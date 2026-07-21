package service

import (
	"errors"

	"github.com/sowmyavejerla13/url-shortener/internal/model"
	"github.com/sowmyavejerla13/url-shortener/internal/repository"
	"github.com/sowmyavejerla13/url-shortener/internal/utils"
)

type URLService struct {
	repo *repository.URLRepository
}

func NewURLService(repo *repository.URLRepository)*URLService{
	return &URLService{
		repo: repo,
	}
}
func (s *URLService)CreateShortURL(userID, originalURL string)(*model.URL, error){
	err := utils.ValidateURL(originalURL)
	if err!=nil{
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
		code, err := utils.GenerateShortCode()
		if err!=nil{
			return nil,err
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
		ShortCode: shortCode,
		OriginalURL: originalURL,
		UserID: userID,
	}

	if err = s.repo.Create(url);err!= nil{
		return nil, err
	}
	return url, nil

}

func (s *URLService) Redirect(shortCode string) (string, error){

	url, err := s.repo.GetByShortCode(shortCode)
	if err!=nil{
		return "",err
	}
	if url == nil{
		return "",errors.New("url not found")
	}
	err = s.repo.IncrementClickCount(url.ID)
	if err!= nil{
		return "", err
	}
	return url.OriginalURL,nil
}

func (s *URLService) GetUserURLs(userID string) ([]model.URL, error){
	urls, err := s.repo.GetByUserID(userID)
	if err!=nil{
		return nil,err
	}

	return urls,nil
}

func (s *URLService) DeleteURL(userID, urlID string) error{
    url, err := s.repo.GetByID(urlID)
	if err!=nil{
		return  err
	}
	if url ==nil{
		return errors.New("no url found")
	}
	if url.UserID != userID {
		return errors.New("forbidden")
	}

	return s.repo.Delete(urlID)
}