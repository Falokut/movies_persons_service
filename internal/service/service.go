package service

import (
	"context"

	"github.com/Falokut/movies_persons_service/internal/models"
	"github.com/Falokut/movies_persons_service/internal/repository"
	"github.com/sirupsen/logrus"
)

type MoviesPersonsService interface {
	GetPersons(ctx context.Context, ids []string) ([]models.Person, error)
}

type MoviesPersonsServiceConfig struct {
	BasePhotoURL     string
	PicturesCategory string
}
type moviesPersonsService struct {
	logger *logrus.Logger
	repo   repository.PersonsRepository
	cfg    MoviesPersonsServiceConfig
}

func NewMoviesPersonsService(logger *logrus.Logger,
	repo repository.PersonsRepository,
	cfg MoviesPersonsServiceConfig) *moviesPersonsService {
	return &moviesPersonsService{
		logger: logger,
		repo:   repo,
		cfg:    cfg,
	}
}

func (s *moviesPersonsService) GetPersons(ctx context.Context, ids []string) (persons []models.Person, err error) {
	repopersons, err := s.repo.GetPersons(ctx, ids)
	persons = make([]models.Person, len(repopersons))
	for i := range repopersons {
		persons[i] = models.Person{
			ID:         repopersons[i].ID,
			FullnameRU: repopersons[i].FullnameRU,
			FullnameEN: repopersons[i].FullnameEN,
			Birthday:   repopersons[i].Birthday,
			Sex:        repopersons[i].Sex,
			PhotoURL:   getPictureURL(repopersons[i].PhotoID, s.cfg.BasePhotoURL, s.cfg.PicturesCategory),
		}
	}
	return
}

func getPictureURL(pictureID, baseURL, category string) string {
	if pictureID == "" || baseURL == "" || category == "" {
		return ""
	}

	return baseURL + "/" + category + "/" + pictureID
}
