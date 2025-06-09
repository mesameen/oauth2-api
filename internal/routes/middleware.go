package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mesameen/oauth2-api/internal/logger"
)

func (h *Handler) AuthMiddleware(c *gin.Context) {
	session, err := h.cookieStore.Get(c.Request, "session-name")
	if err != nil {
		logger.Errorf("Failed to authenticate user. token is not presents")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorised"})
		return
	}
	token, ok := session.Values["auth_token"].(string)
	if !ok || token == "" {
		logger.Errorf("Failed to authenticate user. token is not presents")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorised"})
		return
	}
	c.Set("auth_token", token)
	c.Next()
}
