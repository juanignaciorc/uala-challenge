package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/google/uuid"
	"github.com/juanignaciorc/microbloggin-pltf/internal/domain"
	mock_ports "github.com/juanignaciorc/microbloggin-pltf/mocks"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
)

const userUuidMock = "77dae0ef-658c-44c6-803f-f849854a7033"
const followedUserUuidMock = "88dae0ef-658c-44c6-803f-f849854a7044"

func TestUserHandler_Create(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_ports.NewMockUserService(ctrl)

	// Create the handler with the mock service
	handler := NewUserHandler(mockService)

	tests := []struct {
		name               string
		requestBody        interface{}
		setupMock          func()
		expectedStatusCode int
		expectedResponse   string
	}{
		{
			name:        "Success - User created",
			requestBody: map[string]string{"name": "John Doe", "email": "john@example.com"},
			setupMock: func() {
				mockService.EXPECT().
					CreateUser(gomock.Any(), "John Doe", "john@example.com").
					Return(domain.User{ID: uuid.MustParse(userUuidMock), Name: "John Doe", Email: "john@example.com"}, nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   fmt.Sprintf(`{"message":"User created successfully","data":{"id":"%s","name":"John Doe","email":"john@example.com","followers_count":0,"following_count":0,"tweets_count":0}}`, userUuidMock),
		},
		{
			name:        "Failure - Service error",
			requestBody: map[string]string{"name": "Jane Doe", "email": "jane@example.com"},
			setupMock: func() {
				mockService.EXPECT().
					CreateUser(gomock.Any(), "Jane Doe", "jane@example.com").
					Return(domain.User{}, errors.New("service error"))
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse:   `{"error":"service error"}`,
		},
		{
			name:               "Failure - Invalid JSON",
			requestBody:        `{"name": "John", "email":}`,
			setupMock:          func() {},
			expectedStatusCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up mock expectations
			tt.setupMock()

			// Create a request body
			var requestBody []byte
			if body, ok := tt.requestBody.(map[string]string); ok {
				requestBody, _ = json.Marshal(body)
			} else if str, ok := tt.requestBody.(string); ok {
				requestBody = []byte(str)
			}

			// Create a new HTTP request
			req, err := http.NewRequest(http.MethodPost, "/users/", bytes.NewBuffer(requestBody))
			if err != nil {
				t.Fatal(err)
			}

			req.Header.Set("Content-Type", "application/json")

			// Create a response recorder to capture the response
			rr := httptest.NewRecorder()

			ctx, _ := gin.CreateTestContext(rr)
			ctx.Request = req

			handler.Create(ctx)

			assert.Equal(t, tt.expectedStatusCode, rr.Code)
			if tt.expectedResponse != "" {
				assert.Equal(t, tt.expectedResponse, rr.Body.String())
			}
		})
	}
}

func TestUserHandler_Get(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_ports.NewMockUserService(ctrl)

	// Create the handler with the mock service
	handler := NewUserHandler(mockService)

	tests := []struct {
		name               string
		userID             string
		setupMock          func()
		expectedStatusCode int
		expectedResponse   string
	}{
		{
			name:   "Success - User found",
			userID: userUuidMock,
			setupMock: func() {
				mockService.EXPECT().
					GetUser(gomock.Any(), uuid.MustParse(userUuidMock)).
					Return(domain.User{ID: uuid.MustParse(userUuidMock), Name: "John Doe", Email: "john@example.com"}, nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   fmt.Sprintf(`{"message":"User retrieved successfully","data":{"id":"%s","name":"John Doe","email":"john@example.com","followers_count":0,"following_count":0,"tweets_count":0}}`, userUuidMock),
		},
		{
			name:   "Failure - Service error",
			userID: userUuidMock,
			setupMock: func() {
				mockService.EXPECT().
					GetUser(gomock.Any(), uuid.MustParse(userUuidMock)).
					Return(domain.User{}, errors.New("USER NOT FOUND"))
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse:   `{"error":"USER NOT FOUND"}`,
		},
		{
			name:               "Failure - Invalid UUID",
			userID:             "invalid-uuid",
			setupMock:          func() {},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   `{"error":"Invalid user ID","code":"INVALID_USER_ID"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up mock expectations
			tt.setupMock()

			// Create a new HTTP request
			req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/users/%s", tt.userID), nil)
			if err != nil {
				t.Fatal(err)
			}

			// Create a response recorder to capture the response
			rr := httptest.NewRecorder()

			ctx, _ := gin.CreateTestContext(rr)
			ctx.Request = req
			ctx.Params = gin.Params{
				{Key: "id", Value: tt.userID},
			}

			handler.Get(ctx)

			assert.Equal(t, tt.expectedStatusCode, rr.Code)
			assert.Equal(t, tt.expectedResponse, rr.Body.String())
		})
	}
}

func TestUserHandler_FollowUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_ports.NewMockUserService(ctrl)

	// Create the handler with the mock service
	handler := NewUserHandler(mockService)

	tests := []struct {
		name               string
		userID             string
		followedUserID     string
		setupMock          func()
		expectedStatusCode int
		expectedResponse   string
	}{
		{
			name:           "Success - User followed",
			userID:         userUuidMock,
			followedUserID: followedUserUuidMock,
			setupMock: func() {
				mockService.EXPECT().
					FollowUser(gomock.Any(), uuid.MustParse(userUuidMock), uuid.MustParse(followedUserUuidMock)).
					Return(nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   `{"message":"User followed successfully"}`,
		},
		{
			name:           "Failure - Service error",
			userID:         userUuidMock,
			followedUserID: followedUserUuidMock,
			setupMock: func() {
				mockService.EXPECT().
					FollowUser(gomock.Any(), uuid.MustParse(userUuidMock), uuid.MustParse(followedUserUuidMock)).
					Return(errors.New("follow error"))
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse:   `{"error":"follow error"}`,
		},
		{
			name:               "Failure - Invalid user ID",
			userID:             "invalid-uuid",
			followedUserID:     followedUserUuidMock,
			setupMock:          func() {},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   `{"error":"Invalid user ID","code":"INVALID_USER_ID"}`,
		},
		{
			name:               "Failure - Invalid followed user ID",
			userID:             userUuidMock,
			followedUserID:     "invalid-uuid",
			setupMock:          func() {},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   `{"error":"Invalid followed user ID","code":"INVALID_FOLLOWED_USER_ID"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up mock expectations
			tt.setupMock()

			// Create a new HTTP request
			req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("/users/%s/follow/%s", tt.userID, tt.followedUserID), nil)
			if err != nil {
				t.Fatal(err)
			}

			// Create a response recorder to capture the response
			rr := httptest.NewRecorder()

			ctx, _ := gin.CreateTestContext(rr)
			ctx.Request = req
			ctx.Params = gin.Params{
				{Key: "id", Value: tt.userID},
				{Key: "following_user_id", Value: tt.followedUserID},
			}

			handler.FollowUser(ctx)

			assert.Equal(t, tt.expectedStatusCode, rr.Code)
			assert.Equal(t, tt.expectedResponse, rr.Body.String())
		})
	}
}

func TestUserHandler_GetUserTimeline(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_ports.NewMockUserService(ctrl)

	// Create the handler with the mock service
	handler := NewUserHandler(mockService)

	tests := []struct {
		name               string
		userID             string
		setupMock          func()
		expectedStatusCode int
		expectedResponse   string
	}{
		{
			name:   "Success - Timeline retrieved",
			userID: userUuidMock,
			setupMock: func() {
				tweets := []domain.Tweet{
					{ID: uuid.MustParse(userUuidMock), Message: "Hello World"},
				}
				mockService.EXPECT().
					GetUserTimeline(gomock.Any(), uuid.MustParse(userUuidMock)).
					Return(tweets, nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   fmt.Sprintf(`{"message":"Timeline retrieved successfully","data":[{"id":"%s","message":"Hello World","user":{"id":"00000000-0000-0000-0000-000000000000","name":""}}]}`, userUuidMock),
		},
		{
			name:   "Success - Empty timeline",
			userID: userUuidMock,
			setupMock: func() {
				mockService.EXPECT().
					GetUserTimeline(gomock.Any(), uuid.MustParse(userUuidMock)).
					Return([]domain.Tweet{}, nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   `{"message":"Timeline retrieved successfully","data":[]}`,
		},
		{
			name:   "Failure - Service error",
			userID: userUuidMock,
			setupMock: func() {
				mockService.EXPECT().
					GetUserTimeline(gomock.Any(), uuid.MustParse(userUuidMock)).
					Return(nil, errors.New("timeline error"))
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse:   `{"error":"timeline error"}`,
		},
		{
			name:               "Failure - Invalid UUID",
			userID:             "invalid-uuid",
			setupMock:          func() {},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   `{"error":"Invalid user ID","code":"INVALID_USER_ID"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up mock expectations
			tt.setupMock()

			// Create a new HTTP request
			req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/users/%s/timeline", tt.userID), nil)
			if err != nil {
				t.Fatal(err)
			}

			// Create a response recorder to capture the response
			rr := httptest.NewRecorder()

			ctx, _ := gin.CreateTestContext(rr)
			ctx.Request = req
			ctx.Params = gin.Params{
				{Key: "id", Value: tt.userID},
			}

			handler.GetUserTimeline(ctx)

			assert.Equal(t, tt.expectedStatusCode, rr.Code)
			assert.Equal(t, tt.expectedResponse, rr.Body.String())
		})
	}
}
