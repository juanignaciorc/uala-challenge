package handlers

import (
	"github.com/google/uuid"
	"github.com/juanignaciorc/microbloggin-pltf/internal/domain"
)

// Response DTOs to avoid exposing domain objects

type UserResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
	// Note: Email is intentionally omitted for privacy
}

type UserDetailResponse struct {
	ID             uuid.UUID `json:"id"`
	Name           string    `json:"name"`
	Email          string    `json:"email"` // Only included in detailed responses when appropriate
	FollowersCount int       `json:"followers_count"`
	FollowingCount int       `json:"following_count"`
	TweetsCount    int       `json:"tweets_count"`
}

type TweetResponse struct {
	ID      uuid.UUID    `json:"id"`
	Message string       `json:"message"`
	User    UserResponse `json:"user"` // Nested user info without sensitive data
}

// Error response structures for better error formatting

type ErrorResponse struct {
	Error   string            `json:"error"`
	Code    string            `json:"code,omitempty"`
	Details map[string]string `json:"details,omitempty"`
}

type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Helper functions to convert domain objects to response DTOs

func ToUserResponse(user domain.User) UserResponse {
	return UserResponse{
		ID:   user.ID,
		Name: user.Name,
		// Email is intentionally omitted for privacy
	}
}

func ToUserDetailResponse(user domain.User) UserDetailResponse {
	return UserDetailResponse{
		ID:             user.ID,
		Name:           user.Name,
		Email:          user.Email,
		FollowersCount: len(user.Followers),
		FollowingCount: len(user.Follwing), // Note: keeping the typo from domain for now
		TweetsCount:    len(user.Tweets),
	}
}

func ToTweetResponseSimple(tweet domain.Tweet) TweetResponse {
	return TweetResponse{
		ID:      tweet.ID,
		Message: tweet.Message,
		// User info will be empty in this case
	}
}

// Error helper functions

func NewErrorResponse(message string) ErrorResponse {
	return ErrorResponse{
		Error: message,
	}
}

func NewErrorResponseWithCode(message, code string) ErrorResponse {
	return ErrorResponse{
		Error: message,
		Code:  code,
	}
}

func NewErrorResponseWithDetails(message string, details map[string]string) ErrorResponse {
	return ErrorResponse{
		Error:   message,
		Details: details,
	}
}

func NewSuccessResponse(message string, data interface{}) SuccessResponse {
	return SuccessResponse{
		Message: message,
		Data:    data,
	}
}
