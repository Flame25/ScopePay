package httpx

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
)

func newRequestID() string {
	var b [12]byte

	if _, err := rand.Read(b[:]); err != nil {
		return "req_unknown"
	}

	return "req_" + hex.EncodeToString(b[:])
}

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.GetHeader("X-Request-ID")
		if id == "" {
			id = newRequestID()
		}

		c.Set(ContextKeyRequestID, id)
		c.Header("X-Request-ID", id)
		c.Next()
	}
}

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			r := recover()
			if r != nil {
				log.Printf("[panic] ... %v\n%s", r, debug.Stack())

				if !c.Writer.Written() {
					Error(c, http.StatusInternalServerError, ErrorBody{Code: CodeInternal, Message: "internal error"})
				}
			}

		}()
		c.Next()
	}
}
