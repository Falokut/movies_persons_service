package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

type CacheMetric interface {
	IncCacheHits(method string, times int)
	IncCacheMiss(method string, times int)
}

type RepositoryManager struct {
	logger     *logrus.Logger
	repo       PersonsRepository
	cache      PersonsCache
	personsTTL time.Duration
	metric     CacheMetric
}

func NewPersonsRepositoryManager(logger *logrus.Logger, repo PersonsRepository, cache PersonsCache,
	personsTTL time.Duration, metric CacheMetric) *RepositoryManager {
	return &RepositoryManager{
		logger:     logger,
		repo:       repo,
		cache:      cache,
		personsTTL: personsTTL,
		metric:     metric,
	}
}

func (m *RepositoryManager) GetPersons(ctx context.Context, ids []string) ([]Person, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "RepositoryManager.GetPersonss")
	defer span.Finish()

	var findedPersonss = make([]Person, 0, len(ids))
	persons, notFoundedIds, err := m.cache.GetPersons(ctx, ids)
	if errors.Is(err, redis.Nil) {
		m.metric.IncCacheMiss("GetPersonss", len(ids))
	} else if err != nil {
		m.logger.Error(err)
	}

	if len(persons) == len(ids) {
		m.metric.IncCacheHits("GetPersonss", len(ids))
		return persons, nil
	}

	if len(persons) != 0 && err == nil {
		m.metric.IncCacheHits("GetPersonss", len(ids)-len(notFoundedIds))
		m.metric.IncCacheMiss("GetPersonss", len(notFoundedIds))
		ids = notFoundedIds
	}
	findedPersonss = append(findedPersonss, persons...)

	persons, err = m.repo.GetPersons(ctx, ids)
	if errors.Is(err, sql.ErrNoRows) {
		return []Person{}, ErrNotFound
	}
	if err != nil {
		return []Person{}, err
	}

	go func() {
		m.cache.CachePersons(context.Background(), persons, m.personsTTL)
	}()

	findedPersonss = append(findedPersonss, persons...)
	return findedPersonss, nil
}
