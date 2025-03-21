package utils

import (
	"encoding/json"
	"io"
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

var ToByte = func(data interface{}) ([]byte, error) {
	bytes, err := json.Marshal(data)
	if err != nil {
		logrus.Errorf("Failed to marshal message: %v", err)
		return nil, err
	}
	return bytes, nil
}

var FromByte = func(bytes []byte, result interface{}) error {
	err := json.Unmarshal(bytes, result)
	if err != nil {
		logrus.Errorf("Failed to unmarshal message: %v", err)
		return err
	}
	return nil
}

var Encode = func(w io.Writer, data interface{}) error {
	return json.NewEncoder(w).Encode(data)
}
