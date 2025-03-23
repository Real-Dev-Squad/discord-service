package queue

import (
	"errors"
	"testing"

	"github.com/Real-Dev-Squad/discord-service/config"
	"github.com/Real-Dev-Squad/discord-service/dtos"
	_ "github.com/Real-Dev-Squad/discord-service/tests/setup"
	"github.com/Real-Dev-Squad/discord-service/utils"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/assert"
)

type MockQueue struct {
	DialError    error
	ChannelError error
	QueueError   error
	PublishError error
}

func (m *MockQueue) Dial() error {
	return m.DialError
}

func (m *MockQueue) CreateChannel() error {
	return m.ChannelError
}

func (m *MockQueue) DeclareQueue() error {
	return m.QueueError
}
func (m *MockQueue) PublishMessage(message []byte) error {
	return m.PublishError
}

func TestInitQueueConnection(t *testing.T) {
	config.AppConfig.MAX_RETRIES = 1
	t.Run("should not panic when Dial() returns error", func(t *testing.T) {
		mockQueue := &MockQueue{DialError: errors.New("connection failed")}
		assert.NotPanics(t, func() {
			InitQueueConnection(mockQueue)
		}, "InitQueueConnection should not panic when Dial is unsuccessful")
	})
	t.Run("should not throw error when Dial() is successful", func(t *testing.T) {
		originalFunc := utils.AMQPDial
		defer func() { utils.AMQPDial = originalFunc }()
		utils.AMQPDial = func(url string) (*amqp.Connection, error) {
			return &amqp.Connection{}, nil
		}
		mockQueue := &QueueWrapper{}

		assert.NoError(t, mockQueue.Dial())
	})
	t.Run("should not panic when CreateChannel() returns error", func(t *testing.T) {
		mockQueue := &MockQueue{ChannelError: errors.New("channel failed")}
		assert.NotPanics(t, func() {
			InitQueueConnection(mockQueue)
		}, "InitQueueConnection should not panic when CreateChannel is unsuccessful")
	})
	t.Run("should not panic when DeclareQueue() returns error", func(t *testing.T) {
		mockQueue := &MockQueue{QueueError: errors.New("queue failed")}
		assert.NotPanics(t, func() {
			InitQueueConnection(mockQueue)
		}, "InitQueueConnection should not panic when DeclareQueue is unsuccessful")
	})

	t.Run("should pass when no error is returned", func(t *testing.T) {
		mockQueue := &MockQueue{}
		assert.NotPanics(t, func() {
			InitQueueConnection(mockQueue)
		})
	})

	t.Run("should not return error when PublishMessage() is successful", func(t *testing.T) {
		mockQueue := &QueueWrapper{Publish: func(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error {
			return nil
		}}
		assert.NoError(t, mockQueue.PublishMessage([]byte("message")))
	})

	t.Run("should return error when PublishMessage() returns error", func(t *testing.T) {
		mockQueue := &QueueWrapper{Publish: func(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error {
			return assert.AnError
		}}
		assert.Error(t, mockQueue.PublishMessage([]byte("message")))
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

func TestQueueWrapper(t *testing.T) {
	queueWrapper := &QueueWrapper{}

	t.Run("QueueWrapper should always implement Dial() method", func(t *testing.T) {
		err := queueWrapper.Dial()
		assert.Error(t, err)
	})

	t.Run("QueueWrapper should always implement CreateChannel() method", func(t *testing.T) {
		queueWrapper.Connection = &amqp.Connection{}
		assert.Panics(t, func() {
			queueWrapper.CreateChannel()
		})
	})
	t.Run("CreateChannel should always handle error from ChannelFn() method", func(t *testing.T) {
		mockQueue := &QueueWrapper{ChannelFn: func() (*amqp.Channel, error) {
			return nil, errors.New("channel error")
		}}
		assert.Error(t, mockQueue.CreateChannel())
	})

	t.Run("CreateChannel should not throw error from ChannelFn() method", func(t *testing.T) {
		mockQueue := &QueueWrapper{ChannelFn: func() (*amqp.Channel, error) {
			return &amqp.Channel{}, nil
		}}
		assert.NoError(t, mockQueue.CreateChannel())
	})
	t.Run("QueueWrapper should always implement DeclareQueue() method", func(t *testing.T) {
		queueWrapper.Channel = &amqp.Channel{}
		assert.Panics(t, func() {
			queueWrapper.DeclareQueue()
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
		}, "SendMessage should not panic when SendMessage returns error")
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
		GetQueueInstance = func() *QueueWrapper {
			return &QueueWrapper{
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
