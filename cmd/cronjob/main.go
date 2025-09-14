package main

import (
	"alertly/internal/cronjobs/cjbadgeearn"
	"alertly/internal/cronjobs/cjblockincident"
	"alertly/internal/cronjobs/cjblockuser"
	"alertly/internal/cronjobs/cjcomments"
	"alertly/internal/cronjobs/cjinactivityreminder"
	"alertly/internal/cronjobs/cjincidentexpiration"
	"alertly/internal/cronjobs/cjincidentupdate"
	"alertly/internal/cronjobs/cjnewcluster"
	"alertly/internal/cronjobs/cjuserank"
	"alertly/internal/database" // Use the centralized database package
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	// "github.com/joho/godotenv" // No longer needed
)

// Event is the input structure for the Lambda function
type Event struct {
	Task string `json:"task"`
}

// initDB initializes the database connection.
// This function is now simplified for the Lambda environment.
func initDB() {
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true",
		dbUser, dbPass, dbHost, dbPort, dbName,
	)

	database.InitDB(dsn) // Call the correct InitDB function
}

// HandleRequest is the main Lambda handler function.
func HandleRequest(ctx context.Context, event Event) (string, error) {
	initDB() // Initialize DB connection for each invocation
	defer database.DB.Close()

	log.Printf("Executing cronjob task: %s", event.Task)

	switch event.Task {
	case "new_cluster":
		repo := cjnewcluster.NewRepository(database.DB)
		svc := cjnewcluster.NewService(repo)
		svc.Run()
	case "block_user":
		repo := cjblockuser.NewRepository(database.DB)
		svc := cjblockuser.NewService(repo)
		svc.Run()
	case "block_incident":
		repo := cjblockincident.NewRepository(database.DB)
		svc := cjblockincident.NewService(repo)
		svc.Run()
	case "inactivity_reminder":
		repo := cjinactivityreminder.NewRepository(database.DB)
		svc := cjinactivityreminder.NewService(repo)
		svc.Run()
	case "badge_earn":
		repo := cjbadgeearn.NewRepository(database.DB)
		svc := cjbadgeearn.NewService(repo)
		svc.Run()
	case "user_rank":
		repo := cjuserank.NewRepository(database.DB)
		svc := cjuserank.NewService(repo)
		svc.Run()
	case "comments":
		repo := cjcomments.NewRepository(database.DB)
		svc := cjcomments.NewService(repo)
		svc.Run()
	case "incident_update":
		repo := cjincidentupdate.NewRepository(database.DB)
		svc := cjincidentupdate.NewService(repo)
		svc.Run()
	case "incident_expiration":
		repo := cjincidentexpiration.NewRepository(database.DB)
		svc := cjincidentexpiration.NewService(repo)
		svc.Run()
	default:
		log.Printf("Unknown task: %s", event.Task)
		return fmt.Sprintf("Unknown task: %s", event.Task), nil
	}

	log.Printf("Task %s completed successfully.", event.Task)
	return fmt.Sprintf("Task %s completed successfully.", event.Task), nil
}

func main() {
	lambda.Start(HandleRequest)
}
