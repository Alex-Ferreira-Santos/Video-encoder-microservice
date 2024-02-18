package queue

import (
	"github.com/streadway/amqp"
	"log"
	"os"
)

type RabbitMQ struct {
	User              string
	Password          string
	Host              string
	Port              string
	Vhost             string
	ConsumerQueueName string
	ConsumerName       string
	AutoAck           bool
	Args              amqp.Table
	Channel           *amqp.Channel
}

func NewRabbitMQ() *RabbitMQ {
	rabbitMQArgs := amqp.Table{}
	rabbitMQArgs["x-dead-letter-exchange"] = os.Getenv("RABBITMQ_DLX")

	rabbitMQ := RabbitMQ{
		User: os.Getenv("RABBITMQ_DEFAULT_USER"),
		Password: os.Getenv("RABBITMQ_DEFAULT_PASS"),
		Host: os.Getenv("RABBITMQ_DEFAULT_HOST"),
		Port: os.Getenv("RABBITMQ_DEFAULT_PORT"),
		Vhost: os.Getenv("RABBITMQ_DEFAULT_VHOST"),
		ConsumerQueueName: os.Getenv("RABBITMQ_CONSUMER_QUEUE_NAME"),
		ConsumerName: os.Getenv("RABBITMQ_CONSUMER_NAME"),
		AutoAck: false,
		Args: rabbitMQArgs,
	}

	return &rabbitMQ
}

func (rabbitMQ *RabbitMQ) Connect() *amqp.Channel {
	dsn := "amqp://" + rabbitMQ.User + ":" + rabbitMQ.Password + "@" + rabbitMQ.Host + ":" + rabbitMQ.Port + rabbitMQ.Vhost
	conn, err := amqp.Dial(dsn)
	failOnError(err, "Failed to connect to RabbitMQ")

	rabbitMQ.Channel, err = conn.Channel()
	failOnError(err, "Failed to connect to RabbitMQ")

	return rabbitMQ.Channel
}

func (rabbitMQ *RabbitMQ) Consume(messageChannel chan amqp.Delivery) {
	queue, err := rabbitMQ.Channel.QueueDeclare(
		rabbitMQ.ConsumerQueueName, // name
		true, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		rabbitMQ.Args, // arguments
	)

	failOnError(err, "Failed to declare a queue")


	incomingMessage, err := rabbitMQ.Channel.Consume(
		queue.Name, // queue
		rabbitMQ.ConsumerName, // consumer
		rabbitMQ.AutoAck, // auto-ack
		false, // excluse
		false, // no-local
		false, // no wait
		nil, // args
	)
	failOnError(err, "Failed to register a consumer")


	go func(){
		for message := range incomingMessage {
			log.Println("Incoming new message")
			messageChannel <- message
		}
		log.Println("RabbitMQ channel closed")
		close(messageChannel)
	}()
}

func (rabbitMQ *RabbitMQ) Notify(message string, contentType string, exchange string, routingKey string) error {
	err := rabbitMQ.Channel.Publish(
		exchange, // exchange
		routingKey, // routing key
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: contentType,
			Body: []byte(message),
		},
	)

	if err != nil {
		return err
	}

	return nil
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}