package main

import (
    "log"
    "notes-app/config"
    "notes-app/controllers"
    "notes-app/middleware"
    "notes-app/services"
    "github.com/gin-gonic/gin"
    "github.com/gin-contrib/cors"
    "github.com/joho/godotenv"
)

func main() {
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found, using default values")
    }

    config.ConnectMongoDB()
    config.ConnectRedis()

    authService := services.NewAuthService()
    noteService := services.NewNoteService()

    authController := controllers.NewAuthController(authService)
    noteController := controllers.NewNoteController(noteService)

    router := gin.Default()

    corsConfig := cors.DefaultConfig()
    corsConfig.AllowOrigins = []string{"http://localhost:5173"}
    corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
    corsConfig.AllowHeaders = []string{"Content-Type", "X-Requested-With", "Authorization"} // Explicitly allow required headers
    corsConfig.AllowCredentials = true
    router.Use(cors.New(corsConfig))

    authRoutes := router.Group("/api/auth")
    {
        authRoutes.POST("/register", authController.Register)
        authRoutes.POST("/login", authController.Login)
        authRoutes.POST("/logout", authController.Logout)
        authRoutes.GET("/me", middleware.AuthMiddleware(authService), authController.GetCurrentUser)
        authRoutes.POST("/change-password", middleware.AuthMiddleware(authService), authController.ChangePassword)
    }

    noteRoutes := router.Group("/api/notes")
    noteRoutes.Use(middleware.AuthMiddleware(authService))
    {
        noteRoutes.GET("", noteController.GetAll)
        noteRoutes.POST("", noteController.Create)
        noteRoutes.GET(":id", noteController.GetNote)
        noteRoutes.PUT(":id", noteController.Update)
        noteRoutes.DELETE(":id", noteController.Delete)
        noteRoutes.GET("/trash", noteController.GetTrashed)
        noteRoutes.POST(":id/restore", noteController.Restore)
        noteRoutes.POST(":id/pin", noteController.TogglePin)
        noteRoutes.GET(":id/versions", noteController.GetHistory)
        noteRoutes.POST("/version-restore/:noteId/:versionId", noteController.RestoreVersion)
        noteRoutes.GET("/filter", noteController.FilterByTag)
        noteRoutes.PUT("/autosave/:noteId", noteController.AutoSave)
    }

    log.Println("Server starting on :8080")
    if err := router.Run(":8080"); err != nil {
        log.Fatal("Failed to start server:", err)
    }
}
