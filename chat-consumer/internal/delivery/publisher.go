package delivery

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

// MessagePublisher defines the interface for publishing messages
type MessagePublisher interface {
	PublishToUser(ctx context.Context, userID string, message []byte) (int64, error)
	Close() error
}

// RedisMessagePublisher implements MessagePublisher using Redis pub/sub
type RedisMessagePublisher struct {
	client *redis.Client
}

// NewRedisMessagePublisher creates a new Redis message publisher
func NewRedisMessagePublisher(addr, password string, db int) *RedisMessagePublisher {
	return &RedisMessagePublisher{
		client: redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: password,
			DB:       db,
		}),
	}
}

// PublishToUser publishes a message to a specific user's channel
// Returns the number of subscribers that received the message
func (r *RedisMessagePublisher) PublishToUser(ctx context.Context, userID string, message []byte) (int64, error) {
	channel := fmt.Sprintf("user:%s:messages", userID)

	result := r.client.Publish(ctx, channel, message)
	if result.Err() != nil {
		return 0, fmt.Errorf("failed to publish message to user %s: %w", userID, result.Err())
	}

	subscribers := result.Val()
	publishTime := time.Now()
	log.Printf("[TIMING] [Publisher] Published message to user %s on channel %s at %v (subscribers: %d)", userID, channel, publishTime, subscribers)

	// If no subscribers, the user is not online or events-server is not subscribed
	if subscribers == 0 {
		log.Printf("[TIMING] [Publisher] WARNING: No subscribers for user %s - events-server may not be subscribed for this userId at %v", userID, publishTime)
	} else {
		log.Printf("[TIMING] [Publisher] Success: %d subscriber(s) received message for user %s at %v", subscribers, userID, publishTime)
	}

	return subscribers, nil
}

// Close closes the Redis client
func (r *RedisMessagePublisher) Close() error {
	if err := r.client.Close(); err != nil {
		return fmt.Errorf("failed to close Redis client: %w", err)
	}
	log.Println("Redis message publisher closed")
	return nil
}
