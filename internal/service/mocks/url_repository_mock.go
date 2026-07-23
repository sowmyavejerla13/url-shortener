package mocks

import "github.com/sowmyavejerla13/url-shortener/internal/model"

type URLRepositoryMock struct {
	GetByOriginalURLFunc    func(userID, originalURL string) (*model.URL, error)
	GetByShortCodeFunc      func(shortCode string) (*model.URL, error)
	CreateFunc              func(url *model.URL) error
	GetByUserIDFunc         func(userID string) ([]model.URL, error)
	GetByIDFunc             func(id string) (*model.URL, error)
	IncrementClickCountFunc func(id string) error
	DeleteFunc              func(id string) error
}

func (m *URLRepositoryMock) GetByOriginalURL(userID, originalURL string) (*model.URL, error) {
	return m.GetByOriginalURLFunc(userID, originalURL)
}

func (m *URLRepositoryMock) GetByShortCode(shortCode string) (*model.URL, error) {
	if m.GetByShortCodeFunc != nil {
		return m.GetByShortCodeFunc(shortCode)
	}
	return nil, nil
}

func (m *URLRepositoryMock) Create(url *model.URL) error {
	return m.CreateFunc(url)
}

func (m *URLRepositoryMock) IncrementClickCount(id string) error {
	return m.IncrementClickCountFunc(id)
}
func (m *URLRepositoryMock) GetByUserID(userID string) ([]model.URL, error) {
	if m.GetByUserIDFunc != nil {
		return m.GetByUserIDFunc(userID)
	}
	return nil, nil
}

func (m *URLRepositoryMock) GetByID(id string) (*model.URL, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(id)
	}
	return nil, nil
}

func (m *URLRepositoryMock) Delete(id string) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(id)
	}
	return nil
}
