package signup

import (
	"database/sql"
	"fmt"
)

type Repository interface {
	InsertUser(user User) (int64, error)
	GetUserByID(id int64) (User, error)
}

type mysqlRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &mysqlRepository{db: db}
}

func (repo *mysqlRepository) InsertUser(user User) (int64, error) {
	query := `
		INSERT INTO account (email, first_name, last_name,  password, birth_year, birth_month, birth_day, activation_code)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`
	result, err := repo.db.Exec(query,
		user.Email,
		user.FirstName,
		user.LastName,
		user.Password,
		user.BirthYear,
		user.BirthMonth,
		user.BirthDay,
		user.ActivationCode,
	)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	fmt.Printf("Usuario insertado con ID: %d\n", id)
	return id, nil
}

func (repo *mysqlRepository) GetUserByID(id int64) (User, error) {
	query := `
		SELECT account_id, email, password, activation_code, firstName
		FROM account WHERE account_id = ?
	`
	row := repo.db.QueryRow(query, id)
	var user User
	err := row.Scan(&user.ID, &user.Email, &user.Password, &user.ActivationCode, &user.FirstName)
	return user, err
}
