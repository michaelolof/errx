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
func FindData[T DataType](kind ErrKind, err error) (*T, bool) {
	for err != nil {
		switch t := err.(type) {
		case *errx:
			if kind == t.kind && t.data.isSet && t.data.val != nil {
				if v, ok := t.data.val.(T); ok {
					return &v, true
				} else {
					return nil, false
				}
			} else if kind == t.kind && t.data.isSet && t.data.valStr != "" {
				var d T
				if err := json.Unmarshal([]byte(t.data.valStr), &d); err == nil {
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
