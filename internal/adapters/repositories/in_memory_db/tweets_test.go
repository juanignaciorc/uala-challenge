package in_memory_db

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/juanignaciorc/microbloggin-pltf/internal/domain"
	"github.com/stretchr/testify/assert"
	"testing"
)

const uuidMock = "77dae0ef-658c-44c6-803f-f849854a7033"

func TestInMemoryDB_CreateTweet(t *testing.T) {
	// Test cases
	tests := []struct {
		name    string
		setup   func(*InMemoryDB)
		tweet   domain.Tweet
		wantErr bool
	}{
		{
			name: "Successfully create tweet",
			setup: func(db *InMemoryDB) {
				// Create a user first
				user := domain.User{
					ID:    uuid.MustParse(uuidMock),
					Name:  "testuser",
					Email: "test@example.com",
				}
				userBytes, _ := json.Marshal(user)
				db.data[user.ID] = userBytes
			},
			tweet: domain.Tweet{
				ID:      uuid.MustParse(uuidMock),
				UserID:  uuid.MustParse(uuidMock),
				Message: "test tweet",
			},
			wantErr: false,
		},
		{
			name:  "Fail when user doesn't exist",
			setup: func(db *InMemoryDB) {},
			tweet: domain.Tweet{
				ID:      uuid.New(),
				UserID:  uuid.New(), // This user doesn't exist
				Message: "test tweet",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new InMemoryDB instance for each test
			db := NewInMemoryDB()

			// Setup the test case
			tt.setup(db)

			// Execute the method being tested
			gotTweet, err := db.CreateTweet(context.Background(), tt.tweet)

			// Assert the results
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.tweet.ID, gotTweet.ID)
				assert.Equal(t, tt.tweet.UserID, gotTweet.UserID)
				assert.Equal(t, tt.tweet.Message, gotTweet.Message)

				// Verify the tweet was actually stored in the user's tweets
				user, err := db.GetUser(context.Background(), tt.tweet.UserID)
				assert.NoError(t, err)
				assert.Contains(t, user.Tweets, tt.tweet)
			}
		})
	}
}
