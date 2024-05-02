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
		data := map[string]string{
			"name": "John Doe",
			"age":  "22",
		}
		err := NewData(1713712210034, "something went wrong", data)

		data3, ok := FindData[map[string]string](err)
		assert.Equal(t, ok, true, "1713712218011")
		assert.Equal(t, *data3, data, "1713712224977")
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

		pErr := ParseStampedError(msg)

		data, ok := FindData[int](pErr)

		assert.Equal(t, ok, true, "1713567431691")
		assert.Equal(t, *data, 10, "1713567424980")
	}

	{
		err := NewData(1713705985678, "something went wrong", 1.453)
		msg := err.Error()

		pErr := ParseStampedError(msg)

		data, ok := FindData[float32](pErr)

		assert.Equal(t, ok, true, "1713567431691")
		assert.Equal(t, *data, float32(1.453), "1713567424980")
	}

	{
		data := map[string]string{
			"name": "John",
			"age":  "22",
		}
		err := NewData(1713562521351, "something went wrong", data)
		msg := err.Error()

		pErr := ParseStampedError(msg)

		data3, ok := FindData[map[string]string](pErr)
		assert.Equal(t, ok, true, "1713567384853")
		assert.Equal(t, *data3, data, "1713706354978")
	}

	{
		err := New(1713712738707, "something went wrong")
		msg := err.Error()

		perr := ParseStampedError(msg)

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

		data, ok := FindDataByKind[float64](err, notFoundErr)
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

		perr := ParseStampedError(msg)

		dataInt, ok := FindData[float64](perr)
		assert.True(t, ok)
		assert.Equal(t, int(*dataInt), 30)

		dataStr, ok := FindData[string](perr)
		assert.True(t, ok)
		assert.Equal(t, *dataStr, "https://www.google.com")
	}

	{
		err := errors.New("something went wrong")
		err = fmt.Errorf("second generic error %w", err)
		err = WrapData(1714629786341, err, "https://www.google.com")
		err = WrapErr(1714629791533, err, errors.New("another banger error"), []int{1, 2, 3, 4, 7})
		err = fmt.Errorf("third generic error %w", err)
		err = Wrap(1714629796430, err)
		err = WrapErr(1714629800940, err, errors.New("not found error"), 1.456)
		err = WrapData(1714629805408, err, 30)
		err = fmt.Errorf("another generic error %w", err)
		msg := err.Error()

		perr := ParseStampedError(msg)

		dataStr, ok := FindData[string](perr)
		assert.True(t, ok)
		assert.Equal(t, *dataStr, "https://www.google.com")

		dataflt, ok := FindDataByKind[float32](perr, errors.New("not found error"))
		assert.True(t, ok)
		assert.Equal(t, *dataflt, float32(1.456))

		dataInt, ok := FindData[int](perr)
		assert.True(t, ok)
		assert.Equal(t, *dataInt, 30)

		dataIntL, ok := FindData[[]int](perr)
		assert.True(t, ok)
		assert.Equal(t, *dataIntL, []int{1, 2, 3, 4, 7})
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
		err := newE[unknown](1713691153440, "something went wrong", notFound, nil)
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
			if v, ok := err.(AnyStamper); ok {
				ast := assertions[count]
				assert.Equal(t, v.Stamp(), ast.stamp, 1713701175536)
				assert.Equal(t, v.KindStr(), ast.kindStr, 1713701191837)
				assert.Equal(t, v.AnyData(), ast.data, 1713701188390)
			}
			err = Unwrap(err)
			count = count + 1
		}
	}

	{
		err := errors.New("something went wrong")
		err = fmt.Errorf("error context two: %w", err)
		err = WrapErr(1713721643327, err, notFound, 30)
		err = Wrap(1713721709036, err)
		err = fmt.Errorf("error context three: %w", err)

		assertions := []assertion{
			{stamp: 1713721709036},
			{stamp: 1713721643327, kindStr: notFound.Error(), data: 30},
		}

		count := 0
		for err != nil {
			if v, ok := err.(AnyStamper); ok {
				ast := assertions[count]
				assert.Equal(t, v.Stamp(), ast.stamp)
				assert.Equal(t, v.KindStr(), ast.kindStr)
				assert.Equal(t, v.AnyData(), ast.data)

				count = count + 1
			}
			err = Unwrap(err)
		}
	}
}

func TestXxx(t *testing.T) {

	one := New(1714557123808, "something went wrong")
	two := Wrap(1714557274569, one)
	fmt.Println(two)

	fmt.Println("done")
}
