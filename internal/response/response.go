package response

import "github.com/gin-gonic/gin"

// Response is the shared JSON envelope for HTTP API responses.
type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// Success sends a successful response with the supplied HTTP status and data.
func Success(c *gin.Context, status int, data any) {
	c.JSON(status, Response{Code: 0, Message: "success", Data: data})
}

// Error sends an error response and aborts the remaining Gin handlers.
// Callers must still return to stop executing their current handler.
func Error(c *gin.Context, status int, message string) {
	c.AbortWithStatusJSON(status, Response{Code: 1, Message: message})
}
