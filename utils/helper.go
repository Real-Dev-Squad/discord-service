package utils

import (
	"math"
	"time"

	"github.com/sirupsen/logrus"
)

var ExponentialBackoffRetry = func(maxRetries int, operation func() error) error {
	var err error
	for i := 0; i < maxRetries; i++ {
		err = operation()
		if err == nil {
			return nil
		}
		logrus.Errorf("Attempt %d: Operation failed: %s", i+1, err)
		if i < maxRetries-1 {
			time.Sleep(time.Duration(math.Pow(2, float64(i))) * time.Second)
		}
	}
	return err
}

// TODO: Implement the logic for generating a unique token.
func GenerateToken(uuidToken []byte) (string, error) {
	return "", nil
}
