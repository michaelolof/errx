package errx

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrData_Case1(t *testing.T) {
	{
		err := NewData(1713705663160, "something went wrong", 10)
		data, ok := FindData[int](err)

		assert.Equal(t, ok, true, "1713567431691")
		assert.Equal(t, *data, 10, "1713567424980")
	}

	{
		err := NewData(1713705656185, "something went wrong", 1.453)
		data2, ok := FindData[float64](err)

		assert.Equal(t, ok, true, "1713705739174")
		assert.Equal(t, *data2, 1.453, "1713705744517")
	}

	{
		type obj struct {
			Name string
			Age  int
		}
		err := NewData(1713712210034, "something went wrong", obj{Name: "John Doe", Age: 22})

		data3, ok := FindData[obj](err)
		assert.Equal(t, ok, true, "1713712218011")
		assert.Equal(t, *data3, obj{Name: "John Doe", Age: 22}, "1713712224977")
	}

	{
		err := New(1713712707130, "something went wrong")
		data5, ok := FindData[int](err)

		assert.False(t, ok)
		assert.Nil(t, data5)
	}
}

func TestErrData_Case2(t *testing.T) {
	{
		err := NewData(1713704575804, "something went wrong", 10)
		msg := err.Error()

		pErr, err2 := ParseStampedError(msg)
		assert.Equal(t, err2, nil, 1713704761772)

		data, ok := FindData[float64](pErr)

		assert.Equal(t, ok, true, "1713567431691")
		assert.Equal(t, int(*data), 10, "1713567424980")
	}

	{
		err := NewData(1713705985678, "something went wrong", 1.453)
		msg := err.Error()

		pErr, err2 := ParseStampedError(msg)
		assert.Equal(t, err2, nil, 1713705991420)

		data, ok := FindData[float64](pErr)

		assert.Equal(t, ok, true, "1713567431691")
		assert.Equal(t, *data, 1.453, "1713567424980")
	}

	{
		type obj struct {
			Name string
			Age  string
		}
		err := NewData(1713562521351, "something went wrong", obj{Name: "John Doe", Age: "22"})
		msg := err.Error()

		pErr, err2 := ParseStampedError(msg)
		assert.Equal(t, err2, nil, 1713706300491)

		data3, ok := FindData[obj](pErr)
		assert.Equal(t, ok, true, "1713567384853")
		assert.Equal(t, *data3, obj{Name: "John Doe", Age: "22"}, "1713706354978")
	}

	{
		err := New(1713712738707, "something went wrong")
		msg := err.Error()

		perr, err1 := ParseStampedError(msg)
		assert.Nil(t, err1)

		data5, ok := FindData[int](perr)

		assert.False(t, ok)
		assert.Nil(t, data5)
	}
}

func TestErrData_Case3(t *testing.T) {
	{
		err := NewData(1713712915111, "something went wrong", 30)
		err = WrapData(1713712979885, err, "https://www.google.com")
		err = fmt.Errorf("another generic error: %w", err)
		err = Wrap(1713713041917, err)

		dataInt, ok := FindData[int](err)
		assert.True(t, ok)
		assert.Equal(t, *dataInt, 30)
	}

	{
		err := NewData(1713713974259, "something went wrong", 1.56)
		err = WrapData(1713714078710, err, 1.49)
		err = fmt.Errorf("another generic error: %w", err)
		err = Wrap(1713714083537, err)

		data, ok := FindData[float64](err)
		assert.True(t, ok)
		assert.Equal(t, *data, 1.49)
	}

	{
		notFoundErr := errors.New("not-found-error")
		err := NewData(1713713974259, "something went wrong", 1.56)
		err = WrapErr(1713714078710, err, notFoundErr, 1.90)
		err = WrapData(1713714078710, err, 1.49)
		err = fmt.Errorf("another generic error: %w", err)
		err = Wrap(1713714083537, err)

		data, ok := FindDataOfKind[float64](err, notFoundErr)
		assert.True(t, ok)
		assert.Equal(t, *data, 1.90)
	}
}

func TestErrData_Case4(t *testing.T) {
	{
		err := NewData(1713713979066, "something went wrong", 30)
		err = WrapData(1713712979885, err, "https://www.google.com")
		err = fmt.Errorf("another generic error: %w", err)
		err = Wrap(1713713041917, err)
		msg := err.Error()

		perr, err := ParseStampedError(msg)
		assert.Nil(t, err)

		dataInt, ok := FindData[float64](perr)
		assert.True(t, ok)
		assert.Equal(t, int(*dataInt), 30)

		dataStr, ok := FindData[string](perr)
		assert.True(t, ok)
		assert.Equal(t, *dataStr, "https://www.google.com")
	}
}

func TestUnwrap(t *testing.T) {

	type assertion struct {
		kindStr string
		stamp   int
		data    any
	}

	notFound := errors.New("not-found-error")

	{
		err := NewErr[any](1713691153440, "something went wrong", notFound, nil)
		err = WrapData(1713691164311, err, 70)
		err = Wrap(1713691172952, err)
		err = Wrap(1713691181909, err)

		assertions := []assertion{
			{stamp: 1713691181909},
			{stamp: 1713691172952},
			{stamp: 1713691164311, data: 70},
			{stamp: 1713691153440, kindStr: notFound.Error()},
		}

		count := 0
		for err != nil {
			if v, ok := err.(Stamper); ok {
				ast := assertions[count]
				assert.Equal(t, v.Stamp(), ast.stamp, 1713701175536)
				assert.Equal(t, v.KindStr(), ast.kindStr, 1713701191837)
				assert.Equal(t, v.Data(), ast.data, 1713701188390)
			}
			err = Unwrap(err)
			count = count + 1
		}
	}

	{
		err := errors.New("something went wrong")
		err = fmt.Errorf("error context two: %w", err)
		err = WrapErr(1713721643327, err, notFound, true)
		err = Wrap(1713721709036, err)
		err = fmt.Errorf("error context three: %w", err)

		assertions := []assertion{
			{stamp: 1713721709036},
			{stamp: 1713721643327, kindStr: notFound.Error(), data: true},
		}

		count := 0
		for err != nil {
			if v, ok := err.(Stamper); ok {
				ast := assertions[count]
				assert.Equal(t, v.Stamp(), ast.stamp)
				assert.Equal(t, v.KindStr(), ast.kindStr)
				assert.Equal(t, v.Data(), ast.data)

				count = count + 1
			}
			err = Unwrap(err)
		}
	}
}
