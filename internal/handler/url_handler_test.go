package handler_test

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	appErrors "github.com/sowmyavejerla13/url-shortener/internal/apperrors"
	"github.com/sowmyavejerla13/url-shortener/internal/handler"
	"github.com/sowmyavejerla13/url-shortener/internal/handler/mocks"
	"github.com/sowmyavejerla13/url-shortener/internal/model"
)

func TestCreateShortURL(t *testing.T) {

	gin.SetMode(gin.TestMode)

	tests := []struct {
		name string

		requestBody string

		serviceURL *model.URL
		serviceErr error

		expectedStatus int
		expectedBody   string
	}{
		{
			name: "Success",

			requestBody: `{
				"url":"https://google.com"
			}`,

			serviceURL: &model.URL{
				ID:          "1",
				ShortCode:   "abc123",
				OriginalURL: "https://google.com",
				CreatedAt:   time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC),
			},

			expectedStatus: http.StatusCreated,

			expectedBody: `{
					"id":"1",
					"short_code":"abc123",
					"original_url":"https://google.com",
					"created_at":"2025-01-01T10:00:00Z",
					"click_count":0
			}`,
		},

		{
			name: "Invalid JSON",

			requestBody: `{
				"url":
			}`,

			expectedStatus: http.StatusBadRequest,

			expectedBody: `{
				"error":"Invalid request body"
			}`,
		},

		{
			name: "Validation Error",

			requestBody: `{
				"url":""
			}`,

			expectedStatus: http.StatusBadRequest,
		},

		{
			name: "Service Error",

			requestBody: `{
				"url":"https://google.com"
			}`,

			serviceErr: errors.New("url already exists"),

			expectedStatus: http.StatusBadRequest,

			expectedBody: `{
				"error":"url already exists"
			}`,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			mockService := &mocks.URLServiceMock{}

			mockService.CreateShortURLFunc = func(userID, originalURL string) (*model.URL, error) {
				return tt.serviceURL, tt.serviceErr
			}

			urlHandler := handler.NewURLHandler(mockService)

			req, _ := http.NewRequest(
				http.MethodPost,
				"/shorten",
				bytes.NewBuffer([]byte(tt.requestBody)),
			)

			req.Header.Set("Content-Type", "application/json")

			rec := httptest.NewRecorder()

			c, _ := gin.CreateTestContext(rec)

			c.Request = req

			// fake authenticated user
			c.Set("userID", "user1")

			urlHandler.Create(c)

			assert.Equal(t, tt.expectedStatus, rec.Code)

			if tt.name == "Validation Error" {

				assert.Contains(t, rec.Body.String(), "errors")

			} else {

				assert.JSONEq(
					t,
					tt.expectedBody,
					rec.Body.String(),
				)

			}
		})
	}
}

func TestRedirect(t *testing.T) {

	gin.SetMode(gin.TestMode)

	tests := []struct {
		name string

		shortCode string

		originalURL string
		serviceErr  error

		expectedStatus   int
		expectedLocation string
		expectedBody     string
	}{
		{
			name: "Success",

			shortCode: "abc123",

			originalURL: "https://google.com",

			expectedStatus:   http.StatusFound,
			expectedLocation: "https://google.com",
		},

		{
			name: "URL Not Found",

			shortCode: "abc123",

			serviceErr: appErrors.ErrURLNotFound,

			expectedStatus: http.StatusNotFound,

			expectedBody: `{
				"error":"url not found"
			}`,
		},

		{
			name: "Repository Error",

			shortCode: "abc123",

			serviceErr: assert.AnError,

			expectedStatus: http.StatusNotFound,

			expectedBody: `{
				"error":"assert.AnError general error for testing"
			}`,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			mockService := &mocks.URLServiceMock{}

			mockService.RedirectFunc = func(shortCode string) (string, error) {
				return tt.originalURL, tt.serviceErr
			}

			urlHandler := handler.NewURLHandler(mockService)

			req, _ := http.NewRequest(
				http.MethodGet,
				"/"+tt.shortCode,
				nil,
			)

			rec := httptest.NewRecorder()

			c, _ := gin.CreateTestContext(rec)

			c.Params = gin.Params{
				{
					Key:   "shortCode",
					Value: tt.shortCode,
				},
			}

			c.Request = req

			urlHandler.Redirect(c)

			assert.Equal(t, tt.expectedStatus, rec.Code)

			if tt.name == "Success" {

				assert.Equal(
					t,
					tt.expectedLocation,
					rec.Header().Get("Location"),
				)

			} else {

				assert.JSONEq(
					t,
					tt.expectedBody,
					rec.Body.String(),
				)
			}
		})
	}
}

func TestGetUserURLs(t *testing.T) {

	gin.SetMode(gin.TestMode)

	tests := []struct {
		name string

		serviceURLs []model.URL
		serviceErr  error

		expectedStatus int
		expectedBody   string
	}{
		{
			name: "Success",

			serviceURLs: []model.URL{
				{
					ID:          "1",
					ShortCode:   "abc123",
					OriginalURL: "https://google.com",
				},
				{
					ID:          "2",
					ShortCode:   "xyz789",
					OriginalURL: "https://github.com",
				},
			},

			expectedStatus: http.StatusOK,

			expectedBody: `[
				{
					"id":"1",
					"short_code":"abc123",
					"original_url":"https://google.com",
					"click_count":0,
					"created_at":"0001-01-01T00:00:00Z"
				},
				{
					"id":"2",
					"short_code":"xyz789",
					"original_url":"https://github.com",
					"click_count":0,
					"created_at":"0001-01-01T00:00:00Z"
				}
			]`,
		},

		{
			name: "Service Error",

			serviceErr: assert.AnError,

			expectedStatus: http.StatusInternalServerError,

			expectedBody: `{
				"error":"assert.AnError general error for testing"
			}`,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			mockService := &mocks.URLServiceMock{}

			mockService.GetUserURLsFunc = func(userID string) ([]model.URL, error) {
				return tt.serviceURLs, tt.serviceErr
			}

			urlHandler := handler.NewURLHandler(mockService)

			req, _ := http.NewRequest(http.MethodGet, "/urls", nil)

			rec := httptest.NewRecorder()

			c, _ := gin.CreateTestContext(rec)

			c.Request = req
			c.Set("userID", "user1")

			urlHandler.GetUserURLs(c)

			assert.Equal(t, tt.expectedStatus, rec.Code)

			assert.JSONEq(t, tt.expectedBody, rec.Body.String())
		})
	}
}

func TestDeleteURL(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		urlID          string
		serviceErr     error
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Success",
			urlID:          "1",
			expectedStatus: http.StatusNoContent,
		},
		{
			name:           "URL Not Found",
			urlID:          "1",
			serviceErr:     appErrors.ErrURLNotFound,
			expectedStatus: http.StatusNotFound,
			expectedBody: `{
				"error":"url not found"
			}`,
		},
		{
			name:           "Forbidden",
			urlID:          "1",
			serviceErr:     appErrors.ErrForbidden,
			expectedStatus: http.StatusForbidden,
			expectedBody: `{
				"error":"forbidden"
			}`,
		},
		{
			name:           "Repository Error",
			urlID:          "1",
			serviceErr:     assert.AnError,
			expectedStatus: http.StatusInternalServerError,
			expectedBody: `{
				"error":"assert.AnError general error for testing"
			}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockService := &mocks.URLServiceMock{}

			mockService.DeleteURLFunc = func(userID, urlID string) error {
				return tt.serviceErr
			}

			urlHandler := handler.NewURLHandler(mockService)

			router := gin.New()

			router.Use(func(c *gin.Context) {
				c.Set("userID", "user1")
				c.Next()
			})

			router.DELETE("/urls/:id", urlHandler.Delete)

			req, _ := http.NewRequest(
				http.MethodDelete,
				"/urls/"+tt.urlID,
				nil,
			)

			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)

			if tt.expectedStatus == http.StatusNoContent {
				assert.Empty(t, rec.Body.String())
			} else {
				assert.JSONEq(t, tt.expectedBody, rec.Body.String())
			}
		})
	}
}
