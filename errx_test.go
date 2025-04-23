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
	err0 := errors.New("something one")
	err := wrapErr(1741638522, err0)
	err = wrapErr(1741638536, err)
	match := errors.Is(err, err0)
	assert.True(t, match)

	err1 := newErr(1741638630, "something two")
	err = wrapErr(1741638643, err1)
	err = wrapErr(1741638650, err)
	match = errors.Is(err, err1)
	assert.True(t, match)

	err2 := newErr(1741638827, "something two")
	err = wrapErr(1741638830, err2)
	err = wrapErr(1741638832, err)
	match = Is(err, err2)
	fmt.Println(match)
	assert.True(t, match)
}

func TestFindData(t *testing.T) {
	{
		kind := DataKind[int]("failure")
		err := newErr(1741653892, "something went wrong").WithKind(kind(10))
		data, ok := FindData(kind, err)

		assert.Equal(t, ok, true)
		assert.Equal(t, *data, 10)
	}

	{
		kind := DataKind[float64]("failure")
		err := newErr(1741655302, "something went wrong").WithKind(kind(1.453))
		data2, ok := FindData(kind, err)

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

		data3, ok := FindData(kind, err)
		assert.Equal(t, ok, true)
		assert.Equal(t, *data3, data)
	}

	{
		kind := DataKind[int]("xx")
		err := newErr(1713712707130, "something went wrong")
		data5, ok := FindData(kind, err)

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

		data, ok := FindData(kind, pErr)

		assert.Equal(t, ok, true)
		assert.Equal(t, *data, 10)
	}

	{
		kind := DataKind[float32]("failure")
		err := newErr(1713705985678, "something went wrong").WithKind(kind(1.453))
		msg := err.Error()

		pErr := ParseStampedError(msg)

		data, ok := FindData(kind, pErr)

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

		data3, ok := FindData(kind, pErr)
		assert.Equal(t, ok, true)
		assert.Equal(t, *data3, data)
	}

	{
		err := newErr(1713712738707, "something went wrong")
		msg := err.Error()

		perr := ParseStampedError(msg)

		data5, ok := FindData(DataKind[int]("xxx"), perr)

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

		dataInt, ok := FindData[int](kind, err)
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

		data, ok := FindData(notFoundErr, err)
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

		dataInt, ok := FindData(genericErr, perr)
		assert.True(t, ok)
		assert.Equal(t, *dataInt, 30)

		dataStr, ok := FindData(urlErr, perr)
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

		dataStr, ok := FindData(urlErr, perr)
		assert.True(t, ok)
		assert.Equal(t, *dataStr, "https://www.google.com")

		dataflt, ok := FindData(notFound, perr)
		assert.True(t, ok)
		assert.Equal(t, *dataflt, float32(1.456))

		dataInt, ok := FindData(genericErr, perr)
		assert.True(t, ok)
		assert.Equal(t, *dataInt, 30)

		dataIntL, ok := FindData(listErr, perr)
		assert.True(t, ok)
		assert.Equal(t, *dataIntL, []int{1, 2, 3, 4, 7})
	}
}
