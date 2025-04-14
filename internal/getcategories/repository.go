package getcategories

import (
	"database/sql"
)

type Repository interface {
	GetCategories() ([]Category, error)
}

type mysqlRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &mysqlRepository{db: db}
}

func (r *mysqlRepository) GetCategories() ([]Category, error) {
	query := `SELECT inca_id, name, description, icon, code FROM incident_categories ORDER BY name DESC`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []Category
	for rows.Next() {
		var c Category
		if err := rows.Scan(&c.IncaId, &c.Name, &c.Description, &c.Icon, &c.Code); err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return categories, nil
}
