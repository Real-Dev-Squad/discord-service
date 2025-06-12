package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"errors"
	"testing"

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
		assert.Len(t, token, 64)
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
	t.Run("should generate an RS256 auth token successfully", func(t *testing.T) {
		privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
		assert.NoError(t, err)
		
		name := "Discord Service"
		claims := jwt.MapClaims{
			"name": name,
		}
		
		authToken := &AuthToken{}
		token, err := authToken.GenerateAuthToken(jwt.SigningMethodRS256, claims, privateKey)

		assert.NoError(t, err)
		assert.NotEmpty(t, token)

		parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
			return &privateKey.PublicKey, nil
		})
		assert.NoError(t, err)
		res, _ := parsedToken.Claims.(jwt.MapClaims)["name"]
		assert.Equal(t, name, res)
	})

	t.Run("should return an error for invalid key type", func(t *testing.T) {
		claims := jwt.MapClaims{}
		invalidKey := "<invalid-rsa-key>"
		authToken := &AuthToken{}
		token, err := authToken.GenerateAuthToken(jwt.SigningMethodRS256, claims, invalidKey)
		assert.Error(t, err)
		assert.Empty(t, token)
	})
}
