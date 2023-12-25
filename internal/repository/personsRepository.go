package repository

import (
	"context"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/opentracing/opentracing-go"
)

type personsRepository struct {
	db *sqlx.DB
}

const (
	personsTableName = "persons"
)

func NewPersonsRepository(db *sqlx.DB) *personsRepository {
	return &personsRepository{db: db}
}

func NewPostgreDB(cfg DBConfig) (*sqlx.DB, error) {
	conStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.DBName, cfg.SSLMode)

	return sqlx.Connect("pgx", conStr)
}

func (r *personsRepository) Shutdown() error {
	return r.db.Close()
}

func (r *personsRepository) GetPersons(ctx context.Context, ids []string) ([]Person, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "personsRepository.GetPersons")
	defer span.Finish()

	query := fmt.Sprintf("SELECT * FROM %s WHERE id=ANY($1) ORDER BY id",
		personsTableName)

	var persons []Person
	err := r.db.SelectContext(ctx, &persons, query, ids)
	if err != nil {
		return []Person{}, err
	} else if len(persons) == 0 {
		return []Person{}, ErrNotFound
	}

	return persons, nil
}
