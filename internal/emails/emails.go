package emails

import (
	"bytes"
	"html/template"
	"log"
	"os"
	"path/filepath"

	"github.com/resend/resend-go/v2"
)

var resendClient *resend.Client

// InitEmails inicializa el cliente de Resend.
// Debe ser llamado una vez al iniciar la aplicación.
func InitEmails() {
	apiKey := os.Getenv("RESEND_API_KEY")
	if apiKey == "" {
		log.Println("⚠️ WARNING: RESEND_API_KEY not set. Email sending will be disabled.")
		return
	}

	resendClient = resend.NewClient(apiKey)
	log.Println("✅ Resend Email Client Initialized")
}

// InitSES is an alias for InitEmails (backwards compatibility)
func InitSES() {
	InitEmails()
}

// SendTemplate envía un correo HTML basado en un template y datos dinámicos usando Resend.
func SendTemplate(email, subject, templateName string, data any) {
	if resendClient == nil {
		log.Println("Error: Resend client not initialized. Skipping email send.")
		return
	}

	// Cargar y renderizar el template HTML
	tmplBase := filepath.Join("internal", "emails", "templates", "base.html")
	tmplView := filepath.Join("internal", "emails", "templates", templateName+".html")

	tmpl, err := template.ParseFiles(tmplBase, tmplView)
	if err != nil {
		log.Printf("Error parsing templates (%s + base): %v", templateName, err)
		return
	}

	var body bytes.Buffer
	if err := tmpl.ExecuteTemplate(&body, "base", data); err != nil {
		log.Printf("Error rendering template %s: %v", templateName, err)
		return
	}

	// Enviar email vía Resend
	params := &resend.SendEmailRequest{
		From:    "Alertly <no-reply@alertly.ca>",
		To:      []string{email},
		Subject: subject,
		Html:    body.String(),
	}

	_, err = resendClient.Emails.Send(params)
	if err != nil {
		log.Printf("Error sending email to %s via Resend: %v", email, err)
		return
	}

	log.Printf("Email sent to %s via Resend using template %s", email, templateName)
}
