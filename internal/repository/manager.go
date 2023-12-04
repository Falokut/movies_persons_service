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
	logger    *logrus.Logger
	repo      PeopleRepository
	cache     PeopleCache
	peopleTTL time.Duration
	metric    CacheMetric
}

func NewPeopleRepositoryManager(logger *logrus.Logger, repo PeopleRepository, cache PeopleCache,
	peopleTTL time.Duration, metric CacheMetric) *RepositoryManager {
	return &RepositoryManager{
		logger:    logger,
		repo:      repo,
		cache:     cache,
		peopleTTL: peopleTTL,
		metric:    metric,
	}
}

func (m *RepositoryManager) GetPeople(ctx context.Context, ids []string) ([]People, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "RepositoryManager.GetPeoples")
	defer span.Finish()

	var findedPeoples = make([]People, 0, len(ids))
	people, notFoundedIds, err := m.cache.GetPeople(ctx, ids)
	if errors.Is(err, redis.Nil) {
		m.metric.IncCacheMiss("GetPeoples", len(ids))
	} else if err != nil {
		m.logger.Error(err)
	}

	if len(people) == len(ids) {
		m.metric.IncCacheHits("GetPeoples", len(ids))
		return people, nil
	}

	if len(people) != 0 && err == nil {
		m.metric.IncCacheHits("GetPeoples", len(ids)-len(notFoundedIds))
		m.metric.IncCacheMiss("GetPeoples", len(notFoundedIds))
		ids = notFoundedIds
	}
	findedPeoples = append(findedPeoples, people...)

	people, err = m.repo.GetPeople(ctx, ids)
	if errors.Is(err, sql.ErrNoRows) {
		return people, nil
	}
	if err != nil {
		return []People{}, err
	}

	go func() {
		m.cache.CachePeople(context.Background(), people, m.peopleTTL)
	}()

	findedPeoples = append(findedPeoples, people...)
	return findedPeoples, nil
}
