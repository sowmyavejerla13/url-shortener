package service_test

import (
	"errors"
	"testing"

	"github.com/sowmyavejerla13/url-shortener/internal/config"
	"github.com/sowmyavejerla13/url-shortener/internal/model"
	"github.com/sowmyavejerla13/url-shortener/internal/service"
	"github.com/sowmyavejerla13/url-shortener/internal/service/mocks"
	"github.com/stretchr/testify/assert"
)

var fakeHash = func(password []byte, cost int) ([]byte, error) {
	return []byte("hashed-password"), nil
}

var fakeCompare = func(hash, password []byte) error {
	return nil
}

var fakeToken = func(userID, secret string) (string, error) {
	return "jwt-token", nil
}

func TestRegister(t *testing.T) {

	tests := []struct {
		name string

		nameInput     string
		emailInput    string
		passwordInput string

		existingUser *model.User
		getErr       error

		hashErr error

		createErr error

		expectedErr error
	}{
		{
			name:          "Success",
			nameInput:     "John",
			emailInput:    "john@test.com",
			passwordInput: "password123",
		},

		{
			name:          "Email Already Exists",
			nameInput:     "John",
			emailInput:    "john@test.com",
			passwordInput: "password123",

			existingUser: &model.User{
				ID:    "1",
				Email: "john@test.com",
			},

			expectedErr: errors.New("email already exists"),
		},

		{
			name:          "Repository Error",
			nameInput:     "John",
			emailInput:    "john@test.com",
			passwordInput: "password123",

			getErr: assert.AnError,

			expectedErr: assert.AnError,
		},

		{
			name:          "Hash Password Error",
			nameInput:     "John",
			emailInput:    "john@test.com",
			passwordInput: "password123",

			hashErr: assert.AnError,

			expectedErr: assert.AnError,
		},

		{
			name:          "Create Error",
			nameInput:     "John",
			emailInput:    "john@test.com",
			passwordInput: "password123",

			createErr: assert.AnError,

			expectedErr: assert.AnError,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			mockRepo := &mocks.UserRepositoryMock{}

			mockRepo.GetByEmailFunc = func(email string) (*model.User, error) {
				return tt.existingUser, tt.getErr
			}

			mockRepo.CreateFunc = func(user *model.User) error {
				return tt.createErr
			}

			cfg := &config.Config{
				JWTSecret: "secret",
			}

			hashFn := fakeHash

			if tt.hashErr != nil {
				hashFn = func(password []byte, cost int) ([]byte, error) {
					return nil, tt.hashErr
				}
			}

			userService := service.NewUserServiceWithDependencies(
				mockRepo,
				cfg,
				hashFn,
				fakeCompare,
				fakeToken,
			)
			err := userService.Register(
				tt.nameInput,
				tt.emailInput,
				tt.passwordInput,
			)

			if tt.expectedErr != nil {

				assert.Error(t, err)

				if tt.name == "Email Already Exists" {
					assert.EqualError(t, err, "email already exists")
				} else {
					assert.Equal(t, tt.expectedErr, err)
				}

			} else {

				assert.NoError(t, err)

			}
		})
	}
}

func TestLogin(t *testing.T) {

	tests := []struct {
		name string

		email    string
		password string

		mockUser *model.User
		getErr   error

		compareErr error

		token    string
		tokenErr error

		expectedToken string
		expectedErr   error
	}{
		{
			name:     "Success",
			email:    "john@test.com",
			password: "password123",

			mockUser: &model.User{
				ID:           "1",
				Email:        "john@test.com",
				PasswordHash: "hashed-password",
			},

			token:         "jwt-token",
			expectedToken: "jwt-token",
		},

		{
			name:     "Repository Error",
			email:    "john@test.com",
			password: "password123",

			getErr: assert.AnError,

			expectedErr: assert.AnError,
		},

		{
			name:     "User Not Found",
			email:    "john@test.com",
			password: "password123",

			mockUser: nil,
		},

		{
			name:     "Invalid Password",
			email:    "john@test.com",
			password: "wrong-password",

			mockUser: &model.User{
				ID:           "1",
				Email:        "john@test.com",
				PasswordHash: "hashed-password",
			},

			compareErr: assert.AnError,
		},

		{
			name:     "Generate Token Error",
			email:    "john@test.com",
			password: "password123",

			mockUser: &model.User{
				ID:           "1",
				Email:        "john@test.com",
				PasswordHash: "hashed-password",
			},

			tokenErr: assert.AnError,

			expectedErr: assert.AnError,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			mockRepo := &mocks.UserRepositoryMock{}

			mockRepo.GetByEmailFunc = func(email string) (*model.User, error) {
				return tt.mockUser, tt.getErr
			}

			cfg := &config.Config{
				JWTSecret: "secret",
			}

			compareFn := fakeCompare
			if tt.compareErr != nil {
				compareFn = func(hash, password []byte) error {
					return tt.compareErr
				}
			}

			tokenFn := fakeToken
			if tt.tokenErr != nil {
				tokenFn = func(userID, secret string) (string, error) {
					return "", tt.tokenErr
				}
			} else {
				tokenFn = func(userID, secret string) (string, error) {
					return tt.token, nil
				}
			}

			userService := service.NewUserServiceWithDependencies(
				mockRepo,
				cfg,
				fakeHash,
				compareFn,
				tokenFn,
			)

			token, err := userService.Login(
				tt.email,
				tt.password,
			)

			switch tt.name {

			case "Success":

				assert.NoError(t, err)
				assert.Equal(t, tt.expectedToken, token)

			case "User Not Found":

				assert.Error(t, err)
				assert.EqualError(t, err, "invalid email or password")
				assert.Empty(t, token)

			case "Invalid Password":

				assert.Error(t, err)
				assert.EqualError(t, err, "invalid email and password")
				assert.Empty(t, token)

			default:

				assert.Error(t, err)
				assert.Equal(t, tt.expectedErr, err)
				assert.Empty(t, token)
			}
		})
	}
}
