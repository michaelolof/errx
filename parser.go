package errx

type lexer struct {
	source            string
	currChar          byte
	currPos           int
	nextPos           int
	openBracketsCount []int
	list              []token
}

func newLexer(input string) *lexer {
	l := &lexer{
		source:            input,
		openBracketsCount: make([]int, 0, (len(input)/7)+1),
	}
	l.readChar()
	return l
}

func (l *lexer) nextToken() token {
	var tok token
	currStr := string(l.currChar)
	switch l.currChar {
	case 0:
		tok.typ = eof
	case '[':
		tok.literal = currStr
		l.openBracketsCount = append(l.openBracketsCount, l.currPos)
		if len(l.openBracketsCount) == 1 {
			tok.typ = openBrackets
		} else {
			tok.typ = lBrackets
		}
	case ']':
		tok.literal = currStr
		arrPop(&l.openBracketsCount)
		if len(l.openBracketsCount) == 0 {
			tok.typ = closeBrackets
		} else {
			tok.typ = rBrackets
		}
	case ';':
		// preText := l.peekBehind(1)
		postText := l.peekAhead(2)
		// if preText == "]" && postText == "; " {
		if postText == "; " {
			tok.literal = currStr
			tok.typ = wrapperDelimiter
		}
	case 't':
		preText := l.peekBehind(1)
		postText := l.peekAhead(3)
		if preText == "[" && postText == "ts " {
			tok.typ = stampDirective
			tok.literal = "ts"
			l.moveCursor(3)
		} else {
			tok.literal = currStr
			tok.typ = unknownToken
		}
	case 'k':
		pretext := l.peekBehind(1)
		postText := l.peekAhead(5)
		if pretext == " " && postText == "kind " {
			arrPop(&l.list) // remove pretext space from tree
			tok.typ = kindDirective
			tok.literal = "kind"
			l.moveCursor(5)
		} else {
			tok.literal = currStr
			tok.typ = unknownToken
		}
	case 'd':
		preText := l.peekBehind(1)
		postText := l.peekAhead(5)
		if preText == " " && postText == "data " {
			arrPop(&l.list) // remove pretext space from tree
			tok.typ = dataDirective
			tok.literal = "data"
			l.moveCursor(5)
		} else {
			tok.literal = currStr
			tok.typ = unknownToken
		}
	default:
		tok.literal = currStr
		tok.typ = unknownToken
	}

	l.readChar()
	l.list = append(l.list, tok)
	return tok
}

func (l *lexer) readChar() {
	if l.nextPos >= len(l.source) {
		l.currChar = 0
	} else {
		l.currChar = (l.source)[l.nextPos]
	}

	l.currPos = l.nextPos
	l.nextPos += 1
}

func (l *lexer) peekBehind(n int) string {
	pos := l.currPos
	count := 0
	for pos > 0 {
		if count >= n {
			break
		}

		pos--
		count++
	}

	return l.source[pos:l.currPos]
}

func (l *lexer) peekAhead(n int) string {
	pos := l.currPos
	count := 0
	for l.currChar != 0 {
		if count >= n {
			break
		}

		pos++
		count++
	}

	return l.source[l.currPos:pos]
}

func (l *lexer) moveCursor(n int) {
	l.currPos = l.currPos + (n - 1)
	l.nextPos = l.currPos + 1
}

type token struct {
	literal string
	typ     tokenType
}

type tokenType int

const (
	openBrackets tokenType = iota + 1
	closeBrackets
	lBrackets
	rBrackets
	stampDirective
	kindDirective
	dataDirective
	wrapperDelimiter
	unknownToken
	eof
)
