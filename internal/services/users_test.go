package services

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/juanignaciorc/microbloggin-pltf/internal/domain"
	mock_ports "github.com/juanignaciorc/microbloggin-pltf/mocks"
	"go.uber.org/mock/gomock"
	"reflect"
	"testing"
)

func TestUserService_CreateUser(t *testing.T) {
	mockUUID := uuid.New()

	type testCase struct {
		name       string
		inputName  string
		inputEmail string
		mockInput  domain.User
		mockOutput domain.User
		mockErr    error
		expected   domain.User
		wantErr    bool
	}

	tests := []testCase{
		{
			name:       "Success case",
			inputName:  "John Doe",
			inputEmail: "john@example.com",
			mockInput:  domain.User{Name: "John Doe", Email: "john@example.com"},
			mockOutput: domain.User{ID: mockUUID, Name: "John Doe", Email: "john@example.com"},
			mockErr:    nil,
			expected:   domain.User{ID: mockUUID, Name: "John Doe", Email: "john@example.com"},
			wantErr:    false,
		},
		{
			name:       "Repository error",
			inputName:  "Jane Doe",
			inputEmail: "jane@example.com",
			mockInput:  domain.User{Name: "Jane Doe", Email: "jane@example.com"},
			mockOutput: domain.User{},
			mockErr:    errors.New("repository error"),
			expected:   domain.User{},
			wantErr:    true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockCtx := context.Background()
			mockRepo := mock_ports.NewMockUsersRepository(ctrl)
			s := NewUserService(mockRepo)

			mockRepo.
				EXPECT().
				CreateUser(mockCtx, tc.mockInput).
				Return(tc.mockOutput, tc.mockErr)

			got, err := s.CreateUser(mockCtx, tc.inputName, tc.inputEmail)

			if (err != nil) != tc.wantErr {
				t.Errorf("CreateUser() error = %v, wantErr = %v", err, tc.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tc.expected) {
				t.Errorf("CreateUser() got = %v, want = %v", got, tc.expected)
			}
		})
	}
}

func TestUserService_GetUser(t *testing.T) {
	mockUUID := uuid.New()

	type testCase struct {
		name       string
		inputID    uuid.UUID
		mockOutput domain.User
		mockErr    error
		expected   domain.User
		wantErr    bool
	}

	tests := []testCase{
		{
			name:       "Success case",
			inputID:    mockUUID,
			mockOutput: domain.User{ID: mockUUID, Name: "John Doe", Email: "john@example.com"},
			mockErr:    nil,
			expected:   domain.User{ID: mockUUID, Name: "John Doe", Email: "john@example.com"},
			wantErr:    false,
		},
		{
			name:       "Repository error",
			inputID:    mockUUID,
			mockOutput: domain.User{},
			mockErr:    errors.New("user not found"),
			expected:   domain.User{},
			wantErr:    true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockCtx := context.Background()
			mockRepo := mock_ports.NewMockUsersRepository(ctrl)
			s := NewUserService(mockRepo)

			mockRepo.
				EXPECT().
				GetUser(mockCtx, tc.inputID).
				Return(tc.mockOutput, tc.mockErr)

			got, err := s.GetUser(mockCtx, tc.inputID)

			if (err != nil) != tc.wantErr {
				t.Errorf("GetUser() error = %v, wantErr = %v", err, tc.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tc.expected) {
				t.Errorf("GetUser() got = %v, want = %v", got, tc.expected)
			}
		})
	}
}

func TestUserService_FollowUser(t *testing.T) {
	userID := uuid.New()
	followedID := uuid.New()

	type testCase struct {
		name       string
		userID     uuid.UUID
		followedID uuid.UUID
		mockErr    error
		wantErr    bool
	}

	tests := []testCase{
		{
			name:       "Success case",
			userID:     userID,
			followedID: followedID,
			mockErr:    nil,
			wantErr:    false,
		},
		{
			name:       "Repository error",
			userID:     userID,
			followedID: followedID,
			mockErr:    errors.New("follow operation failed"),
			wantErr:    true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockCtx := context.Background()
			mockRepo := mock_ports.NewMockUsersRepository(ctrl)
			s := NewUserService(mockRepo)

			mockRepo.
				EXPECT().
				FollowUser(mockCtx, tc.userID, tc.followedID).
				Return(tc.mockErr)

			err := s.FollowUser(mockCtx, tc.userID, tc.followedID)

			if (err != nil) != tc.wantErr {
				t.Errorf("FollowUser() error = %v, wantErr = %v", err, tc.wantErr)
			}
		})
	}
}

func TestUserService_GetUserTimeline(t *testing.T) {
	mockUUID := uuid.New()
	mockTweets := []domain.Tweet{
		{ID: uuid.New(), UserID: mockUUID, Message: "First tweet"},
		{ID: uuid.New(), UserID: mockUUID, Message: "Second tweet"},
	}

	type testCase struct {
		name       string
		inputID    uuid.UUID
		mockOutput []domain.Tweet
		mockErr    error
		expected   []domain.Tweet
		wantErr    bool
	}

	tests := []testCase{
		{
			name:       "Success case",
			inputID:    mockUUID,
			mockOutput: mockTweets,
			mockErr:    nil,
			expected:   mockTweets,
			wantErr:    false,
		},
		{
			name:       "Repository error",
			inputID:    mockUUID,
			mockOutput: nil,
			mockErr:    errors.New("timeline fetch failed"),
			expected:   nil,
			wantErr:    true,
		},
		{
			name:       "Empty timeline",
			inputID:    mockUUID,
			mockOutput: []domain.Tweet{},
			mockErr:    nil,
			expected:   []domain.Tweet{},
			wantErr:    false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockCtx := context.Background()
			mockRepo := mock_ports.NewMockUsersRepository(ctrl)
			s := NewUserService(mockRepo)

			mockRepo.
				EXPECT().
				GetUserTimeline(mockCtx, tc.inputID).
				Return(tc.mockOutput, tc.mockErr)

			got, err := s.GetUserTimeline(mockCtx, tc.inputID)

			if (err != nil) != tc.wantErr {
				t.Errorf("GetUserTimeline() error = %v, wantErr = %v", err, tc.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tc.expected) {
				t.Errorf("GetUserTimeline() got = %v, want = %v", got, tc.expected)
			}
		})
	}
}