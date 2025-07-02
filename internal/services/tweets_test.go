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

func TestTweetsService_CreateTweet(t *testing.T) {
	mockUUID := uuid.New()

	type testCase struct {
		name       string
		mockInput  domain.Tweet
		mockOutput domain.Tweet
		mockErr    error
		expected   domain.Tweet
		wantErr    bool
	}

	tests := []testCase{
		{
			name:       "Success case",
			mockInput:  domain.Tweet{UserID: mockUUID, Message: "message"},
			mockOutput: domain.Tweet{ID: mockUUID, UserID: mockUUID, Message: "message"},
			mockErr:    nil,
			expected:   domain.Tweet{ID: mockUUID, UserID: mockUUID, Message: "message"},
			wantErr:    false,
		},
		{
			name:       "Repository error",
			mockInput:  domain.Tweet{UserID: uuid.New(), Message: "message"},
			mockOutput: domain.Tweet{},
			mockErr:    errors.New("repository error"),
			expected:   domain.Tweet{},
			wantErr:    true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockCtx := context.Background()
			mockRepo := mock_ports.NewMockTweetRepository(ctrl)
			s := NewTweetsService(mockRepo)

			mockRepo.
				EXPECT().
				CreateTweet(mockCtx, tc.mockInput).
				Return(tc.mockOutput, tc.mockErr)

			got, err := s.CreateTweet(mockCtx, tc.mockInput.UserID, tc.mockInput.Message)

			if (err != nil) != tc.wantErr {
				t.Errorf("CreateTweet() error = %v, wantErr = %v", err, tc.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tc.expected) {
				t.Errorf("CreateTweet() got = %v, want = %v", got, tc.expected)
			}
		})
	}
}
