package errx

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrorWrappings(t *testing.T) {
	{
		err := NewErr(1741599154, "something went wrong")
		err = WrapErr(1741600103, err)
		err = WrapErr(1741600368, err)
		assert.Equal(t, err.Error(), "[ts 1741600368]; [ts 1741600103]; [ts 1741599154] something went wrong")
	}

	{
		err := errors.New("something went wrong")
		err = WrapErr(1741600621, err)
		err = WrapErr(1741600633, err)
		assert.Equal(t, err.Error(), "[ts 1741600633]; [ts 1741600621]; something went wrong")
	}

	{
		err1 := NewErr(1741601009, "something went wrong")
		err1 = WrapErr(1741601336, err1)
		err := fmt.Errorf("another generic error: %w", err1)
		err = WrapErr(1741601177, err)
		err = WrapErr(1741601190, err)
		assert.Equal(t, err.Error(), "[ts 1741601190]; [ts 1741601177]; another generic error: [ts 1741601336]; [ts 1741601009] something went wrong")
	}

	{
		err := NewErr(1741601666, "something went wrong").WithData(ErrKind("test_error"), Data(30))
		err = WrapErr(1741601699, err)
		err = WrapErr(1741601711, err).WithData(ErrKind("url_error"), Data("www.test.com"))
		assert.Equal(t, err.Error(), `[ts 1741601711 kind url_error data "www.test.com"]; [ts 1741601699]; [ts 1741601666 kind test_error data 30] something went wrong`)
	}

	{
		err := errors.New("something went wrong")
		err = fmt.Errorf("something else broke: %w", err)
		err = WrapErr(1741602329, err)
		err = WrapErr(1741602338, err).WithData(ErrKind("block_err"), Data(40))
		err = WrapErr(1741602379, err)
		assert.Equal(t, err.Error(), "[ts 1741602379]; [ts 1741602338 kind block_err data 40]; [ts 1741602329]; something else broke: something went wrong")
	}

	{
		err := errx{msg: "something went wrong"}
		fmt.Println(err.Error())
		assert.Equal(t, err.Error(), "something went wrong")
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
	err0 := errors.New("something one")
	err := WrapErr(1741638522, err0)
	err = WrapErr(1741638536, err)
	match := errors.Is(err, err0)
	assert.True(t, match)

	err1 := NewErr(1741638630, "something two")
	err = WrapErr(1741638643, err1)
	err = WrapErr(1741638650, err)
	match = errors.Is(err, err1)
	assert.True(t, match)

	err2 := NewErr(1741638827, "something two")
	err = WrapErr(1741638830, err2)
	err = WrapErr(1741638832, err)
	match = Is(err, err2)
	fmt.Println(match)
	assert.True(t, match)
}

func TestFindData(t *testing.T) {
	{
		kind := ErrKind("failure")
		err := NewErr(1741653892, "something went wrong").WithData(kind, Data(10))
		data, ok := FindData[int](kind, err)

		assert.Equal(t, ok, true)
		assert.Equal(t, *data, 10)
	}

	{
		kind := ErrKind("failure")
		err := NewErr(1741655302, "something went wrong").WithData(kind, Data(1.453))
		data2, ok := FindData[float64](kind, err)

		assert.Equal(t, ok, true)
		assert.Equal(t, *data2, 1.453)
	}

	{
		kind := ErrKind("failure")
		data := map[string]string{
			"name": "John Doe",
			"age":  "22",
		}
		err := NewErr(1713712210034, "something went wrong").WithData(kind, Data(data))

		data3, ok := FindData[map[string]string](kind, err)
		assert.Equal(t, ok, true)
		assert.Equal(t, *data3, data)
	}

	{
		err := NewErr(1713712707130, "something went wrong")
		data5, ok := FindData[int](ErrKind("xx"), err)

		assert.False(t, ok)
		assert.Nil(t, data5)
	}
}

func TestFindDataFromParsedError(t *testing.T) {
	{
		kind := ErrKind("failure")
		err := NewErr(1741655898, "something went wrong").WithData(kind, Data(10))
		msg := err.Error()

		pErr := ParseStampedError(msg)

		data, ok := FindData[int](kind, pErr)

		assert.Equal(t, ok, true)
		assert.Equal(t, *data, 10)
	}

	{
		kind := ErrKind("failure")
		err := NewErr(1713705985678, "something went wrong").WithData(kind, Data(1.453))
		msg := err.Error()

		pErr := ParseStampedError(msg)

		data, ok := FindData[float32](kind, pErr)

		assert.Equal(t, ok, true)
		assert.Equal(t, *data, float32(1.453))
	}

	{
		kind := ErrKind("failure")
		data := map[string]string{
			"name": "John",
			"age":  "22",
		}
		err := NewErr(1741656212, "something went wrong").WithData(kind, Data(data))
		msg := err.Error()

		pErr := ParseStampedError(msg)

		data3, ok := FindData[map[string]string](kind, pErr)
		assert.Equal(t, ok, true)
		assert.Equal(t, *data3, data)
	}

	{
		err := NewErr(1713712738707, "something went wrong")
		msg := err.Error()

		perr := ParseStampedError(msg)

		data5, ok := FindData[int](ErrKind("xxx"), perr)

		assert.False(t, ok)
		assert.Nil(t, data5)
	}
}

func TestWrappedFindData(t *testing.T) {
	{
		kind := ErrKind("failure")
		err := NewErr(1741656430, "something went wrong").WithData(kind, Data(30))
		err = WrapErr(1741656433, err).WithData(ErrKind("another_failure"), Data("https://www.google.com"))
		err1 := fmt.Errorf("another generic error: %w", err)
		err = WrapErr(1741656520, err1)

		dataInt, ok := FindData[int](kind, err)
		assert.True(t, ok)
		assert.Equal(t, *dataInt, 30)
	}

	{
		notFoundErr := ErrKind("not-found-error")
		err := NewErr(1713713974259, "something went wrong").WithData(ErrKind("failure_one"), Data(1.56))
		err = WrapErr(1713714078710, err).WithData(notFoundErr, Data(1.90))
		err = WrapErr(1713714078710, err).WithData(ErrKind("failure_three"), Data(1.49))
		err1 := fmt.Errorf("another generic error: %w", err)
		err = WrapErr(1713714083537, err1)

		data, ok := FindData[float64](notFoundErr, err)
		assert.True(t, ok)
		assert.Equal(t, *data, 1.90)
	}
}

func TestWrappedFindDataFromParsedError(t *testing.T) {
	{
		genericErr := ErrKind("generic_error")
		urlErr := ErrKind("url_error")
		err := NewErr(1741657070, "something went wrong").WithData(genericErr, Data(30))
		err = WrapErr(1741657072, err).WithData(urlErr, Data("https://www.google.com"))
		err1 := fmt.Errorf("another generic error %w", err)
		err = WrapErr(1713713041917, err1)

		perr := ParseStampedError(err.Error())

		dataInt, ok := FindData[float32](genericErr, perr)
		assert.True(t, ok)
		assert.Equal(t, int(*dataInt), 30)

		dataStr, ok := FindData[string](urlErr, perr)
		assert.True(t, ok)
		assert.Equal(t, *dataStr, "https://www.google.com")
	}

	{
		urlErr := ErrKind("url_error")
		listErr := ErrKind("list_error")
		notFound := ErrKind("not_found_error")
		genericErr := ErrKind("generic_error")
		err := errors.New("something went wrong")
		err = fmt.Errorf("second generic error %w", err)
		err = WrapErr(1741659127, err).WithData(urlErr, Data("https://www.google.com"))
		err = WrapErr(1741659129, err).WithData(listErr, Data([]int{1, 2, 3, 4, 7}))
		err1 := fmt.Errorf("third generic error %w", err)
		err = WrapErr(1741659230, err1)
		err = WrapErr(1741659238, err).WithData(notFound, Data(1.456))
		err = WrapErr(1741659340, err).WithData(genericErr, Data(30))
		err = fmt.Errorf("another generic error %w", err)

		perr := ParseStampedError(err.Error())

		dataStr, ok := FindData[string](urlErr, perr)
		assert.True(t, ok)
		assert.Equal(t, *dataStr, "https://www.google.com")

		dataflt, ok := FindData[float32](notFound, perr)
		assert.True(t, ok)
		assert.Equal(t, *dataflt, float32(1.456))

		dataInt, ok := FindData[int](genericErr, perr)
		assert.True(t, ok)
		assert.Equal(t, *dataInt, 30)

		dataIntL, ok := FindData[[]int](listErr, perr)
		assert.True(t, ok)
		assert.Equal(t, *dataIntL, []int{1, 2, 3, 4, 7})
	}
}
