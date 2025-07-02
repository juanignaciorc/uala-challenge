package services

import (
	"context"
	"github.com/google/uuid"
	"github.com/juanignaciorc/microbloggin-pltf/internal/domain"
	ports "github.com/juanignaciorc/microbloggin-pltf/internal/ports/repositories"
)

type tweetsServiceImpl struct {
	tweetsRepository ports.TweetRepository
}

// NewTweetsService creates a new TweetService instance.
func NewTweetsService(tweetsRepository ports.TweetRepository) TweetService {
	return &tweetsServiceImpl{
		tweetsRepository: tweetsRepository,
	}
}

func (s *tweetsServiceImpl) CreateTweet(ctx context.Context, userID uuid.UUID, message string) (domain.Tweet, error) {
	tweet := domain.Tweet{
		UserID:  userID,
		Message: message,
	}
	tw, err := s.tweetsRepository.CreateTweet(ctx, tweet)
	if err != nil {
		return domain.Tweet{}, err
	}

	return tw, nil
}
