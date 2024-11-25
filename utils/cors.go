package utils

import (
	"log"

	"github.com/gin-gonic/gin"
)

// EnableCORS sets up CORS middleware for both local development and production use.
func EnableCORS() gin.HandlerFunc {
    return func(c *gin.Context) {
        log.Printf("CORS Request from Origin: %s", c.Request.Header.Get("Origin"))
        c.Header("Access-Control-Allow-Origin", c.Request.Header.Get("Origin"))
        c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
        c.Header("Access-Control-Allow-Credentials", "true")

        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(200)
            return
        }
        c.Next()
    }
}
