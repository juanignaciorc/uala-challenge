package in_memory_db

import (
	"context"
	"github.com/google/uuid"
	"github.com/juanignaciorc/microbloggin-pltf/internal/domain"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInMemoryDB_CreateUser(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(*InMemoryDB)
		user    domain.User
		wantErr bool
	}{
		{
			name:  "Successfully create user",
			setup: func(db *InMemoryDB) {},
			user: domain.User{
				Name:  "testuser",
				Email: "test@example.com",
			},
			wantErr: false,
		},
		{
			name:  "Successfully create user with valid data",
			setup: func(db *InMemoryDB) {},
			user: domain.User{
				Name:  "anotheruser",
				Email: "another@example.com",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := NewInMemoryDB()
			tt.setup(db)

			gotUser, err := db.CreateUser(context.Background(), tt.user)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotEqual(t, uuid.Nil, gotUser.ID) // ID should be generated
				assert.Equal(t, tt.user.Name, gotUser.Name)
				assert.Equal(t, tt.user.Email, gotUser.Email)

				// Verify user was actually stored
				storedUser, err := db.GetUser(context.Background(), gotUser.ID)
				assert.NoError(t, err)
				assert.Equal(t, gotUser, storedUser)
			}
		})
	}
}

func TestInMemoryDB_GetUser(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(*InMemoryDB) uuid.UUID
		wantErr bool
	}{
		{
			name: "Successfully get user",
			setup: func(db *InMemoryDB) uuid.UUID {
				user := domain.User{
					Name:  "testuser",
					Email: "test@example.com",
				}
				createdUser, _ := db.CreateUser(context.Background(), user)
				return createdUser.ID
			},
			wantErr: false,
		},
		{
			name: "User not found",
			setup: func(db *InMemoryDB) uuid.UUID {
				return uuid.New() // Return non-existent user ID
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := NewInMemoryDB()
			userID := tt.setup(db)

			user, err := db.GetUser(context.Background(), userID)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, userID, user.ID)
				assert.NotEmpty(t, user.Name)
				assert.NotEmpty(t, user.Email)
			}
		})
	}
}

func TestInMemoryDB_FollowUser(t *testing.T) {
	tests := []struct {
		name       string
		setup      func(*InMemoryDB) (uuid.UUID, uuid.UUID)
		wantErr    bool
		checkState func(*testing.T, *InMemoryDB, uuid.UUID, uuid.UUID)
	}{
		{
			name: "Successfully follow user",
			setup: func(db *InMemoryDB) (uuid.UUID, uuid.UUID) {
				follower := domain.User{
					Name:  "follower",
					Email: "follower@example.com",
				}
				following := domain.User{
					Name:  "following",
					Email: "following@example.com",
				}
				createdFollower, _ := db.CreateUser(context.Background(), follower)
				createdFollowing, _ := db.CreateUser(context.Background(), following)
				return createdFollower.ID, createdFollowing.ID
			},
			wantErr: false,
			checkState: func(t *testing.T, db *InMemoryDB, followerID, followingID uuid.UUID) {
				follower, err := db.GetUser(context.Background(), followerID)
				assert.NoError(t, err)
				assert.Contains(t, follower.Follwing, followingID)

				following, err := db.GetUser(context.Background(), followingID)
				assert.NoError(t, err)
				assert.Contains(t, following.Followers, followerID)
			},
		},
		{
			name: "Fail when follower doesn't exist",
			setup: func(db *InMemoryDB) (uuid.UUID, uuid.UUID) {
				following := domain.User{
					Name:  "following",
					Email: "following@example.com",
				}
				createdFollowing, _ := db.CreateUser(context.Background(), following)
				return uuid.New(), createdFollowing.ID // Return non-existent follower ID
			},
			wantErr: true,
			checkState: func(t *testing.T, db *InMemoryDB, followerID, followingID uuid.UUID) {
				following, err := db.GetUser(context.Background(), followingID)
				assert.NoError(t, err)
				assert.NotContains(t, following.Followers, followerID)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := NewInMemoryDB()
			followerID, followingID := tt.setup(db)

			err := db.FollowUser(context.Background(), followerID, followingID)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			if tt.checkState != nil {
				tt.checkState(t, db, followerID, followingID)
			}
		})
	}
}

func TestInMemoryDB_GetUserTimeline(t *testing.T) {
	tests := []struct {
		name          string
		setup         func(*InMemoryDB) uuid.UUID
		expectedCount int
		wantErr       bool
	}{
		{
			name: "Successfully get user timeline",
			setup: func(db *InMemoryDB) uuid.UUID {
				// Create a follower user
				follower := domain.User{
					Name:  "follower",
					Email: "follower@example.com",
				}
				createdFollower, _ := db.CreateUser(context.Background(), follower)

				// Create a user to be followed
				followed := domain.User{
					Name:  "followed",
					Email: "followed@example.com",
				}
				createdFollowed, _ := db.CreateUser(context.Background(), followed)

				// Create a tweet for the followed user
				tweet := domain.Tweet{
					ID:      uuid.New(),
					UserID:  createdFollowed.ID,
					Message: "Test tweet",
				}
				_, _ = db.CreateTweet(context.Background(), tweet)

				// Make follower follow the user with tweets
				_ = db.FollowUser(context.Background(), createdFollower.ID, createdFollowed.ID)

				return createdFollower.ID
			},
			expectedCount: 1,
			wantErr:       false,
		},
		{
			name: "Empty timeline",
			setup: func(db *InMemoryDB) uuid.UUID {
				user := domain.User{
					Name:  "testuser",
					Email: "test@example.com",
				}
				createdUser, _ := db.CreateUser(context.Background(), user)
				return createdUser.ID
			},
			expectedCount: 0,
			wantErr:       false,
		},
		{
			name: "User not found",
			setup: func(db *InMemoryDB) uuid.UUID {
				return uuid.New() // Return non-existent user ID
			},
			expectedCount: 0,
			wantErr:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := NewInMemoryDB()
			userID := tt.setup(db)

			timeline, err := db.GetUserTimeline(context.Background(), userID)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, timeline, tt.expectedCount)

				if tt.expectedCount > 0 && len(timeline) > 0 {
					// Timeline contains tweets from followed users, not the user themselves
					assert.NotEqual(t, userID, timeline[0].UserID)
					assert.NotEmpty(t, timeline[0].Message)
				}
			}
		})
	}
}

func TestInMemoryDB_GetFollowedUsers(t *testing.T) {
	tests := []struct {
		name          string
		setup         func(*InMemoryDB) uuid.UUID
		expectedCount int
		wantErr       bool
	}{
		{
			name: "Successfully get followed users",
			setup: func(db *InMemoryDB) uuid.UUID {
				follower := domain.User{
					Name:  "follower",
					Email: "follower@example.com",
				}
				following1 := domain.User{
					Name:  "following1",
					Email: "following1@example.com",
				}
				following2 := domain.User{
					Name:  "following2",
					Email: "following2@example.com",
				}

				createdFollower, _ := db.CreateUser(context.Background(), follower)
				createdFollowing1, _ := db.CreateUser(context.Background(), following1)
				createdFollowing2, _ := db.CreateUser(context.Background(), following2)

				_ = db.FollowUser(context.Background(), createdFollower.ID, createdFollowing1.ID)
				_ = db.FollowUser(context.Background(), createdFollower.ID, createdFollowing2.ID)

				return createdFollower.ID
			},
			expectedCount: 2,
			wantErr:       false,
		},
		{
			name: "User with no followed users",
			setup: func(db *InMemoryDB) uuid.UUID {
				user := domain.User{
					Name:  "loneuser",
					Email: "lone@example.com",
				}
				createdUser, _ := db.CreateUser(context.Background(), user)
				return createdUser.ID
			},
			expectedCount: 0,
			wantErr:       false,
		},
		{
			name: "User not found",
			setup: func(db *InMemoryDB) uuid.UUID {
				return uuid.New() // Return non-existent user ID
			},
			expectedCount: 0,
			wantErr:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := NewInMemoryDB()
			userID := tt.setup(db)

			followedUsers, err := db.GetFollowedUsers(context.Background(), userID)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, followedUsers, tt.expectedCount)
			}
		})
	}
}
