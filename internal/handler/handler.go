package handler

import (
	"context"
	"errors"
	"regexp"
	"strings"

	"github.com/Falokut/movies_persons_service/internal/models"
	"github.com/Falokut/movies_persons_service/internal/service"
	movies_persons_service "github.com/Falokut/movies_persons_service/pkg/movies_persons_service/v1/protos"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type MoviesPersonsServiceHandler struct {
	movies_persons_service.UnimplementedMoviesPersonsServiceV1Server
	s service.MoviesPersonsService
}

func NewMoviesPersonsServiceHandler(s service.MoviesPersonsService) *MoviesPersonsServiceHandler {
	return &MoviesPersonsServiceHandler{s: s}
}

func (h *MoviesPersonsServiceHandler) GetPersons(ctx context.Context,
	in *movies_persons_service.GetMoviePersonsRequest) (res *movies_persons_service.Persons, err error) {
	defer h.handleError(&err)

	in.PersonsIDs = strings.TrimSpace(strings.ReplaceAll(in.PersonsIDs, `"`, ""))
	ok := regexp.MustCompile("^[,0-9]+$").MatchString(in.PersonsIDs)
	if !ok {
		return nil, status.Error(codes.InvalidArgument, "persons_ids must contains only digits and commas")
	}

	ids := strings.Split(in.PersonsIDs, ",")
	persons, err := h.s.GetPersons(ctx, ids)
	if err != nil {
		return
	}

	return convertRepoPersonsToProto(persons), nil
}

func convertRepoPersonsToProto(persons []models.Person) *movies_persons_service.Persons {
	protoPersons := &movies_persons_service.Persons{
		Persons: make([]*movies_persons_service.Person, len(persons)),
	}

	for i := range persons {
		birthday := ""
		if !persons[i].Birthday.IsZero() {
			birthday = persons[i].Birthday.Format("2006-01-02")
		}
		protoPersons.Persons[i] = &movies_persons_service.Person{
			ID:         persons[i].ID,
			FullnameRU: persons[i].FullnameRU,
			FullnameEN: persons[i].FullnameEN,
			Birthday:   birthday,
			Sex:        persons[i].Sex,
			PhotoUrl:   persons[i].PhotoURL,
		}
	}

	return protoPersons
}

func (h *MoviesPersonsServiceHandler) handleError(err *error) {
	if err == nil || *err == nil {
		return
	}

	serviceErr := &models.ServiceError{}
	if errors.As(*err, &serviceErr) {
		*err = status.Error(convertServiceErrCodeToGrpc(serviceErr.Code), serviceErr.Msg)
	} else if _, ok := status.FromError(*err); !ok {
		e := *err
		*err = status.Error(codes.Unknown, e.Error())
	}
}

func convertServiceErrCodeToGrpc(code models.ErrorCode) codes.Code {
	switch code {
	case models.Internal:
		return codes.Internal
	case models.InvalidArgument:
		return codes.InvalidArgument
	case models.Unauthenticated:
		return codes.Unauthenticated
	case models.Conflict:
		return codes.AlreadyExists
	case models.NotFound:
		return codes.NotFound
	case models.Canceled:
		return codes.Canceled
	case models.DeadlineExceeded:
		return codes.DeadlineExceeded
	case models.PermissionDenied:
		return codes.PermissionDenied
	default:
		return codes.Unknown
	}
}
