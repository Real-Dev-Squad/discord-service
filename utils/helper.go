package utils

import (
	"math"
	"time"

	"github.com/sirupsen/logrus"
)

func ExponentialBackoffRetry(maxRetries int, operation func() error) error {
	var err error
	for i := 0; i < maxRetries; i++ {
		err = operation()
		if err == nil {
			return nil
		}
		logrus.Errorf("Attempt %d: Operation failed: %s", i+1, err)
		time.Sleep(time.Duration(math.Pow(3, float64(i))) * time.Second)
	}
	return err
}
