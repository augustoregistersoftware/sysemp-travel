package publisher

import (
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type PaymentCreatedEvent struct {
	PaymentID string `json:"payment_id"`
}

type PaymentPublisher struct {
	channel *amqp.Channel
}

func NewPaymentPublisher(ch *amqp.Channel) *PaymentPublisher {
	return &PaymentPublisher{
		channel: ch,
	}
}

func (p *PaymentPublisher) PublishPaymentCreated(paymentID string) error {
	event := PaymentCreatedEvent{
		PaymentID: paymentID,
	}

	body, err := json.Marshal(event)
	if err != nil {
		return err
	}

	err = p.channel.Publish(
		"",
		"payment.created",
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Body:         body,
		},
	)
	if err != nil {
		return err
	}

	log.Printf("[publisher] mensagem publicada na fila payment.created: payment_id=%s", paymentID)
	return nil
}
