package db_config

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock for the Open function
type MockOpenFunc struct {
	mock.Mock
}

func (m *MockOpenFunc) Open(driverName, dataSourceName string) (*sql.DB, error) {
	args := m.Called(driverName, dataSourceName)
	return args.Get(0).(*sql.DB), args.Error(1)
}

func TestInitDatabaseConfig(t *testing.T) {
	// Save original values
	originalDBHost := DB_HOST
	originalDBPort := DB_PORT
	originalDBName := DB_NAME
	originalDBUser := DB_USER
	originalDBPassword := DB_PASSWORD
	originalDBDriver := DB_DRIVER

	// Restore original values after test
	defer func() {
		DB_HOST = originalDBHost
		DB_PORT = originalDBPort
		DB_NAME = originalDBName
		DB_USER = originalDBUser
		DB_PASSWORD = originalDBPassword
		DB_DRIVER = originalDBDriver
	}()

	// Set environment variables for testing
	t.Setenv("DB_HOST", "test_host")
	t.Setenv("DB_PORT", "test_port")
	t.Setenv("DB_NAME", "test_name")
	t.Setenv("DB_USER", "test_user")
	t.Setenv("DB_PASSWORD", "test_password")
	t.Setenv("DB_DRIVER", "test_driver")

	InitDatabaseConfig()
	assert.Equal(t, "test_host", DB_HOST)
	assert.Equal(t, "test_port", DB_PORT)
	assert.Equal(t, "test_name", DB_NAME)
	assert.Equal(t, "test_user", DB_USER)
	assert.Equal(t, "test_password", DB_PASSWORD)
	assert.Equal(t, "test_driver", DB_DRIVER)
}
func TestConnectDatabase(t *testing.T) {
	t.Setenv("DB_HOST", "test_host")
	t.Setenv("DB_PORT", "test_port")
	t.Setenv("DB_NAME", "test_name")
	t.Setenv("DB_USER", "test_user")
	t.Setenv("DB_PASSWORD", "test_password")

	tests := []struct {
		name        string
		driver      string
		expectedDSN string
		expectedErr bool
		assertExp   bool
	}{
		{
			name:        "MySQL connection",
			driver:      "mysql",
			expectedDSN: "test_user:test_password@tcp(test_host:test_port)/test_name",
			expectedErr: false,
			assertExp:   true,
		},
		{
			name:        "Postgres connection",
			driver:      "postgres",
			expectedDSN: "postgres://test_user:test_password@test_host:test_port/test_name?sslmode=disable&TimeZone=Asia/Jakarta",
			expectedErr: false,
			assertExp:   true,
		},
		{
			name:        "MySQL Failed connection",
			driver:      "mysql",
			expectedDSN: "test_user:test_password@tcp(test_host:test_port)/test_name",
			expectedErr: true,
			assertExp:   true,
		},
		{
			name:        "Postgres Failed connection",
			driver:      "postgres",
			expectedDSN: "postgres://test_user:test_password@test_host:test_port/test_name?sslmode=disable&TimeZone=Asia/Jakarta",
			expectedErr: true,
			assertExp:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("DB_DRIVER", tt.driver)
			InitDatabaseConfig()
			mockOpenFunc := new(MockOpenFunc)
			if tt.expectedErr {
				mockOpenFunc.On("Open", tt.driver, tt.expectedDSN).Return(nil, assert.AnError)
				assert.Panics(t, func() {
					ConnectDatabase(mockOpenFunc.Open)
				}, "The code did not panic")
			} else {
				mockOpenFunc.On("Open", tt.driver, tt.expectedDSN).Return(&sql.DB{}, nil)
				assert.NotPanics(t, func() {
					ConnectDatabase(mockOpenFunc.Open)
				}, "The code did not panic")
			}
			if tt.assertExp {
				mockOpenFunc.AssertExpectations(t)

			}
		})
	}
}
