package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"golang.org/x/exp/maps"
)

type personsCache struct {
	rdb    *redis.Client
	logger *logrus.Logger
}

func (c *personsCache) PingContext(ctx context.Context) error {
	if err := c.rdb.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("error while pinging persons cache: %w", err)
	}

	return nil
}

func (c *personsCache) Shutdown() error {
	return c.rdb.Close()
}

func NewPersonsCache(logger *logrus.Logger, opt *redis.Options) (*personsCache, error) {
	logger.Info("Creating persons cache client")
	rdb := redis.NewClient(opt)
	if rdb == nil {
		return nil, errors.New("can't create new redis client")
	}

	logger.Info("Pinging persons cache client")
	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		return nil, fmt.Errorf("connection is not established: %s", err.Error())
	}

	return &personsCache{rdb: rdb, logger: logger}, nil
}

func (c *personsCache) CachePersons(ctx context.Context, persons []Person, ttl time.Duration) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "personsCache.CachePersons")
	defer span.Finish()

	tx := c.rdb.Pipeline()
	for _, p := range persons {
		toCache, err := json.Marshal(p)
		if err != nil {
			return err
		}
		tx.Set(ctx, fmt.Sprint(p.ID), toCache, ttl)
	}
	_, err := tx.Exec(ctx)
	return err
}

func (c *personsCache) GetPersons(ctx context.Context, ids []string) ([]Person, []string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "personsCache.GetPersons")
	defer span.Finish()

	var personsIDs = make(map[string]struct{}, len(ids))
	for _, id := range ids {
		personsIDs[id] = struct{}{}
	}

	cached, err := c.rdb.MGet(ctx, ids...).Result()
	if err != nil {
		return []Person{}, []string{}, err
	}

	var persons = make([]Person, 0, len(cached))
	for _, cache := range cached {
		if cache == nil {
			continue
		}

		person := Person{}
		err = json.Unmarshal([]byte(cache.(string)), &person)
		if err != nil {
			return []Person{}, []string{}, err
		}
		delete(personsIDs, person.ID)
		persons = append(persons, person)
	}

	return persons, maps.Keys(personsIDs), nil
}
