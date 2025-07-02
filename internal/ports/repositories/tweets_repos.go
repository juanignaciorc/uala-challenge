package ports

import (
	"context"
	"github.com/juanignaciorc/microbloggin-pltf/internal/domain"
)

type TweetRepository interface {
	CreateTweet(ctx context.Context, tweet domain.Tweet) (domain.Tweet, error)
}
