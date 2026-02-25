package errx

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseStapmedError(t *testing.T) {

	notfound := DataKind[string]("not found error")

	err := newErr(1713591173899, "something went wrong").WithKind(DataKind[int]("my_data_1")(30))
	err = wrapErr(1713592763837, err)
	err = wrapErr(1713592780139, err).WithKind(DataKind[int]("my_data_2")(30))
	err = wrapErr(1713591205370, err).WithKind(notfound("something i want to do often"))

	itMatches := IsDataKind(err, notfound)
	assert.Equal(t, itMatches, true)

	msg := err.Error()
	parsed := ParseStampedError(msg)

	itMatches = IsDataKind(parsed, notfound)
	assert.Equal(t, itMatches, true)

	newMsg := parsed.Error()
	assert.Equal(t, newMsg, msg)
}

func TestParseStampedErrorCase2(t *testing.T) {
	err := newErr(1713606995137, "something went wrong").WithKind(DataKind[int]("one")(30))
	err1 := fmt.Errorf("i am a suspect error: %w", err)
	err = wrapErr(1713607005378, err1).WithKind(DataKind[float32]("two")(1.560))
	err = wrapErr(1713607010089, err).WithKind(DataKind[string]("three")("https://www.google.com"))

	msg := err.Error()
	fmt.Println(msg)
	parsed := ParseStampedError(msg)

	newMsg := parsed.Error()
	assert.Equal(t, msg, newMsg)
}

func TestGetStackFrames(t *testing.T) {
	{
		err := errors.New("some thing went wrong")
		err = fmt.Errorf("i am a suspect error: %w", err)
		err = wrapErr(1741664539, err)
		err = wrapErr(1741664541, err)
		err = wrapErr(1741664544, err)

		frames := GetStackFrames(err)
		assert.Len(t, frames, 4)
	}

	{
		err := newErr(1713606995137, "something went wrong")
		err = wrapErr(1713607000211, err)
		err1 := fmt.Errorf("i am a suspect error: %w", err)
		err = wrapErr(1713607005378, err1)
		err = wrapErr(1713607010089, err)

		frames := GetStackFrames(err)
		assert.Len(t, frames, 5)
	}
}

func TestCause_Case1(t *testing.T) {

	err1 := newErr(1715845918044, "something went wrong")
	err := wrapErr(1715845936107, err1)
	err = wrapErr(1715845950562, err)
	err = wrapErr(1715845961777, err)

	errC := Cause(err)

	assert.Equal(t, errC.Error(), err1.Error())
}
func TestMust(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		val := Must(10, nil)
		assert.Equal(t, 10, val)
	})

	t.Run("Panic", func(t *testing.T) {
		assert.Panics(t, func() {
			Must(0, errors.New("fail"))
		})
	})
}

func TestPanic(t *testing.T) {
	assert.Panics(t, func() {
		Panic(errors.New("fail"))
	})
	assert.NotPanics(t, func() {
		Panic(nil)
	})
}

func TestContains(t *testing.T) {
	err := errors.New("hello world")
	assert.True(t, Contains(err, "hello"))
	assert.False(t, Contains(err, "goodbye"))
	assert.False(t, Contains(nil, "any"))
}

func TestSplit(t *testing.T) {
	err1 := errors.New("e1")
	err2 := errors.New("e2")
	joined := Join(err1, err2)

	errs := Split(joined)
	// errors.Join returns an error that implements Unwrap() []error
	assert.Len(t, errs, 2)
}

func TestCauseExtra(t *testing.T) {
	inner := errors.New("inner")
	err := fmt.Errorf("outer: %w", inner)
	assert.Equal(t, inner, Cause(err))
}

func TestCauseMessage(t *testing.T) {
	t.Run("errx tree", func(t *testing.T) {
		err1 := newErr(100, "base failure")
		err := wrapErr(200, err1).WithKind(Kind("some_kind"))
		err = wrapErr(300, err)
		assert.Equal(t, "base failure", CauseMessage(err))
	})

	t.Run("errx tree with wrapped formatted fmt.Errorf", func(t *testing.T) {
		stdErr := fmt.Errorf("formatting user %s failed", "John")
		err1 := fmt.Errorf("wrapped context: %w", stdErr)
		err := wrapErr(100, err1).WithKind(Kind("wrap_kind"))
		assert.Equal(t, "formatting user John failed", CauseMessage(err))
	})

	t.Run("errx tree with nested stamped error containing kind and data", func(t *testing.T) {
		err1 := newErr(100, "deepest stamped failure").WithKind(DataKind[int]("numeric_fault")(42))
		err := wrapErr(200, err1).WithKind(Kind("middle_kind"))
		err = wrapErr(300, err)
		assert.Equal(t, "deepest stamped failure", CauseMessage(err))
	})

	t.Run("errx stringified and parsed back", func(t *testing.T) {
		err1 := newErr(100, "original base message").WithKind(DataKind[string]("db_error")("connection lost"))
		err := wrapErr(200, err1).WithKind(Kind("api_fault"))
		err = wrapErr(300, err)

		formattedStr := err.Error()
		parsedErr := ParseStampedError(formattedStr)

		assert.Equal(t, "original base message", CauseMessage(parsedErr))
	})

	t.Run("errx wrapping fmt.Errorf", func(t *testing.T) {
		stdErr := errors.New("standard error")
		err1 := fmt.Errorf("wrapped context: %w", stdErr)
		err := wrapErr(100, err1).WithKind(Kind("wrap_kind"))
		assert.Equal(t, "standard error", CauseMessage(err))
	})

	t.Run("pure standard error", func(t *testing.T) {
		stdErr := errors.New("just standard")
		assert.Equal(t, "just standard", CauseMessage(stdErr))
	})

	t.Run("nil error", func(t *testing.T) {
		assert.Equal(t, "", CauseMessage(nil))
	})
}
