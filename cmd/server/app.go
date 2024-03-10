package main

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"

	server "github.com/Falokut/grpc_rest_server"
	"github.com/Falokut/healthcheck"
	"github.com/Falokut/movies_persons_service/internal/config"
	"github.com/Falokut/movies_persons_service/internal/handler"
	"github.com/Falokut/movies_persons_service/internal/repository"
	"github.com/Falokut/movies_persons_service/internal/repository/postgresrepository"
	"github.com/Falokut/movies_persons_service/internal/repository/rediscache"
	"github.com/Falokut/movies_persons_service/internal/service"
	jaegerTracer "github.com/Falokut/movies_persons_service/pkg/jaeger"
	"github.com/Falokut/movies_persons_service/pkg/logging"
	"github.com/Falokut/movies_persons_service/pkg/metrics"
	movies_persons_service "github.com/Falokut/movies_persons_service/pkg/movies_persons_service/v1/protos"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/opentracing/opentracing-go"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

func main() {
	logging.NewEntry(logging.ConsoleOutput)
	logger := logging.GetLogger()
	cfg := config.GetConfig()

	logLevel, err := logrus.ParseLevel(cfg.LogLevel)
	if err != nil {
		logger.Fatal(err)
	}
	logger.Logger.SetLevel(logLevel)

	tracer, closer, err := jaegerTracer.InitJaeger(cfg.JaegerConfig)
	if err != nil {
		logger.Errorf("Shutting down, error while creating tracer %v", err)
		return
	}
	logger.Info("Jaeger connected")
	defer closer.Close()
	opentracing.SetGlobalTracer(tracer)

	logger.Info("Metrics initializing")
	metric, err := metrics.CreateMetrics(cfg.PrometheusConfig.Name)
	if err != nil {
		logger.Errorf("Shutting down, error while creating metrics %v", err)
		return
	}

	shutdown := make(chan error, 1)
	go func() {
		logger.Info("Metrics server running")
		if serr := metrics.RunMetricServer(cfg.PrometheusConfig.ServerConfig); serr != nil {
			logger.Errorf("Shutting down, error while running metrics server %v", serr)
			shutdown <- serr
			return
		}
	}()

	logger.Info("Database initializing")
	database, err := postgresrepository.NewPostgreDB(&cfg.DBConfig)
	if err != nil {
		logger.Errorf("Shutting down, connection to the database is not established: %s", err.Error())
		return
	}
	defer database.Close()

	logger.Info("Repository initializing")
	personsRepository := postgresrepository.NewPersonsRepository(logger.Logger, database)

	cache, err := rediscache.NewPersonsCache(logger.Logger, getCacheOptions(cfg), metric)
	if err != nil {
		logger.Errorf("Shutting down, connection to the cache is not established: %s", err.Error())
		return
	}
	defer cache.Shutdown()

	logger.Info("Healthcheck initializing")
	healthcheckManager := healthcheck.NewHealthManager(logger.Logger,
		[]healthcheck.HealthcheckResource{database, cache}, cfg.HealthcheckPort, nil)
	go func() {
		logger.Info("Healthcheck server running")
		if err := healthcheckManager.RunHealthcheckEndpoint(); err != nil {
			logger.Errorf("Shutting down, error while running healthcheck endpoint %s", err.Error())
			shutdown <- err
			return
		}
	}()
	repo := repository.NewPersonsRepository(logger.Logger, personsRepository,
		cache, cfg.MoviesPersonsCache.MoviesPersonsTTL)

	logger.Info("Service initializing")
	s := service.NewMoviesPersonsService(logger.Logger,
		repo, service.MoviesPersonsServiceConfig{
			BasePhotoURL:     cfg.BasePhotoURL,
			PicturesCategory: cfg.PhotoCategory,
		})

	logger.Info("Handler initializing")
	h := handler.NewMoviesPersonsServiceHandler(s)

	logger.Info("Server initializing")
	serv := server.NewServer(logger.Logger, h)
	go func() {
		if serr := serv.Run(getListenServerConfig(cfg), metric, nil, nil); serr != nil {
			logger.Errorf("Shutting down, error while running server %s", serr.Error())
			shutdown <- serr
			return
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGHUP, syscall.SIGTERM)

	select {
	case <-quit:
		break
	case <-shutdown:
		break
	}

	serv.Shutdown()
}

func getListenServerConfig(cfg *config.Config) server.Config {
	return server.Config{
		Mode:        cfg.Listen.Mode,
		Host:        cfg.Listen.Host,
		Port:        cfg.Listen.Port,
		ServiceDesc: &movies_persons_service.MoviesPersonsServiceV1_ServiceDesc,
		RegisterRestHandlerServer: func(ctx context.Context, mux *runtime.ServeMux, service any) error {
			serv, ok := service.(movies_persons_service.MoviesPersonsServiceV1Server)
			if !ok {
				return errors.New("can't convert")
			}
			return movies_persons_service.RegisterMoviesPersonsServiceV1HandlerServer(ctx,
				mux, serv)
		},
	}
}

func getCacheOptions(cfg *config.Config) *redis.Options {
	return &redis.Options{
		Network:  cfg.MoviesPersonsCache.Network,
		Addr:     cfg.MoviesPersonsCache.Addr,
		Password: cfg.MoviesPersonsCache.Password,
		DB:       cfg.MoviesPersonsCache.DB,
	}
}
