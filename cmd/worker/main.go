package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/smtp"
	"strings"

	"sysemp_travel/publisher"
	rabbitmq "sysemp_travel/rabbitMQ"
)

func main() {
	conn, ch := rabbitmq.ConnectRabbitMQ()

	defer conn.Close()
	defer ch.Close()

	_, err := rabbitmq.DeclareQueue(ch, rabbitmq.PaymentCreatedQueue)
	if err != nil {
		log.Fatal(err)
	}

	msgs, err := ch.Consume(
		rabbitmq.PaymentCreatedQueue,
		"payment-worker",
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("[worker] iniciado e aguardando mensagens na fila %s", rabbitmq.PaymentCreatedQueue)

	for msg := range msgs {
		log.Printf("[worker] mensagem recebida: %s", string(msg.Body))

		var event publisher.PaymentCreatedEvent

		err := json.Unmarshal(
			msg.Body,
			&event,
		)

		if err != nil {
			log.Printf("[worker] erro ao decodificar mensagem: %s", err)
			_ = msg.Nack(false, false)
			continue
		}

		processPayment(event.PaymentID)

		err = msg.Ack(false)
		if err != nil {
			log.Printf("[worker] erro ao confirmar mensagem: %s", err)
			continue
		}

		log.Printf("[worker] mensagem consumida e confirmada: payment_id=%s", event.PaymentID)
	}
}

func processPayment(paymentID string) {
	// updateCache(paymentID)
	sendNotification(paymentID)
	updateAnalytics(paymentID)
}

// func updateCache(paymentID string) {
// 	log.Printf("[worker] atualizando cache pagamento %s", paymentID)
// }

func sendNotification(paymentID string) {
	log.Printf("[worker] enviando notificacao pagamento %s", paymentID)

	from := "augusto.valenciano@fabricadecodigos.com.br"
	password := "ziyq givc undp dowt"

	to := []string{"augustovalenciano2004@gmail.com"}

	smtpHost := "smtp.gmail.com"
	smtpPort := "587"
	address := smtpHost + ":" + smtpPort

	auth := smtp.PlainAuth("", from, password, smtpHost)

	subject := "Novo pagamento criado"

	html := BuildPurchaseEmail("Augusto", "Passagem Aérea", paymentID, "1500,00")

	message := []byte(
		"From: " + from + "\r\n" +
			"To: " + strings.Join(to, ",") + "\r\n" +
			"Subject: " + subject + "\r\n" +
			"Content-Type: text/html; charset=\"UTF-8\"\r\n" +
			"\r\n" +
			html + "\r\n",
	)

	err := smtp.SendMail(address, auth, from, to, message)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Email enviado com sucesso!")
}

func updateAnalytics(paymentID string) {
	log.Printf("[worker] atualizando analytics para pagamento... %s", paymentID)
}

func BuildPurchaseEmail(customerName string, product string, transactionID string, price string) string {
	html := fmt.Sprintf(` <!DOCTYPE html> <html lang="pt-BR"> <head> <meta charset="UTF-8" /> <meta name="viewport" content="width=device-width, initial-scale=1.0" /> <style> body { background-color: #f4f4f7; font-family: Arial, sans-serif; margin: 0; padding: 0; } .container { max-width: 600px; margin: 40px auto; background: #ffffff; border-radius: 16px; overflow: hidden; box-shadow: 0 4px 20px rgba(0,0,0,0.08); } .header { background: linear-gradient(135deg, #4f46e5, #7c3aed); color: white; padding: 32px; text-align: center; } .header h1 { margin: 0; font-size: 28px; } .content { padding: 32px; color: #333; } .badge { display: inline-block; background: #dcfce7; color: #166534; padding: 8px 14px; border-radius: 999px; font-size: 14px; font-weight: bold; margin-bottom: 20px; } .info-box { background: #f9fafb; border-radius: 12px; padding: 20px; margin-top: 20px; } .row { display: flex; justify-content: space-between; margin-bottom: 12px; font-size: 15px; } .label { color: #6b7280; } .value { font-weight: bold; color: #111827; } .price { font-size: 28px; color: #4f46e5; font-weight: bold; text-align: center; margin-top: 25px; } .footer { text-align: center; padding: 24px; font-size: 13px; color: #9ca3af; background: #f9fafb; } </style> </head> <body> <div class="container"> <div class="header"> <h1>Pagamento Confirmado</h1> </div> <div class="content"> <div class="badge"> Compra Aprovada </div> <p>Olá <strong>%s</strong>,</p> <p> Seu pagamento foi processado com sucesso. Abaixo estão os detalhes da compra: </p> <div class="info-box"> <div class="row"> <span class="label">Produto</span> <span class="value">%s</span> </div> <div class="row"> <span class="label">Transação</span> <span class="value">%s</span> </div> <div class="row"> <span class="label">Status</span> <span class="value">Pago</span> </div> </div> <div class="price"> R$ %s </div> </div> <div class="footer"> Este email é automático. Não responda. </div> </div> </body> </html> `, customerName, product, transactionID, price)
	return html
}
