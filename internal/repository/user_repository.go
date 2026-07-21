package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sowmyavejerla13/url-shortener/internal/model"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) GetByEmail(email string) (*model.User, error) {

	query := `SELECT 
   				id,
				name,
				email,
				password_hash,
				created_at,
				updated_at
			FROM users
			WHERE email = $1`
	user := &model.User{}
	err := r.db.QueryRow(context.Background(), query, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return user, nil

}

func (r *UserRepository) Create(user *model.User) error {

	query := `
		INSERT INTO  users(
			name,
			email,
			password_hash
		)
		VALUES ($1, $2, $3)`

	_, err := r.db.Exec(context.Background(), query, user.Name, user.Email, user.PasswordHash)
	return err
}
