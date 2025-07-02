package postgre_db

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/juanignaciorc/microbloggin-pltf/internal/domain"
	"log"
)

type TweetsPGRepository struct {
	db *DB
}

func NewTweetRepository(db *DB) *TweetsPGRepository {
	return &TweetsPGRepository{
		db,
	}
}

func (tr *TweetsPGRepository) CreateTweet(ctx context.Context, tweet domain.Tweet) (domain.Tweet, error) {
	tweet.ID = uuid.New()

	result, err := tr.db.connPool.Exec(ctx, "INSERT INTO tweets (id, user_id, message) VALUES ($1, $2, $3)", tweet.ID, tweet.UserID, tweet.Message)
	if err != nil {
		return domain.Tweet{}, err
	}

	rowsAffected := result.RowsAffected()
	if err != nil {
		return domain.Tweet{}, err
	}

	if rowsAffected != 1 {
		log.Printf("Expected to affect 1 row, affected %d", rowsAffected)
		return domain.Tweet{}, fmt.Errorf("expected to affect 1 row, affected %d", rowsAffected)
	}

	return tweet, nil
}
