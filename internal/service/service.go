package service

import (
	"context"
	"strings"

	"github.com/Falokut/grpc_errors"
	"github.com/Falokut/movies_people_service/internal/repository"
	movies_people_service "github.com/Falokut/movies_people_service/pkg/movies_people_service/v1/protos"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

type MoviesPeoplesService struct {
	movies_people_service.UnimplementedMoviesPeopleServiceV1Server
	logger       *logrus.Logger
	repoManager  repository.Manager
	errorHandler errorHandler
}

func NewMoviesPeoplesService(logger *logrus.Logger, repoManager repository.Manager) *MoviesPeoplesService {
	errorHandler := newErrorHandler(logger)
	return &MoviesPeoplesService{
		logger:       logger,
		repoManager:  repoManager,
		errorHandler: errorHandler,
	}
}

func (s *MoviesPeoplesService) GetPeople(ctx context.Context,
	in *movies_people_service.GetMoviePeopleRequest) (*movies_people_service.Humans, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "PeopleService.GetPeople")
	defer span.Finish()
	var err error
	defer span.SetTag("grpc.status", grpc_errors.GetGrpcCode(err))

	if err := validateFilter(in); err != nil {
		return nil, s.errorHandler.createErrorResponce(ErrInvalidFilter, err.Error())
	}

	in.PeopleIDs = strings.TrimSpace(strings.ReplaceAll(in.PeopleIDs, `"`, ""))
	if in.PeopleIDs == "" {
		return &movies_people_service.Humans{}, nil
	}

	ids := strings.Split(in.PeopleIDs, ",")
	people, err := s.repoManager.GetPeople(ctx, ids)
	if err != nil {
		return nil, s.errorHandler.createErrorResponce(ErrInternal, err.Error())
	}

	return convertRepoPeopleToProto(people), err
}

func convertRepoPeopleToProto(people []repository.People) *movies_people_service.Humans {
	protoPeople := &movies_people_service.Humans{}
	protoPeople.People = make(map[string]*movies_people_service.Human, len(people))
	for _, p := range people {
		protoPeople.People[p.ID] = &movies_people_service.Human{
			ID:             p.ID,
			FullnameRU:     p.FullnameRU,
			FullnameEN:     p.FullnameEN.String,
			BirthCountryID: p.BirthCountryID,
		}
	}

	return protoPeople
}
