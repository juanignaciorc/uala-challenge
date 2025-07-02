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
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tweet, err := h.service.CreateTweet(ctx, uuid.MustParse(userID), body.Message)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Tweet created successfully", "tweet": tweet})
}
