package model

import "time"

type URL struct {
	ID            string    `db:"id"`
	ShortCode     string    `db:"short_code"`
	OriginalURL   string    `db:"original_url"`
	ClickCount 	  int    	`db:"click_count"`
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
	UserID 		  string    `db:"user_id"`
}