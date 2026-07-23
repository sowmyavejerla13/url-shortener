package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	appErrors "github.com/sowmyavejerla13/url-shortener/internal/apperrors"
	"github.com/sowmyavejerla13/url-shortener/internal/dto"
	"github.com/sowmyavejerla13/url-shortener/internal/service"
	"github.com/sowmyavejerla13/url-shortener/internal/utils"
)

type URLHandler struct {
	urlService service.URLServiceInterface
}

func NewURLHandler(service service.URLServiceInterface) *URLHandler {
	return &URLHandler{
		urlService: service,
	}
}

// Create godoc
//
// @Summary Create a Short URL
// @Description Creates a short URL
// @Tags Authentication
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.CreateURLRequest true "Create Shorten URL Request"
// @Success 201 {object} dto.CreateURLResponse
// @Failure 400 {object} map[string]string
// @Router /shorten [post]
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

// GetUserURLs godoc
//
// @Summary Get all URLs
// @Description Returns all shortened URLs created by the authenticated user.
// @Tags URLs
// @Security BearerAuth
// @Produce json
// @Success 200 {array} dto.CreateURLResponse
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /urls [get]
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

// DeleteURL godoc
//
// @Summary Delete a URL
// @Description Deletes a shortened URL owned by the authenticated user.
// @Tags URLs
// @Security BearerAuth
// @Produce json
// @Param id path string true "URL ID"
// @Success 204 "No Content"
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /urls/{id} [delete]
func (h *URLHandler) Delete(c *gin.Context) {
	urlID := c.Param("id")
	userID := c.GetString("userID")
	err := h.urlService.DeleteURL(userID, urlID)
	if err != nil {
		switch err {
		case appErrors.ErrURLNotFound:
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})

		case appErrors.ErrForbidden:
			c.JSON(http.StatusForbidden, gin.H{
				"error": err.Error(),
			})

		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		}
		return
	}

	c.Status(http.StatusNoContent)
}
