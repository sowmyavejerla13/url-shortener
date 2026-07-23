package mocks

import "github.com/sowmyavejerla13/url-shortener/internal/model"

type URLServiceMock struct {
	CreateShortURLFunc func(userID, originalURL string) (*model.URL, error)
	RedirectFunc       func(shortCode string) (string, error)
	GetUserURLsFunc    func(userID string) ([]model.URL, error)
	DeleteURLFunc      func(userID, urlID string) error
}

func (m *URLServiceMock) CreateShortURL(userID, originalURL string) (*model.URL, error) {
	if m.CreateShortURLFunc != nil {
		return m.CreateShortURLFunc(userID, originalURL)
	}
	return nil, nil
}

func (m *URLServiceMock) Redirect(shortCode string) (string, error) {
	if m.RedirectFunc != nil {
		return m.RedirectFunc(shortCode)
	}
	return "", nil
}

func (m *URLServiceMock) GetUserURLs(userID string) ([]model.URL, error) {
	if m.GetUserURLsFunc != nil {
		return m.GetUserURLsFunc(userID)
	}
	return nil, nil
}

func (m *URLServiceMock) DeleteURL(userID, urlID string) error {
	if m.DeleteURLFunc != nil {
		return m.DeleteURLFunc(userID, urlID)
	}
	return nil
}
