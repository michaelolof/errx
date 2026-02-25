package errx

import (
	"errors"
)

func Is(err, target error) bool {
	if errors.Is(err, target) {
		return true
	}
	if err == nil || target == nil {
		return false
	}
	return err.Error() == target.Error()
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

func JoinWrap(ts lint, errs ...error) error {
	if len(errs) == 0 {
		return nil
	}

	return Wrap(ts, Join(errs...))
}
