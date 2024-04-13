package errx

import (
	"encoding/json"
	"errors"
	"fmt"
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
	err  error
	data T
}

// NewData returns a timestamped error given the timestamp, message and a data value of type T
func NewData[T comparable](ts int, msg string, data T) error {
	var empD T
	var erro error
	if data == empD {
		erro = fmt.Errorf("STAMP-%d: %s", ts, msg)
	} else {
		erro = fmt.Errorf("STAMP-%d: %s: ERRDATA-%v", ts, msg, data)
	}

	if _logger != nil {
		_logger(erro)
	}

	return &Err[T]{
		err:  erro,
		data: data,
	}
}

// WrapData wraps an existing error given the timestamp, and a data value of type T
func WrapData[T comparable](ts int, err error, data T) error {
	var empD T
	var erro error
	if data == empD {
		erro = fmt.Errorf("STAMP-%d: %w", ts, err)
	} else {
		erro = fmt.Errorf("STAMP-%d: %w: ERRDATA-%v", ts, err, data)
	}

	if _logger != nil {
		_logger(erro)
	}

	return &Err[T]{
		err:  erro,
		data: data,
	}
}

// New returns an error given a timestamp and error message.
func New(ts int, msg string) error {
	return NewData[any](ts, msg, nil)
}

// Wrap formats an existing error based on the timestamp given and returns the string as a value that satisfies error.
func Wrap(ts int, err error) error {
	return WrapData[any](ts, err, nil)
}

// NewFmt returns a timestamped error with the message formatted according to a format specifier.
func NewFmt(ts int, msg string, a ...any) error {
	return New(ts, fmt.Sprintf(msg, a...))
}

// NewKind returns a timestamped error given the timestamp, error kind, message and a data value of type T
func NewKind[T comparable](ts int, kind error, msg string, data T) error {
	return WrapData(ts, fmt.Errorf("%w: %s", kind, msg), data)
}

// WrapKind wraps an existing error given the timestamp, error kind, and a data value of type T
func WrapKind[T comparable](ts int, kind error, err error, data T) error {
	return WrapData(ts, fmt.Errorf("%w: %w", kind, err), data)
}

func NewOnly(ts int, kind error, msg string) error {
	return NewKind[any](ts, kind, msg, nil)
}

func WrapOnly(ts int, kind error, err error) error {
	return WrapKind[any](ts, kind, err, nil)
}

func (e *Err[T]) Error() string {
	return e.err.Error()
}

// Adds support for [errors.Is] function
func (e Err[T]) Is(target error) bool {
	return errors.Is(e.err, target)
}

// Adds support for [errors.Unwrap] function
func (e Err[T]) Unwrap() error {
	return errors.Unwrap(e.err)
}

// Adss support for marshaling error object
func (e *Err[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(newErrObj(e.Error(), e.data))
}

// Adss support for unmarshaling error object
func (e *Err[T]) UnmarshalJSON(data []byte) error {
	var obj ErrObj
	if err := json.Unmarshal(data, &obj); err != nil {
		return err
	}

	var dataVal T
	if obj.Data != nil {
		v, ok := obj.Data.(T)
		if ok {
			dataVal = v
		}
	}
	e.err = &Err[T]{
		err:  errors.New(obj.Msg),
		data: dataVal,
	}
	e.data = dataVal
	return nil
}

// Tries to retrieve and cast a data value on an error object to given type T
func GetData[T comparable](err error) (T, error) {
	var empty T
	var errd *Err[T]
	if errors.As(err, &errd) {
		return errd.data, nil
	}
	return empty, fmt.Errorf("couldn't cast error data to given type [T]: %w", err)
}

type ErrObj struct {
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

func newErrObj(msg string, data any) ErrObj {
	return ErrObj{
		Msg:  msg,
		Data: data,
	}
}
