package repository

import (
	"context"
	"time"

	"github.com/Falokut/movies_persons_service/internal/models"
	"github.com/sirupsen/logrus"
)

type DBConfig struct {
	Host     string `yaml:"host" env:"DB_HOST"`
	Port     string `yaml:"port" env:"DB_PORT"`
	Username string `yaml:"username" env:"DB_USERNAME"`
	Password string `yaml:"password" env:"DB_PASSWORD"`
	DBName   string `yaml:"db_name" env:"DB_NAME"`
	SSLMode  string `yaml:"ssl_mode" env:"DB_SSL_MODE"`
}

type PersonsRepository interface {
	GetPersons(ctx context.Context, ids []string) ([]models.RepositoryPerson, error)
}

type PersonsCache interface {
	CachePersons(ctx context.Context, people []models.RepositoryPerson, TTL time.Duration) error
	GetPersons(ctx context.Context, ids []string) ([]models.RepositoryPerson, []string, error)
}

type personsRepository struct {
	logger     *logrus.Logger
	repo       PersonsRepository
	cache      PersonsCache
	personsTTL time.Duration
}

func NewPersonsRepository(logger *logrus.Logger, repo PersonsRepository, cache PersonsCache,
	personsTTL time.Duration) *personsRepository {
	return &personsRepository{
		logger:     logger,
		repo:       repo,
		cache:      cache,
		personsTTL: personsTTL,
	}
}

func (r *personsRepository) GetPersons(ctx context.Context, ids []string) (persons []models.RepositoryPerson, err error) {
	var findedPersons = make([]models.RepositoryPerson, 0, len(ids))
	persons, notFoundedIds, err := r.cache.GetPersons(ctx, ids)
	if err != nil {
		r.logger.Error(err)
	}

	if len(persons) == len(ids) {
		return persons, nil
	}

	if len(persons) != 0 && err == nil {
		ids = notFoundedIds
	}

	findedPersons = append(findedPersons, persons...)

	persons, err = r.repo.GetPersons(ctx, ids)
	if err != nil {
		return
	}

	go func() {
		err := r.cache.CachePersons(context.Background(), persons, r.personsTTL)
		if err != nil {
			r.logger.Errorf("error while cachin persons: %v", err)
		}
	}()

	findedPersons = append(findedPersons, persons...)
	return findedPersons, nil
}
