package emails

import (
	"crypto/tls"
	"log"

	"github.com/go-mail/mail"
)

func Send(email, subject, body string) {
	m := mail.NewMessage()
	// Configura los encabezados.
	m.SetHeader("From", "remitente@example.com")
	m.SetHeader("To", email)
	m.SetHeader("Subject", subject)
	// Configura el cuerpo del mensaje, en este caso HTML.
	m.SetBody("text/html", body)

	d := mail.NewDialer("127.0.0.1", 1025, "", "")
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Envía el correo.
	if err := d.DialAndSend(m); err != nil {
		log.Printf("Error sending emails/emails.go")
		// log.Fatal(err)
	}
}
