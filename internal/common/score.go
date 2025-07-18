package common

import (
	"fmt"
)

func SaveScore(dbExec DBExecutor, accountID int64, score uint8) error {
	updateQuery := `UPDATE account SET score = score + ? WHERE account_id = ?`
	_, err := dbExec.Exec(updateQuery, score, accountID)
	if err != nil {
		return fmt.Errorf("failed to update score: %w", err)
	}
	return nil
}
