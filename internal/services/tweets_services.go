package services

import (
	"context"
	"github.com/google/uuid"
	"github.com/juanignaciorc/microbloggin-pltf/internal/domain"
)

type TweetService interface {
	CreateTweet(ctx context.Context, userID uuid.UUID, message string) (domain.Tweet, error)
}
