//go:build ignore

package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("🔍 CHECKING CRONJOB ACTIVITY IN PRODUCTION DATABASE")
	fmt.Println("=" + repeat("=", 79))

	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("⚠️  No .env file found, using environment variables")
	}

	// Connect to database
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	if dbPort == "" {
		dbPort = "3306"
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		dbUser, dbPass, dbHost, dbPort, dbName)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("❌ Database connection failed: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("❌ Database ping failed: %v", err)
	}

	fmt.Println("✅ Connected to database:", dbHost)
	fmt.Println()

	// Check 1: Recent badge activity (badge_earn cronjob)
	fmt.Println("📋 CHECK 1: Recent Badge Activity (badge_earn cronjob)")
	fmt.Println("-" + repeat("-", 79))

	var badgeCount int
	var lastBadgeTime *time.Time
	err = db.QueryRow(`
		SELECT COUNT(*) as count, MAX(created_at) as last_activity
		FROM account_badges
		WHERE created_at > DATE_SUB(NOW(), INTERVAL 7 DAY)
	`).Scan(&badgeCount, &lastBadgeTime)

	if err != nil {
		fmt.Printf("❌ Error checking badges: %v\n", err)
	} else {
		fmt.Printf("   Badges awarded in last 7 days: %d\n", badgeCount)
		if lastBadgeTime != nil {
			fmt.Printf("   Last badge awarded: %s\n", lastBadgeTime.Format("2006-01-02 15:04:05"))
		} else {
			fmt.Println("   Last badge awarded: NEVER")
		}
	}

	// Check 2: Incident expiration activity
	fmt.Println("\n📋 CHECK 2: Expired Incidents (incident_expiration cronjob)")
	fmt.Println("-" + repeat("-", 79))

	var expiredCount int
	var lastExpiredTime *time.Time
	err = db.QueryRow(`
		SELECT COUNT(*) as count, MAX(end_time) as last_activity
		FROM incident_clusters
		WHERE is_active = 0 AND end_time > DATE_SUB(NOW(), INTERVAL 7 DAY)
	`).Scan(&expiredCount, &lastExpiredTime)

	if err != nil {
		fmt.Printf("❌ Error checking expired incidents: %v\n", err)
	} else {
		fmt.Printf("   Incidents expired in last 7 days: %d\n", expiredCount)
		if lastExpiredTime != nil {
			fmt.Printf("   Last incident expired: %s\n", lastExpiredTime.Format("2006-01-02 15:04:05"))
		} else {
			fmt.Println("   Last incident expired: NEVER")
		}
	}

	// Check 3: Recent notifications (new_cluster cronjob)
	fmt.Println("\n📋 CHECK 3: Recent Notifications (new_cluster cronjob)")
	fmt.Println("-" + repeat("-", 79))

	var notifCount int
	var lastNotifTime *time.Time
	err = db.QueryRow(`
		SELECT COUNT(*) as count, MAX(created_at) as last_activity
		FROM notifications
		WHERE created_at > DATE_SUB(NOW(), INTERVAL 7 DAY)
	`).Scan(&notifCount, &lastNotifTime)

	if err != nil {
		fmt.Printf("❌ Error checking notifications: %v\n", err)
	} else {
		fmt.Printf("   Notifications created in last 7 days: %d\n", notifCount)
		if lastNotifTime != nil {
			fmt.Printf("   Last notification created: %s\n", lastNotifTime.Format("2006-01-02 15:04:05"))
		} else {
			fmt.Println("   Last notification created: NEVER")
		}
	}

	// Check 4: User rank updates
	fmt.Println("\n📋 CHECK 4: User Rank Updates (user_rank cronjob)")
	fmt.Println("-" + repeat("-", 79))

	var rankedUsersCount int
	var lastRankUpdate *time.Time
	err = db.QueryRow(`
		SELECT COUNT(*) as count, MAX(updated_at) as last_activity
		FROM account
		WHERE rank_score > 0
	`).Scan(&rankedUsersCount, &lastRankUpdate)

	if err != nil {
		fmt.Printf("❌ Error checking user ranks: %v\n", err)
	} else {
		fmt.Printf("   Users with rank score > 0: %d\n", rankedUsersCount)
		if lastRankUpdate != nil {
			fmt.Printf("   Last rank update: %s\n", lastRankUpdate.Format("2006-01-02 15:04:05"))
		} else {
			fmt.Println("   Last rank update: NEVER")
		}
	}

	// Check 5: Bot incidents (if bot_creator is running)
	fmt.Println("\n📋 CHECK 5: Bot-Created Incidents (bot_creator cronjob)")
	fmt.Println("-" + repeat("-", 79))

	var botIncidentCount int
	var lastBotIncident *time.Time
	err = db.QueryRow(`
		SELECT COUNT(*) as count, MAX(created_at) as last_activity
		FROM incident_reports
		WHERE account_id = 1 AND created_at > DATE_SUB(NOW(), INTERVAL 7 DAY)
	`).Scan(&botIncidentCount, &lastBotIncident)

	if err != nil {
		fmt.Printf("❌ Error checking bot incidents: %v\n", err)
	} else {
		fmt.Printf("   Bot incidents in last 7 days: %d\n", botIncidentCount)
		if lastBotIncident != nil {
			fmt.Printf("   Last bot incident: %s\n", lastBotIncident.Format("2006-01-02 15:04:05"))
		} else {
			fmt.Println("   Last bot incident: NEVER")
		}
	}

	// Summary
	fmt.Println("\n" + repeat("=", 80))
	fmt.Println("📊 SUMMARY")
	fmt.Println(repeat("=", 80))

	cronActivity := false
	if badgeCount > 0 || expiredCount > 0 || notifCount > 0 || rankedUsersCount > 0 || botIncidentCount > 0 {
		cronActivity = true
	}

	if cronActivity {
		fmt.Println("✅ CRONJOBS ARE RUNNING - Evidence found in database:")
		if badgeCount > 0 {
			fmt.Printf("   ✅ badge_earn: %d badges awarded\n", badgeCount)
		}
		if expiredCount > 0 {
			fmt.Printf("   ✅ incident_expiration: %d incidents expired\n", expiredCount)
		}
		if notifCount > 0 {
			fmt.Printf("   ✅ new_cluster: %d notifications sent\n", notifCount)
		}
		if rankedUsersCount > 0 {
			fmt.Printf("   ✅ user_rank: %d users have ranks\n", rankedUsersCount)
		}
		if botIncidentCount > 0 {
			fmt.Printf("   ✅ bot_creator: %d bot incidents created\n", botIncidentCount)
		}
	} else {
		fmt.Println("❌ NO CRONJOB ACTIVITY DETECTED")
		fmt.Println("   No evidence of cronjobs running in the last 7 days")
	}
}

func repeat(s string, count int) string {
	result := ""
	for i := 0; i < count; i++ {
		result += s
	}
	return result
}
