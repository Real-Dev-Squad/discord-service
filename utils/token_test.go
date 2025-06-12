package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"errors"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

type faultyReader struct{}

func (r *faultyReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("reader error")
}

func TestGenerateUniqueToken(t *testing.T) {
	t.Run("should generate a unique token successfully", func(t *testing.T) {
		uniqueToken := &UniqueToken{}
		token, err := uniqueToken.GenerateUniqueToken()
		assert.NoError(t, err)
		assert.NotEmpty(t, token)
		assert.Len(t, token, 64) // SHA-256 hash is 64 hex characters
	})

	t.Run("should return an error when random number generation fails", func(t *testing.T) {
		originalReader := rand.Reader
		rand.Reader = &faultyReader{}
		defer func() {
			rand.Reader = originalReader
		}()
		
		uniqueToken := &UniqueToken{}
		token, err := uniqueToken.GenerateUniqueToken()
		assert.Error(t, err)
		assert.Empty(t, token)
		assert.Contains(t, err.Error(), "failed to generate random number")
	})
}

func TestGenerateAuthToken(t *testing.T) {
	t.Run("should generate an HS256 auth token successfully", func(t *testing.T) {
		claims := jwt.MapClaims{
			"username": "test",
			"exp":      time.Now().Add(time.Hour * 1).Unix(),
		}
		secret := []byte("secret-key")

		authToken := &AuthToken{}
		token, err := authToken.GenerateAuthToken(jwt.SigningMethodHS256, claims, secret)
		assert.NoError(t, err)
		assert.NotEmpty(t, token)

		parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
			return secret, nil
		})
		assert.NoError(t, err)
		assert.True(t, parsedToken.Valid)
	})

	t.Run("should generate an RS256 auth token successfully", func(t *testing.T) {
		privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
		assert.NoError(t, err)

		claims := jwt.MapClaims{
			"username": "test",
			"exp":      time.Now().Add(time.Hour * 1).Unix(),
		}

		authToken := &AuthToken{}
		token, err := authToken.GenerateAuthToken(jwt.SigningMethodRS256, claims, privateKey)
		assert.NoError(t, err)
		assert.NotEmpty(t, token)

		parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
			return &privateKey.PublicKey, nil
		})
		assert.NoError(t, err)
		assert.True(t, parsedToken.Valid)
	})

	t.Run("should return an error for invalid key type", func(t *testing.T) {
		claims := jwt.MapClaims{
			"username": "test",
			"exp":      time.Now().Add(time.Hour * 1).Unix(),
		}
		// Using a string as a key for RS256, which expects *rsa.PrivateKey
		invalidKey := "not-a-real-key"

		authToken := &AuthToken{}
		token, err := authToken.GenerateAuthToken(jwt.SigningMethodRS256, claims, invalidKey)
		assert.Error(t, err)
		assert.Empty(t, token)
	})
}
