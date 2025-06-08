package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type tokenHelper struct{}

var TokenHelper = &tokenHelper{}

// GenerateUniqueToken creates a secure, unique token by hashing a combination of a UUID,
// a cryptographically secure random number, and the current time.
func (t *tokenHelper) GenerateUniqueToken() (string, error) {
	// 1. Generate a new UUID
	uuidToken := uuid.NewString()

	// 2. Generate a random number up to 1,000,000
	randNum, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		logrus.Errorf("Error generating random number: %v", err)
		return "", fmt.Errorf("failed to generate random number: %w", err)
	}

	// 3. Get the current timestamp in milliseconds
	generationTime := time.Now().UnixMilli()

	// 4. Concatenate the parts into a single string
	combinedString := fmt.Sprintf("%s%d%d", uuidToken, randNum, generationTime)
	fmt.Println("Combined string: ", combinedString)
	// 5. Calculate the SHA-256 hash
	hasher := sha256.New()
	hasher.Write([]byte(combinedString))
	hashBytes := hasher.Sum(nil)

	// 6. Encode the hash to a hexadecimal string
	token := hex.EncodeToString(hashBytes)
	fmt.Println("Token: ", token)
	return token, nil
}
