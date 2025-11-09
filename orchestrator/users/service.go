package users

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	httplib "github.com/kunal768/cmpe202/http-lib"
	"github.com/kunal768/cmpe202/orchestrator/internal/queue"
	"github.com/kunal768/cmpe202/orchestrator/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type svc struct {
	repo        Repository
	mongoClient *mongo.Client
	publisher   queue.Publisher
}

type Service interface {
	Signup(ctx context.Context, req SignupRequest) (*SignupResponse, error)
	Login(ctx context.Context, req LoginRequest) (*LoginResponse, error)
	RefreshToken(ctx context.Context, req RefreshTokenRequest) (*RefreshTokenResponse, error)
	GetUserByID(ctx context.Context, userID string) (*models.User, error)
	// FetchUndeliveredMessages returns undelivered messages for recipient (optional - requires mongo client)
	FetchUndeliveredMessages(ctx context.Context, recipientID string) ([]map[string]interface{}, error)
}

func NewService(repo Repository, publisher queue.Publisher, mongoClient ...*mongo.Client) Service {
	var mc *mongo.Client
	if len(mongoClient) > 0 {
		mc = mongoClient[0]
	}
	return &svc{
		repo:        repo,
		mongoClient: mc,
		publisher:   publisher,
	}
}

// Signup creates a new user account
func (s *svc) Signup(ctx context.Context, req SignupRequest) (*SignupResponse, error) {
	// Check if user already exists
	existingUser, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err == nil && existingUser != nil {
		return nil, fmt.Errorf("user with email %s already exists", req.Email)
	}

	// Generate user ID
	userID, err := generateUserID()
	if err != nil {
		return nil, fmt.Errorf("failed to generate user ID: %w", err)
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	now := time.Now()

	// Create user
	user := &models.User{
		UserId:   userID,
		UserName: req.UserName,
		Email:    req.Email,
		Role:     models.USER, // Default role
		Contact: models.Contact{
			Email: req.Email,
			Phone: req.Phone,
		},
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Create user authentication
	userAuth := &models.UserAuth{
		UserId:    userID,
		Password:  string(hashedPassword),
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Save user to database
	if err := s.repo.CreateUser(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Save user authentication to database
	if err := s.repo.CreateUserAuth(ctx, userAuth); err != nil {
		return nil, fmt.Errorf("failed to create user authentication: %w", err)
	}

	// Generate access token
	accessToken, err := httplib.GenerateJWT(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	// Generate refresh token
	refreshToken, err := httplib.GenerateRefreshToken(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Create user login authentication
	userLoginAuth := &models.UserLoginAuth{
		UserId:       userID,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    now.Add(24 * time.Hour), // Token expires in 24 hours
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	// Create login auth record
	if err := s.repo.CreateUserLoginAuth(ctx, userLoginAuth); err != nil {
		return nil, fmt.Errorf("failed to create user login authentication: %w", err)
	}

	return &SignupResponse{
		Message:      "User created successfully",
		Token:        accessToken,
		RefreshToken: refreshToken,
		User:         *user,
	}, nil
}

// Login authenticates a user and returns a JWT token
func (s *svc) Login(ctx context.Context, req LoginRequest) (*LoginResponse, error) {
	// Get user by email
	user, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("invalid email or password")
	}

	// Get user authentication
	userAuth, err := s.repo.GetUserAuthByUserID(ctx, user.UserId)
	if err != nil {
		return nil, fmt.Errorf("invalid email or password")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(userAuth.Password), []byte(req.Password)); err != nil {
		return nil, fmt.Errorf("invalid email or password")
	}

	// Generate access token
	accessToken, err := httplib.GenerateJWT(user.UserId)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	// Generate new refresh token
	refreshToken, err := httplib.GenerateRefreshToken(user.UserId)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Create or update user login authentication
	now := time.Now()
	userLoginAuth := &models.UserLoginAuth{
		UserId:       user.UserId,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    now.Add(24 * time.Hour), // Token expires in 24 hours
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	// Check if login auth already exists for this user
	existingLoginAuth, err := s.repo.GetUserLoginAuthByUserID(ctx, user.UserId)
	if err != nil {
		// Create new login auth
		if err := s.repo.CreateUserLoginAuth(ctx, userLoginAuth); err != nil {
			return nil, fmt.Errorf("failed to create user login authentication: %w", err)
		}
	} else {
		// Update existing login auth
		userLoginAuth.CreatedAt = existingLoginAuth.CreatedAt
		if err := s.repo.UpdateUserLoginAuth(ctx, userLoginAuth); err != nil {
			return nil, fmt.Errorf("failed to update user login authentication: %w", err)
		}
	}

	return &LoginResponse{
		Message:      "Login successful",
		Token:        accessToken,
		RefreshToken: refreshToken,
		User:         *user,
	}, nil
}

// RefreshToken handles token refresh
func (s *svc) RefreshToken(ctx context.Context, req RefreshTokenRequest) (*RefreshTokenResponse, error) {
	// Validate refresh token
	_, err := httplib.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	// Get user login auth by refresh token from database
	userLoginAuth, err := s.repo.GetUserLoginAuthByRefreshToken(ctx, req.RefreshToken)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	// Check if token is expired
	if time.Now().After(userLoginAuth.ExpiresAt) {
		return nil, fmt.Errorf("refresh token expired")
	}

	// Get user details
	user, err := s.repo.GetUserByID(ctx, userLoginAuth.UserId)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Generate new access token
	accessToken, err := httplib.GenerateJWT(user.UserId)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	// Generate new refresh token
	newRefreshToken, err := httplib.GenerateRefreshToken(user.UserId)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Update user login authentication
	now := time.Now()
	userLoginAuth.AccessToken = accessToken
	userLoginAuth.RefreshToken = newRefreshToken
	userLoginAuth.ExpiresAt = now.Add(24 * time.Hour) // Token expires in 24 hours
	userLoginAuth.UpdatedAt = now

	if err := s.repo.UpdateUserLoginAuth(ctx, userLoginAuth); err != nil {
		return nil, fmt.Errorf("failed to update user login authentication: %w", err)
	}

	return &RefreshTokenResponse{
		Message:      "Token refreshed successfully",
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		User:         *user,
	}, nil
}

// GetUserByID retrieves a user by ID
func (s *svc) GetUserByID(ctx context.Context, userID string) (*models.User, error) {
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	return user, nil
}

// FetchUndeliveredMessages returns undelivered messages for a recipient using mongo client if available
func (s *svc) FetchUndeliveredMessages(ctx context.Context, recipientID string) ([]map[string]interface{}, error) {
	if s.mongoClient == nil {
		return nil, fmt.Errorf("mongo client not configured")
	}
	coll := s.mongoClient.Database("chatdb").Collection("chatmessages")
	filter := bson.M{"recipientId": recipientID, "status": "UNDELIVERED"}
	cur, err := coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	var results []map[string]interface{}
	now := time.Now()
	// Skip messages that were updated within the last 2 seconds to prevent race conditions
	// This gives chat-consumer time to update status from UNDELIVERED to DELIVERED
	recentThreshold := now.Add(-2 * time.Second)
	for cur.Next(ctx) {
		var doc map[string]interface{}
		if err := cur.Decode(&doc); err != nil {
			continue
		}
		// Check if message was recently updated (might be in process of being delivered)
		// updatedAt can be stored as primitive.DateTime or time.Time depending on MongoDB driver version
		var updatedTime time.Time
		if updatedAt, ok := doc["updatedAt"].(primitive.DateTime); ok {
			updatedTime = updatedAt.Time()
		} else if updatedAt, ok := doc["updatedAt"].(time.Time); ok {
			updatedTime = updatedAt
		} else {
			// If we can't determine updatedAt, include the message (better to republish than skip)
			results = append(results, doc)
			continue
		}
		
		if updatedTime.After(recentThreshold) {
			// Skip recently updated messages to prevent republishing messages that are being processed
			fmt.Printf("[TIMING] [Orchestrator] Skipping message %v - updated %v ago (likely being processed by chat-consumer)\n", doc["messageId"], time.Since(updatedTime))
			continue
		}
		results = append(results, doc)
	}

	// Best-effort: republish undelivered messages to the queue if publisher is configured
	// Normalize message format to match what chat-consumer expects (not the full MongoDB document)
	if s.publisher != nil && len(results) > 0 {
		republishStart := time.Now()
		fmt.Printf("[TIMING] [Orchestrator] Starting republish of %d undelivered messages for recipient %s at %v\n", len(results), recipientID, republishStart)
		for i, m := range results {
			// Extract only the fields that chat-consumer expects
			// Remove MongoDB-specific fields like _id, status, createdAt, updatedAt
			messageForQueue := map[string]interface{}{
				"messageId":   m["messageId"],
				"senderId":    m["senderId"],
				"recipientId": m["recipientId"],
				"content":     m["content"],
				"timestamp":   m["timestamp"],
				"type":        m["type"],
			}
			
			msgStart := time.Now()
			b, err := json.Marshal(messageForQueue)
			if err != nil {
				// skip malformed document
				fmt.Printf("[TIMING] [Orchestrator] failed to marshal undelivered message %v for recipient %s: %v (took %v)\n", m["messageId"], recipientID, err, time.Since(msgStart))
				continue
			}
			// publish with context but do not fail the whole operation if publish fails
			publishStart := time.Now()
			if err := s.publisher.Publish(ctx, b); err != nil {
				// Log publish failure using fmt for minimal dependencies
				fmt.Printf("[TIMING] [Orchestrator] failed to publish undelivered message %v for recipient %s: %v (took %v)\n", m["messageId"], recipientID, err, time.Since(publishStart))
			} else {
				publishDuration := time.Since(publishStart)
				fmt.Printf("[TIMING] [Orchestrator] Republished undelivered message %v to queue for recipient %s (message %d/%d, took %v)\n", m["messageId"], recipientID, i+1, len(results), publishDuration)
			}
		}
		republishDuration := time.Since(republishStart)
		fmt.Printf("[TIMING] [Orchestrator] Completed republish of %d messages for recipient %s (total time: %v)\n", len(results), recipientID, republishDuration)
	}

	return results, nil
}

// generateUserID generates a unique user ID
func generateUserID() (string, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	return id.String(), nil
}
