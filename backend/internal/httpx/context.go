package httpx

import "github.com/gin-gonic/gin"

const ContextKeyUserID = "userID"

func UserID(c *gin.Context) (uint, bool) {
	v, exists := c.Get(ContextKeyUserID)

	if !exists {
		return 0, false
	}

	id, ok := v.(uint)
	return id, ok
}
