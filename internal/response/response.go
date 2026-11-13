package response

import (
	"Finance/internal/errs"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// Body is the shared envelope for JSON responses and SSE error events.
type Body struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// Success sends a successful response with HTTP 200.
func Success(c *gin.Context, data any) {
	c.JSON(http.StatusOK, Body{Code: int(errs.Success), Message: "ok", Data: data})
}

// Created sends a successful response with HTTP 201.
func Created(c *gin.Context, data any) {
	c.JSON(http.StatusCreated, Body{Code: int(errs.Success), Message: "ok", Data: data})
}

// SSE starts the stream on its first event and writes one complete event.
func SSE(c *gin.Context, event string, payload any) error {
	if err := c.Request.Context().Err(); err != nil {
		return err
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("encode SSE event %q: %w", event, err)
	}
	if !c.Writer.Written() {
		c.Header("Content-Type", "text/event-stream; charset=utf-8")
		c.Header("Cache-Control", "no-cache")
		c.Header("X-Accel-Buffering", "no")
		c.Status(http.StatusOK)
	}
	if _, err := fmt.Fprintf(c.Writer, "event: %s\ndata: %s\n\n", event, data); err != nil {
		return fmt.Errorf("write SSE event %q: %w", event, err)
	}
	c.Writer.Flush()
	return nil
}

// Fail sends an error response and aborts the remaining Gin handlers.
// Callers must still return to stop executing their current handler.
// Once SSE starts, errors use the same envelope in an error event.
func Fail(c *gin.Context, err error) {
	c.Abort()
	if errors.Is(c.Request.Context().Err(), context.Canceled) {
		return
	}

	e := errs.From(err)

	if e.HTTPCode >= http.StatusInternalServerError {
		log.Printf("request failed: method=%s path=%s error=%v", c.Request.Method, c.Request.URL.Path, err)
	}

	body := Body{Code: int(e.Code), Message: e.Msg}
	if c.Writer.Written() {
		if strings.HasPrefix(c.Writer.Header().Get("Content-Type"), "text/event-stream") {
			if sendErr := SSE(c, "error", body); sendErr != nil {
				log.Printf("send stream error failed: method=%s path=%s error=%v", c.Request.Method, c.Request.URL.Path, sendErr)
			}
		}
		return
	}
	c.AbortWithStatusJSON(e.HTTPCode, body)
}
