package queue

import (
	"sync"

	"github.com/Real-Dev-Squad/discord-service/config"
	"github.com/Real-Dev-Squad/discord-service/dtos"
	"github.com/Real-Dev-Squad/discord-service/utils"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
)

type Queue struct {
	Connection *amqp.Connection
	Queue      amqp.Queue
	Name       string
	Channel    *amqp.Channel
}

func (q *Queue) dial() error {
	var err error
	q.Connection, err = amqp.Dial("amqp://localhost")
	return err
}

func (q *Queue) createChannel() error {
	var err error
	q.Channel, err = q.Connection.Channel()
	return err
}

func (q *Queue) declareQueue() error {
	var err error
	q.Queue, err = q.Channel.QueueDeclare(
		config.AppConfig.QUEUE_NAME,     // name
		true,                            // durable
		false,                           // delete when unused
		false,                           // exclusive
		false,                           // no-wait
		amqp.Table{"x-max-priority": 2}, // arguments
	)
	return err
}

var (
	queueInstance *Queue
	once          sync.Once
)

type sessionInterface interface {
	dial() error
	createChannel() error
	declareQueue() error
}

func InitQueueConnection(openSession sessionInterface) {
	var err error
	f := func() error {
		err = openSession.dial()
		if err != nil {
			return err
		}
		err = openSession.createChannel()
		if err != nil {
			return err
		}
		err = openSession.declareQueue()
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
	queueInstance := &Queue{}
	InitQueueConnection(queueInstance)
}

var GetQueueInstance = func() *Queue {
	once.Do(queueHandler)
	return queueInstance
}

func SendMessage(message dtos.TextMessage) {
	queue := GetQueueInstance()
	err := queue.Channel.Publish(
		"",
		queue.Name,
		true,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message.Text),
			Priority:    message.Priority,
		})
	if err != nil {
		logrus.Errorf("Failed to publish a message: %s", err)
	}
	logrus.Info("Message Sent Successfully")
}
