package service

import (
	"context"
	"strings"

	"github.com/Falokut/grpc_errors"
	"github.com/Falokut/movies_persons_service/internal/repository"
	movies_persons_service "github.com/Falokut/movies_persons_service/pkg/movies_persons_service/v1/protos"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

type MoviesPeoplesService struct {
	movies_persons_service.UnimplementedMoviesPersonsServiceV1Server
	logger        *logrus.Logger
	imagesService *imageService
	repoManager   repository.Manager
	errorHandler  errorHandler
}

func NewMoviesPeoplesService(logger *logrus.Logger, repoManager repository.Manager,
	imagesService *imageService) *MoviesPeoplesService {
	errorHandler := newErrorHandler(logger)
	return &MoviesPeoplesService{
		logger:        logger,
		repoManager:   repoManager,
		errorHandler:  errorHandler,
		imagesService: imagesService,
	}
}

func (s *MoviesPeoplesService) GetPeople(ctx context.Context,
	in *movies_persons_service.GetMoviePersonsRequest) (*movies_persons_service.Persons, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "PeopleService.GetPeople")
	defer span.Finish()
	var err error
	defer span.SetTag("grpc.status", grpc_errors.GetGrpcCode(err))

	if err := validateFilter(in); err != nil {
		return nil, s.errorHandler.createErrorResponce(ErrInvalidFilter, err.Error())
	}

	in.PersonsIDs = strings.TrimSpace(strings.ReplaceAll(in.PersonsIDs, `"`, ""))
	if in.PersonsIDs == "" {
		return &movies_persons_service.Persons{}, nil
	}

	ids := strings.Split(in.PersonsIDs, ",")
	people, err := s.repoManager.GetPersons(ctx, ids)
	if err != nil {
		return nil, s.errorHandler.createErrorResponce(ErrInternal, err.Error())
	}

	return s.convertRepoPeopleToProto(ctx, people), err
}

func (s *MoviesPeoplesService) convertRepoPeopleToProto(ctx context.Context,
	people []repository.Person) *movies_persons_service.Persons {
	protoPersons := &movies_persons_service.Persons{}
	protoPersons.Persons = make(map[string]*movies_persons_service.Person, len(people))
	for _, p := range people {
		protoPersons.Persons[p.ID] = &movies_persons_service.Person{
			ID:         p.ID,
			FullnameRU: p.FullnameRU,
			FullnameEN: p.FullnameEN.String,
			Birthday:   p.Birthday.Time.Format("2006-01-02"),
			Sex:        p.Sex.String,
			PhotoUrl:   s.imagesService.GetPictureURL(ctx, p.PhotoID.String),
		}
	}

	return protoPersons
}
