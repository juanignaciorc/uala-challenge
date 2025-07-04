package api

import (
	"context"
	ports "github.com/juanignaciorc/microbloggin-pltf/internal/ports/repositories"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/juanignaciorc/microbloggin-pltf/internal/adapters/handlers"
	"github.com/juanignaciorc/microbloggin-pltf/internal/adapters/repositories/in_memory_db"
	"github.com/juanignaciorc/microbloggin-pltf/internal/adapters/repositories/postgre_db"
	"github.com/juanignaciorc/microbloggin-pltf/internal/services"
)

const basePath = "/api/v1"

func createHandlers(userRepo ports.UsersRepository, tweetRepo ports.TweetRepository) (*handlers.UserHandler, *handlers.TweetHandler) {
	userService := services.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService)

	tweetService := services.NewTweetsService(tweetRepo)
	tweetHandler := handlers.NewTweetHandler(tweetService)

	return userHandler, tweetHandler
}

func setupRoutes(router *gin.Engine, userHandler *handlers.UserHandler, tweetHandler *handlers.TweetHandler) {
	router.GET("/ping", handlers.PingHandler)
	router.POST(basePath+"/users", userHandler.Create)
	router.GET(basePath+"/users/:id", userHandler.Get)
	router.POST(basePath+"/users/:id/tweet", tweetHandler.CreateTweet)
	router.POST(basePath+"/users/:id/follow/:following_user_id", userHandler.FollowUser)
	router.GET(basePath+"/users/:id/timeline", userHandler.GetUserTimeline)
}

func SetupEngine() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Println("DATABASE_URL not set, using in-memory database")
		repoIMDB := in_memory_db.NewInMemoryDB()
		userHandler, tweetHandler := createHandlers(repoIMDB, repoIMDB)

		setupRoutes(router, userHandler, tweetHandler)
		return router
	}

	ctx := context.Background()
	db, err := postgre_db.NewDB(ctx, databaseURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	userRepo := postgre_db.NewUserRepository(db)
	tweetRepo := postgre_db.NewTweetRepository(db)
	userHandler, tweetHandler := createHandlers(userRepo, tweetRepo)

	setupRoutes(router, userHandler, tweetHandler)
	return router
}
