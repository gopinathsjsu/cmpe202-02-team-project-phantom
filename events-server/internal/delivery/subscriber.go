package delivery

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

// MessageSubscriber handles Redis pub/sub for message delivery
type MessageSubscriber interface {
	Subscribe(ctx context.Context, userID string, messageHandler func([]byte) error, onConfirmed func()) error
	Unsubscribe(ctx context.Context, userID string) error
	Close() error
}

// RedisMessageSubscriber implements MessageSubscriber using Redis pub/sub
type RedisMessageSubscriber struct {
	client   *redis.Client
	subs     map[string]*redis.PubSub
	mu       sync.RWMutex
	closed   bool
	stopChan chan struct{}
}

// NewRedisMessageSubscriber creates a new Redis message subscriber
func NewRedisMessageSubscriber(addr, password string, db int) *RedisMessageSubscriber {
	return &RedisMessageSubscriber{
		client: redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: password,
			DB:       db,
		}),
		subs:     make(map[string]*redis.PubSub),
		stopChan: make(chan struct{}),
	}
}

// Subscribe starts listening to messages for a specific user
// onConfirmed callback is called when subscription is confirmed (after pubsub.Receive succeeds)
func (r *RedisMessageSubscriber) Subscribe(ctx context.Context, userID string, messageHandler func([]byte) error, onConfirmed func()) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.closed {
		return fmt.Errorf("subscriber is closed")
	}

	// Check if already subscribed - if so, unsubscribe the old one first
	if oldPubsub, exists := r.subs[userID]; exists {
		log.Printf("[Subscriber] User %s already has active subscription, closing old subscription before creating new one", userID)
		if err := oldPubsub.Close(); err != nil {
			log.Printf("[Subscriber] Error closing old subscription for user %s: %v", userID, err)
		}
		delete(r.subs, userID)
		// Give a brief moment for the old goroutine to exit
		time.Sleep(50 * time.Millisecond)
	}

	channel := fmt.Sprintf("user:%s:messages", userID)
	pubsub := r.client.Subscribe(ctx, channel)
	r.subs[userID] = pubsub

	// Start goroutine to handle messages
	go r.handleMessages(ctx, userID, pubsub, messageHandler, onConfirmed)

	log.Printf("Subscribing to messages for user %s on channel %s (waiting for confirmation...)", userID, channel)
	return nil
}

// Unsubscribe stops listening to messages for a specific user
func (r *RedisMessageSubscriber) Unsubscribe(ctx context.Context, userID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.closed {
		return nil
	}

	pubsub, exists := r.subs[userID]
	if !exists {
		return nil // Already unsubscribed
	}

	if err := pubsub.Close(); err != nil {
		log.Printf("Error closing subscription for user %s: %v", userID, err)
	}

	delete(r.subs, userID)
	log.Printf("Unsubscribed from messages for user %s", userID)
	return nil
}

// Close closes all subscriptions and the Redis client
func (r *RedisMessageSubscriber) Close() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.closed {
		return nil
	}

	r.closed = true
	close(r.stopChan)

	// Close all subscriptions
	for userID, pubsub := range r.subs {
		if err := pubsub.Close(); err != nil {
			log.Printf("Error closing subscription for user %s: %v", userID, err)
		}
	}
	r.subs = make(map[string]*redis.PubSub)

	// Close Redis client
	if err := r.client.Close(); err != nil {
		return fmt.Errorf("failed to close Redis client: %w", err)
	}

	log.Println("Redis message subscriber closed")
	return nil
}

// handleMessages processes incoming messages from Redis pub/sub
func (r *RedisMessageSubscriber) handleMessages(ctx context.Context, userID string, pubsub *redis.PubSub, messageHandler func([]byte) error, onConfirmed func()) {
	defer func() {
		if err := pubsub.Close(); err != nil {
			log.Printf("Error closing pubsub for user %s: %v", userID, err)
		}
	}()

	// Wait for subscription confirmation
	confirmStart := time.Now()
	log.Printf("[TIMING] [%s] Waiting for Redis subscription confirmation at %v", userID, confirmStart)
	_, err := pubsub.Receive(ctx)
	if err != nil {
		log.Printf("[TIMING] [%s] Failed to receive subscription confirmation for user %s: %v (took %v)", userID, userID, err, time.Since(confirmStart))
		return
	}

	// Subscription confirmed - call the callback
	confirmDuration := time.Since(confirmStart)
	log.Printf("[TIMING] [%s] Subscription confirmed for user %s on channel user:%s:messages at %v (took %v)", userID, userID, userID, time.Now(), confirmDuration)
	if onConfirmed != nil {
		onConfirmed()
	}

	// Channel to receive messages
	chStart := time.Now()
	log.Printf("[TIMING] [%s] Getting pubsub.Channel() at %v", userID, chStart)
	ch := pubsub.Channel()
	chDuration := time.Since(chStart)
	log.Printf("[TIMING] [%s] pubsub.Channel() ready at %v (took %v), now listening for messages", userID, time.Now(), chDuration)

	for {
		select {
		case <-r.stopChan:
			return
		case <-ctx.Done():
			return
		case msg, ok := <-ch:
			if !ok {
				log.Printf("[TIMING] [%s] Message channel closed for user %s at %v", userID, userID, time.Now())
				return
			}

			if msg == nil {
				continue
			}

			receiveTime := time.Now()
			log.Printf("[TIMING] [%s] Received message from Redis channel at %v", userID, receiveTime)

			// Parse the message payload
			var messageData struct {
				MessageID   string    `json:"messageId"`
				SenderID    string    `json:"senderId"`
				RecipientID string    `json:"recipientId"`
				Content     string    `json:"content"`
				Timestamp   time.Time `json:"timestamp"`
				Type        string    `json:"type"`
			}

			parseStart := time.Now()
			if err := json.Unmarshal([]byte(msg.Payload), &messageData); err != nil {
				log.Printf("[TIMING] [%s] Failed to unmarshal message for user %s: %v (took %v)", userID, userID, err, time.Since(parseStart))
				continue
			}
			parseDuration := time.Since(parseStart)
			log.Printf("[TIMING] [%s] Parsed message %s for user %s (took %v)", userID, messageData.MessageID, userID, parseDuration)

			// Call the message handler
			handlerStart := time.Now()
			if err := messageHandler([]byte(msg.Payload)); err != nil {
				log.Printf("[TIMING] [%s] Error handling message %s for user %s: %v (took %v)", userID, messageData.MessageID, userID, err, time.Since(handlerStart))
			} else {
				handlerDuration := time.Since(handlerStart)
				log.Printf("[TIMING] [%s] Successfully handled message %s for user %s (took %v, total time from receive: %v)", userID, messageData.MessageID, userID, handlerDuration, time.Since(receiveTime))
			}
		}
	}
}
