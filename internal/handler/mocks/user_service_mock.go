package mocks

type UserServiceMock struct {
	RegisterFunc func(name, email, password string) error
	LoginFunc    func(email, password string) (string, error)
}

func (m *UserServiceMock) Register(name, email, password string) error {
	if m.RegisterFunc != nil {
		return m.RegisterFunc(name, email, password)
	}
	return nil
}

func (m *UserServiceMock) Login(email, password string) (string, error) {
	if m.LoginFunc != nil {
		return m.LoginFunc(email, password)
	}
	return "", nil
}
