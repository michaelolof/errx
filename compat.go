package errx

import (
	"errors"
)

func IsKind(err error, kind errKind) bool {
	for err != nil {
		if e, ok := err.(interface{ Kind() string }); ok {
			return e.Kind() == kind.kind
		}
		err = Unwrap(err)
	}
	return false
}

func IsDataKind[T DataType](err error, kind func(d T) errKind) bool {
	var d T
	for err != nil {
		if e, ok := err.(interface{ Kind() string }); ok {
			return e.Kind() == kind(d).kind
		}
		err = Unwrap(err)
	}
	return false
}

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
