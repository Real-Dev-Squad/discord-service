package queue

import (
	"errors"
	"sync"

	"github.com/Real-Dev-Squad/discord-service/config"
	"github.com/Real-Dev-Squad/discord-service/utils"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
)

var (
	queueInstance *QueueWrapper
	once          sync.Once
)

type QueueInterface interface {
	Dial() error
	CreateChannel() error
	DeclareQueue() error
	PublishMessage(message []byte) error
}

type QueueWrapper struct {
	Connection *amqp.Connection
	Queue      amqp.Queue
	Name       string
	ChannelFn  func() (*amqp.Channel, error)
	Channel    *amqp.Channel
	Publish    func(exchange string, key string, mandatory bool, immediate bool, msg amqp.Publishing) error
}

func (q *QueueWrapper) Dial() error {
	var err error
	q.Connection, err = utils.AMQPDial(config.AppConfig.QUEUE_URL)
	if err != nil {
		logrus.Errorf("Failed to establish connection to RabbitMQ: %v", err)
		return err
	}
	q.ChannelFn = q.Connection.Channel
	return nil
}

func (q *QueueWrapper) CreateChannel() error {
	var err error
	q.Channel, err = q.ChannelFn()
	if err != nil {
		return err
	}
	q.Publish = q.Channel.Publish
	return nil
}

func (q *QueueWrapper) DeclareQueue() error {
	var err error
	q.Queue, err = q.Channel.QueueDeclare(
		config.AppConfig.QUEUE_NAME,
		true,
		false,
		false,
		false,
		amqp.Table{"x-max-priority": 2},
	)
	return err
}

func (q *QueueWrapper) PublishMessage(message []byte) error {
	if err := q.Publish(
		"",
		q.Queue.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        message,
		},
	); err != nil {
		logrus.Errorf("Failed to publish message: %v", err)
		return err
	}
	logrus.Info("Message sent successfully")
	return nil
}

func InitQueueConnection(queue QueueInterface) {
	var err error
	f := func() error {
		err = queue.Dial()
		if err != nil {
			return err
		}
		err = queue.CreateChannel()
		if err != nil {
			return err
		}
		err = queue.DeclareQueue()
		return err
	}

	err = utils.ExponentialBackoffRetry(config.AppConfig.MAX_RETRIES, f)
	if err != nil {
		logrus.Errorf("Failed to initialize queue after %d attempts: %s", config.AppConfig.MAX_RETRIES, err)
		return
	}
	logrus.Infof("Established a connection to RabbitMQ named %s", config.AppConfig.QUEUE_NAME)
}

func queueHandler() {
	queueInstance = &QueueWrapper{}
	InitQueueConnection(queueInstance)
}

var GetQueueInstance = func() *QueueWrapper {
	once.Do(queueHandler)
	return queueInstance
}

var SendMessage = func(message []byte) error {
	queue := GetQueueInstance()

	if queue.Channel == nil {
		logrus.Errorf("Queue channel is not initialized")
		return errors.New("Queue channel is not initialized")
	}

	return queue.PublishMessage(message)

}
