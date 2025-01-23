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

type mockSession struct {
	dialError    error
	channelError error
	queueError   error
}

func (m *mockSession) dial() error {
	return m.dialError
}

func (m *mockSession) createChannel() error {
	return m.channelError
}
func (m *mockSession) declareQueue() error {
	return m.queueError
}

func TestInitQueueConnection(t *testing.T) {
	config.AppConfig.MAX_RETRIES = 1
	t.Run("should not panic when Dial() returns error", func(t *testing.T) {
		mockSess := &mockSession{dialError: errors.New("connection failed")}
		assert.NotPanics(t, func() {
			InitQueueConnection(mockSess)
		}, "InitQueueConnection should not panic when Dial is unsuccessful")

	})

	t.Run("should not panic when CreateChannel() returns error", func(t *testing.T) {
		mockSess := &mockSession{channelError: errors.New("channel failed")}
		assert.NotPanics(t, func() {
			InitQueueConnection(mockSess)
		}, "InitQueueConnection should not panic when CreateChannel is unsuccessful")

	})

	t.Run("should not panic when DeclareQueue() returns error", func(t *testing.T) {
		mockSess := &mockSession{queueError: errors.New("queue failed")}
		assert.NotPanics(t, func() {
			InitQueueConnection(mockSess)
		}, "InitQueueConnection should not when DeclareQueue is unsuccessful")

	})

	t.Run("should pass when no error is returned", func(t *testing.T) {
		mockSess := &mockSession{}
		assert.NotPanics(t, func() {
			InitQueueConnection(mockSess)
		})
	})
}

func TestGetQueueInstance(t *testing.T) {
	t.Run("Should use ExponentialBackoffRetry via GetQueueInstance", func(t *testing.T) {
		attempt := 0
		originalDef := utils.ExponentialBackoffRetry
		utils.ExponentialBackoffRetry = func(maxRetries int, operation func() error) error {
			attempt++
			return errors.New("error")
		}
		defer func() { utils.ExponentialBackoffRetry = originalDef }()
		assert.Nil(t, GetQueueInstance())
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
	t.Run("Should panic when SendMessage returns error", func(t *testing.T) {
		config.AppConfig.MAX_RETRIES = 1
		message := dtos.TextMessage{
			Text:     "test",
			Priority: 1,
		}
		assert.Panics(t, func() {
			SendMessage(message)
		}, "SendMessage should panic when SendMessage returns error")
	})

	t.Run("Should panic when SendMessage returns error", func(t *testing.T) {
		config.AppConfig.MAX_RETRIES = 1
		message := dtos.TextMessage{
			Text:     "test",
			Priority: 1,
		}
		mockQueueSession := &Queue{
			Connection: &amqp.Connection{},
			Name:       "Testing",
			Channel:    &amqp.Channel{},
			Queue:      amqp.Queue{},
		}
		originalFunc := GetQueueInstance
		defer func() { GetQueueInstance = originalFunc }()
		GetQueueInstance = func() *Queue {
			return mockQueueSession
		}
		assert.Panics(t, func() {
			SendMessage(message)
		}, "SendMessage should panic when SendMessage returns error")
	})
}
