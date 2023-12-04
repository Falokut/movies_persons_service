package service

import (
	"errors"
	"regexp"

	movies_people_service "github.com/Falokut/movies_people_service/pkg/movies_people_service/v1/protos"
)

var ErrInvalidFilter = errors.New("invalid filter value, filter must contain only digits and commas")

func validateFilter(filter *movies_people_service.GetMoviePeopleRequest) error {
	if filter.GetPeopleIDs() != "" {
		if err := checkFilterParam(filter.PeopleIDs); err != nil {
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
