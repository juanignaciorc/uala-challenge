package in_memory_db

import (
	"context"
	"github.com/google/uuid"
	"github.com/juanignaciorc/microbloggin-pltf/internal/domain"
)

type InMemoryDBTweetsInterface interface {
	CreateTweet(ctx context.Context, tweet domain.Tweet) (domain.Tweet, error)
}

type InMemoryDB struct {
	data map[uuid.UUID][]byte
}

func NewInMemoryDB() *InMemoryDB {
	return &InMemoryDB{
		data: make(map[uuid.UUID][]byte),
	}
}
