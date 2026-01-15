package getsubcategoriesbycategoryid

import (
	"database/sql"
)

type Repository interface {
	GetSubcategoriesByCategoryId(id int) ([]Subcategory, error)
}

type pgRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &pgRepository{db: db}
}

func (r *pgRepository) GetSubcategoriesByCategoryId(id int) ([]Subcategory, error) {
	query := `SELECT insu_id, inca_id, name, description, icon, code, min_circle_range, max_circle_range, default_circle_range, category_code, subcategory_code FROM incident_subcategories WHERE inca_id = $1 ORDER BY name DESC`
	rows, err := r.db.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []Subcategory
	for rows.Next() {
		var c Subcategory
		if err := rows.Scan(&c.InsuId, &c.IncaId, &c.Name, &c.Description, &c.Icon, &c.Code, &c.MinCircleRange, &c.MaxCircleRange, &c.DefaultCircleRange, &c.CategoryCode, &c.SubcategoryCode); err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return categories, nil
}
