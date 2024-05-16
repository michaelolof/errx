package errx

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseStamedError(t *testing.T) {

	notfound := errors.New("not found error")

	err := NewData(1713591173899, "something went wrong", 30)
	err = Wrap(1713592763837, err)
	err = WrapData(1713592780139, err, 30)
	err = WrapErr(1713591205370, err, notfound, "something i want to do often")

	itMatches := Is(err, notfound)
	assert.Equal(t, itMatches, true, 1713680824055)

	msg := err.Error()
	parsed := ParseStampedError(msg)

	itMatches = Is(parsed, notfound)
	assert.Equal(t, itMatches, true, 1713680824055)

	newMsg := parsed.Error()
	assert.Equal(t, newMsg, msg, 1713680554000)
}

func TestParseStamedError_Case2(t *testing.T) {
	err := errors.New("some thing went wrong")
	err = fmt.Errorf("i am a suspect error: %w", err)
	err = WrapData(1713607000211, err, 30)
	err = WrapData(1713607005378, err, 1.560)
	err = WrapData(1713607010089, err, "https://www.google.com")

	msg := err.Error()
	parsed := ParseStampedError(msg)

	newMsg := parsed.Error()
	assert.Equal(t, newMsg, msg, 1713685783298)
}

func TestParseStamedError_Case3(t *testing.T) {
	err := NewData(1713606995137, "something went wrong", 30)
	err = fmt.Errorf("i am a suspect error: %w", err)
	err = WrapData(1713607005378, err, 1.560)
	err = WrapData(1713607010089, err, "https://www.google.com")

	rep := Report(err)
	fmt.Println(rep)

	msg := err.Error()
	parsed := ParseStampedError(msg)

	newMsg := parsed.Error()
	assert.Equal(t, newMsg, msg, 1713685783298)
}

func TestGetStackFrames_Case1(t *testing.T) {

	err := errors.New("some thing went wrong")
	err = fmt.Errorf("i am a suspect error: %w", err)
	err = WrapData(1713607000211, err, 30)
	err = WrapData(1713607005378, err, 1.560)
	err = WrapData(1713607010089, err, "https://www.google.com")

	frames := GetStackFrames(err)
	assert.Len(t, frames, 4)
}

func TestGetStackFrames_Case2(t *testing.T) {

	err := NewData(1713606995137, "something went wrong", 30)
	err = WrapData(1713607000211, err, []float32{1.5, 2.5})
	err = fmt.Errorf("i am a suspect error: %w", err)
	err = WrapData(1713607005378, err, 1.560)
	err = WrapData(1713607010089, err, "https://www.google.com")

	frames := GetStackFrames(err)
	assert.Len(t, frames, 5)

}

func TestCause_Case1(t *testing.T) {

	err1 := New(1715845918044, "something went wrong")
	err := Wrap(1715845936107, err1)
	err = Wrap(1715845950562, err)
	err = Wrap(1715845961777, err)

	err = Cause(err)

	assert.Equal(t, err.Error(), err1.Error())
}

func TestReport_Case1(t *testing.T) {
	err := New(1715845918044, "something went wrong")
	err = Wrap(1715845936107, err)
	err = Wrap(1715845950562, err)
	err = Wrap(1715845961777, err)

	rep := Report(err)

	assert.Equal(t, rep.Msg, "something went wrong")
	assert.Equal(t, rep.Traces, []int{1715845961777, 1715845950562, 1715845936107, 1715845918044})
}
