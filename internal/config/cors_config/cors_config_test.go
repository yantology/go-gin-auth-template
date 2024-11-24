package cors_config

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCorsConfig(t *testing.T) {
	tests := []struct {
		name            string
		envValue        string
		origin          string
		wantStatus      int
		wantAllowOrigin string
	}{
		{
			name:            "Default Origin",
			envValue:        "",
			origin:          "http://example.com",
			wantStatus:      http.StatusOK,
			wantAllowOrigin: "*",
		},
		{
			name:            "Allowed Origin",
			envValue:        "http://example.com,http://another-example.com",
			origin:          "http://example.com",
			wantStatus:      http.StatusOK,
			wantAllowOrigin: "http://example.com",
		},
		{
			name:            "Wildcard Origin",
			envValue:        "http://example.com,http://another-example.com",
			origin:          "http://another-example.com",
			wantStatus:      http.StatusOK,
			wantAllowOrigin: "http://another-example.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set or unset the environment variable
			if tt.envValue != "" {
				os.Setenv("CORS_ALLOW_ORIGINS", tt.envValue)
			} else {
				os.Unsetenv("CORS_ALLOW_ORIGINS")
			}

			// Setup Gin router with CORS middleware
			router := gin.Default()
			router.Use(CorsConfig())

			// Define a simple handler for testing
			router.GET("/test", func(c *gin.Context) {
				c.String(http.StatusOK, "CORS test")
			})

			// Create a new HTTP request
			req, _ := http.NewRequest("GET", "/test", nil)
			req.Header.Set("Origin", tt.origin)

			// Record the response
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Check the response status code
			assert.Equal(t, tt.wantStatus, w.Code)

			// Check the CORS headers
			assert.Equal(t, tt.wantAllowOrigin, w.Header().Get("Access-Control-Allow-Origin"))

			// Clean up environment variable
			os.Unsetenv("CORS_ALLOW_ORIGINS")
		})
	}
}
