package in_memory_db

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"

	"github.com/juanignaciorc/microbloggin-pltf/internal/domain"
)

func (db *InMemoryDB) CreateUser(ctx context.Context, user domain.User) (domain.User, error) {
	id := uuid.New()
	user.ID = id

	userBytes, err := json.Marshal(user)
	if err != nil {
		return domain.User{}, err
	}

	db.data[user.ID] = userBytes

	var createdUser domain.User
	err = json.Unmarshal(userBytes, &createdUser)
	return createdUser, err
}

func (db *InMemoryDB) GetUser(ctx context.Context, id uuid.UUID) (domain.User, error) {
	userBytes, ok := db.data[id]
	if !ok {
		return domain.User{}, fmt.Errorf("user with id %v not found", id)
	}

	var user domain.User
	err := json.Unmarshal(userBytes, &user)
	if err != nil {
		return domain.User{}, err
	}

	return user, nil
}

func (db *InMemoryDB) FollowUser(ctx context.Context, userID uuid.UUID, followedID uuid.UUID) error {
	user, err := db.GetUser(ctx, userID)
	if err != nil {
		return err
	}

	followedUser, err := db.GetUser(ctx, followedID)
	if err != nil {
		return err
	}

	followedUser.Followers = append(user.Followers, userID)
	user.Follwing = append(user.Follwing, followedID)

	userBytes, err := json.Marshal(user)
	if err != nil {
		return err
	}

	followedUserBytes, err := json.Marshal(followedUser)
	if err != nil {
		return err
	}

	db.data[userID] = userBytes
	db.data[followedID] = followedUserBytes

	return nil
}

func (db *InMemoryDB) GetUserTimeline(ctx context.Context, userID uuid.UUID) ([]domain.Tweet, error) {
	followedUsers, err := db.GetFollowedUsers(ctx, userID)
	if err != nil {
		return nil, err
	}

	var userTimeline []domain.Tweet
	for _, followedID := range followedUsers {
		user, err := db.GetUser(ctx, followedID)
		if err != nil {
			return nil, err
		}

		userTimeline = append(userTimeline, user.Tweets...)
	}

	return userTimeline, nil
}

func (db *InMemoryDB) GetFollowedUsers(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error) {
	user, err := db.GetUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	return user.Follwing, nil
}
