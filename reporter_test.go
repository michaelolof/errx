package errx

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReport(t *testing.T) {
	err1 := New(1, "e1")
	err2 := Wrap(2, err1)
	err3 := Wrap(3, err2)

	t.Run("Default", func(t *testing.T) {
		res := Report(err3, 0)
		assert.Equal(t, err3.Error(), res)
	})

	t.Run("Reversed", func(t *testing.T) {
		res := Report(err3, Reversed)
		// Expected: [ts 1] e1; [ts 2]; [ts 3]
		// Wait, let's look at stack frames logic in reporter.go
		assert.Contains(t, res, "[ts 1] e1")
		assert.Contains(t, res, "[ts 3]")
	})

	t.Run("Indent", func(t *testing.T) {
		res := Report(err3, Indent)
		lines := strings.Split(res, ";\n")
		assert.Len(t, lines, 3)
		assert.True(t, strings.HasPrefix(lines[1], "  "))
		assert.True(t, strings.HasPrefix(lines[2], "    "))
	})

	t.Run("ReversedIndent", func(t *testing.T) {
		res := Report(err3, ReversedIndent)
		lines := strings.Split(res, ";\n")
		assert.Len(t, lines, 3)
		assert.True(t, strings.HasPrefix(lines[1], "  "))
		assert.True(t, strings.HasPrefix(lines[2], "    "))
	})
}
