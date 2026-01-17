package comments

import (
	"alertly/internal/common"
	"alertly/internal/dbtypes"
	"database/sql"
	"fmt"
	"log"
)

type Repository interface {
	Save(comment InComment) (int64, error)
	GetClusterCommentsByID(inclID int64) ([]Comment, error)
	GetCommentById(incoID int64) (Comment, error)
}

type pgRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &pgRepository{db: db}
}

func (r *pgRepository) Save(comment InComment) (int64, error) {

	// 1. Insertar comentario
	query := `INSERT INTO incident_comments (account_id, comment, created_at, incl_id) VALUES ($1, $2, NOW(), $3) RETURNING inco_id`
	var commentID int64
	err := r.db.QueryRow(query, comment.AccountID, comment.Comment, comment.InclID).Scan(&commentID)
	if err != nil {
		return 0, fmt.Errorf("failed to insert comment: %w", err)
	}

	query = `UPDATE incident_clusters SET counter_total_comments = counter_total_comments + 1 WHERE incl_id = $1`
	_, err = r.db.Exec(query, comment.InclID)

	if err != nil {
		log.Printf("Error updating total comments count: %v", err)
	}

	// ----------------------CITIZEN SCORE----------------------- //
	// Save to DB
	err = common.SaveScore(r.db, comment.AccountID, 5)
	if err != nil {
		fmt.Println("error saving score") // It's not necesary to stop the server
	}

	return commentID, nil
}

func (r *pgRepository) GetClusterCommentsByID(inclID int64) ([]Comment, error) {
	query := `SELECT
	t1.inco_id,
	t1.account_id,
	t1.comment,
	t1.created_at,
	t1.comment_status,
	t1.counter_flags,
	t2.nickname,
	COALESCE(t2.thumbnail_url, '') as thumbnail_url
	FROM incident_comments t1 INNER JOIN account t2 ON t1.account_id = t2.account_id
	WHERE t1.incl_id = $1
	ORDER BY t1.inco_id DESC`
	rows, err := r.db.Query(query, inclID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var comments []Comment

	for rows.Next() {
		var c Comment
		var commentStatus dbtypes.NullBool

		if err := rows.Scan(
			&c.IncoID,
			&c.AccountID,
			&c.Comment,
			&c.CreatedAt,
			&commentStatus,
			&c.CounterFlags,
			&c.Nickname,
			&c.ThumbnailUrl,
		); err != nil {
			return nil, err
		}
		c.CommentStatus = commentStatus.Valid && commentStatus.Bool
		comments = append(comments, c)
	}

	if err != nil {
		return nil, err
	}

	return comments, nil
}

func (r *pgRepository) GetCommentById(incoID int64) (Comment, error) {
	query := `SELECT
	t1.inco_id,
	t1.account_id,
	t1.comment,
	t1.created_at,
	t1.comment_status,
	t1.counter_flags,
	t2.nickname,
	COALESCE(t2.thumbnail_url, '') as thumbnail_url
	FROM incident_comments t1 INNER JOIN account t2 ON t1.account_id = t2.account_id
	WHERE t1.inco_id = $1`

	var c Comment
	var commentStatus dbtypes.NullBool
	err := r.db.QueryRow(query, incoID).Scan(
		&c.IncoID,
		&c.AccountID,
		&c.Comment,
		&c.CreatedAt,
		&commentStatus,
		&c.CounterFlags,
		&c.Nickname,
		&c.ThumbnailUrl,
	)

	if err != nil {
		return Comment{}, err
	}
	c.CommentStatus = commentStatus.Valid && commentStatus.Bool

	return c, nil
}
