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

type peopleCache struct {
	rdb    *redis.Client
	logger *logrus.Logger
}

func (c *peopleCache) PingContext(ctx context.Context) error {
	if err := c.rdb.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("error while pinging genres cache: %w", err)
	}

	return nil
}

func (c *peopleCache) Shutdown() {
	c.rdb.Close()
}

func NewPeopleCache(logger *logrus.Logger, opt *redis.Options) (*peopleCache, error) {
	logger.Info("Creating genres cache client")
	rdb := redis.NewClient(opt)
	if rdb == nil {
		return nil, errors.New("can't create new redis client")
	}

	logger.Info("Pinging genres cache client")
	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		return nil, fmt.Errorf("connection is not established: %s", err.Error())
	}

	return &peopleCache{rdb: rdb, logger: logger}, nil
}

func (c *peopleCache) CachePeople(ctx context.Context, people []People, TTL time.Duration) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "peopleCache.CachePeople")
	defer span.Finish()

	tx := c.rdb.Pipeline()
	for _, p := range people {
		toCache, err := json.Marshal(p)
		if err != nil {
			return err
		}
		tx.Set(ctx, fmt.Sprint(p.ID), toCache, TTL)
	}
	_, err := tx.Exec(ctx)
	return err
}

func (c *peopleCache) GetPeople(ctx context.Context, ids []string) ([]People, []string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "peopleCache.GetPeople")
	defer span.Finish()

	var peopleIDs = make(map[string]struct{}, len(ids))
	for _, id := range ids {
		peopleIDs[id] = struct{}{}
	}

	cached, err := c.rdb.MGet(ctx, ids...).Result()
	if err != nil {
		return []People{}, []string{}, err
	}

	var Peoples = make([]People, 0, len(cached))
	for _, cache := range cached {
		if cache == nil {
			continue
		}

		people := People{}
		err = json.Unmarshal([]byte(cache.(string)), &people)
		if err != nil {
			return []People{}, []string{}, err
		}
		delete(peopleIDs, people.ID)
		Peoples = append(Peoples, people)
	}

	return Peoples, maps.Keys(peopleIDs), nil
}
