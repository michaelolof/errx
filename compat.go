package errx

import (
	"errors"
)

func Is(err, target error) bool {
	itis := errors.Is(err, target)
	if itis {
		return itis
	} else {
		return err.Error() == target.Error()
	}
}

func As(err error, taget any) bool {
	return errors.As(err, taget)
}

func Join(errs ...error) error {
	return errors.Join(errs...)
}

func Unwrap(err error) error {
	return errors.Unwrap(err)
}
