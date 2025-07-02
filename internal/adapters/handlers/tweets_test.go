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

const uuidMock = "77dae0ef-658c-44c6-803f-f849854a7033"

func TestTweetHandler_CreateTweet(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_ports.NewMockTweetService(ctrl)

	// Create the handler with the mock service
	handler := NewTweetHandler(mockService)

	tests := []struct {
		name               string
		userID             uuid.UUID
		requestBody        interface{}
		setupMock          func()
		expectedStatusCode int
		expectedResponse   string
	}{
		{
			name:        "Success - Tweet created",
			userID:      uuid.MustParse(uuidMock),
			requestBody: map[string]string{"message": "Tweet created successfully"},
			setupMock: func() {
				// Set up the expected behavior for the mock
				mockService.EXPECT().
					CreateTweet(gomock.Any(), gomock.Any(), "Tweet created successfully").
					Return(domain.Tweet{ID: uuid.MustParse(uuidMock), Message: "Tweet created successfully"}, nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   fmt.Sprintf(`{"message":"Tweet created successfully","tweet":{"id":"%s","user_id":"00000000-0000-0000-0000-000000000000","message":"Tweet created successfully"}}`, uuidMock),
		},
		{
			name:        "Failure - Service error",
			userID:      uuid.MustParse(uuidMock),
			requestBody: map[string]string{"message": "Hello, Error!"},
			setupMock: func() {
				mockService.EXPECT().
					CreateTweet(gomock.Any(), gomock.Any(), "Hello, Error!").
					Return(domain.Tweet{}, errors.New("service error"))
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse:   `{"error":"service error"}`,
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
			req, err := http.NewRequest(http.MethodPost, "/tweets/", bytes.NewBuffer(requestBody))
			if err != nil {
				t.Fatal(err)
			}

			req.Header.Set("Content-Type", "application/json")

			// Create a response recorder to capture the response
			rr := httptest.NewRecorder()

			ctx, _ := gin.CreateTestContext(rr)
			ctx.Request = req
			ctx.Params = gin.Params{
				{Key: "id", Value: uuidMock},
			}

			handler.CreateTweet(ctx)

			assert.Equal(t, tt.expectedStatusCode, rr.Code)
			assert.Equal(t, tt.expectedResponse, rr.Body.String())

		})
	}

}
