package scripts

import (
	"fmt"
	"habba/models"
	"net/smtp"
	"os"
)

func EmailVerifymail(to string, token string) error {
    from := os.Getenv("Email")
    pass := os.Getenv("Emailpass")
    smtpHost := "smtp.gmail.com"
    smtpPort := "587"

    msg := []byte("Subject: Verify your email\n\n" +
        "Click the link to verify: https://booking-backend-csn1.onrender.com/user/verify?token=" + token)

    auth := smtp.PlainAuth("", from, pass, smtpHost)

    return smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, msg)
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
