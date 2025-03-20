package queue

import (
	"errors"
	"testing"

	"github.com/Real-Dev-Squad/discord-service/config"
	"github.com/Real-Dev-Squad/discord-service/dtos"
	_ "github.com/Real-Dev-Squad/discord-service/tests/helpers"
	"github.com/Real-Dev-Squad/discord-service/utils"
	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/stretchr/testify/assert"
)

type mockQueue struct {
	dialError    error
	channelError error
	queueError   error
}

func (m *mockQueue) dial() error {
	return m.dialError
}

func (m *mockQueue) createChannel() error {
	return m.channelError
}
func (m *mockQueue) declareQueue() error {
	return m.queueError
}

func TestInitQueueConnection(t *testing.T) {
	config.AppConfig.MAX_RETRIES = 1
	t.Run("should not panic when Dial() returns error", func(t *testing.T) {
		mockQueue := &mockQueue{dialError: errors.New("connection failed")}
		assert.NotPanics(t, func() {
			InitQueueConnection(mockQueue)
		}, "InitQueueConnection should not panic when Dial is unsuccessful")

	})

	t.Run("should not panic when CreateChannel() returns error", func(t *testing.T) {
		mockQueue := &mockQueue{channelError: errors.New("channel failed")}
		assert.NotPanics(t, func() {
			InitQueueConnection(mockQueue)
		}, "InitQueueConnection should not panic when CreateChannel is unsuccessful")

	})

	t.Run("should not panic when DeclareQueue() returns error", func(t *testing.T) {
		mockQueue := &mockQueue{queueError: errors.New("queue failed")}
		assert.NotPanics(t, func() {
			InitQueueConnection(mockQueue)
		}, "InitQueueConnection should not when DeclareQueue is unsuccessful")

	})

	t.Run("should pass when no error is returned", func(t *testing.T) {
		mockQueue := &mockQueue{}
		assert.NotPanics(t, func() {
			InitQueueConnection(mockQueue)
		})
	})
}

func TestGetQueueInstance(t *testing.T) {
	t.Run("Should use ExponentialBackoffRetry via GetQueueInstance", func(t *testing.T) {
		attempt := 0
		originalFunc := utils.ExponentialBackoffRetry
		utils.ExponentialBackoffRetry = func(maxRetries int, operation func() error) error {
			attempt++
			return errors.New("error")
		}
		defer func() { utils.ExponentialBackoffRetry = originalFunc }()
		assert.NotNil(t, GetQueueInstance())
		assert.Equal(t, 1, attempt)
	})
}

func TestSessionWrapper(t *testing.T) {
	sessionWrapper := &Queue{}

	t.Run("SessionWrapper should always implement dial() method", func(t *testing.T) {
		err := sessionWrapper.dial()
		assert.Error(t, err)
	})

	t.Run("SessionWrapper should always implement createChannel() method", func(t *testing.T) {
		sessionWrapper.Connection = &amqp.Connection{}
		assert.Panics(t, func() {
			sessionWrapper.createChannel()
		})

	})

	t.Run("SessionWrapper should always implement declareQueue() method", func(t *testing.T) {
		sessionWrapper.Channel = &amqp.Channel{}
		assert.Panics(t, func() {
			sessionWrapper.declareQueue()
		})
	})

}

func TestSendMessage(t *testing.T) {
	t.Run("Should not panic when SendMessage returns error", func(t *testing.T) {
		config.AppConfig.MAX_RETRIES = 1
		message := dtos.DataPacket{
			UserID:      "1",
			CommandName: utils.CommandNames.Listening,
		}
		bytes, err := utils.ToByte(message)
		assert.NoError(t, err)
		assert.NotPanics(t, func() {
			SendMessage(bytes)
		}, "SendMessage should panic when SendMessage returns error")
	})
	t.Run("Should panic when trying to send message from queue without proper connection", func(t *testing.T) {
		config.AppConfig.MAX_RETRIES = 1
		message := dtos.DataPacket{
			UserID:      "1",
			CommandName: utils.CommandNames.Listening,
		}
		bytes, err := utils.ToByte(message)
		assert.NoError(t, err)
		originalFunc := utils.ExponentialBackoffRetry
		originalQueueInstance := GetQueueInstance
		utils.ExponentialBackoffRetry = func(maxRetries int, operation func() error) error {
			return nil
		}
		GetQueueInstance = func() *Queue {
			return &Queue{
				Channel: &amqp.Channel{},
			}
		}
		defer func() {
			utils.ExponentialBackoffRetry = originalFunc
			GetQueueInstance = originalQueueInstance
		}()

		assert.Panics(t, func() {
			SendMessage(bytes)
		}, "SendMessage should panic when SendMessage returns error")
	})
}
