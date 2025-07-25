package invitefriend

import (
	"alertly/internal/common"
	"database/sql"
)

type Repository interface {
	Save(invitation Invitation) error
}

type mysqlRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &mysqlRepository{db: db}
}

func (r *mysqlRepository) Save(invitation Invitation) error {
	err := common.SaveScore(r.db, invitation.AccountID, 20)
	return err
}
