package service

import (
	"context"
	"errors"
	"strings"

	"github.com/Falokut/movies_persons_service/internal/repository"
	movies_persons_service "github.com/Falokut/movies_persons_service/pkg/movies_persons_service/v1/protos"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
)

type MoviesPersonsService struct {
	movies_persons_service.UnimplementedMoviesPersonsServiceV1Server
	logger        *logrus.Logger
	imagesService *imageService
	repoManager   repository.Manager
	errorHandler  errorHandler
}

func NewMoviesPersonsService(logger *logrus.Logger, repoManager repository.Manager,
	imagesService *imageService) *MoviesPersonsService {
	errorHandler := newErrorHandler(logger)
	return &MoviesPersonsService{
		logger:        logger,
		repoManager:   repoManager,
		errorHandler:  errorHandler,
		imagesService: imagesService,
	}
}

func (s *MoviesPersonsService) GetPersons(ctx context.Context,
	in *movies_persons_service.GetMoviePersonsRequest) (*movies_persons_service.Persons, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "PeopleService.GetPeople")
	defer span.Finish()

	if err := validateFilter(in); err != nil {
		return nil, s.errorHandler.createErrorResponceWithSpan(span, ErrInvalidFilter, err.Error())
	}

	in.PersonsIDs = strings.TrimSpace(strings.ReplaceAll(in.PersonsIDs, `"`, ""))
	if in.PersonsIDs == "" {
		return &movies_persons_service.Persons{}, nil
	}

	ids := strings.Split(in.PersonsIDs, ",")
	people, err := s.repoManager.GetPersons(ctx, ids)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, s.errorHandler.createErrorResponceWithSpan(span, ErrNotFound, "")
	}
	if err != nil {
		return nil, s.errorHandler.createErrorResponceWithSpan(span, ErrInternal, err.Error())
	}

	span.SetTag("grpc.status", codes.OK)
	return s.convertRepoPeopleToProto(ctx, people), nil
}

func (s *MoviesPersonsService) convertRepoPeopleToProto(ctx context.Context,
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
