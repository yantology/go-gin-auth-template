package jwt_config

import (
	"os"
	"strconv"
)

var JWT_ACCESS_SECRET = "access_secret"
var JWT_REFRESH_SECRET = "refresh_secret"
var JWT_ACCESS_TIMEOUT = 15
var JWT_REFRESH_TIMEOUT = 10080

func InitJWTConfig() {
	env_JWT_ACCESS_SECRET := os.Getenv("JWT_ACCESS_SECRET")
	if env_JWT_ACCESS_SECRET != "" {
		JWT_ACCESS_SECRET = env_JWT_ACCESS_SECRET
	}
	env_JWT_REFRESH_SECRET := os.Getenv("JWT_REFRESH_SECRET")
	if env_JWT_REFRESH_SECRET != "" {
		JWT_REFRESH_SECRET = env_JWT_REFRESH_SECRET
	}
	env_JWT_ACCESS_TIMEOUT := os.Getenv("JWT_ACCESS_TIMEOUT")
	if env_JWT_ACCESS_TIMEOUT != "" {
		if timeout, err := strconv.Atoi(env_JWT_ACCESS_TIMEOUT); err == nil {
			JWT_ACCESS_TIMEOUT = timeout
		}
	}
	env_JWT_REFRESH_TIMEOUT := os.Getenv("JWT_REFRESH_TIMEOUT")
	if env_JWT_REFRESH_TIMEOUT != "" {
		if timeout, err := strconv.Atoi(env_JWT_REFRESH_TIMEOUT); err == nil {
			JWT_REFRESH_TIMEOUT = timeout
		}
	}
}
