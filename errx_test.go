package errx

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrorWrappings(t *testing.T) {
	{
		err := newErr(1741599154, "something went wrong")
		err = wrapErr(1741600103, err)
		err = wrapErr(1741600368, err)
		assert.Equal(t, err.Error(), "[ts 1741600368]; [ts 1741600103]; [ts 1741599154] something went wrong")
	}

	{
		err := errors.New("something went wrong")
		err = wrapErr(1741600621, err)
		err = wrapErr(1741600633, err)
		assert.Equal(t, err.Error(), "[ts 1741600633]; [ts 1741600621]; something went wrong")
	}

	{
		err1 := newErr(1741601009, "something went wrong")
		err1 = wrapErr(1741601336, err1)
		err := fmt.Errorf("another generic error: %w", err1)
		err = wrapErr(1741601177, err)
		err = wrapErr(1741601190, err)
		assert.Equal(t, err.Error(), "[ts 1741601190]; [ts 1741601177]; another generic error: [ts 1741601336]; [ts 1741601009] something went wrong")
	}

	{
		err := newErr(1741601666, "something went wrong").WithKind(DataKind[int]("test_error")(30))
		err = wrapErr(1741601699, err)
		err = wrapErr(1741601711, err).WithKind(DataKind[string]("url_error")("www.test.com"))
		assert.Equal(t, err.Error(), `[ts 1741601711 kind url_error data "www.test.com"]; [ts 1741601699]; [ts 1741601666 kind test_error data 30] something went wrong`)
	}

	{
		err := errors.New("something went wrong")
		err = fmt.Errorf("something else broke: %w", err)
		err = wrapErr(1741602329, err)
		err = wrapErr(1741602338, err).WithKind(DataKind[int]("block_err")(40))
		err = wrapErr(1741602379, err)
		assert.Equal(t, err.Error(), "[ts 1741602379]; [ts 1741602338 kind block_err data 40]; [ts 1741602329]; something else broke: something went wrong")
	}

	{
		err := errx{msg: "something went wrong"}
		fmt.Println(err.Error())
		assert.Equal(t, err.Error(), "something went wrong")
	}
}

func TestErrorWrapFormating(t *testing.T) {
	{
		err := New(1745412853, "something went wrong")
		err = Wrapf(1745413114, "i said %v", err)
		assert.Equal(t, err.Error(), "[ts 1745413114]; [ts 1745412853] i said something went wrong")
	}

	{
		err := errors.New("something went wrong")
		err = Wrapf(1745413538, "i said %v", err)
		assert.Equal(t, err.Error(), "[ts 1745413538]; i said something went wrong")
	}
}

func TestErrorUnwrapping(t *testing.T) {
	err := New(1741636590, "something went wrong")
	err = Wrap(1741636604, err)
	err = Wrap(1741636616, err)
	fmt.Println(err)

	err1 := errors.Unwrap(err)
	assert.Equal(t, err1.Error(), "[ts 1741636604]; [ts 1741636590] something went wrong")

	err1 = Unwrap(err1)
	assert.Equal(t, err1.Error(), "[ts 1741636590] something went wrong")

	err1 = Unwrap(err1)
	assert.Equal(t, err1, nil)
}

func TestErrorIs(t *testing.T) {
	t.Run("Basic string matching", func(t *testing.T) {
		err0 := errors.New("something one")
		err := wrapErr(1741638522, err0)
		err = wrapErr(1741638536, err)
		assert.True(t, errors.Is(err, err0))
		assert.True(t, Is(err, err0))
	})

	t.Run("Exact structural match", func(t *testing.T) {
		err1 := newErr(1741638630, "something two")
		err := wrapErr(1741638643, err1)
		err = wrapErr(1741638650, err)
		assert.True(t, errors.Is(err, err1))
		assert.True(t, Is(err, err1))
	})

	t.Run("Stamps don't have to match. They are not part of error signature", func(t *testing.T) {
		err1 := newErr(100, "same message")
		err2 := newErr(200, "same message")
		assert.True(t, Is(err1, err2))
	})

	t.Run("Mismatching kinds", func(t *testing.T) {
		err1 := newErr(100, "same message").WithKind(Kind("kind_a"))
		err2 := newErr(100, "same message").WithKind(Kind("kind_b"))
		assert.False(t, Is(err1, err2))
	})

	t.Run("Mismatching string data", func(t *testing.T) {
		err1 := newErr(100, "msg").WithKind(DataKind[string]("k")("val1"))
		err2 := newErr(100, "msg").WithKind(DataKind[string]("k")("val2"))
		assert.False(t, Is(err1, err2))
	})

	t.Run("Matching slice data natively against parsed string", func(t *testing.T) {
		err1 := newErr(100, "msg").WithKind(DataKind[[]int]("k")([]int{1, 2, 3}))
		err2 := ParseStampedError(err1.Error())
		// Although err1 holds native []int and err2 holds a JSON string "[1,2,3]",
		// they should evaluate equivalently under Is() due to the .String() conversion matching
		assert.True(t, Is(err1, err2))
	})

	t.Run("Unrelated custom error type", func(t *testing.T) {
		type customErr struct{ msg string }
		// Need to declare a function for it but we can't attach methods to types defined inside functions.
		// So instead I'll just rely on a simple errors.New test, or use exactly what was requested.
		// To avoid scope issues, I'll use a local type that delegates, but there's a simpler way: just use a standard error interface mock or errors.New
		// To emulate an external custom error:
		err1 := fmt.Errorf("custom failure") // acts as a non-errx error

		err2 := newErr(100, "custom failure") // An errx error with a stamp
		err3 := newErr(0, "custom failure")   // An errx error with a zero stamp

		// errors.Is should find a match if the string messages are identical,
		// as fmt.Errorf does not have an Unwrap method that matches errx.
		// Wait, errors.Is checks Unwraps, which won't match strings unless we use Is

		// Our custom Is() should match based on string for non-errx types
		// when the errx error has a zero stamp, or if the target is not an *errx.
		assert.False(t, Is(err2, err1)) // 100 vs clean
		assert.True(t, Is(err3, err1))  // 0 clean string matches exactly
	})

	t.Run("Nil comparisons", func(t *testing.T) {
		var err error = nil
		assert.True(t, Is(err, nil))

		var errTyped *errx = nil
		assert.False(t, Is(errTyped, nil)) // typed nils inside error interfaces are technically not untyped nil

		err1 := newErr(100, "msg")
		assert.False(t, Is(err1, nil))
		assert.False(t, Is(nil, err1))
	})
}

func TestFindData(t *testing.T) {
	{
		kind := DataKind[int]("failure")
		err := newErr(1741653892, "something went wrong").WithKind(kind(10))
		data, ok := FindData(err, kind)

		assert.Equal(t, ok, true)
		assert.Equal(t, *data, 10)
	}

	{
		kind := DataKind[float64]("failure")
		err := newErr(1741655302, "something went wrong").WithKind(kind(1.453))
		data2, ok := FindData(err, kind)

		assert.Equal(t, ok, true)
		assert.Equal(t, *data2, 1.453)
	}

	{
		kind := DataKind[map[string]string]("failure")
		data := map[string]string{
			"name": "John Doe",
			"age":  "22",
		}
		err := newErr(1713712210034, "something went wrong").WithKind(kind(data))

		data3, ok := FindData(err, kind)
		assert.Equal(t, ok, true)
		assert.Equal(t, *data3, data)
	}

	{
		kind := DataKind[int]("xx")
		err := newErr(1713712707130, "something went wrong")
		data5, ok := FindData(err, kind)

		assert.False(t, ok)
		assert.Nil(t, data5)
	}
}

func TestFindDataFromParsedError(t *testing.T) {
	{
		kind := DataKind[int]("failure")
		err := newErr(1741655898, "something went wrong").WithKind(kind(10))
		msg := err.Error()

		pErr := ParseStampedError(msg)

		data, ok := FindData(pErr, kind)

		assert.Equal(t, ok, true)
		assert.Equal(t, *data, 10)
	}

	{
		kind := DataKind[float32]("failure")
		err := newErr(1713705985678, "something went wrong").WithKind(kind(1.453))
		msg := err.Error()

		pErr := ParseStampedError(msg)

		data, ok := FindData(pErr, kind)

		assert.Equal(t, ok, true)
		assert.Equal(t, *data, float32(1.453))
	}

	{
		kind := DataKind[map[string]string]("failure")
		data := map[string]string{
			"name": "John",
			"age":  "22",
		}
		err := newErr(1741656212, "something went wrong").WithKind(kind(data))
		msg := err.Error()

		pErr := ParseStampedError(msg)

		data3, ok := FindData(pErr, kind)
		assert.Equal(t, ok, true)
		assert.Equal(t, *data3, data)
	}

	{
		err := newErr(1713712738707, "something went wrong")
		msg := err.Error()

		perr := ParseStampedError(msg)

		data5, ok := FindData(perr, DataKind[int]("xxx"))

		assert.False(t, ok)
		assert.Nil(t, data5)
	}
}

func TestWrappedFindData(t *testing.T) {
	{
		kind := DataKind[int]("failure")
		err := newErr(1741656430, "something went wrong").WithKind(kind(30))
		err = wrapErr(1741656433, err).WithKind(DataKind[string]("another_failure")("https://www.google.com"))
		err1 := fmt.Errorf("another generic error: %w", err)
		err = wrapErr(1741656520, err1)

		dataInt, ok := FindData(err, kind)
		assert.True(t, ok)
		assert.Equal(t, *dataInt, 30)
	}

	{
		notFoundErr := DataKind[float32]("not-found-error")
		err := newErr(1713713974259, "something went wrong").WithKind(DataKind[float32]("failure_one")(1.56))
		err = wrapErr(1713714078710, err).WithKind(notFoundErr(1.90))
		err = wrapErr(1713714078710, err).WithKind(DataKind[float32]("failure_three")(1.49))
		err1 := fmt.Errorf("another generic error: %w", err)
		err = wrapErr(1713714083537, err1)

		data, ok := FindData(err, notFoundErr)
		assert.True(t, ok)
		assert.Equal(t, *data, float32(1.90))
	}
}

func TestWrappedFindDataFromParsedError(t *testing.T) {
	{
		genericErr := DataKind[int]("generic_error")
		urlErr := DataKind[string]("url_error")
		err := newErr(1741657070, "something went wrong").WithKind(genericErr(30))
		err = wrapErr(1741657072, err).WithKind(urlErr("https://www.google.com"))
		err1 := fmt.Errorf("another generic error %w", err)
		err = wrapErr(1713713041917, err1)

		perr := ParseStampedError(err.Error())

		dataInt, ok := FindData(perr, genericErr)
		assert.True(t, ok)
		assert.Equal(t, *dataInt, 30)

		dataStr, ok := FindData(perr, urlErr)
		assert.True(t, ok)
		assert.Equal(t, *dataStr, "https://www.google.com")
	}

	{
		urlErr := DataKind[string]("url_error")
		listErr := DataKind[[]int]("list_error")
		notFound := DataKind[float32]("not_found_error")
		genericErr := DataKind[int]("generic_error")
		err := errors.New("something went wrong")
		err = fmt.Errorf("second generic error %w", err)
		err = wrapErr(1741659127, err).WithKind(urlErr("https://www.google.com"))
		err = wrapErr(1741659129, err).WithKind(listErr([]int{1, 2, 3, 4, 7}))
		err1 := fmt.Errorf("third generic error %w", err)
		err = wrapErr(1741659230, err1)
		err = wrapErr(1741659238, err).WithKind(notFound(1.456))
		err = wrapErr(1741659340, err).WithKind(genericErr(30))
		err = fmt.Errorf("another generic error %w", err)

		perr := ParseStampedError(err.Error())

		dataStr, ok := FindData(perr, urlErr)
		assert.True(t, ok)
		assert.Equal(t, *dataStr, "https://www.google.com")

		dataflt, ok := FindData(perr, notFound)
		assert.True(t, ok)
		assert.Equal(t, *dataflt, float32(1.456))

		dataInt, ok := FindData(perr, genericErr)
		assert.True(t, ok)
		assert.Equal(t, *dataInt, 30)

		dataIntL, ok := FindData(perr, listErr)
		assert.True(t, ok)
		assert.Equal(t, *dataIntL, []int{1, 2, 3, 4, 7})
	}
}
func TestBuilders(t *testing.T) {
	t.Run("NewBuild", func(t *testing.T) {
		err := NewBuild(123, "test error").WithKind(Kind("test_kind"))
		assert.Equal(t, "[ts 123 kind test_kind] test error", err.Error())
		assert.Equal(t, 123, err.Stamp())
		assert.Equal(t, "test_kind", err.Kind())
	})

	t.Run("BuildFrom", func(t *testing.T) {
		baseErr := errors.New("base error")
		err := BuildFrom(456, baseErr).WithKind(Kind("wrap_kind"))
		assert.Equal(t, "[ts 456 kind wrap_kind]; base error", err.Error())
	})
}

func TestFormatting(t *testing.T) {
	t.Run("Newf", func(t *testing.T) {
		err := Newf(789, "error %d: %s", 1, "failed")
		assert.Equal(t, "[ts 789] error 1: failed", err.Error())
	})

	t.Run("Wrapf", func(t *testing.T) {
		baseErr := errors.New("original")
		err := Wrapf(101, "wrapped %s: %v", baseErr, "context")
		assert.Equal(t, "[ts 101]; wrapped original: context", err.Error())
	})

	t.Run("Wrapf with errx", func(t *testing.T) {
		baseErr := New(202, "inner")
		err := Wrapf(303, "outer %s: %v", baseErr, "msg")
		assert.Equal(t, "[ts 303]; [ts 202] outer inner: msg", err.Error())
	})
}

func TestKindFunctions(t *testing.T) {
	kind := Kind("auth")

	t.Run("NewKind", func(t *testing.T) {
		err := NewKind(404, kind, "not authorized")
		assert.Equal(t, "[ts 404 kind auth] not authorized", err.Error())
		assert.True(t, IsKind(err, kind))
	})

	t.Run("WrapKind", func(t *testing.T) {
		base := errors.New("db error")
		err := WrapKind(505, kind, base)
		assert.Equal(t, "[ts 505 kind auth]; db error", err.Error())
		assert.True(t, IsKind(err, kind))
	})

	t.Run("NewKindf", func(t *testing.T) {
		err := NewKindf(606, kind, "user %d failed", 123)
		assert.Equal(t, "[ts 606 kind auth] user 123 failed", err.Error())
	})

	t.Run("WrapKindf", func(t *testing.T) {
		base := errors.New("io error")
		err := WrapKindf(707, kind, "wrapped: %v", base)
		assert.Equal(t, "[ts 707 kind auth]; wrapped: io error", err.Error())
	})
}

func TestStamps(t *testing.T) {
	err := New(1, "e1")
	err = Wrap(2, err)
	err = Wrap(3, err)

	ex, ok := err.(*errx)
	assert.True(t, ok)
	stamps := ex.Stamps()
	assert.Equal(t, []int{3, 2, 1}, stamps)
}

func TestEdgeCases(t *testing.T) {
	t.Run("Wrap nil", func(t *testing.T) {
		err := Wrap(123, nil)
		assert.Equal(t, "[ts 123]", err.Error())
	})

	t.Run("Zero stamp", func(t *testing.T) {
		err := New(0, "nothing")
		assert.Equal(t, "nothing", err.Error())
	})

	t.Run("Empty message", func(t *testing.T) {
		err := New(123, "")
		assert.Equal(t, "[ts 123]", err.Error())
	})
}
func TestJoinWrap(t *testing.T) {
	err1 := errors.New("e1")
	err2 := errors.New("e2")
	joined := JoinWrap(123, err1, err2)
	assert.Contains(t, joined.Error(), "[ts 123]")
	assert.Contains(t, joined.Error(), "e1")
	assert.Contains(t, joined.Error(), "e2")
}

func TestUseLogger(t *testing.T) {
	UseLogger(nil)
}
