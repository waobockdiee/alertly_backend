package cronjob

import (
	"database/sql"
	"log"
	"time"
)

type PremiumExpirationService struct {
	db *sql.DB
}

func NewPremiumExpirationService(db *sql.DB) *PremiumExpirationService {
	return &PremiumExpirationService{db: db}
}

// ExpiredAccount represents an account with expired premium
type ExpiredAccount struct {
	AccountID          int64     `json:"account_id"`
	Email              string    `json:"email"`
	PremiumExpiredDate time.Time `json:"premium_expired_date"`
}

// CheckAndExpirePremiumAccounts checks for expired premium accounts and deactivates them
// This should be run once daily (recommended: 2 AM) via cronjob
func (s *PremiumExpirationService) CheckAndExpirePremiumAccounts() error {
	log.Println("🔍 Starting daily premium expiration check...")

	// Find accounts with expired premium
	expiredAccounts, err := s.findExpiredAccounts()
	if err != nil {
		log.Printf("❌ Error finding expired accounts: %v", err)
		return err
	}

	if len(expiredAccounts) == 0 {
		log.Println("✅ No expired premium accounts found")
		return nil
	}

	log.Printf("⚠️ Found %d expired premium accounts", len(expiredAccounts))

	// Process each expired account
	for _, account := range expiredAccounts {
		err := s.expirePremiumAccount(account.AccountID, account.Email)
		if err != nil {
			log.Printf("❌ Error expiring premium for account %d (%s): %v",
				account.AccountID, account.Email, err)
			continue
		}

		log.Printf("✅ Expired premium for account %d (%s)",
			account.AccountID, account.Email)
	}

	log.Printf("✅ Premium expiration check completed. Processed %d accounts", len(expiredAccounts))
	return nil
}

// findExpiredAccounts finds all accounts with expired premium subscriptions
func (s *PremiumExpirationService) findExpiredAccounts() ([]ExpiredAccount, error) {
	query := `
		SELECT account_id, email, premium_expired_date 
		FROM account 
		WHERE is_premium = 1 
		AND premium_expired_date IS NOT NULL 
		AND premium_expired_date <= NOW()
		ORDER BY premium_expired_date ASC
	`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var expiredAccounts []ExpiredAccount

	for rows.Next() {
		var account ExpiredAccount
		err := rows.Scan(
			&account.AccountID,
			&account.Email,
			&account.PremiumExpiredDate,
		)
		if err != nil {
			log.Printf("Error scanning expired account: %v", err)
			continue
		}
		expiredAccounts = append(expiredAccounts, account)
	}

	return expiredAccounts, nil
}

// expirePremiumAccount deactivates premium for a specific account
func (s *PremiumExpirationService) expirePremiumAccount(accountID int64, email string) error {
	// Start transaction
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 1. Update account to set is_premium = 0
	updateAccountQuery := `
		UPDATE account 
		SET is_premium = 0, premium_expired_date = NULL 
		WHERE account_id = ?
	`
	_, err = tx.Exec(updateAccountQuery, accountID)
	if err != nil {
		return err
	}

	// 2. Log the expiration in payment history
	insertHistoryQuery := `
		INSERT INTO account_premium_payment_history 
		(account_id, type_plan, description, created) 
		VALUES (?, ?, ?, NOW())
	`
	description := "Premium subscription expired automatically by system cronjob"
	_, err = tx.Exec(insertHistoryQuery, accountID, "free", description)
	if err != nil {
		return err
	}

	// Commit transaction
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

// GetPremiumExpirationStats returns statistics about premium expirations
func (s *PremiumExpirationService) GetPremiumExpirationStats() (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// Count active premium users
	var activePremium int
	err := s.db.QueryRow("SELECT COUNT(*) FROM account WHERE is_premium = 1").Scan(&activePremium)
	if err != nil {
		return nil, err
	}
	stats["active_premium_users"] = activePremium

	// Count users expiring in next 7 days
	var expiringSoon int
	expiringQuery := `
		SELECT COUNT(*) FROM account 
		WHERE is_premium = 1 
		AND premium_expired_date IS NOT NULL 
		AND premium_expired_date BETWEEN NOW() AND DATE_ADD(NOW(), INTERVAL 7 DAY)
	`
	err = s.db.QueryRow(expiringQuery).Scan(&expiringSoon)
	if err != nil {
		return nil, err
	}
	stats["expiring_in_7_days"] = expiringSoon

	// Count users already expired but still marked as premium (should be 0 after cronjob)
	var alreadyExpired int
	expiredQuery := `
		SELECT COUNT(*) FROM account 
		WHERE is_premium = 1 
		AND premium_expired_date IS NOT NULL 
		AND premium_expired_date <= NOW()
	`
	err = s.db.QueryRow(expiredQuery).Scan(&alreadyExpired)
	if err != nil {
		return nil, err
	}
	stats["already_expired"] = alreadyExpired

	return stats, nil
}

// SendExpirationWarnings sends warnings to users whose premium will expire soon
func (s *PremiumExpirationService) SendExpirationWarnings() error {
	log.Println("📧 Checking for users who need expiration warnings...")

	// Find users expiring in 3 days
	warningQuery := `
		SELECT account_id, email, premium_expired_date 
		FROM account 
		WHERE is_premium = 1 
		AND premium_expired_date IS NOT NULL 
		AND premium_expired_date BETWEEN DATE_ADD(NOW(), INTERVAL 2 DAY) AND DATE_ADD(NOW(), INTERVAL 4 DAY)
	`

	rows, err := s.db.Query(warningQuery)
	if err != nil {
		return err
	}
	defer rows.Close()

	var warningCount int
	for rows.Next() {
		var accountID int64
		var email string
		var expirationDate time.Time

		err := rows.Scan(&accountID, &email, &expirationDate)
		if err != nil {
			log.Printf("Error scanning warning account: %v", err)
			continue
		}

		// Here you would integrate with your notification system
		// For now, just log
		log.Printf("⚠️ Premium expiring soon for account %d (%s) on %s",
			accountID, email, expirationDate.Format("2006-01-02"))

		warningCount++
	}

	log.Printf("📧 Sent %d expiration warnings", warningCount)
	return nil
}
