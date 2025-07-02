package ports

import (
	"context"
	"github.com/google/uuid"
	"github.com/juanignaciorc/microbloggin-pltf/internal/domain"
)

type UsersRepository interface {
	CreateUser(ctx context.Context, user domain.User) (domain.User, error)
	GetUser(ctx context.Context, id uuid.UUID) (domain.User, error)
	FollowUser(ctx context.Context, userID uuid.UUID, followedID uuid.UUID) error
	GetUserTimeline(ctx context.Context, userID uuid.UUID) ([]domain.Tweet, error)
}
