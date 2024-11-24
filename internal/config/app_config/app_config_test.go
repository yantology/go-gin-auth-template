package app_config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitAppConfig(t *testing.T) {
	tests := []struct {
		name     string
		envValue string
		wantPort string
	}{
		{
			name:     "Environment variable not set",
			envValue: "",
			wantPort: ":8000",
		},
		{
			name:     "Environment variable set",
			envValue: ":9000",
			wantPort: ":9000",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Backup original PORT value
			originalPort := PORT

			// Restore original PORT value after test
			defer func() {
				PORT = originalPort
			}()

			// Set or unset the environment variable
			if tt.envValue != "" {
				os.Setenv("APP_PORT", tt.envValue)
			} else {
				os.Unsetenv("APP_PORT")
			}

			// Initialize app config
			InitAppConfig()

			// Check the result
			assert.Equal(t, tt.wantPort, PORT)

			// Clean up environment variable
			os.Unsetenv("APP_PORT")
		})
	}
}
