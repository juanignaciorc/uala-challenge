package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/juanignaciorc/microbloggin-pltf/internal/services"
	"net/http"
)

type TweetHandler struct {
	service services.TweetService
}

func NewTweetHandler(service services.TweetService) *TweetHandler {
	return &TweetHandler{
		service: service,
	}
}

func (h *TweetHandler) CreateTweet(ctx *gin.Context) {
	userID := ctx.Param("id")

	var body CreateTweetBody
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, NewErrorResponseWithCode(err.Error(), "INVALID_REQUEST_BODY"))
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

	response := NewSuccessResponse("Tweet created successfully", ToTweetResponseSimple(tweet))
	ctx.JSON(http.StatusOK, response)
}
