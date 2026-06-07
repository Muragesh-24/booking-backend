package scripts

import (
	"fmt"
	"habba/models"
	"net/smtp"
	"os"
	"encoding/json"
	"strings"
	"log"
	amqp "github.com/rabbitmq/amqp091-go"
)

func SendEmailVerificationMail(to string, token string) error {
	smtpUser := os.Getenv("Emailuser")
	smtpPass := os.Getenv("EmailPass")

	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	from := smtpUser

	frontendURL := strings.TrimSuffix(os.Getenv("FRONTEND_URL"), "/")
	if frontendURL == "" {
		frontendURL = "https://kannaddaganeshiitk.vercel.app"
	}

	verificationURL := frontendURL + "/auth/verified?query=" + token

	msg := []byte("From: EngiGrow <" + from + ">\r\n" +
		"To: " + to + "\r\n" +
		"Subject: Verify your email\r\n" +
		"\r\n" +
		"Click the link to verify: " + verificationURL)

	auth := smtp.PlainAuth("", smtpUser, smtpPass, smtpHost)

	return smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, msg)
} 
func EmailVerifyMessage(to string, token string) amqp.Publishing {
	payload := map[string]string{
		"to":    to,
		"token": token,
	}

	body, _ := json.Marshal(payload)

	return amqp.Publishing{
		ContentType: "application/json",
		Body:        body,
	}
}

func EmailInvitation(to string, book models.Booking) error {
	from := os.Getenv("Email")
	pass := os.Getenv("Emailpass")
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	subject := "🎉 Your Ganesh Chaturthi Order is Confirmed!"
	body := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
  <style>
    body { font-family: 'Arial', sans-serif; background: #fff8e1; margin:0; padding:0; }
    .container { max-width: 600px; margin: auto; background: #ffffff; padding: 20px; border-radius: 10px; box-shadow: 0 4px 6px rgba(0,0,0,0.1); text-align: center; }
    h2 { color: #d84315; }
    p { font-size: 16px; color: #444; }
    .ganesha { font-size: 40px; margin-bottom: 10px; }
    .highlight { font-weight: bold; color: #2e7d32; }
  </style>
</head>
<body>
  <div class="container">
    <div class="ganesha">🐘🙏🌸</div>
    <h2>Dear %s,</h2>
    <p>Your order has been <span class="highlight">successfully completed</span>.</p>
    <p>We are eagerly waiting to welcome you on the <span class="highlight">Ganesh Chaturthi event day</span>.</p>
    <p>✨ Ganpati Bappa Morya! ✨</p>
  </div>
</body>
</html>
`, book.Name)

	msg := []byte("To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/html; charset=\"UTF-8\"\r\n\r\n" +
		body)

	auth := smtp.PlainAuth("", from, pass, smtpHost)
	return smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, msg)
}

type EmailVerifyPayload struct {
	To    string `json:"to"`
	Token string `json:"token"`
}

func StartEmailVerificationConsumer() {
	msgs, err := RabbitMQChannel.Consume(
		"mail_queue",
		"",
		true,  // autoAck
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		log.Fatal("Failed to register email consumer:", err)
	}

	go func() {
		for msg := range msgs {
			var payload EmailVerifyPayload

			if err := json.Unmarshal(msg.Body, &payload); err != nil {
				log.Println("Invalid email payload:", err)
				continue
			}

			err := SendEmailVerificationMail(payload.To, payload.Token)
			if err != nil {
				log.Println("Failed to send verification email:", err)
				continue
			}

			log.Println("Verification email sent to:", payload.To)
		}
	}()
}