package errx

import (
	"strings"
)

type ReportMode int

const (
	Reversed       ReportMode = 2
	Indent         ReportMode = 3
	ReversedIndent ReportMode = 4
)

func Report(err error, mode ReportMode) string {
	switch mode {
	case Reversed:
		frames := splitToFrames(err, 10)
		reversed := make([]string, 0, len(frames))
		for i := len(frames) - 1; i >= 0; i-- {
			v := strings.TrimSpace(frames[i].err().Error())
			if len(v) > 0 {
				reversed = append(reversed, v)
			}
		}

		return strings.Join(reversed, "; ")

	case Indent:
		frames := splitToFrames(err, 10)
		indented := make([]string, 0, len(frames))
		for idx, frame := range frames {
			v := strings.TrimSpace(frame.err().Error())
			if len(v) > 0 {
				indented = append(indented, leftPad(v, idx*2))
			}
		}

		return strings.Join(indented, ";\n")

	case ReversedIndent:
		frames := splitToFrames(err, 10)
		reversed := make([]string, 0, len(frames))
		count := 0
		for i := len(frames) - 1; i >= 0; i-- {
			v := strings.TrimSpace(frames[i].err().Error())
			if len(v) > 0 {
				reversed = append(reversed, leftPad(v, count*2))
			}
			count++
		}

		return strings.Join(reversed, ";\n")
	}

	return err.Error()
}

func splitToFrames(err error, cap int) []stackFrame {
	if cap == 0 {
		cap = 10
	}
	frames := make([]stackFrame, 0, cap)
	for err != nil {
		uerr := Unwrap(err)
		if uerr == nil {
			fms := getStackFrames(strings.TrimSpace(err.Error()))
			frames = append(frames, fms[0])
		} else if v := strings.Split(err.Error(), uerr.Error()); len(v) > 0 {
			if v[0] != "" {
				fms := getStackFrames(strings.TrimSpace(v[0]))
				frames = append(frames, fms[0])
			} else {
				err = Unwrap(err)
				continue
			}
		}
		err = uerr
	}
	return frames
}

func leftPad(s string, length int) string {
	// if len(s) >= length {
	// 	return s
	// }
	padding := ""
	for range length {
		padding += " "
	}
	return padding + s
}
