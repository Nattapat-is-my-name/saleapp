package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"saleapp/pkg/response"
)

func Recovery(logger zerolog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				requestID, _ := c.Get("request_id")

				logger.Error().
					Interface("request_id", requestID).
					Interface("error", err).
					Str("stack", string(debug.Stack())).
					Msg("Panic recovered")

				response.Error(
					c,
					http.StatusInternalServerError,
					"INTERNAL_ERROR",
					"An unexpected error occurred",
				)
				c.Abort()
			}
		}()
		c.Next()
	}
}

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last()
			status := c.Writer.Status()
			if status == 0 {
				status = http.StatusInternalServerError
			}

			response.Error(
				c,
				status,
				"ERROR",
				err.Error(),
			)
		}
	}
}

func NotFoundHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"error": gin.H{
				"code":    "NOT_FOUND",
				"message": fmt.Sprintf("Route %s %s not found", c.Request.Method, c.Request.URL.Path),
			},
		})
	}
}
