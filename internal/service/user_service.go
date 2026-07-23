package service

import (
	"errors"

	"github.com/sowmyavejerla13/url-shortener/internal/config"
	"github.com/sowmyavejerla13/url-shortener/internal/model"
	"github.com/sowmyavejerla13/url-shortener/internal/repository"
	"github.com/sowmyavejerla13/url-shortener/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo   repository.UserRepositoryInterface
	config *config.Config

	hashPassword    func([]byte, int) ([]byte, error)
	comparePassword func([]byte, []byte) error
	generateToken   func(string, string) (string, error)
}

// Production constructor
func NewUserService(
	repo repository.UserRepositoryInterface,
	cfg *config.Config,
) *UserService {

	return &UserService{
		repo:            repo,
		config:          cfg,
		hashPassword:    bcrypt.GenerateFromPassword,
		comparePassword: bcrypt.CompareHashAndPassword,
		generateToken:   utils.GenerateToken,
	}
}

// Test constructor
func NewUserServiceWithDependencies(
	repo repository.UserRepositoryInterface,
	cfg *config.Config,
	hashPassword func([]byte, int) ([]byte, error),
	comparePassword func([]byte, []byte) error,
	generateToken func(string, string) (string, error),
) *UserService {

	return &UserService{
		repo:            repo,
		config:          cfg,
		hashPassword:    hashPassword,
		comparePassword: comparePassword,
		generateToken:   generateToken,
	}
}

type UserServiceInterface interface {
	Register(name, email, password string) error
	Login(email, password string) (string, error)
}

func (s *UserService) Register(name, email, password string) error {

	existingUser, err := s.repo.GetByEmail(email)
	if err != nil {
		return err
	}

	if existingUser != nil {
		return errors.New("email already exists")
	}

	hashedPassword, err := s.hashPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return err
	}

	user := &model.User{
		Name:         name,
		Email:        email,
		PasswordHash: string(hashedPassword),
	}

	return s.repo.Create(user)
}

func (s *UserService) Login(email, password string) (string, error) {

	user, err := s.repo.GetByEmail(email)
	if err != nil {
		return "", err
	}

	if user == nil {
		return "", errors.New("invalid email or password")
	}

	err = s.comparePassword(
		[]byte(user.PasswordHash),
		[]byte(password),
	)
	if err != nil {
		return "", errors.New("invalid email and password")
	}

	token, err := s.generateToken(
		user.ID,
		s.config.JWTSecret,
	)
	if err != nil {
		return "", err
	}

	return token, nil
}
