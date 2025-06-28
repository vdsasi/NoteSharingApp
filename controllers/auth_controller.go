package controllers

import (
    "net/http"
    "notes-app/models"
    "notes-app/services"
    "github.com/gin-gonic/gin"
)

type AuthController struct {
    authService *services.AuthService
}

func NewAuthController(authService *services.AuthService) *AuthController {
    return &AuthController{authService: authService}
}

func (ac *AuthController) Register(c *gin.Context) {
    var req models.RegisterRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    err := ac.authService.Register(req)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "User registered successfully!"})
}

func (ac *AuthController) Login(c *gin.Context) {
    var req models.LoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    response, err := ac.authService.Login(req)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
        return
    }

    // Set sessionId as HTTP-only, Secure cookie
    c.SetCookie("sessionId", response.SessionID, 24*3600, "/", "", true, true)

    c.JSON(http.StatusOK, response.User)
}

func (ac *AuthController) Logout(c *gin.Context) {
    sessionID, err := c.Cookie("sessionId")
    if err == nil && sessionID != "" {
        ac.authService.Logout(sessionID)
        // Clear the cookie
        c.SetCookie("sessionId", "", -1, "/", "", true, true)
    }
    c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully!"})
}

func (ac *AuthController) GetCurrentUser(c *gin.Context) {
    user, exists := c.Get("user")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
        return
    }

    u := user.(*models.User)
    profile := models.UserProfileDto{
        ID:       u.ID.Hex(),
        Username: u.Username,
        Email:    u.Email,
    }

    c.JSON(http.StatusOK, profile)
}

func (ac *AuthController) ChangePassword(c *gin.Context) {
    user, exists := c.Get("user")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
        return
    }

    var req models.ChangePasswordRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    u := user.(*models.User)
    err := ac.authService.ChangePassword(u.ID, req)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Password changed successfully."})
}
