package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// TestCORSMiddleware tests the CORSMiddleware to ensure it sets correct CORS headers and handles preflight requests correctly.
func TestCORSMiddleware(t *testing.T) {
	tests := []struct {
		name            string
		method          string
		path            string
		expectedStatus  int
		expectedHeaders map[string]string
	}{
		{
			name:           "OPTIONS request",
			method:         http.MethodOptions,
			path:           "/",
			expectedStatus: http.StatusNoContent,
			expectedHeaders: map[string]string{
				"Access-Control-Allow-Origin":      "*",
				"Access-Control-Allow-Methods":     "POST, HEAD, PATCH, OPTIONS, GET, PUT, DELETE",
				"Access-Control-Allow-Headers":     "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With",
				"Access-Control-Allow-Credentials": "true",
			},
		},
		{
			name:           "GET request",
			method:         http.MethodGet,
			path:           "/test-get",
			expectedStatus: http.StatusOK,
			expectedHeaders: map[string]string{
				"Access-Control-Allow-Origin":      "*",
				"Access-Control-Allow-Methods":     "POST, HEAD, PATCH, OPTIONS, GET, PUT, DELETE",
				"Access-Control-Allow-Headers":     "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With",
				"Access-Control-Allow-Credentials": "true",
			},
		},
		{
			name:           "POST request",
			method:         http.MethodPost,
			path:           "/test-post",
			expectedStatus: http.StatusOK,
			expectedHeaders: map[string]string{
				"Access-Control-Allow-Origin":      "*",
				"Access-Control-Allow-Methods":     "POST, HEAD, PATCH, OPTIONS, GET, PUT, DELETE",
				"Access-Control-Allow-Headers":     "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With",
				"Access-Control-Allow-Credentials": "true",
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				router := gin.New()
				router.Use(CORSMiddleware())

				switch tt.method {
				case http.MethodGet:
					router.GET(
						tt.path, func(c *gin.Context) {
							c.String(http.StatusOK, "OK")
						},
					)
				case http.MethodPost:
					router.POST(
						tt.path, func(c *gin.Context) {
							c.String(http.StatusOK, "OK")
						},
					)
				case http.MethodOptions:
					router.OPTIONS(tt.path, func(c *gin.Context) {})
				}

				req, _ := http.NewRequest(tt.method, tt.path, nil)
				resp := httptest.NewRecorder()

				router.ServeHTTP(resp, req)

				assert.Equal(t, tt.expectedStatus, resp.Code)

				for key, expectedValue := range tt.expectedHeaders {
					assert.Equal(t, expectedValue, resp.Header().Get(key))
				}

				if tt.method == http.MethodOptions {
					assert.Empty(t, resp.Body.String())
				}
			},
		)
	}
}

// TestCORSMiddleware_WithOrigin tests if the CORSMiddleware correctly sets the Access-Control-Allow-Origin header for a valid origin.
func TestCORSMiddleware_WithOrigin(t *testing.T) {
	router := gin.New()
	router.Use(CORSMiddleware())
	router.GET(
		"/test", func(c *gin.Context) {
			c.String(http.StatusOK, "OK")
		},
	)

	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Origin", "http://example.com")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, "*", resp.Header().Get("Access-Control-Allow-Origin"))
	assert.Equal(t, http.StatusOK, resp.Code)
}
