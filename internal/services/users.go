package services

import (
	"context"
	"github.com/google/uuid"
	"github.com/juanignaciorc/microbloggin-pltf/internal/domain"
	ports "github.com/juanignaciorc/microbloggin-pltf/internal/ports/repositories"
)

type userServiceImpl struct {
	userRepository ports.UsersRepository
}

func NewUserService(userRepository ports.UsersRepository) userServiceImpl {
	return userServiceImpl{
		userRepository: userRepository,
	}
}

func (s userServiceImpl) CreateUser(ctx context.Context, name, mail string) (domain.User, error) {
	user := domain.User{
		Name:  name,
		Email: mail,
	}

	createdUser, err := s.userRepository.CreateUser(ctx, user)
	if err != nil {
		return domain.User{}, err
	}

	return createdUser, nil
}

func (s userServiceImpl) GetUser(ctx context.Context, id uuid.UUID) (domain.User, error) {
	user, err := s.userRepository.GetUser(ctx, id)
	if err != nil {
		return domain.User{}, err
	}

	return user, nil
}

func (s userServiceImpl) FollowUser(ctx context.Context, userID, followedID uuid.UUID) error {
	if err := s.userRepository.FollowUser(ctx, userID, followedID); err != nil {
		return err
	}

	return nil
}

func (s userServiceImpl) GetUserTimeline(ctx context.Context, userID uuid.UUID) ([]domain.Tweet, error) {
	tweets, err := s.userRepository.GetUserTimeline(ctx, userID)
	if err != nil {
		return nil, err
	}

	return tweets, nil
}
