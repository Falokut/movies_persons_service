package postgresrepository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Falokut/movies_persons_service/internal/models"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type PersonsRepository struct {
	db     *sqlx.DB
	logger *logrus.Logger
}

const (
	personsTableName = "persons"
)

func NewPersonsRepository(logger *logrus.Logger, db *sqlx.DB) *PersonsRepository {
	return &PersonsRepository{db: db, logger: logger}
}

type Person struct {
	ID         string         `db:"id" json:"id"`
	FullnameRU string         `db:"fullname_ru"`
	FullnameEN sql.NullString `db:"fullname_en"`
	Birthday   sql.NullTime   `db:"birthday"`
	Sex        sql.NullString `db:"sex"`
	PhotoID    sql.NullString `db:"photo_id"`
}

func (r *PersonsRepository) GetPersons(ctx context.Context, ids []string) (persons []models.RepositoryPerson, err error) {
	defer handleError(ctx, &err)
	defer r.logError(err, "GetPersons")

	query := fmt.Sprintf("SELECT * FROM %s WHERE id=ANY($1) ORDER BY id",
		personsTableName)
	var p []Person
	err = r.db.SelectContext(ctx, &p, query, ids)
	if err != nil {
		return
	}

	persons = make([]models.RepositoryPerson, len(p))
	for i := range p {
		persons[i] = models.RepositoryPerson{
			ID:         p[i].ID,
			FullnameRU: p[i].FullnameRU,
			FullnameEN: p[i].FullnameEN.String,
			Birthday:   p[i].Birthday.Time,
			Sex:        p[i].Sex.String,
			PhotoID:    p[i].PhotoID.String,
		}
	}
	return
}

func (r *PersonsRepository) logError(err error, functionName string) {
	if err == nil {
		return
	}

	var repoErr = &models.ServiceError{}
	if errors.As(err, &repoErr) {
		r.logger.WithFields(
			logrus.Fields{
				"error.function.name": functionName,
				"error.msg":           repoErr.Msg,
				"error.code":          repoErr.Code,
			},
		).Error("persons repository error occurred")
	} else {
		r.logger.WithFields(
			logrus.Fields{
				"error.function.name": functionName,
				"error.msg":           err.Error(),
			},
		).Error("persons repository error occurred")
	}
}
