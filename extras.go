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

func GetStackFrames(err error) []stackFrame {
	return getStackFrames(err.Error())
}

func ParseStampedError(errString string) *errx {
	return stacksToErr(getStackFrames(errString))
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

type stackFrame struct {
	IsStamped bool
	Stamp     lint
	Kind      errKind
	Msg       string
}

func (s stackFrame) err() *errx {
	return stacksToErr([]stackFrame{s})
}

func newStackFrame(stampStr string, kindStr, dataStr, msg string) stackFrame {
	isUnstamped := false
	ts, err := strconv.Atoi(stampStr)
	if err != nil {
		isUnstamped = true
	}

	var data dataValue
	if dataStr != "" {
		data = dataValue{valStr: dataStr, isSet: true}
	}

	return stackFrame{
		IsStamped: !isUnstamped,
		Stamp:     lint(ts),
		Kind:      errKind{kind: kindStr, data: data},
		Msg:       strings.TrimSpace(msg),
	}
}

func getStackFrames(errStr string) []stackFrame {
	lex := newLexer(errStr)
	tok := lex.nextToken()
	for lex.hasNext() && tok.typ != eof {
		tok = lex.nextToken()
	}

	tree := lex.list
	frames := make([]stackFrame, 0, (len(errStr)/7)+1)

	type state struct {
		stamp   string
		kind    string
		data    string
		msg     string
		inStamp bool
	}

	buff := state{}

	clearBuffer := func() {
		buff = state{}
	}

	for i := 0; i < len(tree); i++ {
		tok := tree[i]
		switch tok.typ {
		case openBrackets:
			// A generic error message wrapping a stamped error
			if len(buff.msg) > 0 {
				frames = append(frames, newStackFrame("", "", "", strings.TrimSpace(buff.msg)))
				clearBuffer()
			}

			buff.inStamp = true

		case wrapperDelimiter:
			// A stamped error message wrapping another error
			if buff.inStamp {
				frames = append(frames, newStackFrame(buff.stamp, buff.kind, buff.data, strings.TrimSpace(buff.msg)))
				clearBuffer()
			}

		case stampDirective:
			// Bring out stamp id
			for {
				nextToken := tree[i+1]
				if nextToken.typ == closeBrackets || nextToken.typ == kindDirective || nextToken.typ == dataDirective {
					break
				}
				buff.stamp = buff.stamp + nextToken.literal
				i++
			}

		case kindDirective:
			// Bring out kind error
			for {
				nextToken := tree[i+1]
				if nextToken.typ == closeBrackets || nextToken.typ == dataDirective {
					break
				}
				buff.kind = buff.kind + nextToken.literal
				i++
			}

		case dataDirective:
			// Bring out data value
			for {
				nextToken := tree[i+1]
				if nextToken.typ == closeBrackets {
					break
				}
				buff.data = buff.data + nextToken.literal
				i++
			}

		case unknownToken:
			// Append unknown characters
			buff.msg = buff.msg + tok.literal

		case emptySpace:
			if len(buff.msg) > 0 {
				buff.msg = buff.msg + tok.literal
			}

		}
	}

	if buff.inStamp {
		frames = append(frames, newStackFrame(buff.stamp, buff.kind, buff.data, buff.msg))
	} else if buff.msg != "" {
		frames = append(frames, newStackFrame("", "", "", buff.msg))
	}

	return frames
}

func stacksToErr(frames []stackFrame) *errx {
	var existinge error
	var existingErr *errx
	l := len(frames)
	for i := l - 1; i >= 0; i-- {
		frame := frames[i]
		isWrapper := i < (l - 1)

		switch true {
		case frame.IsStamped && !isWrapper:
			existingErr = newErr(frame.Stamp, frame.Msg).WithKind(frame.Kind)
			existinge = nil
		case frame.IsStamped && isWrapper:
			if existinge != nil {
				existingErr = wrapErr(frame.Stamp, existinge).WithKind(frame.Kind)
				existinge = nil
			} else {
				existingErr = wrapErr(frame.Stamp, existingErr).WithKind(frame.Kind)
			}
		case !frame.IsStamped && isWrapper:
			existinge = fmt.Errorf("%s %w", frame.Msg, existingErr)
			existingErr = nil
		case !frame.IsStamped && !isWrapper:
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
