package emails

import (
	"bytes"
	"crypto/tls"
	"html/template"
	"log"
	"path/filepath"

	"github.com/go-mail/mail"
)

// SendTemplate envía un correo HTML basado en un template y datos dinámicos.
func SendTemplate(email, subject, templateName string, data any) {
	// Cargar tanto el template base como el específico (verify_email.html, alert_nearby.html, etc.)
	tmplBase := filepath.Join("internal", "emails", "templates", "base.html")
	tmplView := filepath.Join("internal", "emails", "templates", templateName+".html")

	tmpl, err := template.ParseFiles(tmplBase, tmplView)
	if err != nil {
		log.Printf("Error parsing templates (%s + base): %v", templateName, err)
		return
	}

	// Renderizar el contenido final en un buffer
	var body bytes.Buffer
	if err := tmpl.ExecuteTemplate(&body, "base", data); err != nil {
		log.Printf("Error rendering template %s: %v", templateName, err)
		return
	}

	// Configurar el email
	m := mail.NewMessage()
	m.SetHeader("From", "no-reply@alertly.ca")
	m.SetHeader("To", email)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body.String())

	// Configurar y enviar por SMTP
	d := mail.NewDialer("localhost", 1025, "", "") // Cambiar si usás otro servicio SMTP real
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	if err := d.DialAndSend(m); err != nil {
		log.Printf("Error sending email to %s: %v", email, err)
		return
	}

	log.Printf("Email sent to %s using template %s", email, templateName)
}
