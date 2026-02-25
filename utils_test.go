package errx

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToString_Int(t *testing.T) {
	oval := 10
	val := toStr(oval)

	ovalBytes, _ := json.Marshal(oval)
	assert.Equal(t, val, string(ovalBytes))

	cval, err := fromStr[int](val)
	assert.Nil(t, err)
	assert.Equal(t, *cval, oval)
}

func TestToString_float(t *testing.T) {
	oval := 1.3440
	val := toStr(oval)

	ovalBytes, _ := json.Marshal(oval)
	assert.Equal(t, val, string(ovalBytes))

	cval, err := fromStr[float64](val)
	assert.Nil(t, err)
	assert.Equal(t, *cval, oval)
}

func TestToString_floa64(t *testing.T) {
	var oval float64 = 1.34401
	val := toStr(oval)

	ovalBytes, _ := json.Marshal(oval)
	assert.Equal(t, val, string(ovalBytes))

	cval, err := fromStr[float64](val)
	assert.Nil(t, err)
	assert.Equal(t, *cval, oval)
}

func TestToString_floa32(t *testing.T) {
	var oval float32 = 1.34401
	val := toStr(oval)

	ovalBytes, _ := json.Marshal(oval)
	assert.Equal(t, val, string(ovalBytes))

	cval, err := fromStr[float32](val)
	assert.Nil(t, err)
	assert.Equal(t, *cval, oval)
}

func TestToString_String(t *testing.T) {
	var oval string = "https://www.google.com/golang?wehere='at\"intelligence\"'"
	val := toStr(oval)

	ovalBytes, _ := json.Marshal(oval)
	assert.Equal(t, val, string(ovalBytes))

	cval, err := fromStr[string](val)
	assert.Nil(t, err)
	assert.Equal(t, *cval, oval)
}

func TestToString_IntArr(t *testing.T) {
	oval := []int{1, 2, 4, 5, 10, 20}
	val := toStr(oval)

	ovalBytes, _ := json.Marshal(oval)
	assert.Equal(t, val, string(ovalBytes))

	cval, err := fromStr[[]int](val)
	assert.Nil(t, err)
	assert.Equal(t, *cval, oval)
}

func TestToString_Float32Arr(t *testing.T) {
	oval := []float32{1.4500, 2.211, 4.0193, 5.101, 10, 20.0001}
	val := toStr(oval)

	ovalBytes, _ := json.Marshal(oval)
	assert.Equal(t, val, string(ovalBytes))

	cval, err := fromStr[[]float32](val)
	assert.Nil(t, err)
	assert.Equal(t, *cval, oval)
}

func TestToString_Float64Arr(t *testing.T) {
	oval := []float64{1.4500, 2.211, 4.0193, 5.101, 10, 20.0001}
	val := toStr(oval)

	ovalBytes, _ := json.Marshal(oval)
	assert.Equal(t, val, string(ovalBytes))

	cval, err := fromStr[[]float64](val)
	assert.Nil(t, err)
	assert.Equal(t, *cval, oval)
}

func TestToString_StringArr(t *testing.T) {
	oval := []string{"one", "two", "three", "four", "five", "six"}
	val := toStr(oval)

	ovalBytes, _ := json.Marshal(oval)
	assert.Equal(t, val, string(ovalBytes))

	cval, err := fromStr[[]string](val)
	assert.Nil(t, err)
	assert.Equal(t, *cval, oval)
}

func TestToString_MapInt(t *testing.T) {
	oval := map[string]int{"one": 1, "two": 2, "three": 3, "four": 4, "five": 5}
	val := toStr(oval)

	ovalBytes, _ := json.Marshal(oval)
	assert.Equal(t, val, string(ovalBytes))

	cval, err := fromStr[map[string]int](val)
	assert.Nil(t, err)
	assert.Equal(t, *cval, oval)
}
