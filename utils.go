package errx

import (
	"encoding/json"
	"fmt"
)

func toStr(data any) string {
	bs, err := json.Marshal(data)
	if err != nil {
		return fmt.Sprintf("%v", data)
	}
	return string(bs)
}
