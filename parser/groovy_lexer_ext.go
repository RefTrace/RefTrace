package parser

import (
	"fmt"
	"sort"
	"unicode"
	"unicode/utf16"

	"github.com/antlr4-go/antlr/v4"
)

type Paren struct {
	text          string
	lastTokenType int
	line          int
	column        int
}

type MyGroovyLexer struct {
	*antlr.BaseLexer
	errorIgnored      bool
	tokenIndex        int64
	lastTokenType     int
	invalidDigitCount int
	parenStack        []Paren
}

// isJavaIdentifierStart checks if a given code point is a valid start character for a Java identifier.
// https://docs.oracle.com/javase%2F8%2Fdocs%2Fapi%2F%2F/java/lang/Character.html#isJavaIdentifierStart-char-
func isJavaIdentifierStart(codePoint rune) bool {
	return unicode.IsLetter(codePoint) || unicode.Is(unicode.Lm, codePoint) || unicode.Is(unicode.Nl, codePoint) || unicode.Is(unicode.Pc, codePoint)
}

// isIdentifierIgnorable checks if a given rune is an ignorable character in a Java identifier or a Unicode identifier.
// https://docs.oracle.com/javase%2F8%2Fdocs%2Fapi%2F%2F/java/lang/Character.html#isIdentifierIgnorable-char-
func isIdentifierIgnorable(ch rune) bool {
	// Check if the character is an ISO control character that is not whitespace
	if (ch >= '\u0000' && ch <= '\u0008') || (ch >= '\u000E' && ch <= '\u001B') || (ch >= '\u007F' && ch <= '\u009F') {
		return true
	}
	// Check if the character has the FORMAT general category value
	return unicode.Is(unicode.Cf, ch)
}

// isJavaIdentifierStartAndNotIdentifierIgnorable checks if a given rune is a valid start character for a Java identifier and not ignorable.
func isJavaIdentifierStartAndNotIdentifierIgnorable(ch rune) bool {
	return isJavaIdentifierStart(ch) && !isIdentifierIgnorable(ch)
}

func isJavaIdentifierPartAndNotIdentifierIgnorable(ch rune) bool {
	return isJavaIdentifierPart(ch) && !isIdentifierIgnorable(ch)
}

// isJavaIdentifierStartFromSurrogatePair checks if the characters at positions laMinus2 and laMinus1 form a valid surrogate pair and if the resulting code point is a valid start character for a Java identifier.
func isJavaIdentifierStartFromSurrogatePair(laMinus2, laMinus1 int) bool {
	if laMinus2 >= 0xD800 && laMinus2 <= 0xDBFF && laMinus1 >= 0xDC00 && laMinus1 <= 0xDFFF {
		codePoint := utf16.DecodeRune(rune(laMinus2), rune(laMinus1))
		return isJavaIdentifierStart(codePoint)
	}
	return false
}

// isJavaIdentifierPart checks if a given code point is a valid part character for a Java identifier.
// https://docs.oracle.com/javase%2F8%2Fdocs%2Fapi%2F%2F/java/lang/Character.html#isJavaIdentifierPart-char-
func isJavaIdentifierPart(codePoint rune) bool {
	return unicode.IsLetter(codePoint) ||
		unicode.IsDigit(codePoint) ||
		unicode.Is(unicode.Lm, codePoint) ||
		unicode.Is(unicode.Nl, codePoint) ||
		unicode.Is(unicode.Pc, codePoint) ||
		unicode.Is(unicode.Mn, codePoint) ||
		unicode.Is(unicode.Mc, codePoint) ||
		isIdentifierIgnorable(codePoint)
}

// isJavaIdentifierPartFromSurrogatePair checks if the characters at positions laMinus2 and laMinus1 form a valid surrogate pair and if the resulting code point is a valid part character for a Java identifier.
func isJavaIdentifierPartFromSurrogatePair(laMinus2, laMinus1 int) bool {
	if laMinus2 >= 0xD800 && laMinus2 <= 0xDBFF && laMinus1 >= 0xDC00 && laMinus1 <= 0xDFFF {
		codePoint := utf16.DecodeRune(rune(laMinus2), rune(laMinus1))
		return isJavaIdentifierPart(codePoint)
	}
	return false
}

func require(condition bool, message string, offset int, lexer *GroovyLexer) {
	if !condition {
		line := lexer.GetLine()
		column := lexer.GetCharPositionInLine() + offset
		errorMsg := fmt.Sprintf("line %d:%d %s", line, column, message)
		panic(antlr.NewBaseRecognitionException(errorMsg, lexer, lexer.GetInputStream(), nil))
	}
}

func (l *GroovyLexer) enterParenCallback(text string) {
	// This method is intended to be overridden
}

func (l *GroovyLexer) enterParen() {
	text := l.GetText()
	l.enterParenCallback(text)
	l.parenStack = append(l.parenStack, Paren{text, l.lastTokenType, l.GetLine(), l.GetCharPositionInLine()})
}

func (l *GroovyLexer) exitParenCallback(text string) {
	// This method is intended to be overridden
}

func (l *GroovyLexer) exitParen() {
	text := l.GetText()
	l.exitParenCallback(text)
	if len(l.parenStack) > 0 {
		l.parenStack = l.parenStack[:len(l.parenStack)-1]
	}
}

func (l *GroovyLexer) isInsideParens() bool {
	if len(l.parenStack) == 0 {
		return false
	}
	paren := l.parenStack[len(l.parenStack)-1]
	text := paren.text
	return (text == "(" && paren.lastTokenType != GroovyLexerTRY) || text == "[" || text == "?["
}

func (l *GroovyLexer) ignoreTokenInsideParens() {
	if !l.isInsideParens() {
		return
	}
	l.SetChannel(antlr.TokenHiddenChannel)
}

func (l *GroovyLexer) addComment(_type int) {
	// TODO: implement this
	//text := l.GetInputStream().GetText(antlr.NewInterval(l.GetTokenStartCharIndex(), l.GetCharIndex()-1))
	// Handle the comment text as needed
}

func (l *GroovyLexer) isFollowedByWhiteSpaces() bool {
	input := l.GetInputStream()
	for i := l.GetCharIndex(); i < input.Size(); i++ {
		ch := input.LA(i + 1)
		if ch == antlr.TokenEOF {
			break
		}
		if !unicode.IsSpace(rune(ch)) {
			return false
		}
		if unicode.IsSpace(rune(ch)) {
			return true
		}
	}
	return false
}

func (l *GroovyLexer) ignoreMultiLineCommentConditionally() {
	if !l.isInsideParens() && l.isFollowedByWhiteSpaces() {
		return
	}
	l.SetChannel(antlr.TokenHiddenChannel)
}

var REGEX_CHECK_ARRAY = []int{
	GroovyLexerDEC, GroovyLexerINC, GroovyLexerTHIS, GroovyLexerRBRACE, GroovyLexerRBRACK, GroovyLexerRPAREN, GroovyLexerGStringEnd, GroovyLexerNullLiteral,
	GroovyLexerStringLiteral, GroovyLexerBooleanLiteral, GroovyLexerIntegerLiteral, GroovyLexerFloatingPointLiteral,
	GroovyLexerIdentifier, GroovyLexerCapitalizedIdentifier,
}

func (l *GroovyLexer) isRegexAllowed() bool {
	return sort.SearchInts(REGEX_CHECK_ARRAY, l.lastTokenType) < 0
}
