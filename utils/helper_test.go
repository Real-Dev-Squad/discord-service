package utils

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExponentialBackoffRetry_Success(t *testing.T) {
	attempts := 0
	operation := func() error {
		attempts++
		if attempts < 3 {
			return errors.New("temporary error")
		}
		return nil
	}

	err := ExponentialBackoffRetry(5, operation)
	assert.NoError(t, err)
	assert.Equal(t, 3, attempts)
}

func TestExponentialBackoffRetry_Failure(t *testing.T) {
	attempts := 0
	operation := func() error {
		attempts++
		return errors.New("permanent error")
	}

	err := ExponentialBackoffRetry(3, operation)
	assert.Error(t, err)
	assert.Equal(t, 3, attempts)
}

func TestExponentialBackoffRetry_NoRetries(t *testing.T) {
	attempts := 0
	operation := func() error {
		attempts++
		return errors.New("error")
	}

	err := ExponentialBackoffRetry(0, operation)
	assert.Nil(t, err)
	assert.Equal(t, 0, attempts)
}

func TestExponentialBackoffRetry_ImmediateSuccess(t *testing.T) {
	attempts := 0
	operation := func() error {
		attempts++
		return nil
	}
	err := ExponentialBackoffRetry(5, operation)
	assert.NoError(t, err)
	assert.Equal(t, 1, attempts)
}
func TestToByte(t *testing.T) {
	t.Run("should convert string to byte", func(t *testing.T) {
		data := "test"
		_, err := ToByte(data)
		assert.NoError(t, err)
	})
	t.Run("should return error if conversion fails", func(t *testing.T) {
		data := make(chan int)
		_, err := ToByte(data)
		assert.Error(t, err)
	})
}

func TestFromByte(t *testing.T) {
	t.Run("should convert string from byte", func(t *testing.T) {
		data := "test"
		bytes, err := ToByte(data)
		assert.NoError(t, err)

		var result string
		err = FromByte(bytes, &result)
		assert.NoError(t, err)
		assert.Equal(t, data, result)
	})
	t.Run("should return error if conversion fails", func(t *testing.T) {
		var result interface{}
		data := "test"
		err := FromByte([]byte(data), &result)
		assert.Error(t, err)
	})
}
