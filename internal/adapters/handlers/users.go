package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/juanignaciorc/microbloggin-pltf/internal/services"
	"net/http"
)

type UserHandler struct {
	service services.UserService
}

func NewUserHandler(service services.UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

type CreateUserBody struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type CreateTweetBody struct {
	Message string `json:"message" binding:"required,max=280"`
}

func (h UserHandler) Create(ctx *gin.Context) {
	var body CreateUserBody
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.service.CreateUser(ctx, body.Name, body.Email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "User created successfully", "user": user})
}

func (h UserHandler) Get(ctx *gin.Context) {
	userIDStr := ctx.Param("id")

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := h.service.GetUser(ctx, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"user": user})
}

func (h UserHandler) FollowUser(ctx *gin.Context) {
	userIDStr := ctx.Param("id")

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	followedUserIDStr := ctx.Param("following_user_id")

	followedUserID, err := uuid.Parse(followedUserIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid followed user ID"})
		return
	}

	if err := h.service.FollowUser(ctx, userID, followedUserID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "User followed successfully"})
}

func (h UserHandler) GetUserTimeline(ctx *gin.Context) {
	userIDStr := ctx.Param("id")

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	tweets, err := h.service.GetUserTimeline(ctx, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"tweets": tweets})
}
