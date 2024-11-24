package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/yantology/go-gin-auth-template/internal/config"
	"github.com/yantology/go-gin-auth-template/internal/config/db_config"
	"github.com/yantology/go-gin-auth-template/internal/handlers"
	"github.com/yantology/go-gin-auth-template/internal/middleware"
	"github.com/yantology/go-gin-auth-template/internal/repository"
	"github.com/yantology/go-gin-auth-template/internal/services"
	"github.com/yantology/go-gin-auth-template/internal/utils"
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize configurations
	config.InitConfig()

	// Initialize repositories, services, and handlers
	userRepo := repository.NewUserRepository(db_config.DB)
	jwtUtil := utils.NewJWTUtil(config.JWT_ACCESS_SECRET(), config.JWT_REFRESH_SECRET(), config.JWT_ACCESS_TIMEOUT(), config.JWT_REFRESH_TIMEOUT())
	authService := services.NewAuthService(userRepo, jwtUtil)
	authHandler := handlers.NewAuthHandler(authService)

	// Initialize Gin router
	router := gin.Default()

	// Serve static files
	router.Static(config.PUBLIC_ROUTE(), config.PUBLIC_ASSETS_DIR())
	router.StaticFile("/", "./public/index.html")

	// Public routes
	router.POST("/register", authHandler.Register)
	router.POST("/login", authHandler.Login)
	router.POST("/refresh-token", authHandler.RefreshToken)

	// Protected routes
	authMiddleware := middleware.AuthMiddleware(jwtUtil)
	protected := router.Group("/auth").Use(authMiddleware)
	{
		protected.POST("/change-password", authHandler.ChangePassword)
	}

	// Start the server
	log.Fatal(router.Run(config.PORT()))
}
