package controllers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Real-Dev-Squad/discord-service/controllers"
	utils "github.com/Real-Dev-Squad/discord-service/tests"
	_ "github.com/Real-Dev-Squad/discord-service/tests/helpers"
	"github.com/stretchr/testify/assert"
)

func TestWSHandler_UpgradeError(t *testing.T) {
	// Use an invalid HTTP method to trigger upgrade failure
	r := httptest.NewRequest(http.MethodPost, "/ws", nil)
	w := httptest.NewRecorder()

	controllers.WSHandler(w, r)

	result := w.Result()
	assert.Equal(t, http.StatusBadRequest, result.StatusCode, "Expected status code 400 on upgrade failure")
}

func TestWSHandler_Success(t *testing.T) {
	conn, resp, err, defFunc := utils.SocketConnection(controllers.WSHandler)
	defer defFunc()

	assert.NoError(t, err, "Error should be nil when establishing connection")
	assert.NotNil(t, conn, "Connection should be established")
	assert.Equal(t, http.StatusSwitchingProtocols, resp.StatusCode, "Status code should be 101")
}

func TestWSHandler_Disconnect(t *testing.T) {
	session, _, err, defFunc := utils.SocketConnection(controllers.WSHandler)
	defer defFunc()
	assert.NoError(t, err, "Error should be nil when establishing connection")
	assert.NotNil(t, session, "Session should be established")
	err = session.Close()
	assert.NoError(t, err, "Closing connection should not produce an error")
}
