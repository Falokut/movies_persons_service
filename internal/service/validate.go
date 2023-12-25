package service

import (
	"errors"
	"regexp"
)

var ErrInvalidFilter = errors.New("invalid filter value, filter must contain only digits and commas")

func checkParam(val string) error {
	exp, err := regexp.Compile("^[!-&!+,0-9]+$")
	if err != nil {
		return err
	}
	if !exp.Match([]byte(val)) {
		return ErrInvalidFilter
	}

	return nil
}
