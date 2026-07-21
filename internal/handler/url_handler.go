package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sowmyavejerla13/url-shortener/internal/dto"
	"github.com/sowmyavejerla13/url-shortener/internal/service"
	"github.com/sowmyavejerla13/url-shortener/internal/utils"
)

type URLHandler struct {
	urlService *service.URLService
	baseURL    string
}

func NewURLHandler(service *service.URLService, baseURL string) *URLHandler {
	return &URLHandler{
		urlService: service,
		baseURL:    baseURL,
	}
}

func (h *URLHandler) Create(c *gin.Context) {
	var req dto.CreateURLRequest
	userID := c.GetString("userID")

	if err := c.ShouldBindJSON(&req); err != nil {
		if _, ok := err.(validator.ValidationErrors); ok {
			c.JSON(http.StatusBadRequest, gin.H{
				"errors": utils.FormatValidationErrors(err),
			})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	url, err := h.urlService.CreateShortURL(userID, req.URL)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, dto.CreateURLResponse{
		ID:          url.ID,
		ShortCode:   url.ShortCode,
		OriginalURL: url.OriginalURL,
		CreatedAt:   url.CreatedAt.Format(time.RFC3339),
		ClickCount:  url.ClickCount,
	})

}

func (h *URLHandler) Redirect(c *gin.Context) {
	shortCode := c.Param("shortCode")
	originalURL, err := h.urlService.Redirect(shortCode)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.Redirect(http.StatusFound, originalURL)
}

func (h *URLHandler) GetUserURLs(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}
	urls, err := h.urlService.GetUserURLs(userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	responses := make([]dto.CreateURLResponse, 0, len(urls))

	for _, url := range urls {
		responses = append(responses, dto.CreateURLResponse{
			ID:          url.ID,
			ShortCode:   url.ShortCode,
			OriginalURL: url.OriginalURL,
			ClickCount:  url.ClickCount,
			CreatedAt:   url.CreatedAt.Format(time.RFC3339),
		})
	}
	c.JSON(http.StatusOK, responses)
}

func (h *URLHandler) Delete(c *gin.Context) {
	urlID := c.Param("id")
	userID := c.GetString("userID")
	err := h.urlService.DeleteURL(userID, urlID)
	if err != nil {
		switch err.Error() {
		case "url not found":
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case "forbidden":
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.Status(http.StatusNoContent)
}
