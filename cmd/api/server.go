package api

import (
	"log"

	"github.com/gin-gonic/gin"
)

func StartServer(router *gin.Engine) {
	if err := router.Run(":8080"); err != nil {
		log.Fatal("error running server", []error{err})
	}
}
