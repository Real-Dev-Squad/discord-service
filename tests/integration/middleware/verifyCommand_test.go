package middleware_test

import (
	"crypto/ed25519"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Real-Dev-Squad/discord-service/config"
	"github.com/Real-Dev-Squad/discord-service/middleware"
	_ "github.com/Real-Dev-Squad/discord-service/tests/helpers"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
)

func setup() *httprouter.Router {
	router := httprouter.New()
	router.POST("/", middleware.VerifyCommand(func(response http.ResponseWriter, request *http.Request, params httprouter.Params) {
		response.Header().Set("Content-Type", "application/json; charset=UTF-8")
		response.WriteHeader(http.StatusOK)
	}))
	return router
}

func TestVerifyCommand(t *testing.T) {

	subtests := []struct {
		name           string
		setupMock      func()
		expectedStatus int
	}{
		{
			name: "Should verify the verification with valid public key",
			setupMock: func() {
				middleware.VerifyInteraction = func(r *http.Request, key ed25519.PublicKey) bool {
					return true
				}
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Should fail the verification with invalid public key",
			setupMock: func() {
				middleware.VerifyInteraction = func(r *http.Request, key ed25519.PublicKey) bool {
					return false
				}
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "Should fail with Internal Server Error when public key is malformed",
			setupMock: func() {
				originalKey := config.AppConfig.DISCORD_PUBLIC_KEY
				config.AppConfig.DISCORD_PUBLIC_KEY = "invalid hex string"
				t.Cleanup(func() {
					config.AppConfig.DISCORD_PUBLIC_KEY = originalKey
				})
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, subtest := range subtests {
		t.Run(subtest.name, func(t *testing.T) {
			var originalFunc = middleware.VerifyInteraction
			subtest.setupMock()
			router := setup()
			w := httptest.NewRecorder()
			r, _ := http.NewRequest("POST", "/", nil)
			router.ServeHTTP(w, r)
			assert.Equal(t, subtest.expectedStatus, w.Code)
			t.Cleanup(func() {
				middleware.VerifyInteraction = originalFunc
			})
		})
	}
}
