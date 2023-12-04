package repository

import (
	"context"
	"database/sql"
	"time"
)

type DBConfig struct {
	Host     string `yaml:"host" env:"DB_HOST"`
	Port     string `yaml:"port" env:"DB_PORT"`
	Username string `yaml:"username" env:"DB_USERNAME"`
	Password string `yaml:"password" env:"DB_PASSWORD"`
	DBName   string `yaml:"db_name" env:"DB_NAME"`
	SSLMode  string `yaml:"ssl_mode" env:"DB_SSL_MODE"`
}

type People struct {
	ID             string         `db:"id" json:"id"`
	FullnameRU     string         `db:"fullname_ru" json:"fullname_ru"`
	FullnameEN     sql.NullString `db:"fullname_en" json:"fullname_en,omitempty"`
	BirthCountryID int32          `db:"birth_country_id" json:"birth_country_id"`
}

type Manager interface {
	GetPeople(ctx context.Context, ids []string) ([]People, error)
}

type PeopleRepository interface {
	GetPeople(ctx context.Context, ids []string) ([]People, error)
}

type PeopleCache interface {
	CachePeople(ctx context.Context, people []People, TTL time.Duration) error
	GetPeople(ctx context.Context, ids []string) ([]People, []string, error)
}
