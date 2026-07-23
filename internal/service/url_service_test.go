package service_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	appErrors "github.com/sowmyavejerla13/url-shortener/internal/apperrors"
	"github.com/sowmyavejerla13/url-shortener/internal/model"
	"github.com/sowmyavejerla13/url-shortener/internal/service"
	"github.com/sowmyavejerla13/url-shortener/internal/service/mocks"
)

func TestRedirect(t *testing.T) {

	tests := []struct {
		name          string
		shortCode     string
		mockURL       *model.URL
		mockGetErr    error
		mockClickErr  error
		expectedURL   string
		expectedError error
	}{
		{
			name:      "Success",
			shortCode: "abc123",
			mockURL: &model.URL{
				ID:          "1",
				ShortCode:   "abc123",
				OriginalURL: "https://google.com",
			},
			expectedURL: "https://google.com",
		},
		{
			name:          "URL Not Found",
			shortCode:     "invalid",
			mockURL:       nil,
			expectedError: appErrors.ErrURLNotFound,
		},
		{
			name:      "Increment Click Count Error",
			shortCode: "abc123",
			mockURL: &model.URL{
				ID:          "1",
				ShortCode:   "abc123",
				OriginalURL: "https://google.com",
			},
			mockClickErr:  assert.AnError,
			expectedError: assert.AnError,
		},
		{
			name:          "Repository Error",
			shortCode:     "abc123",
			mockGetErr:    assert.AnError,
			expectedError: assert.AnError,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			mockRepo := &mocks.URLRepositoryMock{}

			mockRepo.GetByShortCodeFunc = func(shortCode string) (*model.URL, error) {
				return tt.mockURL, tt.mockGetErr
			}

			mockRepo.IncrementClickCountFunc = func(id string) error {
				return tt.mockClickErr
			}

			urlService := service.NewURLService(
				mockRepo,
				func() (string, error) {
					return "abc123", nil
				},
			)

			url, err := urlService.Redirect(tt.shortCode)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
				assert.Equal(t, "", url)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedURL, url)
			}
		})
	}
}

func TestCreateShortURL(t *testing.T) {

	tests := []struct {
		name string

		userID      string
		originalURL string

		mockExistingURL *model.URL
		mockOriginalErr error

		generatedCodes []string
		generateErr    error

		shortCodeExists map[string]bool
		getShortCodeErr error

		createErr error

		expectedURL *model.URL
		expectedErr error
	}{
		{
			name:        "Success",
			userID:      "user1",
			originalURL: "https://google.com",

			generatedCodes: []string{"abc123"},

			shortCodeExists: map[string]bool{},

			expectedURL: &model.URL{
				UserID:      "user1",
				OriginalURL: "https://google.com",
				ShortCode:   "abc123",
			},
		},

		{
			name:        "Invalid URL",
			userID:      "user1",
			originalURL: "invalid-url",

			expectedErr: nil,
		},

		{
			name:        "URL Already Exists",
			userID:      "user1",
			originalURL: "https://google.com",

			mockExistingURL: &model.URL{
				ID:          "1",
				UserID:      "user1",
				OriginalURL: "https://google.com",
				ShortCode:   "abc123",
			},
		},

		{
			name:        "GetByOriginalURL Error",
			userID:      "user1",
			originalURL: "https://google.com",

			mockOriginalErr: assert.AnError,

			expectedErr: assert.AnError,
		},

		{
			name:        "Generate Short Code Error",
			userID:      "user1",
			originalURL: "https://google.com",

			generateErr: assert.AnError,

			expectedErr: assert.AnError,
		},

		{
			name:        "GetByShortCode Error",
			userID:      "user1",
			originalURL: "https://google.com",

			generatedCodes: []string{"abc123"},

			getShortCodeErr: assert.AnError,

			expectedErr: assert.AnError,
		},

		{
			name:        "Create Error",
			userID:      "user1",
			originalURL: "https://google.com",

			generatedCodes: []string{"abc123"},

			shortCodeExists: map[string]bool{},

			createErr: assert.AnError,

			expectedErr: assert.AnError,
		},

		{
			name:        "Short Code Collision",
			userID:      "user1",
			originalURL: "https://google.com",

			generatedCodes: []string{
				"abc123",
				"xyz789",
			},

			shortCodeExists: map[string]bool{
				"abc123": true,
				"xyz789": false,
			},

			expectedURL: &model.URL{
				UserID:      "user1",
				OriginalURL: "https://google.com",
				ShortCode:   "xyz789",
			},
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			mockRepo := &mocks.URLRepositoryMock{}

			mockRepo.GetByOriginalURLFunc = func(userID, originalURL string) (*model.URL, error) {
				return tt.mockExistingURL, tt.mockOriginalErr
			}

			mockRepo.GetByShortCodeFunc = func(code string) (*model.URL, error) {

				if tt.getShortCodeErr != nil {
					return nil, tt.getShortCodeErr
				}

				if tt.shortCodeExists[code] {
					return &model.URL{
						ShortCode: code,
					}, nil
				}

				return nil, nil
			}

			mockRepo.CreateFunc = func(url *model.URL) error {
				return tt.createErr
			}

			index := 0

			generator := func() (string, error) {

				if tt.generateErr != nil {
					return "", tt.generateErr
				}

				code := tt.generatedCodes[index]
				index++

				return code, nil
			}

			urlService := service.NewURLService(
				mockRepo,
				generator,
			)
			url, err := urlService.CreateShortURL(
				tt.userID,
				tt.originalURL,
			)

			switch tt.name {

			case "Invalid URL":
				assert.Error(t, err)
				assert.Nil(t, url)
				assert.EqualError(t, err, "invalid url")

			case "URL Already Exists":
				assert.NoError(t, err)
				assert.Equal(t, tt.mockExistingURL, url)

			case "Success":
				assert.NoError(t, err)
				assert.NotNil(t, url)

				assert.Equal(t, tt.expectedURL.UserID, url.UserID)
				assert.Equal(t, tt.expectedURL.OriginalURL, url.OriginalURL)
				assert.Equal(t, tt.expectedURL.ShortCode, url.ShortCode)

			case "Short Code Collision":
				assert.NoError(t, err)
				assert.NotNil(t, url)

				assert.Equal(t, tt.expectedURL.ShortCode, url.ShortCode)
				assert.Equal(t, tt.expectedURL.OriginalURL, url.OriginalURL)
				assert.Equal(t, tt.expectedURL.UserID, url.UserID)

			default:
				assert.Error(t, err)
				assert.Equal(t, tt.expectedErr, err)
				assert.Nil(t, url)
			}
		})
	}
}

func TestGetUserURLs(t *testing.T) {

	tests := []struct {
		name string

		userID string

		mockURLs []model.URL
		mockErr  error

		expectedURLs []model.URL
		expectedErr  error
	}{
		{
			name:   "Success",
			userID: "user1",

			mockURLs: []model.URL{
				{
					ID:          "1",
					UserID:      "user1",
					ShortCode:   "abc123",
					OriginalURL: "https://google.com",
				},
				{
					ID:          "2",
					UserID:      "user1",
					ShortCode:   "xyz789",
					OriginalURL: "https://github.com",
				},
			},

			expectedURLs: []model.URL{
				{
					ID:          "1",
					UserID:      "user1",
					ShortCode:   "abc123",
					OriginalURL: "https://google.com",
				},
				{
					ID:          "2",
					UserID:      "user1",
					ShortCode:   "xyz789",
					OriginalURL: "https://github.com",
				},
			},
		},

		{
			name:        "Repository Error",
			userID:      "user1",
			mockErr:     assert.AnError,
			expectedErr: assert.AnError,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			mockRepo := &mocks.URLRepositoryMock{}

			mockRepo.GetByUserIDFunc = func(userID string) ([]model.URL, error) {
				return tt.mockURLs, tt.mockErr
			}

			urlService := service.NewURLService(
				mockRepo,
				nil,
			)

			urls, err := urlService.GetUserURLs(tt.userID)

			if tt.expectedErr != nil {

				assert.Error(t, err)
				assert.Equal(t, tt.expectedErr, err)
				assert.Nil(t, urls)

			} else {

				assert.NoError(t, err)
				assert.Equal(t, tt.expectedURLs, urls)

			}
		})
	}
}

func TestDeleteURL(t *testing.T) {

	tests := []struct {
		name string

		userID string
		urlID  string

		mockURL    *model.URL
		mockGetErr error

		mockDeleteErr error

		expectedErr error
	}{
		{
			name:   "Success",
			userID: "user1",
			urlID:  "1",

			mockURL: &model.URL{
				ID:     "1",
				UserID: "user1",
			},
		},

		{
			name:   "Repository Error",
			userID: "user1",
			urlID:  "1",

			mockGetErr: assert.AnError,

			expectedErr: assert.AnError,
		},

		{
			name:   "URL Not Found",
			userID: "user1",
			urlID:  "1",

			mockURL: nil,

			expectedErr: appErrors.ErrURLNotFound,
		},

		{
			name:   "Forbidden",
			userID: "user1",
			urlID:  "1",

			mockURL: &model.URL{
				ID:     "1",
				UserID: "another-user",
			},

			expectedErr: appErrors.ErrForbidden,
		},

		{
			name:   "Delete Error",
			userID: "user1",
			urlID:  "1",

			mockURL: &model.URL{
				ID:     "1",
				UserID: "user1",
			},

			mockDeleteErr: assert.AnError,

			expectedErr: assert.AnError,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			mockRepo := &mocks.URLRepositoryMock{}

			mockRepo.GetByIDFunc = func(id string) (*model.URL, error) {
				return tt.mockURL, tt.mockGetErr
			}

			mockRepo.DeleteFunc = func(id string) error {
				return tt.mockDeleteErr
			}

			urlService := service.NewURLService(
				mockRepo,
				nil,
			)

			err := urlService.DeleteURL(
				tt.userID,
				tt.urlID,
			)

			if tt.expectedErr != nil {

				assert.Error(t, err)
				assert.Equal(t, tt.expectedErr, err)

			} else {

				assert.NoError(t, err)

			}
		})
	}
}
