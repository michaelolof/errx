package errx

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmptySource(t *testing.T) {
	lex := newLexer("")
	tok := lex.nextToken()

	assert.False(t, lex.hasNext())
	assert.Equal(t, tok.typ, eof)
}

func TestSampleStrings(t *testing.T) {
	{
		const errStr string = " [ts 1746358853 kind fileopen data users.txt]; "

		frames := getStackFrames(errStr)

		assert.Len(t, frames, 1)
		assert.Equal(t, frames[0].Kind.kind, "fileopen")
		assert.Equal(t, frames[0].Kind.data.valStr, "users.txt")
		assert.Equal(t, int(frames[0].Stamp), 1746358853)
		assert.Equal(t, frames[0].IsStamped, true)

	}

	{
		const errStr string = " [ts 1746358853 kind fileopen data users.txt] stuff went wrong"

		frames := getStackFrames(errStr)

		assert.Len(t, frames, 1)
		assert.Equal(t, frames[0].Kind.kind, "fileopen")
		assert.Equal(t, frames[0].Kind.data.valStr, "users.txt")
		assert.Equal(t, int(frames[0].Stamp), 1746358853)
		assert.Equal(t, frames[0].IsStamped, true)
		assert.Equal(t, frames[0].Msg, "stuff went wrong")
	}
}

func TestSampleString2(t *testing.T) {
	{
		const errString string = "[ts 1740258221]; [ts 1734858838]; [ts 1734805746] hash doesn't match"

		frames := getStackFrames(errString)

		fmt.Println("done", frames)
	}

	{
		const errStr string = `[ts 1713607010089 kind three data https://www.google.com]; [ts 1713607005378 kind two data 1.56]; i am a suspect error: [ts 1713606995137 kind one data 30] something went wrong`

		frames := getStackFrames(errStr)

		err := stacksToErr(frames)
		fmt.Println(err)

		assert.Len(t, frames, 4)
		fm1, fm2, fm3, fm4 := frames[0], frames[1], frames[2], frames[3]

		assert.Equal(t, fm1.IsStamped, true)
		assert.Equal(t, fm1.Kind.data.valStr, "https://www.google.com")
		assert.Equal(t, fm1.Kind.kind, "three")
		assert.Equal(t, fm2.IsStamped, true)
		assert.Equal(t, fm2.Kind.data.valStr, "1.56")
		assert.Equal(t, fm2.Kind.kind, "two")
		assert.Equal(t, fm3.IsStamped, false)
		assert.Equal(t, fm3.Kind.kind, "")
		assert.Equal(t, fm3.Msg, "i am a suspect error:")
		assert.Equal(t, fm4.IsStamped, true)
		assert.Equal(t, fm4.Kind.kind, "one")
		assert.Equal(t, fm4.Msg, "something went wrong")
	}

	{
		const errStr string = " [ts 1746358853]; stuff went wrong"

		frames := getStackFrames(errStr)

		assert.Len(t, frames, 2)
		frm1, frm2 := frames[0], frames[1]

		assert.Equal(t, frm1.Kind.kind, "")
		assert.Equal(t, frm1.Kind.data.valStr, "")
		assert.Equal(t, int(frm1.Stamp), 1746358853)
		assert.Equal(t, frm1.IsStamped, true)
		assert.Equal(t, frm1.Msg, "")
		assert.Equal(t, frm2.IsStamped, false)
		assert.Equal(t, frm2.Msg, "stuff went wrong")
	}

	{
		const errStr string = " [ts 1746358853]; [ts 1746358854 kind notfound]; [ts 1746358855] stuff went wrong "

		frames := getStackFrames(errStr)

		fmt.Println(stacksToErr(frames))

		assert.Len(t, frames, 3)
		frm1, frm2, frm3 := frames[0], frames[1], frames[2]

		assert.Equal(t, frm1.Kind.kind, "")
		assert.Equal(t, frm1.Kind.data.valStr, "")
		assert.Equal(t, int(frm1.Stamp), 1746358853)
		assert.Equal(t, frm1.IsStamped, true)
		assert.Equal(t, frm1.Msg, "")
		assert.Equal(t, frm2.IsStamped, true)
		assert.Equal(t, frm2.Kind.kind, "notfound")
		assert.Equal(t, frm3.IsStamped, true)
		assert.Equal(t, frm3.Kind.kind, "")
		assert.Equal(t, frm3.Msg, "stuff went wrong")
	}
}
