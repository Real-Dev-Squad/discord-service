package middleware_test

import (
	"crypto/ed25519"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Real-Dev-Squad/discord-service/config"
	"github.com/Real-Dev-Squad/discord-service/middleware"
	_ "github.com/Real-Dev-Squad/discord-service/tests/setup"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
)

var testControllerCalled bool = false

func testController(response http.ResponseWriter, request *http.Request, params httprouter.Params) {
	testControllerCalled = true
	response.Header().Set("Content-Type", "application/json; charset=UTF-8")
	response.WriteHeader(http.StatusOK)
}
func setup() *httprouter.Router {
	testControllerCalled = false
	router := httprouter.New()
	router.POST("/", middleware.VerifyCommand(testController))
	return router
}

func TestVerifyCommand(t *testing.T) {

	t.Run("Should verify the verification with valid public key", func(t *testing.T) {
		router := setup()
		var originalFunc = middleware.VerifyInteraction
		middleware.VerifyInteraction = func(r *http.Request, key ed25519.PublicKey) bool {
			return true
		}
		defer func() { middleware.VerifyInteraction = originalFunc }()

		w := httptest.NewRecorder()
		r, _ := http.NewRequest("POST", "/", nil)
		router.ServeHTTP(w, r)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, true, testControllerCalled)
	})

	t.Run("Should fail with Unauthorized error when public key is invalid", func(t *testing.T) {
		router := setup()
		var originalFunc = middleware.VerifyInteraction
		middleware.VerifyInteraction = func(r *http.Request, key ed25519.PublicKey) bool {
			return false
		}
		defer func() { middleware.VerifyInteraction = originalFunc }()

		w := httptest.NewRecorder()
		r, _ := http.NewRequest("POST", "/", nil)
		router.ServeHTTP(w, r)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Should fail with Internal Server Error when public key is malformed", func(t *testing.T) {
		router := setup()
		var originalKey = config.AppConfig.DISCORD_PUBLIC_KEY
		config.AppConfig.DISCORD_PUBLIC_KEY = "malformed_key"
		defer func() { config.AppConfig.DISCORD_PUBLIC_KEY = originalKey }()

		w := httptest.NewRecorder()
		r, _ := http.NewRequest("POST", "/", nil)
		router.ServeHTTP(w, r)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
