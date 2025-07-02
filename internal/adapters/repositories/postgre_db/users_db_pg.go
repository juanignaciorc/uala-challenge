package postgre_db

import (
	"context"
	"github.com/google/uuid"
	"github.com/juanignaciorc/microbloggin-pltf/internal/domain"
	"log"
)

/**
 * UsersPGRepository implements port.UsersRepository interface
 * and provides an access to the postgres database
 */
type UsersPGRepository struct {
	db *DB
}

// NewUserRepository creates a new user repository instance
func NewUserRepository(db *DB) *UsersPGRepository {
	return &UsersPGRepository{
		db,
	}
}

func (ur *UsersPGRepository) CreateUser(ctx context.Context, user domain.User) (domain.User, error) {
	user.ID = uuid.New()

	result, err := ur.db.connPool.Exec(ctx, "INSERT INTO users (id, name, email) VALUES ($1, $2, $3)", user.ID, user.Name, user.Email)
	if err != nil {
		return domain.User{}, err
	}

	rowsAffected := result.RowsAffected()
	if err != nil {
		return domain.User{}, err
	}

	if rowsAffected != 1 {
		log.Printf("Expected to affect 1 row, affected %d", rowsAffected)
		return domain.User{}, err
	}

	return user, nil
}

func (ur *UsersPGRepository) GetUser(ctx context.Context, id uuid.UUID) (domain.User, error) {
	var user domain.User

	err := ur.db.connPool.QueryRow(ctx, "SELECT id, name, email FROM users WHERE id = $1", id).Scan(&user.ID, &user.Name, &user.Email)
	if err != nil {
		return domain.User{}, err
	}

	return user, nil
}

func (ur *UsersPGRepository) FollowUser(ctx context.Context, userID uuid.UUID, followedID uuid.UUID) error {
	_, err := ur.db.connPool.Exec(ctx, "INSERT INTO followers (user_id, followed_id) VALUES ($1, $2)", userID, followedID)
	if err != nil {
		return err
	}

	return nil
}

func (ur *UsersPGRepository) GetUserTimeline(ctx context.Context, userID uuid.UUID) ([]domain.Tweet, error) {
	var userTimeLine []domain.Tweet

	followedUsers, err := ur.GetFollowedUserIDs(ctx, userID)
	if err != nil {
		return nil, err
	}

	for _, followedUser := range followedUsers {
		rows, err := ur.db.connPool.Query(ctx, "SELECT id, user_id, message, created_at FROM tweets WHERE user_id = $1", followedUser)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		for rows.Next() {
			var tweet domain.Tweet
			if err := rows.Scan(&tweet.ID, &tweet.UserID, &tweet.Message); err != nil {
				return nil, err
			}
			userTimeLine = append(userTimeLine, tweet)
		}

		if err := rows.Err(); err != nil {
			return nil, err
		}
	}

	return userTimeLine, nil
}

func (ur *UsersPGRepository) GetFollowedUserIDs(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error) {
	var followedUserIDs []uuid.UUID

	rows, err := ur.db.connPool.Query(ctx, "SELECT followed_id FROM followers WHERE user_id = $1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var followedID uuid.UUID
		if err := rows.Scan(&followedID); err != nil {
			return nil, err
		}
		followedUserIDs = append(followedUserIDs, followedID)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return followedUserIDs, nil
}

func (ur *UsersPGRepository) GetFollowerIDs(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error) {
	var followerIDs []uuid.UUID

	rows, err := ur.db.connPool.Query(ctx, "SELECT user_id FROM followers WHERE followed_id = $1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var followerID uuid.UUID
		if err := rows.Scan(&followerID); err != nil {
			return nil, err
		}
		followerIDs = append(followerIDs, followerID)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return followerIDs, nil
}
