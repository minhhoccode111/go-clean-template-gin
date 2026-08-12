package middleware_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/minhhoccode111/go-clean-template-gin/internal/controller/restapi/middleware"
	"github.com/minhhoccode111/go-clean-template-gin/pkg/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestApp(t *testing.T) (*gin.Engine, *jwt.Manager) {
	t.Helper()

	gin.SetMode(gin.TestMode)

	jwtManager := jwt.New("test-secret", time.Hour)

	app := gin.New()
	app.Use(middleware.Auth(jwtManager))
	app.GET("/test", func(c *gin.Context) {
		v, ok := c.Get("userID")
		if !ok {
			c.Status(http.StatusUnauthorized)

			return
		}

		userID, ok := v.(string)
		if !ok {
			c.Status(http.StatusUnauthorized)

			return
		}

		c.String(http.StatusOK, userID)
	})

	return app, jwtManager
}

func TestAuthMiddleware(t *testing.T) {
	t.Parallel()

	app, jwtManager := newTestApp(t)

	validToken, err := jwtManager.GenerateToken("user-id-123")
	require.NoError(t, err)

	tests := []struct {
		name           string
		authHeader     string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "missing header",
			authHeader:     "",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "invalid format",
			authHeader:     "Basic xxx",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "invalid token",
			authHeader:     "Bearer invalid",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "valid token",
			authHeader:     "Bearer " + validToken,
			expectedStatus: http.StatusOK,
			expectedBody:   "user-id-123",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/test", http.NoBody)
			if tc.authHeader != "" {
				req.Header.Set("Authorization", tc.authHeader)
			}

			w := httptest.NewRecorder()
			app.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatus, w.Code)

			if tc.expectedBody != "" {
				body, readErr := io.ReadAll(w.Result().Body)
				require.NoError(t, readErr)
				assert.Equal(t, tc.expectedBody, string(body))
			}
		})
	}
}
