package middleware

import (
	"bicycle-store/internal/models"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Check if there are any errors
		if len(c.Errors) > 0 {
			err := c.Errors.Last()
			log.Printf("Error: %v", err)

			c.JSON(http.StatusInternalServerError, models.APIResponse{
				Success: false,
				Error:   err.Error(),
			})
		}
	}
}

func RecoveryHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Panic recovered: %v", r)
				c.JSON(http.StatusInternalServerError, models.APIResponse{
					Success: false,
					Error:   "Internal server error",
				})
				c.Abort()
			}
		}()
		c.Next()
	}
}
