package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/juanignaciorc/microbloggin-pltf/internal/services"
	"net/http"
)

type TweetHandler struct {
	service     services.TweetService
	userService services.UserService
}

func NewTweetHandler(service services.TweetService, userService services.UserService) *TweetHandler {
	return &TweetHandler{
		service:     service,
		userService: userService,
	}
}

func (h *TweetHandler) CreateTweet(ctx *gin.Context) {
	userID := ctx.Param("id")

	var body CreateTweetBody
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, NewErrorResponseWithCode(err.Error(), "EXCEEDED_MAX_TWEET_CHARACTERS"))
		return
	}

	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, NewErrorResponseWithCode("Invalid user ID", "INVALID_USER_ID"))
		return
	}

	tweet, err := h.service.CreateTweet(ctx, parsedUserID, body.Message)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, NewErrorResponse(err.Error()))
		return
	}

	// Fetch user information to include in the response
	user, err := h.userService.GetUser(ctx, parsedUserID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, NewErrorResponse(err.Error()))
		return
	}

	response := NewSuccessResponse("Tweet created successfully", ToTweetResponseWithUser(tweet, user))
	ctx.JSON(http.StatusCreated, response)
}
