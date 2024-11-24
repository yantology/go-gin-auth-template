package config

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/yantology/go-gin-auth-template/internal/config/app_config"
	"github.com/yantology/go-gin-auth-template/internal/config/cors_config"
	"github.com/yantology/go-gin-auth-template/internal/config/db_config"
	"github.com/yantology/go-gin-auth-template/internal/config/jwt_config"
)

func InitConfig() {
	app_config.InitAppConfig()
	db_config.InitDatabaseConfig()
	db_config.ConnectDatabase(sql.Open)
	jwt_config.InitJWTConfig()
}

// variable app_config
// var (
// 	PORT             = app_config.PORT
// 	PUBLIC_ROUTE      = app_config.PUBLIC_ROUTE
// 	PUBLIC_ASSETS_DIR = &app_config.PUBLIC_ASSETS_DIR
// )

func PORT() string {
	return app_config.PORT
}

func PUBLIC_ROUTE() string {
	return app_config.PUBLIC_ROUTE
}

func PUBLIC_ASSETS_DIR() string {
	return app_config.PUBLIC_ASSETS_DIR
}

// variable db_config

func DB_DRIVER() string {
	return db_config.DB_DRIVER
}

func DB_USER() string {
	return db_config.DB_USER
}

func DB_PASSWORD() string {
	return db_config.DB_PASSWORD
}

func DB_NAME() string {
	return db_config.DB_NAME
}

func DB_HOST() string {
	return db_config.DB_HOST
}

func DB_PORT() string {
	return db_config.DB_PORT
}

func DB() *sql.DB {
	return db_config.DB
}

// variable cors_config
func CORS_ALLOW_ORIGINS() gin.HandlerFunc {
	return cors_config.CorsConfig()
}

// variable jwt_config
func JWT_ACCESS_SECRET() string {
	return jwt_config.JWT_ACCESS_SECRET
}

func JWT_REFRESH_SECRET() string {
	return jwt_config.JWT_REFRESH_SECRET
}

func JWT_ACCESS_TIMEOUT() int {
	return jwt_config.JWT_ACCESS_TIMEOUT
}

func JWT_REFRESH_TIMEOUT() int {
	return jwt_config.JWT_REFRESH_TIMEOUT
}
