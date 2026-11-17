package auth

import (
	"net/http"
	"strings"

	"github.com/gadzz/scopepay/internal/httpx"
	"github.com/gin-gonic/gin"
)

func RequireAuth(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		h := c.GetHeader("Authorization")
		const prefix = "Bearer "

		if !strings.HasPrefix(h, prefix) {
			httpx.Error(c, http.StatusUnauthorized, httpx.ErrorBody{Code: httpx.CodeUnauthorized, Message: "missing bearer token"})
			return
		}

		id, err := ParseToken(secret, strings.TrimPrefix(h, prefix))

		if err != nil {
			httpx.Error(c, http.StatusUnauthorized, httpx.ErrorBody{Code: httpx.CodeUnauthorized, Message: "invalid or expired token"})
			return
		}

		c.Set(httpx.ContextKeyUserID, id)
		c.Next()
	}
}
