package utils

import (
    "github.com/gin-gonic/gin"
  
)

// EnableCORS middleware function for Gin
func EnableCORS() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Set CORS headers
        c.Header("Access-Control-Allow-Origin", "https://quickz.onrender.com") // Allow the frontend's deployed URL
        c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

        // Handle preflight request (OPTIONS request)
        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(200) // Respond with 200 OK to OPTIONS requests
            return
        }

        // Continue with the request
        c.Next()
    }
}
