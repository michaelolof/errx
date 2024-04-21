package errx

import (
	"errors"
	"fmt"
	"strings"
)

func Is(err error, target error) bool {
	itis := errors.Is(err, target)
	if itis {
		return itis
	} else {
		return strings.Contains(err.Error(), fmt.Sprintf(" kind %s", target.Error()))
	}
}

func Join(errs ...error) error {
	return errors.Join(errs...)
}

func Unwrap(err error) error {
	return errors.Unwrap(err)
}
