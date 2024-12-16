package handlers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"codeinstyle.io/captain/config"
	"codeinstyle.io/captain/db"
	"codeinstyle.io/captain/models"
	"codeinstyle.io/captain/repository"
	"codeinstyle.io/captain/system"
	"codeinstyle.io/captain/utils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func setupAuthTest(t *testing.T, database *gorm.DB) (*gin.Engine, *AuthHandlers) {
	repositories := repository.NewRepositories(database)
	router := setupTestRouter(repositories)

	// Create minimal templates for testing
	authHandlers := NewAuthHandlers(repositories, &config.Config{})

	// Create test user
	password := "Test1234!"
	hashedPassword, err := utils.HashPassword(password)
	assert.NoError(t, err)

	testUser := models.User{
		Email:    "test@example.com",
		Password: hashedPassword,
	}

	err = database.Create(&testUser).Error
	assert.NoError(t, err)

	return router, authHandlers
}

func TestAuthHandlers_Login(t *testing.T) {
	database := db.SetupTestDB()
	router, authHandlers := setupAuthTest(t, database)

	// Add routes
	router.GET("/login", authHandlers.Login)
	router.POST("/login", authHandlers.PostLogin)

	tests := []struct {
		name          string
		email         string
		password      string
		next          string
		expectedCode  int
		expectedPath  string
		sessionCookie bool
	}{
		{
			name:          "Valid login",
			email:         "test@example.com",
			password:      "Test1234!",
			expectedCode:  http.StatusFound,
			expectedPath:  "/admin",
			sessionCookie: true,
		},
		{
			name:          "Invalid password",
			email:         "test@example.com",
			password:      "wrong",
			expectedCode:  http.StatusUnauthorized,
			expectedPath:  "",
			sessionCookie: false,
		},
		{
			name:          "Invalid email",
			email:         "wrong@example.com",
			password:      "Test1234!",
			expectedCode:  http.StatusUnauthorized,
			expectedPath:  "",
			sessionCookie: false,
		},
		{
			name:          "Custom return path",
			email:         "test@example.com",
			password:      "Test1234!",
			next:          "/admin/posts",
			expectedCode:  http.StatusFound,
			expectedPath:  "/admin/posts",
			sessionCookie: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			form := url.Values{}
			form.Add("email", tt.email)
			form.Add("password", tt.password)

			// Add return path if specified
			if tt.next != "" {
				form.Add("next", tt.next)
			}

			// Create request
			req, err := http.NewRequest("POST", "/login", strings.NewReader(form.Encode()))
			assert.NoError(t, err)

			req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

			// Create response recorder
			w := httptest.NewRecorder()

			// Serve request
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)
			if tt.expectedPath != "" {
				assert.Equal(t, tt.expectedPath, w.Header().Get("Location"))
			}

			// Check session cookie
			if tt.sessionCookie {
				assert.Contains(t, w.Header().Get("Set-Cookie"), system.CookieName+"=")
			} else {
				assert.NotContains(t, w.Header().Get("Set-Cookie"), system.CookieName+"=")
			}
		})
	}
}

func TestAuthHandlers_Logout(t *testing.T) {
	database := db.SetupTestDB()
	router, authHandlers := setupAuthTest(t, database)

	// Add route
	router.GET("/logout", authHandlers.Logout)

	// Test GET /logout
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/logout", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusFound, w.Code)
	assert.Equal(t, "/login", w.Header().Get("Location"))

	// Check that session cookie is cleared
	assert.Contains(t, w.Header().Get("Set-Cookie"), system.CookieName+"=;")
	assert.Contains(t, w.Header().Get("Set-Cookie"), "Max-Age=0")
}
