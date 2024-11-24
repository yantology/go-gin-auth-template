package db_config

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql" // MySQL driver
	_ "github.com/lib/pq"
)

var DB_HOST = "127.0.0.1"
var DB_PORT = "3306"
var DB_NAME = "gin_gonic"
var DB_USER = "root"
var DB_PASSWORD = ""
var DB_DRIVER = "mysql"

// Define a type for the Open function
type OpenFunc func(driverName, dataSourceName string) (*sql.DB, error)

func InitDatabaseConfig() {
	env_DB_HOST := os.Getenv("DB_HOST")
	if env_DB_HOST != "" {
		DB_HOST = env_DB_HOST
	}
	env_DB_PORT := os.Getenv("DB_PORT")
	if env_DB_PORT != "" {
		DB_PORT = env_DB_PORT
	}
	env_DB_NAME := os.Getenv("DB_NAME")
	if env_DB_NAME != "" {
		DB_NAME = env_DB_NAME
	}
	env_DB_USER := os.Getenv("DB_USER")
	if env_DB_USER != "" {
		DB_USER = env_DB_USER
	}
	env_DB_PASSWORD := os.Getenv("DB_PASSWORD")
	if env_DB_PASSWORD != "" {
		DB_PASSWORD = env_DB_PASSWORD
	}
	env_DB_DRIVER := os.Getenv("DB_DRIVER")
	if env_DB_DRIVER != "" {
		DB_DRIVER = env_DB_DRIVER
	}
}

var DB *sql.DB

func ConnectDatabase(openFunc OpenFunc) {
	var errConnection error

	if DB_DRIVER == "mysql" {
		dsnMysql := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", DB_USER, DB_PASSWORD, DB_HOST, DB_PORT, DB_NAME)
		DB, errConnection = openFunc(DB_DRIVER, dsnMysql)
	} else if DB_DRIVER == "postgres" {
		fmt.Printf("postgres://%s:%s@%s:%s/%s?sslmode=disable&TimeZone=Asia/Jakarta", DB_USER, DB_PASSWORD, DB_HOST, DB_PORT, DB_NAME)
		dsnPgSql := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable&TimeZone=Asia/Jakarta", DB_USER, DB_PASSWORD, DB_HOST, DB_PORT, DB_NAME)
		DB, errConnection = openFunc(DB_DRIVER, dsnPgSql)
	} else {
		panic("Invalid database driver")
	}

	if errConnection != nil {
		panic(fmt.Sprintf("Failed to connect to database: %v", errConnection))
	} else {
		fmt.Println("Database connected")
	}
}
