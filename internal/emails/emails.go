package emails

import (
	"bytes"
	"context"
	"html/template"
	"log"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/aws-sdk-go-v2/service/sesv2/types"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

var sesClient *sesv2.Client

// InitSES inicializa el cliente de AWS SES v2.
// Debe ser llamado una vez al iniciar la aplicación.
func InitSES() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Printf("🚨 CRITICAL: Failed to load AWS SDK config for SES. Email sending will be disabled. Error: %v", err)
		return // No hacer crash, solo registrar y continuar.
	}

	// En Fargate, el Task Role ARN no se asume automáticamente por el SDK.
	// Debemos crear un proveedor de credenciales que asuma explícitamente el rol.
	taskRoleARN := os.Getenv("AWS_IAM_ROLE_ARN") // ECS inyecta esta variable si hay un Task Role
	if taskRoleARN != "" {
		log.Printf("Found Task Role ARN: %s. Attempting to assume role for SES.", taskRoleARN)
		stsClient := sts.NewFromConfig(cfg)
		provider := stscreds.NewAssumeRoleProvider(stsClient, taskRoleARN)
		cfg.Credentials = aws.NewCredentialsCache(provider)
	} else {
		log.Println("No AWS_IAM_ROLE_ARN found. Using default credential chain for SES.")
	}

	sesClient = sesv2.NewFromConfig(cfg)
	log.Println("✅ AWS SES Client Initialized")
}

// SendTemplate envía un correo HTML basado en un template y datos dinámicos usando AWS SES.
func SendTemplate(email, subject, templateName string, data any) {
	if sesClient == nil {
		log.Println("Error: SES client not initialized. Skipping email send.")
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

	// Configurar y enviar el email vía AWS SES
	input := &sesv2.SendEmailInput{
		FromEmailAddress: aws.String("no-reply@alertly.ca"),
		Destination: &types.Destination{
			ToAddresses: []string{email},
		},
		Content: &types.EmailContent{
			Simple: &types.Message{
				Subject: &types.Content{
					Data:    aws.String(subject),
					Charset: aws.String("UTF-8"),
				},
				Body: &types.Body{
					Html: &types.Content{
						Data:    aws.String(body.String()),
						Charset: aws.String("UTF-8"),
					},
				},
			},
		},
	}

	_, err = sesClient.SendEmail(context.TODO(), input)
	if err != nil {
		log.Printf("Error sending email to %s via SES: %v", email, err)
		return
	}

	log.Printf("Email sent to %s via SES using template %s", email, templateName)
}
