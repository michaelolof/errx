package errx

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

var (
	ErrUnsupported = errors.ErrUnsupported
)

func Split(err error) []error {
	if uw, ok := err.(interface{ Unwrap() []error }); ok {
		return uw.Unwrap()
	}
	return []error{err}
}

type Err[T comparable] struct {
	msg     string
	stamp   int
	kindStr string
	data    T
	err     error
}

// New returns an error given a timestamp and error message.
func New(ts int, msg string) error {
	return NewErr[any](ts, msg, nil, nil)
}

// Wrap formats an existing error based on the timestamp given and returns the string as a value that satisfies error.
func Wrap(ts int, err error) error {
	return WrapErr[any](ts, err, nil, nil)
}

// NewData returns a timestamped error given the timestamp, message and a data value of type T
func NewData[T comparable](ts int, msg string, data T) error {
	return NewErr(ts, msg, nil, data)
}

// WrapData wraps an existing error given the timestamp, and a data value of type T
func WrapData[T comparable](ts int, err error, data T) error {
	return WrapErr(ts, err, nil, data)
}

// NewF returns a timestamped error with the message formatted according to a format specifier.
func Newf(ts int, msg string, a ...any) error {
	return New(ts, fmt.Sprintf(msg, a...))
}

// NewErr returns a timestamped error given the timestamp, error kind, message and a data value of type T
func NewErr[T comparable](ts int, msg string, kind error, data T) error {
	var kindStr string
	if kind != nil {
		kindStr = kind.Error()
	}

	err := &Err[T]{
		stamp:   ts,
		msg:     msg,
		kindStr: kindStr,
		data:    data,
	}

	err.err = err.buildErr(nil)

	go func() {
		if _logger != nil {
			_logger(err)
		}
	}()

	return err
}

// WrapErr wraps an existing error given the timestamp, error kind, and a data value of type T
func WrapErr[T comparable](ts int, err error, kind error, data T) error {
	var kindStr string
	if kind != nil {
		kindStr = kind.Error()
	}

	inst := &Err[T]{
		stamp:   ts,
		kindStr: kindStr,
		data:    data,
	}

	inst.err = inst.buildErr(err)

	go func() {
		if _logger != nil {
			_logger(inst)
		}
	}()

	return inst
}

// NewKind returns a timestamped error with a message and given error kind which can be used to provide context or error matching
func NewKind(ts int, msg string, kind error) error {
	return NewErr[any](ts, msg, kind, nil)
}

// WrapKind wraps an existing error given the timestamp and a given error kind which can be used to provide context or error matching
func WrapKind(ts int, err error, kind error) error {
	return WrapErr[any](ts, err, kind, nil)
}

func (e *Err[T]) buildErr(wrapErr error) error {
	var ed T
	var details string
	data, err1 := json.Marshal(e.data)
	if err1 != nil {
		return fmt.Errorf("[stamp %d] %w", e.stamp, err1)
	}

	if e.kindStr != "" && e.data != ed {
		details = fmt.Sprintf("[stamp %d kind %s data %v]", e.stamp, e.kindStr, string(data))
	} else if e.kindStr != "" && e.data == ed {
		details = fmt.Sprintf("[stamp %d kind %s]", e.stamp, e.kindStr)
	} else if e.data != ed && e.kindStr == "" {
		details = fmt.Sprintf("[stamp %d data %v]", e.stamp, string(data))
	} else {
		details = fmt.Sprintf("[stamp %d]", e.stamp)
	}

	if wrapErr == nil && e.msg != "" {
		return fmt.Errorf("%s %s", details, e.msg)
	} else if wrapErr != nil {
		return fmt.Errorf("%s; %w", details, wrapErr)
	}
	return nil
}

func (e *Err[T]) Error() string {
	return e.err.Error()
}

func (e *Err[T]) Is(target error) bool {
	itis := errors.Is(e.err, target)
	if itis {
		return itis
	} else {
		return e.err.Error() == target.Error()
	}
}

func (e *Err[T]) Data() any {
	return e.data
}

func (e *Err[T]) KindStr() string {
	if e.kindStr != "" {
		return e.kindStr
	}
	return ""
}

func (e *Err[T]) Stamp() int {
	return e.stamp
}

// Adds support for [errors.Unwrap] function
func (e Err[T]) Unwrap() error {
	return errors.Unwrap(e.err)
}

// Unwraps the error retrieves the data values and returns the first one that matches the givent type
func FindData[T comparable](err error) (*T, bool) {
	for err != nil {
		if se, ok := err.(Stamper); ok {
			data := se.Data()

			if v, vok := data.(T); vok {
				return &v, true
			} else if v, vok := data.(map[string]any); vok {
				val, err1 := parseMapToStruct[T](v)
				if err1 != nil {
					err = Unwrap(err)
					continue
				}

				return val, true
			}
		}
		err = Unwrap(err)
	}

	return nil, false
}

// Unwraps the error and retrieves the data values and returns the first one that matches the specified error kind and given type
func FindDataOfKind[T comparable](err error, kind error) (*T, bool) {
	for err != nil {
		if se, ok := err.(Stamper); ok {
			kindStr := se.KindStr()

			if kindStr == kind.Error() {
				data := se.Data()

				if v, vok := data.(T); vok {
					return &v, true
				} else if v, vok := data.(map[string]any); vok {
					val, err1 := parseMapToStruct[T](v)
					if err1 != nil {
						err = Unwrap(err)
						continue
					}

					return val, true
				}
			}
		}
		err = Unwrap(err)
	}

	return nil, false
}

type Stamper interface {
	Stamp() int
	Data() any
	KindStr() string
}

func parseMapToStruct[T any](dic map[string]any) (*T, error) {

	var structT T
	sv := reflect.ValueOf(&structT)

	for _, f := range reflect.VisibleFields(reflect.TypeOf(structT)) {
		tags := strings.Split(f.Tag.Get("json"), ",")

		var name string
		if len(tags) > 0 && tags[0] != "" {
			name = tags[0]
		} else {
			name = f.Name
		}

		for key, val := range dic {
			if name == key {
				sv.Elem().FieldByName(f.Name).Set(reflect.ValueOf(val))
				break
			}
		}
	}

	return &structT, nil
}
