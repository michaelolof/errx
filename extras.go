package errx

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func Must[T any](obj T, err error) T {
	if err != nil {
		panic(err)
	}
	return obj
}

func Split(err error) []error {
	if uw, ok := err.(interface{ Unwrap() []error }); ok {
		return uw.Unwrap()
	}
	return []error{err}
}

func Contains(err error, substr string) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), substr)
}

func GetStackFrames(err error) []StackFrame {
	return getStackFrames(err.Error())
}

func ParseStampedError(errString string) *errx {

	frames := getStackFrames(errString)

	var existinge error
	var existingErr *errx
	for i := len(frames) - 1; i >= 0; i-- {
		frame := frames[i]

		switch true {
		case frame.IsStamped && !frame.IsWrapper:
			existingErr = newErr(frame.Stamp, frame.Msg).WithKind(frame.Kind)
			existinge = nil
		case frame.IsStamped && frame.IsWrapper:
			if existinge != nil {
				existingErr = wrapErr(frame.Stamp, existinge).WithKind(frame.Kind)
				existinge = nil
			} else {
				existingErr = wrapErr(frame.Stamp, existingErr).WithKind(frame.Kind)
			}
		case !frame.IsStamped && frame.IsWrapper:
			existinge = fmt.Errorf("%s %w", frame.Msg, existingErr)
			existingErr = nil
		case !frame.IsStamped && frame.IsWrapper:
			existinge = errors.New(frame.Msg)
			existingErr = nil
		}
	}

	if existingErr != nil {
		return existingErr
	} else {
		return &errx{err: existinge}
	}
}

func Cause(err error) error {
	var e error
	for err != nil {
		e = Unwrap(err)
		if e == nil {
			return err
		} else {
			err = e
		}
	}
	return err
}

func Panic(err error) {
	if err != nil {
		panic(err)
	}
}

func IsKind(err error, kind errKind) bool {
	for err != nil {
		if e, ok := err.(interface{ Kind() string }); ok {
			return e.Kind() == kind.kind
		}
		err = Unwrap(err)
	}
	return false
}

func IsDataKind[T DataType](err error, kind func(d T) errKind) bool {
	var d T
	for err != nil {
		if e, ok := err.(interface{ Kind() string }); ok {
			return e.Kind() == kind(d).kind
		}
		err = Unwrap(err)
	}
	return false
}

type report struct {
	Msg    string
	Traces []int
	Kind   string
}

func (r *report) Error() string {
	return r.Msg
}

func newReport(msg string, traces []int, kind string) report {
	return report{
		Msg:    msg,
		Traces: traces,
		Kind:   kind,
	}
}

func Report(err error) report {
	var msg string
	traces := make([]int, 0, 10)
	var kind string

	for err != nil {
		if v, ok := err.(StampedErr); ok {
			msg = v.Msg()
			traces = append(traces, v.Stamp())
			k := v.Kind()
			if k != "" && kind == "" {
				kind = k
			}

		} else {
			m := err.Error()
			if !strings.HasPrefix(m, "[ts ") {
				msg = m
			}
		}

		err = Unwrap(err)
	}

	return newReport(msg, traces, kind)
}

type StackFrame struct {
	IsWrapper bool
	IsStamped bool
	Stamp     literalInt
	Kind      errKind
	Msg       string
}

func newStackFrame(isWrapper bool, stampStr string, kindStr, dataStr, msg string) StackFrame {
	isUnstamped := false
	ts, err := strconv.Atoi(stampStr)
	if err != nil {
		isUnstamped = true
	}

	var data dataValue
	if dataStr != "" {
		data = dataValue{valStr: dataStr, isSet: true}
	}

	return StackFrame{
		IsWrapper: isWrapper,
		IsStamped: !isUnstamped,
		Stamp:     literalInt(ts),
		Kind:      errKind{kind: kindStr, data: data},
		Msg:       strings.TrimSpace(msg),
	}
}

func getStackFrames(errStr string) []StackFrame {

	lex := newLexer(errStr)
	tok := lex.nextToken()
	for tok.typ != eof {
		tok = lex.nextToken()
	}

	tree := lex.list
	frames := make([]StackFrame, 0, (len(errStr)/7)+1)
	stampStr := ""
	kindStr := ""
	dataStr := ""
	msg := ""

	clearBuffer := func() {
		stampStr = ""
		kindStr = ""
		dataStr = ""
		msg = ""
	}

	for i := 0; i < len(tree); i++ {
		tok := tree[i]
		if ms := strings.TrimSpace(msg); len(ms) > 0 && tok.typ == openBrackets {
			// A generic error message wrapping a stamped error
			frames = append(frames, newStackFrame(true, "", "", "", ms))
			clearBuffer()
		}

		if tok.typ == wrapperDelimiter {
			// A stamped error message wrapping another error
			frames = append(frames, newStackFrame(true, stampStr, kindStr, dataStr, strings.TrimSpace(msg)))
			clearBuffer()
		}

		if tok.typ == stampDirective {
			// Bring out stamp id
			for {
				nextToken := tree[i+1]
				if nextToken.typ == closeBrackets || nextToken.typ == kindDirective || nextToken.typ == dataDirective {
					break
				}
				stampStr = stampStr + nextToken.literal
				i++
			}
		}

		if tok.typ == kindDirective {
			// Bring out kind error
			for {
				nextToken := tree[i+1]
				if nextToken.typ == closeBrackets || nextToken.typ == dataDirective {
					break
				}
				kindStr = kindStr + nextToken.literal
				i++
			}
		}

		if tok.typ == dataDirective {
			// Bring out data value
			for {
				nextToken := tree[i+1]
				if nextToken.typ == closeBrackets {
					break
				}
				dataStr = dataStr + nextToken.literal
				i++
			}
		}

		if tok.typ == unknownToken {
			// Append unknown characters
			msg = msg + tok.literal
		}
	}

	frames = append(frames, newStackFrame(false, stampStr, kindStr, dataStr, msg))
	return frames
}

func fromStr[T any](str string) (*T, error) {
	var result any
	var typeName string = fmt.Sprintf("%T", *new(T))

	switch typeName {
	case "int":
		val, err := strconv.Atoi(str)
		if err != nil {
			return nil, err
		}
		result = val

	case "float32":
		l, err := strconv.ParseFloat(str, 32)
		if err != nil {
			return nil, err
		}
		result = float32(l)

	case "float64":
		l, err := strconv.ParseFloat(str, 64)
		if err != nil {
			return nil, err
		}
		result = l

	case "string":
		l, err := strconv.Unquote(str)
		if err != nil {
			return nil, err
		}
		result = l

	case "[]int":
		lstr := strings.Split(str[1:len(str)-1], ", ")
		lint := make([]int, 0, len(lstr))
		for _, v := range lstr {
			val, err := strconv.Atoi(v)
			if err != nil {
				return nil, err
			}
			lint = append(lint, val)
		}
		result = lint

	case "[]float32":
		lstr := strings.Split(str[1:len(str)-1], ", ")
		lint := make([]float32, 0, len(lstr))
		for _, v := range lstr {
			val, err := strconv.ParseFloat(v, 32)
			if err != nil {
				return nil, err
			}
			lint = append(lint, float32(val))
		}
		result = lint

	case "[]float64":
		lstr := strings.Split(str[1:len(str)-1], ", ")
		lint := make([]float64, 0, len(lstr))
		for _, v := range lstr {
			val, err := strconv.ParseFloat(v, 64)
			if err != nil {
				return nil, err
			}
			lint = append(lint, val)
		}
		result = lint

	case "[]string":
		var l []string
		err := json.Unmarshal([]byte(str), &l)
		if err != nil {
			return nil, err
		}
		result = l

	case "map[string]int":
		ml := make(map[string]int)
		err := json.Unmarshal([]byte(str), &ml)
		if err != nil {
			return nil, err
		}
		result = ml

	case "map[string]float32":
		ml := make(map[string]float32)
		err := json.Unmarshal([]byte(str), &ml)
		if err != nil {
			return nil, err
		}
		result = ml

	case "map[string]float64":
		ml := make(map[string]float64)
		err := json.Unmarshal([]byte(str), &ml)
		if err != nil {
			return nil, err
		}
		result = ml

	case "map[string]string":
		ml := make(map[string]string)
		err := json.Unmarshal([]byte(str), &ml)
		if err != nil {
			return nil, err
		}
		result = ml

	default:
		return nil, fmt.Errorf("match not found for type '%s'", typeName)
	}

	rtn, ok := result.(T)
	if !ok {
		return nil, fmt.Errorf("error converting match value '%v' to type '%s'", result, typeName)
	}

	return &rtn, nil
}
