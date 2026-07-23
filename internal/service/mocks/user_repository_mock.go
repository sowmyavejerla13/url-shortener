package mocks

import "github.com/sowmyavejerla13/url-shortener/internal/model"

type UserRepositoryMock struct {
	GetByEmailFunc func(email string) (*model.User, error)
	CreateFunc     func(user *model.User) error
}

func (m *UserRepositoryMock) GetByEmail(email string) (*model.User, error) {
	if m.GetByEmailFunc != nil {
		return m.GetByEmailFunc(email)
	}
	return nil, nil
}

func (m *UserRepositoryMock) Create(user *model.User) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(user)
	}
	return nil
}
