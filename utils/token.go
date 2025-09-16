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
)

type UniqueTokenI interface {
	GenerateUniqueToken() (string, error)
}
type AuthTokenI interface {
    GenerateAuthToken(method jwt.SigningMethod, claims jwt.Claims, privateKey any) (string, error)
}

type UniqueToken struct{}
type AuthToken struct{}

func (ut *UniqueToken) GenerateUniqueToken() (string, error) {
	uuidToken := uuid.NewString()
	
	randNum, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		return "", fmt.Errorf("failed to generate random number: %w", err)
	}

	generationTime := time.Now().UnixMilli()
	combinedString := fmt.Sprintf("%s%d%d", uuidToken, randNum, generationTime)

	hasher := sha256.New()
	hasher.Write([]byte(combinedString))
	hashBytes := hasher.Sum(nil)
	token := hex.EncodeToString(hashBytes)
	
	return token, nil
}

func (at *AuthToken) GenerateAuthToken(method jwt.SigningMethod, claims jwt.Claims, privateKey any) (string, error) {
	token := jwt.NewWithClaims(method, claims)
	
	tokenString, err := token.SignedString(privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %v", err)
	}

	return tokenString, nil
}