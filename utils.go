package errx

import (
	"fmt"
	"strings"
)

func arrPop[T any](alist *[]T) T {
	f := len(*alist)
	rv := (*alist)[f-1]
	*alist = (*alist)[:f-1]
	return rv
}

func toStr(data any) string {
	convs := strings.Split(fmt.Sprintf("%T <<<]@@@^&!@@@[>>> %#v", data, data), " <<<]@@@^&!@@@[>>> ")

	switch convs[0] {
	case "[]int":
		return "[" + convs[1][6:len(convs[1])-1] + "]"
	case "[]float32", "[]float64":
		return "[" + convs[1][10:len(convs[1])-1] + "]"
	case "[]string":
		return "[" + convs[1][9:len(convs[1])-1] + "]"
	case "map[string]int":
		return convs[1][14:]
	case "map[string]string":
		return convs[1][17:]
	case "map[string]float32", "map[string]float64":
		return convs[1][18:]
	default:
		return convs[1]
	}
}
