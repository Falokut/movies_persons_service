package config

import (
	"sync"
	"time"

	"github.com/Falokut/movies_persons_service/internal/repository"
	"github.com/Falokut/movies_persons_service/pkg/jaeger"
	"github.com/Falokut/movies_persons_service/pkg/metrics"
	logging "github.com/Falokut/online_cinema_ticket_office.loggerwrapper"
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	LogLevel        string `yaml:"log_level" env:"LOG_LEVEL"`
	HealthcheckPort string `yaml:"healthcheck_port" env:"HEALTHCHECK_PORT"`

	ImagesService struct {
		BasePhotoUrl   string `yaml:"base_photo_url" env:"BASE_PHOTO_URL"`
		ImagesCategory string `yaml:"pictures_category" env:"PICTURES_CATEGORY"`
	} `yaml:"images_service"`
	Listen struct {
		Host string `yaml:"host" env:"HOST"`
		Port string `yaml:"port" env:"PORT"`
		Mode string `yaml:"server_mode" env:"SERVER_MODE"` // support GRPC, REST, BOTH
	} `yaml:"listen"`

	PrometheusConfig struct {
		Name         string                      `yaml:"service_name" ENV:"PROMETHEUS_SERVICE_NAME"`
		ServerConfig metrics.MetricsServerConfig `yaml:"server_config"`
	} `yaml:"prometheus"`

	MoviesPeoplesCache struct {
		Network          string        `yaml:"network" env:"MOVIES_PERSONS_CACHE_NETWORK"`
		Addr             string        `yaml:"addr" env:"MOVIES_PERSONS_CACHE_ADDR"`
		Password         string        `yaml:"password" env:"MOVIES_PERSONS_CACHE_PASSWORD"`
		DB               int           `yaml:"db" env:"MOVIES_PERSONS_CACHE_DB"`
		MoviesPersonsTTL time.Duration `yaml:"movies_persons_ttl"`
	} `yaml:"movies_persons_cache"`

	DBConfig     repository.DBConfig `yaml:"db_config"`
	JaegerConfig jaeger.Config       `yaml:"jaeger"`
}

var instance *Config
var once sync.Once

const configsPath = "configs/"

func GetConfig() *Config {
	once.Do(func() {
		logger := logging.GetLogger()
		instance = &Config{}

		if err := cleanenv.ReadConfig(configsPath+"config.yml", instance); err != nil {
			help, _ := cleanenv.GetDescription(instance, nil)
			logger.Fatal(help, " ", err)
		}
	})

	return instance
}
