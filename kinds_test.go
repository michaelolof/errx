package errx

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	NotFound  = Kind("notfound")
	FileOpen  = DataKind[string]("fileopen")
	PageLoad  = DataKind[string]("pageload")
	UserLogin = DataKind[UserInfo]("userlogin")
)

type UserInfo map[string]string

func (u UserInfo) Id() int {
	v, err := strconv.Atoi(u["id"])
	if err != nil {
		panic(err)
	}

	return v
}

func (u UserInfo) Email() string {
	return u["email"]
}

func TestDataKindsExtensive(t *testing.T) {
	t.Run("Int Slice", func(t *testing.T) {
		kind := DataKind[[]int]("ids")
		data := []int{1, 2, 3}
		err := NewKind(123, kind(data), "msg")
		res, ok := FindData(err, kind)
		assert.True(t, ok)
		assert.Equal(t, data, *res)
	})

	t.Run("String Map", func(t *testing.T) {
		kind := DataKind[map[string]string]("meta")
		data := map[string]string{"key": "value"}
		err := NewKind(456, kind(data), "msg")
		res, ok := FindData(err, kind)
		assert.True(t, ok)
		assert.Equal(t, data, *res)
	})

	t.Run("Float64", func(t *testing.T) {
		kind := DataKind[float64]("score")
		data := 0.95
		err := NewKind(789, kind(data), "msg")
		res, ok := FindData(err, kind)
		assert.True(t, ok)
		assert.Equal(t, data, *res)
	})

	t.Run("Missing Data", func(t *testing.T) {
		kind := DataKind[int]("missing")
		err := New(111, "no kind here")
		_, ok := FindData(err, kind)
		assert.False(t, ok)
	})

	t.Run("Wrong Kind same type", func(t *testing.T) {
		k1 := DataKind[int]("k1")
		k2 := DataKind[int]("k2")
		err := NewKind(222, k1(10), "msg")
		_, ok := FindData(err, k2)
		assert.False(t, ok)
	})
}

func TestIsDataKind(t *testing.T) {
	k1 := DataKind[int]("k1")
	err := NewKind(123, k1(10), "msg")
	err = Wrap(456, err)

	assert.True(t, IsDataKind(err, k1))
	assert.False(t, IsDataKind(err, DataKind[int]("other")))
}

func TestOtherTypes(t *testing.T) {
	t.Run("Other Types", func(t *testing.T) {
		k1 := DataKind[float32]("f32")
		k2 := DataKind[[]string]("ss")
		k3 := DataKind[map[string]int]("msi")

		err := NewKind(1, k1(1.2), "msg")
		res, _ := FindData(err, k1)
		assert.Equal(t, float32(1.2), *res)

		err = NewKind(2, k2([]string{"a", "b"}), "msg")
		res2, _ := FindData(err, k2)
		assert.Equal(t, []string{"a", "b"}, *res2)

		err = NewKind(3, k3(map[string]int{"a": 1}), "msg")
		res3, _ := FindData(err, k3)
		assert.Equal(t, map[string]int{"a": 1}, *res3)
	})
}
