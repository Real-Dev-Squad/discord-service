package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type UniqueTokenHelper interface {
	GenerateUniqueToken() (string, error)
}
type AuthTokenHelper interface {
	GenerateAuthToken(method jwt.SigningMethod, claims jwt.Claims, privateKey any) (string, error)
}

type UniqueToken struct{}
type AuthToken struct{}

func (t *UniqueToken) GenerateUniqueToken() (string, error) {
	uuidToken := uuid.NewString()
	randNum, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		logrus.Errorf("Error generating random number: %v", err)
		return "", fmt.Errorf("failed to generate random number: %w", err)
	}

	generationTime := time.Now().UnixMilli()
	combinedString := fmt.Sprintf("%s%d%d", uuidToken, randNum, generationTime)

	hasher := sha256.New()
	if _, err := hasher.Write([]byte(combinedString)); err != nil {
		return "", fmt.Errorf("failed to write to hasher: %w", err)
	}

	hashBytes := hasher.Sum(nil)
	token := hex.EncodeToString(hashBytes)
	return token, nil
}

func (t *AuthToken) GenerateAuthToken(method jwt.SigningMethod, claims jwt.Claims, privateKey any) (string, error) {
	token := jwt.NewWithClaims(method, claims)
	tokenString, err := token.SignedString(privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}
	return tokenString, nil
}
