package scripts;

import (
	"fmt"
	"net/url"
	"time"
"os"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
	
)

var (
	RabbitMQConn    *amqp.Connection
	RabbitMQChannel *amqp.Channel
	MailQueueName   string
	ModerationQueue string
)

func RabbitMQConnection() {
	host :=  os.Getenv("rabbithost")
	port := os.Getenv("rabbitport")
	user :=  os.Getenv("rabbituser")
	password := os.Getenv("rabbitpassword")

	MailQueueName = os.Getenv("rabbitmailqueue")


	logrus.Infof("RabbitMQ host=%q user=%q port=%q", host, user, port)

	var connErr error
	for attempt := 1; attempt <= 10; attempt++ {
		connErr = connectRabbitMQ(host, port, user, password)
		if connErr == nil {
			logrus.Info("Connected to RabbitMQ")
			return
		}

		logrus.WithError(connErr).WithField("attempt", attempt).Warn("RabbitMQ connection failed, retrying")
		time.Sleep(2 * time.Second)
	}

	logrus.Fatal("Failed to connect to RabbitMQ: ", connErr)
}

func connectRabbitMQ(host, port, user, password string) error {

	rabbitURL := url.URL{
		Scheme: "amqp",
		User:   url.UserPassword(user, password),
		Host:   fmt.Sprintf("%s:%s", host, port),
		Path:   "/",
	}

	conn, err := amqp.Dial(rabbitURL.String())
	if err != nil {
		return err
	}

	channel, err := conn.Channel()
	if err != nil {
		_ = conn.Close()
		return err
	}

	if _, err := channel.QueueDeclare(MailQueueName, true, false, false, false, nil); err != nil {
		_ = channel.Close()
		_ = conn.Close()
		return err
	}



	if err := channel.Qos(1, 0, false); err != nil {
		_ = channel.Close()
		_ = conn.Close()
		return err
	}

	RabbitMQConn = conn
	RabbitMQChannel = channel
	return nil
}

