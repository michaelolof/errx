package errx

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func TryPanic(err error) {
	if err != nil {
		panic(err)
	}
}

func Contains(err error, substr string) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), substr)
}

type StackFrame struct {
	IsWrapper   bool
	IsUnstamped bool
	Stamp       int
	Kind        error
	Data        any
	Msg         string
	FullMsg     string
}

func GetStackFrames(err error) ([]StackFrame, error) {
	return getStackFrames(err.Error())
}

func ParseStampedError(errString string) (error, error) {

	frames, err := getStackFrames(errString)
	if err != nil {
		return nil, fmt.Errorf("couldn't get stack frames: %w", err)
	}

	var existingErr error
	for i := len(frames) - 1; i >= 0; i-- {
		frame := frames[i]

		if !frame.IsWrapper && frame.IsUnstamped {
			existingErr = errors.New(frame.FullMsg)
		} else if !frame.IsWrapper {
			existingErr = NewErr(frame.Stamp, frame.Msg, frame.Kind, frame.Data)
		} else if frame.IsWrapper && frame.IsUnstamped {
			existingErr = fmt.Errorf("%s %w", frame.FullMsg, existingErr)
		} else if existingErr != nil {
			existingErr = WrapErr(frame.Stamp, existingErr, frame.Kind, frame.Data)
		}
	}

	return existingErr, nil
}

func getStackFrames(errString string) ([]StackFrame, error) {

	extractFrames := func(str string) (*StackFrame, error) {
		var stampStr, kindStr, dataStr string
		parts := strings.Split(str, "[stamp ")
		parts = strings.Split(parts[1], "]")
		parts = strings.Split(parts[0], " kind ")
		if len(parts) == 1 {
			parts = strings.Split(parts[0], " data ")
			stampStr = parts[0]
			if len(parts) > 1 {
				dataStr = parts[1]
			}
		} else {
			stampStr = parts[0]
			parts = strings.Split(parts[1], " data ")
			if len(parts) > 0 {
				kindStr = parts[0]
			}
			if len(parts) > 1 {
				dataStr = parts[1]
			}
		}

		stamp, err := strconv.Atoi(stampStr)
		if err != nil {
			return nil, fmt.Errorf("couldn't parse stamp text to integer: %w", err)
		}

		var kind error
		if kindStr != "" {
			kind = errors.New(kindStr)
		}

		var data any
		if dataStr != "" {
			err = json.Unmarshal([]byte(dataStr), &data)
			if err != nil {
				return nil, fmt.Errorf("couldn't parse data string to any data: %w", err)
			}
		}

		return &StackFrame{
			Stamp: stamp,
			Kind:  kind,
			Data:  data,
		}, nil
	}

	if !strings.HasPrefix(errString, "[stamp ") {
		return nil, errors.New("unsupported error string. only stamped errors can be parsed")
	}

	strl := len(errString)
	errs := make([]error, 0, (strl/21)+1)
	frames := make([]StackFrame, 0, (strl/21)+1)

	// We already know errString starts with [stamp
	for pos := 0; pos < strl; pos++ {

		wkstr := errString[pos:]

		if pos > 0 && !strings.HasPrefix(wkstr, "[stamp ") {
			// Look for the next stamp
			nidx := strings.Index(wkstr, " [stamp ")
			if nidx < 0 {
				f := StackFrame{
					IsWrapper:   false,
					IsUnstamped: true,
					FullMsg:     wkstr,
				}
				frames = append(frames, f)
				break
			} else {
				f := StackFrame{
					IsWrapper:   true,
					IsUnstamped: true,
					FullMsg:     wkstr[0:nidx],
				}
				frames = append(frames, f)
				pos = pos + nidx
				continue
			}

		}

		// Find the next closing bracket index
		idx := pos + strings.IndexAny(wkstr, "]")
		rem := errString[pos : idx+1]

		// check if this is the end of an error instance
		if strl > idx+2 {
			next := errString[idx+1 : idx+3]
			if next == "; " {
				p, err := extractFrames(rem)
				if err != nil {
					errs = append(errs, fmt.Errorf("couldn't parse wrapped error: %w", err))
					continue
				}
				p.FullMsg = rem
				p.IsWrapper = true
				frames = append(frames, *p)
				pos = idx + 2
				continue
			} else {
				idx2 := strings.Index(errString[idx+1:], "; ")
				if idx2 == -1 {
					f, err := extractFrames(rem)
					if err != nil {
						errs = append(errs, fmt.Errorf("couldn't parse new error: %w", err))
						break
					}
					msg := errString[idx+1:]
					f.Msg = strings.TrimLeft(msg, " ")
					f.IsWrapper = false
					f.FullMsg = rem + msg
					frames = append(frames, *f)
					break
				}
			}
		}
	}

	return frames, errors.Join(errs...)
}
