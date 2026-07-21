package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sowmyavejerla13/url-shortener/internal/model"
)

type URLRepository struct {
	db *pgxpool.Pool
}

func NewURLRepository(db *pgxpool.Pool) *URLRepository {
	return &URLRepository{
		db: db,
	}
}

func (r *URLRepository) Create(url *model.URL) error {
	query := `
		INSERT INTO urls(
			short_code,
			original_url,
			user_id
		)
		VALUES($1,$2,$3)
		RETURNING
			id,
			click_count,
			created_at,
			updated_at
	`

	return r.db.QueryRow(
		context.Background(),
		query,
		url.ShortCode,
		url.OriginalURL,
		url.UserID,
	).Scan(
		&url.ID,
		&url.ClickCount,
		&url.CreatedAt,
		&url.UpdatedAt,
	)
}

func (r *URLRepository) GetByShortCode(code string) (*model.URL, error) {
	query := `SELECT 
   				id,
				short_code,
				original_url,
				click_count,
				user_id,
				created_at,
				updated_at
			FROM urls
			WHERE short_code = $1`
	url := &model.URL{}
	err := r.db.QueryRow(context.Background(), query, code).Scan(
		&url.ID,
		&url.ShortCode,
		&url.OriginalURL,
		&url.ClickCount,
		&url.UserID,
		&url.CreatedAt,
		&url.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return url, nil
}

func (r *URLRepository) GetByUserID(userID string) ([]model.URL, error) {
	query := `SELECT 
   				id,
				short_code,
				original_url,
				click_count,
				user_id,
				created_at,
				updated_at
			FROM urls
			WHERE user_id = $1
			ORDER BY created_at DESC`
	rows, err := r.db.Query(
		context.Background(),
		query,
		userID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var urls []model.URL

	for rows.Next() {
		var url model.URL
		err := rows.Scan(
			&url.ID,
			&url.ShortCode,
			&url.OriginalURL,
			&url.ClickCount,
			&url.UserID,
			&url.CreatedAt,
			&url.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		urls = append(urls, url)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return urls, nil
}

func (r *URLRepository) GetByOriginalURL(userID, originalURL string) (*model.URL, error) {
	query := `	SELECT 
   				 	id,
					short_code,
					original_url,
					user_id,
					click_count,
					created_at,
					updated_at
				FROM urls
				WHERE user_id = $1
  				AND original_url = $2;
			`
	url := &model.URL{}
	err := r.db.QueryRow(context.Background(), query, userID, originalURL).Scan(
		&url.ID,
		&url.ShortCode,
		&url.OriginalURL,
		&url.UserID,
		&url.ClickCount,
		&url.CreatedAt,
		&url.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return url, nil
}

func (r *URLRepository) IncrementClickCount(id string) error {
	query := `
		UPDATE urls
		SET click_count = click_count + 1
			WHERE id = $1
	`
	cmd, err := r.db.Exec(
		context.Background(),
		query,
		id,
	)
	if cmd.RowsAffected() == 0 {
		return errors.New("url not found")
	}
	return err
}
func (r *URLRepository) GetByID(id string) (*model.URL, error) {
	query := `	SELECT 
   				 	id,
					short_code,
					original_url,
					user_id,
					click_count,
					created_at,
					updated_at
				FROM urls
				WHERE id = $1;
			`
	url := &model.URL{}
	err := r.db.QueryRow(context.Background(), query, id).Scan(
		&url.ID,
		&url.ShortCode,
		&url.OriginalURL,
		&url.UserID,
		&url.ClickCount,
		&url.CreatedAt,
		&url.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return url, nil
}
func (r *URLRepository) Delete(id string) error {
	query := `DELETE FROM 
				urls where id =$1`
	cmd, err := r.db.Exec(
		context.Background(),
		query,
		id)

	if err != nil {
		return err
	}

	if cmd.RowsAffected() == 0 {
		return errors.New("url not found")
	}

	return nil
}
