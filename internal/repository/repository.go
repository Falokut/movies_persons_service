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

type Person struct {
	ID         string         `db:"id" json:"id"`
	FullnameRU string         `db:"fullname_ru" json:"fullname_ru"`
	FullnameEN sql.NullString `db:"fullname_en" json:"fullname_en,omitempty"`
	Birthday   sql.NullTime   `db:"birthday" json:"birthday,omitempty"`
	Sex        sql.NullString `db:"sex" json:"sex,omitempty"`
	PhotoID    sql.NullString `db:"photo_id" json:"photo_id,omitempty"`
}

type Manager interface {
	GetPersons(ctx context.Context, ids []string) ([]Person, error)
}

type PersonsRepository interface {
	GetPersons(ctx context.Context, ids []string) ([]Person, error)
}

type PersonsCache interface {
	CachePersons(ctx context.Context, people []Person, TTL time.Duration) error
	GetPersons(ctx context.Context, ids []string) ([]Person, []string, error)
}
