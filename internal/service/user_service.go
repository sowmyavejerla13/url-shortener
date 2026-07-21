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
	repo *repository.UserRepository
	config * config.Config
}

func NewUserService(repo *repository.UserRepository,cfg *config.Config)*UserService  {
	return &UserService{
		repo: repo,
		config: cfg,
	}
}

func (s *UserService)Register(name, email, password string)error{
	existingUser, err:= s.repo.GetByEmail(email)
	if err !=nil{
		return err
	}
	if existingUser !=nil{
		return errors.New("email already exists")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)

	if err!=nil{
		return err
	}
	user := &model.User{
		Name: name,
		Email: email,
		PasswordHash: string(hashedPassword),
	}
	return s.repo.Create(user)

}

func (s *UserService)Login(email , password string)(string,error){
		user, err := s.repo.GetByEmail(email)
		if err!=nil{
			return "",err
		}

		if user == nil{
			return "",errors.New("invalid email or password")
		}
		err = bcrypt.CompareHashAndPassword(
			[]byte(user.PasswordHash),
			[]byte(password),
		)

		if err !=nil{
			return "", errors.New("invalid email and password")
		}
		token, err := utils.GenerateToken(user.ID, s.config.JWTSecret)
		if err!=nil{
			return "", err
		}
		return token, nil

}