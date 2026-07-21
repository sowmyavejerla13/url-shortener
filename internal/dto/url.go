package dto

type CreateURLRequest struct {
	URL string `json:"url" binding:"required,url"`
}

type CreateURLResponse struct {
	ID          string `json:"id"`
	ShortCode   string `json:"short_code"`
	OriginalURL string `json:"original_url"`
	ClickCount  int    `json:"click_count"`
	CreatedAt   string `json:"created_at"`
}
