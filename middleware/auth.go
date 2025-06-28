package middleware

import (
    "net/http"
    "notes-app/services"
    "github.com/gin-gonic/gin"
)

func AuthMiddleware(authService *services.AuthService) gin.HandlerFunc {
    return func(c *gin.Context) {
        sessionID, err := c.Cookie("sessionId")
        if err != nil || sessionID == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "No session provided"})
            c.Abort()
            return
        }

        user, err := authService.GetUserBySession(sessionID)
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid session"})
            c.Abort()
            return
        }

        c.Set("user", user)
        c.Next()
    }
}
