package api

import (
	"context"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/juanignaciorc/microbloggin-pltf/internal/adapters/handlers"
	"github.com/juanignaciorc/microbloggin-pltf/internal/adapters/repositories/in_memory_db"
	"github.com/juanignaciorc/microbloggin-pltf/internal/adapters/repositories/postgre_db"
	"github.com/juanignaciorc/microbloggin-pltf/internal/services"
)

const basePath = "/api/v1"

func SetupEngine() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Println("DATABASE_URL not set, using in-memory database")
		// Fallback to in-memory database
		repoIMDB := in_memory_db.NewInMemoryDB()

		//Dependency Injection
		userService := services.NewUserService(repoIMDB)
		userHandler := handlers.NewUserHandler(userService)

		tweetService := services.NewTweetsService(repoIMDB)
		tweetHandler := handlers.NewTweetHandler(tweetService)

		router.GET("/ping", handlers.PingHandler)
		router.POST(basePath+"/users", userHandler.Create)
		router.GET(basePath+"/users/:id", userHandler.Get)
		router.POST(basePath+"/users/:id/tweet", tweetHandler.CreateTweet)
		router.POST(basePath+"/users/:id/follow/:following_user_id", userHandler.FollowUser)
		router.GET(basePath+"/users/:id/timeline", userHandler.GetUserTimeline)

		return router
	} else {
		// Initialize PostgreSQL connection
		ctx := context.Background()
		db, err := postgre_db.NewDB(ctx, databaseURL)
		if err != nil {
			log.Fatal("Failed to connect to database:", err)
		}

		// Initialize PostgreSQL repositories
		userRepo := postgre_db.NewUserRepository(db)
		tweetRepo := postgre_db.NewTweetRepository(db)

		//Dependency Injection
		userService := services.NewUserService(userRepo)
		userHandler := handlers.NewUserHandler(userService)

		tweetService := services.NewTweetsService(tweetRepo)
		tweetHandler := handlers.NewTweetHandler(tweetService)

		router.GET("/ping", handlers.PingHandler)
		router.POST(basePath+"/users", userHandler.Create)
		router.GET(basePath+"/users/:id", userHandler.Get)
		router.POST(basePath+"/users/:id/tweet", tweetHandler.CreateTweet)
		router.POST(basePath+"/users/:id/follow/:following_user_id", userHandler.FollowUser)
		router.GET(basePath+"/users/:id/timeline", userHandler.GetUserTimeline)

		return router
	}
}
