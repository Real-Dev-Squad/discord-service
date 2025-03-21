package utils

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDiscordResponseWithMessage(t *testing.T) {
	response := httptest.NewRecorder()
	message := "Test message"
	Success.NewDiscordResponse(response, message, nil)

	if status := response.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := `{"success": true, "status": 200, "message": "Test message"}`
	if strings.TrimSpace(response.Body.String()) != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", response.Body.String(), expected)
	}
}

func TestNewDiscordResponseWithEncodeError(t *testing.T) {
	response := httptest.NewRecorder()
	data := make(chan int)

	Success.NewDiscordResponse(response, "hello", data)
	assert.Equal(t, http.StatusInternalServerError, response.Code)

}

func TestNewDiscordResponseWithNoError(t *testing.T) {
	response := httptest.NewRecorder()
	Success.NewDiscordResponse(response, "hello", "data")
	assert.Equal(t, http.StatusOK, response.Code)

}
