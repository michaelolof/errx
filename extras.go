package errx

import (
	"errors"
	"fmt"
	"strings"
)

func TryPanic(err error) {
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}

func Contains(err error, substr string) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), substr)
}

// ScrubMsg returns a new string from the error message with the stamp information and the error data taken removed
func ScrubMsg(err error) string {
	msg := err.Error()
	msgs := strings.Split(msg, ": ")

	rtns := make([]string, 0, len(msgs))
	for _, m := range msgs {
		if strings.HasPrefix(m, "STAMP-") || strings.HasPrefix(m, "ERRDATA-") {
			continue
		}
		rtns = append(rtns, m)
	}

	return strings.Join(rtns, ": ")
}

// NewScrub returns a new error object with the stamp information and the error data taken removed
func NewScrub(err error) error {
	return errors.New(ScrubMsg(err))
}
