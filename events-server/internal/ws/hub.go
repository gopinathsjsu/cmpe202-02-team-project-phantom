package ws

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/kunal768/cmpe202/events-server/internal/delivery"
)

type Hub struct {
	mu              sync.RWMutex
	clients         map[string]*Client
	subscriber      delivery.MessageSubscriber
	recentMessages  map[string]time.Time // Track recently sent messages by messageId to prevent duplicates
	recentMessagesMu sync.RWMutex        // Mutex for recentMessages map
}

func NewHub(subscriber delivery.MessageSubscriber) *Hub {
	return &Hub{
		clients:        make(map[string]*Client),
		subscriber:     subscriber,
		recentMessages: make(map[string]time.Time), // Initialize map to prevent nil map panic
	}
}

func (h *Hub) Register(c *Client) {
	h.mu.Lock()
	h.clients[c.ID] = c
	h.mu.Unlock()

	// Subscribe to messages for this user and wait for confirmation
	if h.subscriber != nil {
		// Channel to signal subscription confirmation
		confirmed := make(chan struct{})
		
		// Subscribe with confirmation callback
		if err := h.subscriber.Subscribe(context.Background(), c.ID, h.handleMessage, func() {
			// Subscription confirmed - signal that we're ready
			close(confirmed)
		}); err != nil {
			log.Printf("Failed to subscribe to messages for user %s: %v", c.ID, err)
			return
		}

		// Wait for subscription confirmation (with timeout)
		waitStart := time.Now()
		select {
		case <-confirmed:
			waitDuration := time.Since(waitStart)
			log.Printf("[TIMING] [%s] Subscription confirmed for user %s, ready to receive messages (waited %v)", c.ID, c.ID, waitDuration)
		case <-time.After(5 * time.Second):
			waitDuration := time.Since(waitStart)
			log.Printf("[TIMING] [%s] Warning: Subscription confirmation timeout for user %s (proceeding anyway, waited %v)", c.ID, c.ID, waitDuration)
		}
	}
}

func (h *Hub) Unregister(userID string) {
	h.mu.Lock()
	delete(h.clients, userID)
	h.mu.Unlock()

	// Unsubscribe from messages for this user
	if h.subscriber != nil {
		if err := h.subscriber.Unsubscribe(context.Background(), userID); err != nil {
			log.Printf("Failed to unsubscribe from messages for user %s: %v", userID, err)
		}
	}
}

func (h *Hub) Get(userID string) (*Client, bool) {
	h.mu.RLock()
	c, ok := h.clients[userID]
	h.mu.RUnlock()
	return c, ok
}

// SendMessageToUser sends a message to a specific user via WebSocket if they are connected
func (h *Hub) SendMessageToUser(userID string, message []byte) error {
	client, exists := h.Get(userID)
	if !exists {
		return fmt.Errorf("user %s is not connected", userID)
	}

	return client.SendMessage(message)
}

// handleMessage processes incoming messages from Redis pub/sub
func (h *Hub) handleMessage(msg []byte) error {
	hubReceiveTime := time.Now()
	log.Printf("[TIMING] [Hub] Received message from Redis pub/sub at %v: %s", hubReceiveTime, string(msg))
	
	// First, check if this is a notification message
	var notificationCheck struct {
		Type        string `json:"type"`
		SubType     string `json:"subType"`
		Count       int    `json:"count"`
		RecipientID string `json:"recipientId"`
	}
	
	if err := json.Unmarshal(msg, &notificationCheck); err == nil && notificationCheck.Type == "notification" {
		// This is a notification message
		log.Printf("[Hub] Processing notification message for user %s (count: %d)", notificationCheck.RecipientID, notificationCheck.Count)
		client, exists := h.Get(notificationCheck.RecipientID)
		if !exists {
			log.Printf("[Hub] Client %s not found for notification delivery", notificationCheck.RecipientID)
			return nil
		}
		// Send as notification
		notifMsg := NotificationMessage{
			Type:    notificationCheck.Type,
			SubType: notificationCheck.SubType,
			Count:   notificationCheck.Count,
		}
		if err := client.SendNotification(notifMsg); err != nil {
			log.Printf("[Hub] Failed to send notification to user %s: %v", notificationCheck.RecipientID, err)
			return err
		}
		log.Printf("[Hub] Notification sent to user %s via WebSocket", notificationCheck.RecipientID)
		return nil
	}

	// Parse the message to get recipient ID (regular message)
	var messageData struct {
		MessageID   string `json:"messageId"`
		RecipientID string `json:"recipientId"`
	}

	if err := json.Unmarshal(msg, &messageData); err != nil {
		log.Printf("[Hub] Failed to unmarshal message: %v", err)
		return fmt.Errorf("failed to unmarshal message: %w", err)
	}

	log.Printf("[Hub] Processing regular message %s for user %s", messageData.MessageID, messageData.RecipientID)

	// Deduplication: Check if we've recently sent this message (within last 5 seconds)
	// This prevents duplicate delivery when user reconnects and old subscription is still active
	h.recentMessagesMu.RLock()
	lastSent, recentlySent := h.recentMessages[messageData.MessageID]
	h.recentMessagesMu.RUnlock()
	
	if recentlySent && time.Since(lastSent) < 5*time.Second {
		log.Printf("[Hub] Skipping duplicate message %s for user %s (sent %v ago)", messageData.MessageID, messageData.RecipientID, time.Since(lastSent))
		return nil // Not an error, just a duplicate
	}

	// Find the client
	client, exists := h.Get(messageData.RecipientID)
	if !exists {
		log.Printf("[Hub] Client %s not found for message delivery (may have disconnected)", messageData.RecipientID)
		return nil // Not an error, client might have disconnected
	}

	// Send the message to the client
	sendStart := time.Now()
	if err := client.SendMessage(msg); err != nil {
		log.Printf("[TIMING] [Hub] Failed to send message %s to user %s via WebSocket: %v (took %v, total from hub receive: %v)", messageData.MessageID, messageData.RecipientID, err, time.Since(sendStart), time.Since(hubReceiveTime))
		return err
	}
	
	// Mark message as recently sent to prevent duplicates
	h.recentMessagesMu.Lock()
	h.recentMessages[messageData.MessageID] = time.Now()
	h.recentMessagesMu.Unlock()
	
	// Clean up old entries periodically (messages older than 10 seconds)
	// This prevents memory leak from the deduplication map
	go func() {
		h.recentMessagesMu.Lock()
		defer h.recentMessagesMu.Unlock()
		cutoff := time.Now().Add(-10 * time.Second)
		for msgID, sentTime := range h.recentMessages {
			if sentTime.Before(cutoff) {
				delete(h.recentMessages, msgID)
			}
		}
	}()
	
	sendDuration := time.Since(sendStart)
	totalDuration := time.Since(hubReceiveTime)
	log.Printf("[TIMING] [Hub] Message %s sent to user %s via WebSocket (send took %v, total from hub receive: %v)", messageData.MessageID, messageData.RecipientID, sendDuration, totalDuration)
	return nil
}
