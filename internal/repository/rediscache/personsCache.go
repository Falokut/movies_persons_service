package rediscache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/Falokut/movies_persons_service/internal/models"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"golang.org/x/exp/maps"
)

type PersonsCache struct {
	rdb     *redis.Client
	logger  *logrus.Logger
	metrics Metrics
}

func (c *PersonsCache) PingContext(ctx context.Context) error {
	if err := c.rdb.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("error while pinging persons cache: %w", err)
	}

	return nil
}

func (c *PersonsCache) Shutdown() {
	if err := c.rdb.Close(); err != nil {
		c.logger.Errorf("error while shutting down persons cache: %v", err)
	}
}

func NewPersonsCache(logger *logrus.Logger, opt *redis.Options, metrics Metrics) (*PersonsCache, error) {
	rdb, err := NewRedisClient(opt)
	if err != nil {
		return nil, err
	}

	return &PersonsCache{rdb: rdb, logger: logger, metrics: metrics}, nil
}

func (c *PersonsCache) CachePersons(ctx context.Context, persons []models.RepositoryPerson, ttl time.Duration) (err error) {
	defer c.updateMetrics(err, "CachePersons")
	defer handleError(ctx, &err)
	defer c.logError(err, "CachePersons")

	tx := c.rdb.Pipeline()
	for _, p := range persons {
		toCache, merr := json.Marshal(p)
		if merr != nil {
			err = merr
			return
		}
		err = tx.Set(ctx, p.ID, toCache, ttl).Err()
		if err != nil {
			return
		}
	}
	_, err = tx.Exec(ctx)
	return err
}

func (c *PersonsCache) GetPersons(ctx context.Context,
	ids []string) (persons []models.RepositoryPerson, notFoundIds []string, err error) {
	defer c.updateMetrics(err, "GetPersons")
	defer handleError(ctx, &err)
	defer c.logError(err, "GetPersons")

	var personsIDs = make(map[string]struct{}, len(ids))
	for _, id := range ids {
		personsIDs[id] = struct{}{}
	}

	cached, err := c.rdb.MGet(ctx, ids...).Result()
	if err != nil {
		return
	}

	persons = make([]models.RepositoryPerson, 0, len(cached))
	for _, cache := range cached {
		if cache == nil {
			continue
		}

		person := models.RepositoryPerson{}
		err = json.Unmarshal([]byte(cache.(string)), &persons)
		if err != nil {
			return
		}
		delete(personsIDs, person.ID)
		persons = append(persons, person)
	}

	return persons, maps.Keys(personsIDs), nil
}

func (c *PersonsCache) logError(err error, functionName string) {
	if err == nil {
		return
	}

	var repoErr = &models.ServiceError{}
	if errors.As(err, &repoErr) {
		c.logger.WithFields(
			logrus.Fields{
				"error.function.name": functionName,
				"error.msg":           repoErr.Msg,
				"error.code":          repoErr.Code,
			},
		).Error("persons cache error occurred")
	} else {
		c.logger.WithFields(
			logrus.Fields{
				"error.function.name": functionName,
				"error.msg":           err.Error(),
			},
		).Error("persons cache error occurred")
	}
}

func (c *PersonsCache) updateMetrics(err error, functionName string) {
	if err == nil {
		c.metrics.IncCacheHits(functionName, 1)
		return
	}
	if models.Code(err) == models.NotFound {
		c.metrics.IncCacheMiss(functionName, 1)
	}
}
