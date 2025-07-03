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
		ctx.JSON(http.StatusBadRequest, NewErrorResponseWithCode(err.Error(), "INVALID_REQUEST_BODY"))
		return
	}

	user, err := h.service.CreateUser(ctx, body.Name, body.Email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, NewErrorResponse(err.Error()))
		return
	}

	response := NewSuccessResponse("User created successfully", ToUserDetailResponse(user))
	ctx.JSON(http.StatusOK, response)
}

func (h UserHandler) Get(ctx *gin.Context) {
	userIDStr := ctx.Param("id")

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, NewErrorResponseWithCode("Invalid user ID", "INVALID_USER_ID"))
		return
	}

	user, err := h.service.GetUser(ctx, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, NewErrorResponse("USER NOT FOUND"))
		return
	}

	response := NewSuccessResponse("User retrieved successfully", ToUserDetailResponse(user))
	ctx.JSON(http.StatusCreated, response)
}

func (h UserHandler) FollowUser(ctx *gin.Context) {
	userIDStr := ctx.Param("id")

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, NewErrorResponseWithCode("Invalid user ID", "INVALID_USER_ID"))
		return
	}

	followedUserIDStr := ctx.Param("following_user_id")

	followedUserID, err := uuid.Parse(followedUserIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, NewErrorResponseWithCode("Invalid followed user ID", "INVALID_FOLLOWED_USER_ID"))
		return
	}

	if err := h.service.FollowUser(ctx, userID, followedUserID); err != nil {
		ctx.JSON(http.StatusInternalServerError, NewErrorResponse(err.Error()))
		return
	}

	response := NewSuccessResponse("User followed successfully", nil)
	ctx.JSON(http.StatusCreated, response)
}

func (h UserHandler) GetUserTimeline(ctx *gin.Context) {
	userIDStr := ctx.Param("id")

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, NewErrorResponseWithCode("Invalid user ID", "INVALID_USER_ID"))
		return
	}

	tweets, err := h.service.GetUserTimeline(ctx, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, NewErrorResponse(err.Error()))
		return
	}

	// Convert tweets to response DTOs
	tweetResponses := make([]TweetResponse, len(tweets))
	for i, tweet := range tweets {
		tweetResponses[i] = ToTweetResponseSimple(tweet)
	}

	response := NewSuccessResponse("Timeline retrieved successfully", tweetResponses)
	ctx.JSON(http.StatusOK, response)
}
