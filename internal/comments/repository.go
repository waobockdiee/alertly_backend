package comments

import (
	"alertly/internal/common"
	"database/sql"
	"fmt"
	"log"
)

type Repository interface {
	Save(comment InComment) (int64, error)
	GetClusterCommentsByID(inclID int64) ([]Comment, error)
	GetCommentById(incoID int64) (Comment, error)
}

type mysqlRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &mysqlRepository{db: db}
}

func (r *mysqlRepository) Save(comment InComment) (int64, error) {

	// 1. Insertar comentario
	query := `INSERT INTO incident_comments (account_id, comment, created_at, incl_id) VALUES (?, ?, NOW(), ?)`
	result, err := r.db.Exec(query, comment.AccountID, comment.Comment, comment.InclID)
	if err != nil {
		return 0, fmt.Errorf("failed to insert comment: %w", err)
	}

	query = `UPDATE incident_clusters SET counter_total_comments = counter_total_comments + 1 WHERE incl_id = ?`
	_, err = r.db.Exec(query, comment.InclID)

	if err != nil {
		log.Printf("Error updating total comments count: %v", err)
	}

	commentID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert id: %w", err)
	}

	// ----------------------CITIZEN SCORE----------------------- //
	// Save to DB
	err = common.SaveScore(r.db, comment.AccountID, 5)
	if err != nil {
		fmt.Println("error saving score") // It's not necesary to stop the server
	}

	return commentID, nil
}

func (r *mysqlRepository) GetClusterCommentsByID(inclID int64) ([]Comment, error) {
	query := `SELECT 
	t1.inco_id,
	t1.account_id, 
	t1.comment,
	t1.created_at,
	t1.comment_status,
	t1.counter_flags,
	t2.nickname,
	t2.thumbnail_url
	FROM incident_comments t1 INNER JOIN account t2 ON t1.account_id = t2.account_id
	WHERE t1.incl_id = ?
	ORDER BY t1.inco_id DESC`
	rows, err := r.db.Query(query, inclID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var comments []Comment

	for rows.Next() {
		var c Comment

		if err := rows.Scan(
			&c.IncoID,
			&c.AccountID,
			&c.Comment,
			&c.CreatedAt,
			&c.CommentStatus,
			&c.CounterFlags,
			&c.Nickname,
			&c.ThumbnailUrl,
		); err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}

	if err != nil {
		return nil, err
	}

	return comments, nil
}

func (r *mysqlRepository) GetCommentById(incoID int64) (Comment, error) {
	query := `SELECT 
	t1.inco_id,
	t1.account_id, 
	t1.comment,
	t1.created_at,
	t1.comment_status,
	t1.counter_flags,
	t2.nickname,
	t2.thumbnail_url
	FROM incident_comments t1 INNER JOIN account t2 ON t1.account_id = t2.account_id
	WHERE t1.inco_id = ?`

	var c Comment
	err := r.db.QueryRow(query, incoID).Scan(
		&c.IncoID,
		&c.AccountID,
		&c.Comment,
		&c.CreatedAt,
		&c.CommentStatus,
		&c.CounterFlags,
		&c.Nickname,
		&c.ThumbnailUrl,
	)

	if err != nil {
		return Comment{}, err
	}

	return c, nil
}
