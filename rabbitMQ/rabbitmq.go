package rabbitmq

import (
	"log"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const PaymentCreatedQueue = "payment.created"

func ConnectRabbitMQ() (
	*amqp.Connection,
	*amqp.Channel,
) {
	rabbitURL := os.Getenv("RABBITMQ_URL")
	if rabbitURL == "" {
		rabbitURL = "amqp://guest:guest@localhost:5672/"
	}

	var lastErr error
	for attempt := 1; attempt <= 30; attempt++ {
		conn, err := amqp.Dial(rabbitURL)
		if err != nil {
			lastErr = err
			log.Printf("[rabbitmq] aguardando conexao (%d/30): %v", attempt, err)
			time.Sleep(2 * time.Second)
			continue
		}

		ch, err := conn.Channel()
		if err != nil {
			lastErr = err
			_ = conn.Close()
			log.Printf("[rabbitmq] aguardando canal (%d/30): %v", attempt, err)
			time.Sleep(2 * time.Second)
			continue
		}

		log.Println("[rabbitmq] conectado")
		return conn, ch
	}

	log.Fatalf("[rabbitmq] nao foi possivel conectar: %v", lastErr)
	return nil, nil
}

func DeclareQueue(ch *amqp.Channel, queueName string) (amqp.Queue, error) {
	return ch.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
}
