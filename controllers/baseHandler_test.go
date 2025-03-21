package controllers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Real-Dev-Squad/discord-service/controllers"
	_ "github.com/Real-Dev-Squad/discord-service/tests/helpers"
	"github.com/stretchr/testify/assert"
)

func TestDiscordBaseHandler(t *testing.T) {
	t.Run("Should return 400 when request body is empty", func(t *testing.T) {
		w := httptest.NewRecorder()
		r, _ := http.NewRequest("POST", "/", http.NoBody)
		controllers.DiscordBaseHandler(w, r, nil)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
