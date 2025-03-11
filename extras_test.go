package errx

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseStapmedError(t *testing.T) {

	notfound := ErrKind("not found error")

	err := NewErr(1713591173899, "something went wrong").WithData(ErrKind("my_data_1"), Data(30))
	err = WrapErr(1713592763837, err)
	err = WrapErr(1713592780139, err).WithData(ErrKind("my_data_2"), Data(30))
	err = WrapErr(1713591205370, err).WithData(notfound, Data("something i want to do often"))

	itMatches := KindOf(err, notfound)
	assert.Equal(t, itMatches, true)

	msg := err.Error()
	parsed := ParseStampedError(msg)

	itMatches = KindOf(parsed, notfound)
	assert.Equal(t, itMatches, true)

	newMsg := parsed.Error()
	assert.Equal(t, newMsg, msg)
}

func TestParseStampedErrorCase2(t *testing.T) {
	err := NewErr(1713606995137, "something went wrong").WithData(ErrKind("one"), Data(30))
	err1 := fmt.Errorf("i am a suspect error: %w", err)
	err = WrapErr(1713607005378, err1).WithData(ErrKind("two"), Data(1.560))
	err = WrapErr(1713607010089, err).WithData(ErrKind("three"), Data("https://www.google.com"))

	rep := Report(err)
	fmt.Println(rep)

	msg := err.Error()
	parsed := ParseStampedError(msg)

	newMsg := parsed.Error()
	assert.Equal(t, msg, newMsg)
}

func TestGetStackFrames(t *testing.T) {
	{
		err := errors.New("some thing went wrong")
		err = fmt.Errorf("i am a suspect error: %w", err)
		err = WrapErr(1741664539, err)
		err = WrapErr(1741664541, err)
		err = WrapErr(1741664544, err)

		frames := GetStackFrames(err)
		assert.Len(t, frames, 4)
	}

	{
		err := NewErr(1713606995137, "something went wrong")
		err = WrapErr(1713607000211, err)
		err1 := fmt.Errorf("i am a suspect error: %w", err)
		err = WrapErr(1713607005378, err1)
		err = WrapErr(1713607010089, err)

		frames := GetStackFrames(err)
		assert.Len(t, frames, 5)
	}
}

func TestCause_Case1(t *testing.T) {

	err1 := NewErr(1715845918044, "something went wrong")
	err := WrapErr(1715845936107, err1)
	err = WrapErr(1715845950562, err)
	err = WrapErr(1715845961777, err)

	errC := Cause(err)

	assert.Equal(t, errC.Error(), err1.Error())
}

func TestReport_Case1(t *testing.T) {
	err := NewErr(1715845918044, "something went wrong")
	err = WrapErr(1715845936107, err)
	err = WrapErr(1715845950562, err)
	err = WrapErr(1715845961777, err)

	rep := Report(err)

	assert.Equal(t, rep.Msg, "something went wrong")
	assert.Equal(t, rep.Traces, []int{1715845961777, 1715845950562, 1715845936107, 1715845918044})
}
