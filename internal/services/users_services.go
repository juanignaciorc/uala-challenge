package services

import (
	"context"
	"github.com/google/uuid"
	"github.com/juanignaciorc/microbloggin-pltf/internal/domain"
)

type UserService interface {
	CreateUser(ctx context.Context, name, mail string) (domain.User, error)
	GetUser(ctx context.Context, id uuid.UUID) (domain.User, error)
	FollowUser(ctx context.Context, userID, followedID uuid.UUID) error
	GetUserTimeline(ctx context.Context, userID uuid.UUID) ([]domain.Tweet, error)
}
