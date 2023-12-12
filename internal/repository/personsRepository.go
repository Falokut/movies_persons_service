package repository

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/opentracing/opentracing-go"
)

type postgreRepository struct {
	db *sqlx.DB
}

const (
	personsTableName = "persons"
)

func NewPersonsRepository(db *sqlx.DB) *postgreRepository {
	return &postgreRepository{db: db}
}

func NewPostgreDB(cfg DBConfig) (*sqlx.DB, error) {
	conStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.DBName, cfg.SSLMode)
	db, err := sqlx.Connect("pgx", conStr)

	if err != nil {
		return nil, err
	}

	return db, nil
}

func (r *postgreRepository) Shutdown() {
	r.db.Close()
}

func (r *postgreRepository) GetPersons(ctx context.Context, ids []string) ([]Person, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "postgreRepository.GetPersons")
	defer span.Finish()

	query, args, err := sqlx.In(fmt.Sprintf("SELECT * FROM %s WHERE id IN(?) ORDER BY id",
		personsTableName), ids)
	if err != nil {
		return []Person{}, err
	}
	query = sqlx.Rebind(sqlx.DOLLAR, query)
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return []Person{}, err
	}

	var Persons = make([]Person, 0, len(ids))
	for rows.Next() {
		id, nameRU := "", ""
		nameEN := sql.NullString{}

		if err := rows.Scan(&id, &nameRU, &nameEN); err != nil {
			return []Person{}, err
		}
		Persons = append(Persons, Person{
			ID:         id,
			FullnameRU: nameRU,
			FullnameEN: nameEN,
		})
	}

	return Persons, nil
}
