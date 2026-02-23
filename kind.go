package errx

import (
	"encoding/json"
)

// Define a basic error kind
func Kind(k string) errKind {
	return errKind{
		kind: k,
		data: dataValue{isSet: false},
	}
}

// Define an error kind with acceptable data types
func DataKind[T DataType](k string) func(d T) errKind {
	return func(d T) errKind {
		return errKind{
			kind: k,
			data: dataValue{isSet: true, val: d},
		}
	}
}

type dataValue struct {
	isSet  bool
	val    any
	valStr string
}

func (d *dataValue) String() string {
	if d.isSet && d.valStr != "" {
		return d.valStr
	} else if d.isSet && d.val != nil {
		return toStr(d.val)
	} else {
		return ""
	}
}

type unknown struct{}

type DataType interface {
	unknown |
		~int | ~float32 | ~float64 | ~string |
		~[]int | ~[]float32 | ~[]float64 | ~[]string |
		~map[string]int | ~map[string]float32 | ~map[string]float64 | ~map[string]string |
		~map[int]int | ~map[int]float32 | ~map[int]float64 | ~map[int]string
}

func IsKind(err error, kind errKind) bool {
	for err != nil {
		if e, ok := err.(interface{ Kind() string }); ok {
			if e.Kind() == kind.kind {
				return true
			}
		}
		err = Unwrap(err)
	}
	return false
}

func IsDataKind[T DataType](err error, kind func(d T) errKind) bool {
	var d T
	for err != nil {
		if e, ok := err.(interface{ Kind() string }); ok {
			if e.Kind() == kind(d).kind {
				return true
			}
		}
		err = Unwrap(err)
	}
	return false
}

// Unwraps the error and retrieves the data values and returns the first one that matches the specified error kind and given type
func FindData[T DataType](err error, kind func(T) errKind) (*T, bool) {
	var dv T
	k := kind(dv)
	for err != nil {
		switch t := err.(type) {
		case *errx:
			if k.kind == t.kind.kind && t.kind.data.isSet && t.kind.data.val != nil {
				if v, ok := t.kind.data.val.(T); ok {
					return &v, true
				} else {
					return nil, false
				}
			} else if k.kind == t.kind.kind && t.kind.data.isSet && t.kind.data.valStr != "" {
				var d T
				if err := json.Unmarshal([]byte(t.kind.data.valStr), &d); err == nil {
					return &d, true
				} else {
					return nil, false
				}
			}
		}

		err = Unwrap(err)
	}

	return nil, false
}
