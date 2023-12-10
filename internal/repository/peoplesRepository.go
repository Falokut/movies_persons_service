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
	peopleTableName = "people"
)

func NewPeopleRepository(db *sqlx.DB) *postgreRepository {
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

func (r *postgreRepository) GetPeople(ctx context.Context, ids []string) ([]People, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "postgreRepository.GetPeoples")
	defer span.Finish()

	query, args, err := sqlx.In(fmt.Sprintf("SELECT * FROM %s WHERE id IN(?) ORDER BY id",
		peopleTableName), ids)
	if err != nil {
		return []People{}, err
	}
	query = sqlx.Rebind(sqlx.DOLLAR, query)
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return []People{}, err
	}

	var people = make([]People, 0, len(ids))
	for rows.Next() {
		id, nameRU := "", ""
		nameEN := sql.NullString{}

		if err := rows.Scan(&id, &nameRU, &nameEN); err != nil {
			return []People{}, err
		}
		people = append(people, People{
			ID:         id,
			FullnameRU: nameRU,
			FullnameEN: nameEN,
		})
	}

	return people, nil
}
