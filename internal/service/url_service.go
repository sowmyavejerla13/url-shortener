package service

import (
	"github.com/sowmyavejerla13/url-shortener/internal/model"
	"github.com/sowmyavejerla13/url-shortener/internal/repository"
	"github.com/sowmyavejerla13/url-shortener/internal/utils"

	appErrors "github.com/sowmyavejerla13/url-shortener/internal/apperrors"
)

type URLService struct {
	repo              repository.URLRepositoryInterface
	generateShortCode func() (string, error)
}

func NewURLService(
	repo repository.URLRepositoryInterface,
	generateShortCode func() (string, error),
) *URLService {
	return &URLService{
		repo:              repo,
		generateShortCode: generateShortCode,
	}
}

type URLServiceInterface interface {
	CreateShortURL(userID, originalURL string) (*model.URL, error)
	Redirect(shortCode string) (string, error)
	GetUserURLs(userID string) ([]model.URL, error)
	DeleteURL(userID, urlID string) error
}

func (s *URLService) CreateShortURL(userID, originalURL string) (*model.URL, error) {
	err := utils.ValidateURL(originalURL)
	if err != nil {
		return nil, err
	}
	existingUrl, err := s.repo.GetByOriginalURL(userID, originalURL)
	if err != nil {
		return nil, err
	}

	if existingUrl != nil {
		return existingUrl, nil
	}

	var shortCode string

	for {
		code, err := s.generateShortCode()
		if err != nil {
			return nil, err
		}

		existing, err := s.repo.GetByShortCode(code)
		if err != nil {
			return nil, err
		}

		if existing == nil {
			shortCode = code
			break
		}
	}

	url := &model.URL{
		ShortCode:   shortCode,
		OriginalURL: originalURL,
		UserID:      userID,
	}

	if err = s.repo.Create(url); err != nil {
		return nil, err
	}
	return url, nil

}

func (s *URLService) Redirect(shortCode string) (string, error) {

	url, err := s.repo.GetByShortCode(shortCode)
	if err != nil {
		return "", err
	}
	if url == nil {
		return "", appErrors.ErrURLNotFound
	}
	err = s.repo.IncrementClickCount(url.ID)
	if err != nil {
		return "", err
	}
	return url.OriginalURL, nil
}

func (s *URLService) GetUserURLs(userID string) ([]model.URL, error) {
	urls, err := s.repo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	return urls, nil
}

func (s *URLService) DeleteURL(userID, urlID string) error {
	url, err := s.repo.GetByID(urlID)
	if err != nil {
		return err
	}
	if url == nil {
		return appErrors.ErrURLNotFound
	}
	if url.UserID != userID {
		return appErrors.ErrForbidden
	}

	return s.repo.Delete(urlID)
}
