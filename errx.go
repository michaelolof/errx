package errx

import (
	"encoding/json"
	"errors"
	"fmt"
)

var (
	ErrUnsupported = errors.ErrUnsupported
)

type dataVals interface {
}

type Map[T dataVals] struct {
	mp map[string]T
}

func NewMap[T dataVals](mp map[string]T) Map[T] {
	return Map[T]{mp: mp}
}

func (m *Map[T]) String() string {
	return fmt.Sprintf("%v", m.mp)
}

type unknown struct{}

type DataType interface {
	unknown |
		int | float32 | float64 | string |
		[]int | []float32 | []float64 | []string |
		map[string]int | map[string]float32 | map[string]float64 | map[string]string
}

type Err[T DataType] struct {
	msg     string
	stamp   int
	kindStr string
	data    *T
	dataStr string
	err     error
}

func newE[T DataType](ts int, msg string, kind error, data *T) error {
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

func wrapE[T DataType](ts int, err error, kind error, data *T) error {
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

func newS(ts int, msg string, kind error, dataStr string) error {
	var kindStr string
	if kind != nil {
		kindStr = kind.Error()
	}

	err := &Err[unknown]{
		stamp:   ts,
		msg:     msg,
		kindStr: kindStr,
		dataStr: dataStr,
	}

	err.err = err.buildErr(nil)

	go func() {
		if _logger != nil {
			_logger(err)
		}
	}()

	return err
}

func wrapS(ts int, err error, kind error, dataStr string) error {
	var kindStr string
	if kind != nil {
		kindStr = kind.Error()
	}

	inst := &Err[unknown]{
		stamp:   ts,
		kindStr: kindStr,
		dataStr: dataStr,
	}

	inst.err = inst.buildErr(err)

	go func() {
		if _logger != nil {
			_logger(inst)
		}
	}()

	return inst
}

// New returns an error given a timestamp and error message.
func New(ts int, msg string) error {
	return newE[unknown](ts, msg, nil, nil)
}

// Wrap formats an existing error based on the timestamp given and returns the string as a value that satisfies error.
func Wrap(ts int, err error) error {
	return wrapE[unknown](ts, err, nil, nil)
}

// NewData returns a timestamped error given the timestamp, message and a data value of type T
func NewData[T DataType](ts int, msg string, data T) error {
	return newE(ts, msg, nil, &data)
}

// WrapData wraps an existing error given the timestamp, and a data value of type T
func WrapData[T DataType](ts int, err error, data T) error {
	return wrapE(ts, err, nil, &data)
}

// NewF returns a timestamped error with the message formatted according to a format specifier.
func Newf(ts int, msg string, a ...any) error {
	return newE[unknown](ts, fmt.Sprintf(msg, a...), nil, nil)
}

// NewErr returns a timestamped error given the timestamp, error kind, message and a data value of type T
func NewErr[T DataType](ts int, msg string, kind error, data T) error {
	return newE(ts, msg, kind, &data)
}

// WrapErr wraps an existing error given the timestamp, error kind, and a data value of type T
func WrapErr[T DataType](ts int, err error, kind error, data T) error {
	return wrapE(ts, err, kind, &data)
}

// NewKind returns a timestamped error with a message and given error kind which can be used to provide context or error matching
func NewKind(ts int, msg string, kind error) error {
	return newE[unknown](ts, msg, kind, nil)
}

// WrapKind wraps an existing error given the timestamp and a given error kind which can be used to provide context or error matching
func WrapKind(ts int, err error, kind error) error {
	return wrapE[unknown](ts, err, kind, nil)
}

func (e *Err[T]) buildErr(wrapErr error) error {
	var details string

	if e.kindStr != "" && e.data != nil {
		details = fmt.Sprintf("[stamp %d kind %s data %s]", e.stamp, e.kindStr, toStr(*e.data))
	} else if e.kindStr != "" && e.dataStr != "" {
		details = fmt.Sprintf("[stamp %d kind %s data %s]", e.stamp, e.kindStr, e.dataStr)
	} else if e.kindStr != "" && e.data == nil {
		details = fmt.Sprintf("[stamp %d kind %s]", e.stamp, e.kindStr)
	} else if e.data != nil && e.kindStr == "" {
		details = fmt.Sprintf("[stamp %d data %s]", e.stamp, toStr(*e.data))
	} else if e.dataStr != "" && e.kindStr == "" {
		details = fmt.Sprintf("[stamp %d data %s]", e.stamp, e.dataStr)
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

func (e *Err[T]) DataStr() string {
	return e.dataStr
}

func (e *Err[T]) Data() *T {
	return e.data
}

func (e *Err[T]) AnyData() *any {
	var val any
	data := e.Data()
	if data != nil {
		val = *data
	}

	return &val
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
func FindData[T DataType](err error) (*T, bool) {
	for err != nil {
		switch t := err.(type) {
		case Stamper[T]:
			data := t.Data()
			if data != nil {
				return data, true
			}
		case Stamper[unknown]:
			data := t.DataStr()
			if data != "" {
				val, err1 := fromStr[T](data)
				if err1 == nil {
					return val, true
				}
			}
		}
		err = Unwrap(err)
	}

	return nil, false
}

// Unwraps the error and retrieves the data values and returns the first one that matches the specified error kind and given type
func FindDataByKind[T DataType](err error, kind error) (*T, bool) {
	for err != nil {
		switch t := err.(type) {
		case Stamper[T]:
			kindStr := t.KindStr()
			if kindStr == kind.Error() {
				data := t.Data()
				if data != nil {
					return data, true
				}
			}
		case Stamper[unknown]:
			kindStr := t.KindStr()
			if kindStr == kind.Error() {
				data := t.DataStr()
				if data != "" {
					var d T
					err := json.Unmarshal([]byte(data), &d)
					if err == nil {
						return &d, true
					}
				}
			}
		}

		err = Unwrap(err)
	}

	return nil, false
}

type AnyStamper interface {
	Stamp() int
	AnyData() any
	DataStr() string
	KindStr() string
}

type Stamper[T DataType] interface {
	Stamp() int
	Data() *T
	DataStr() string
	KindStr() string
}
