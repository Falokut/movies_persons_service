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
	"github.com/Falokut/movies_persons_service/internal/repository"
	"github.com/Falokut/movies_persons_service/internal/service"
	jaegerTracer "github.com/Falokut/movies_persons_service/pkg/jaeger"
	"github.com/Falokut/movies_persons_service/pkg/metrics"
	movies_persons_service "github.com/Falokut/movies_persons_service/pkg/movies_persons_service/v1/protos"
	logging "github.com/Falokut/online_cinema_ticket_office.loggerwrapper"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/opentracing/opentracing-go"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

func main() {
	logging.NewEntry(logging.ConsoleOutput)
	logger := logging.GetLogger()

	appCfg := config.GetConfig()
	logLevel, err := logrus.ParseLevel(appCfg.LogLevel)
	if err != nil {
		logger.Fatal(err)
	}
	logger.Logger.SetLevel(logLevel)

	tracer, closer, err := jaegerTracer.InitJaeger(appCfg.JaegerConfig)
	if err != nil {
		logger.Fatal("cannot create tracer", err)
	}
	logger.Info("Jaeger connected")
	defer closer.Close()

	opentracing.SetGlobalTracer(tracer)

	logger.Info("Metrics initializing")
	metric, err := metrics.CreateMetrics(appCfg.PrometheusConfig.Name)
	if err != nil {
		logger.Fatal(err)
	}

	go func() {
		logger.Info("Metrics server running")
		if err := metrics.RunMetricServer(appCfg.PrometheusConfig.ServerConfig); err != nil {
			logger.Fatal(err)
		}
	}()

	logger.Info("Database initializing")
	database, err := repository.NewPostgreDB(appCfg.DBConfig)
	if err != nil {
		logger.Fatalf("Shutting down, connection to the database is not established: %s", err.Error())
	}

	logger.Info("Repository initializing")
	repo := repository.NewPersonsRepository(database)
	defer repo.Shutdown()

	cache, err := repository.NewPersonsCache(logger.Logger, getCacheOptions(appCfg))
	if err != nil {
		logger.Fatalf("Shutting down, connection to the cache is not established: %s", err.Error())
	}
	defer cache.Shutdown()

	logger.Info("Healthcheck initializing")
	healthcheckManager := healthcheck.NewHealthManager(logger.Logger,
		[]healthcheck.HealthcheckResource{database, cache}, appCfg.HealthcheckPort, nil)
	go func() {
		logger.Info("Healthcheck server running")
		if err := healthcheckManager.RunHealthcheckEndpoint(); err != nil {
			logger.Fatalf("Shutting down, can't run healthcheck endpoint %s", err.Error())
		}
	}()
	imgService := service.NewImageService(service.ImageServiceConfig{
		BasePhotoUrl:     appCfg.ImagesService.BasePhotoUrl,
		PicturesCategory: appCfg.ImagesService.ImagesCategory,
	}, logger.Logger)
	repoManager := repository.NewPersonsRepositoryManager(logger.Logger, repo,
		cache, appCfg.MoviesPeoplesCache.MoviesPersonsTTL, metric)
	logger.Info("Service initializing")
	service := service.NewMoviesPersonsService(logger.Logger, repoManager, imgService)

	logger.Info("Server initializing")
	s := server.NewServer(logger.Logger, service)
	s.Run(getListenServerConfig(appCfg), metric, nil, nil)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGHUP, syscall.SIGTERM)

	<-quit
	s.Shutdown()
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
			return movies_persons_service.RegisterMoviesPersonsServiceV1HandlerServer(context.Background(),
				mux, serv)
		},
	}
}

func getCacheOptions(cfg *config.Config) *redis.Options {
	return &redis.Options{
		Network:  cfg.MoviesPeoplesCache.Network,
		Addr:     cfg.MoviesPeoplesCache.Addr,
		Password: cfg.MoviesPeoplesCache.Password,
		DB:       cfg.MoviesPeoplesCache.DB,
	}
}
