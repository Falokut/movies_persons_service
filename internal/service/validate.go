package service

import (
	"errors"
	"regexp"

	movies_persons_service "github.com/Falokut/movies_persons_service/pkg/movies_persons_service/v1/protos"
)

var ErrInvalidFilter = errors.New("invalid filter value, filter must contain only digits and commas")

func validateFilter(filter *movies_persons_service.GetMoviePersonsRequest) error {
	if filter.GetPersonsIDs() != "" {
		if err := checkFilterParam(filter.PersonsIDs); err != nil {
			return err
		}
	}

	return nil
}

func checkFilterParam(val string) error {
	exp, err := regexp.Compile("^[!-&!+,0-9]+$")
	if err != nil {
		return err
	}
	if !exp.Match([]byte(val)) {
		return ErrInvalidFilter
	}

	return nil
}
