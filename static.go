package errx

import "encoding/json"

type unknown struct{}

type DataType interface {
	unknown |
		int | float32 | float64 | string |
		[]int | []float32 | []float64 | []string |
		map[string]int | map[string]float32 | map[string]float64 | map[string]string
}

// Unwraps the error and retrieves the data values and returns the first one that matches the specified error kind and given type
func FindData[T DataType](kind func(T) errKind, err error) (*T, bool) {
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
