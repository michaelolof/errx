package errx

import (
	"fmt"
)

// A literal int
type lint int
type errKind struct {
	kind string
	data dataValue
}

type StampedErr interface {
	Msg() string
	Stamp() int
	Kind() string
}

type errx struct {
	ts   lint
	kind errKind
	msg  string
	err  error
	errx *errx
}

// Implements the error interface by returning the error string
func (e *errx) Error() string {
	return buildErrx(e).Error()
}

// Returns the string representation of the errx object.
func (e *errx) String() string {
	return e.Error()
}

// Returns the error message.
func (e *errx) Msg() string {
	return e.msg
}

// Unwraps the error object.
func (e *errx) Unwrap() error {
	if e.errx != nil {
		return e.errx
	} else if e.err != nil {
		return e.err
	} else {
		return nil
	}
}

// Returns the list of stamp traces for a given error.
func (e *errx) Stamps() []int {
	rtn := make([]int, 0, 15)
	for {
		if err := e.Unwrap(); err != nil {
			if v, ok := err.(interface{ Stamp() int }); ok {
				rtn = append(rtn, v.Stamp())
			}
		} else {
			break
		}
	}
	return rtn
}

// Returns the error interface for the errx instance
func (e *errx) Err() error {
	return e
}

// Returns the error stamp for the given error
func (e *errx) Stamp() int {
	return int(e.ts)
}

// Returns the kind of error it is
func (e *errx) Kind() string {
	return e.kind.kind
}

// Add an error kind to your error object.
func (e *errx) WithKind(kind errKind) *errx {
	e.kind = kind
	return e
}

// Create a new errx instance and add properties to it using the builder pattern.
func NewBuild(ts int, msg string) *errx {
	return newErr(lint(ts), msg)
}

// Wraps am existing error into a new errx instance and add properties to it using the builder pattern.
func BuildFrom(ts int, err error) *errx {
	return wrapErr(lint(ts), err)
}

// New returns an error given a timestamp and error message.
func New(ts lint, msg string) error {
	return newErr(ts, msg)
}

// Wrap formats an existing error based on the timestamp given and returns the string as a value that satisfies error.
func Wrap(ts lint, err error) error {
	return wrapErr(ts, err)
}

// NewF returns a timestamped error with the message formatted according to a format specifier.
func Newf(ts lint, pattern string, a ...any) error {
	return newErrf(ts, pattern, a...)
}

// Wrapf formats an existing error based on the timestamp and formats the existing error message according to the format specifier defined
func Wrapf(ts lint, pattern string, err error, a ...any) error {
	return wrapErrf(ts, pattern, err, a...)
}

// NewKind returns a timestamped error with a message and given error kind which can be used to provide context or error matching
func NewKind(ts lint, kind errKind, msg string) error {
	return newErr(ts, msg).WithKind(kind)
}

// WrapKind wraps an existing error given the timestamp and a given error kind which can be used to provide context or error matching
func WrapKind(ts lint, kind errKind, err error) error {
	return wrapErr(ts, err).WithKind(kind)
}

// NewKind returns a timestamped error with a message and given error kind which can be used to provide context or error matching
func NewKindf(ts lint, kind errKind, msg string, a ...any) error {
	return newErrf(ts, msg, a...).WithKind(kind)
}

// WrapKind wraps an existing error given the timestamp and a given error kind which can be used to provide context or error matching
func WrapKindf(ts lint, kind errKind, pattern string, err error, a ...any) error {
	return wrapErrf(ts, pattern, err, a...).WithKind(kind)
}

func newErr(ts lint, msg string) *errx {
	return &errx{ts: ts, msg: msg}
}

func wrapErr(ts lint, err error) *errx {
	switch e := err.(type) {
	case *errx:
		return &errx{ts: ts, errx: e}
	default:
		return &errx{ts: ts, err: err}
	}
}

func newErrf(ts lint, pattern string, a ...any) *errx {
	return newErr(ts, fmt.Sprintf(pattern, a...))
}

func wrapErrf(ts lint, pattern string, err error, a ...any) *errx {
	arr := make([]any, 0, len(a)+1)
	switch v := err.(type) {
	case *errx:
		arr = append(arr, v.msg)
		arr = append(arr, a...)
		v.msg = fmt.Sprintf(pattern, arr...)
		return wrapErr(ts, v)
	default:
		arr = append(arr, err)
		arr = append(arr, a...)
		err2 := fmt.Errorf(pattern, arr...)
		return wrapErr(ts, err2)
	}
}

func buildErrx(e *errx) error {
	var details string

	if e.kind.kind != "" && e.kind.data.isSet {
		details = fmt.Sprintf("[ts %d kind %s data %s]", e.ts, e.kind.kind, e.kind.data.String())
	} else if e.kind.kind != "" && !e.kind.data.isSet {
		details = fmt.Sprintf("[ts %d kind %s]", e.ts, e.kind.kind)
	} else if e.kind.data.isSet && e.kind.kind == "" {
		details = fmt.Sprintf("[ts %d data %s]", e.ts, e.kind.data.String())
	} else if e.ts != 0 {
		details = fmt.Sprintf("[ts %d]", e.ts)
	}

	if e.errx != nil {
		return fmt.Errorf("%s; %s", details, buildErrx(e.errx).Error())
	} else if e.err != nil {
		return fmt.Errorf("%s; %s", details, e.err.Error())
	} else if details != "" {
		return fmt.Errorf("%s %s", details, e.msg)
	} else {
		return fmt.Errorf("%s", e.msg)
	}

}
