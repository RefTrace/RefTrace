package lexer

import (
	"fmt"
	"strings"
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

func (l *GroovyLexer) Emit() antlr.Token {
	l.tokenIndex++
	token := l.BaseLexer.Emit()
	tokenType := token.GetTokenType()
	if token.GetChannel() == antlr.TokenDefaultChannel {
		l.lastTokenType = tokenType
	}
	if tokenType == GroovyLexerRollBackOne {
		l.rollbackOneChar()
	}
	return token
}

func (b *GroovyLexer) GetAllTokens() []antlr.Token {
	vl := b // the base class uses b.Virt which resolves wrong
	tokens := make([]antlr.Token, 0)
	t := vl.NextToken()
	for t.GetTokenType() != antlr.TokenEOF {
		if t.GetChannel() == antlr.TokenDefaultChannel {
			tokens = append(tokens, t)
		}
		t = vl.NextToken()
	}
	return tokens
}

func (g *GroovyLexer) NextToken() antlr.Token {
	input := g.GetInput()
	if input == nil {
		panic("NextToken requires a non-nil input stream.")
	}

	tokenStartMarker := input.Mark()

	// previously in finally block
	defer func() {
		// make sure we release marker after Match or
		// unbuffered char stream will keep buffering
		input.Release(tokenStartMarker)
	}()

	for {
		if g.GetHitEOF() {
			g.EmitEOF()
			return g.GetToken()
		}
		g.SetToken(nil)
		g.SetChannel(antlr.TokenDefaultChannel)
		g.SetTokenStartCharIndex(input.Index())
		g.SetTokenStartColumn(g.Interpreter.GetCharPositionInLine())
		g.SetTokenStartLine(g.Interpreter.GetLine())
		g.SetText("")
		continueOuter := false
		for {
			g.SetTheType(antlr.TokenInvalidType)

			ttype := g.BaseLexer.SafeMatch() // Defaults to LexerSkip

			if input.LA(1) == antlr.TokenEOF {
				g.SetHitEOF(true)
			}
			if g.GetTheType() == antlr.TokenInvalidType {
				g.SetTheType(ttype)
			}
			if g.GetTheType() == antlr.LexerSkip {
				continueOuter = true
				break
			}
			if g.GetTheType() != antlr.LexerMore {
				break
			}
		}

		if continueOuter {
			continue
		}
		if g.GetToken() == nil {
			g.Emit()
		}
		return g.GetToken()
	}
}

func (l *GroovyLexer) rollbackOneChar() {
	interpreter := l.GetInterpreter().(*antlr.LexerATNSimulator)
	resetAcceptPosition(interpreter, l.GetInputStream(), l.TokenStartCharIndex-1, l.TokenStartLine, l.GetInterpreter().GetCharPositionInLine()-1)
}

func (l *GroovyLexer) handleRollBackOne() {
	istream := l.GetInputStream()
	readChar := istream.LA(-1)
	if l.GetInputStream().LA(1) == antlr.TokenEOF && (readChar == '"' || readChar == '/') {
		l.SetType(GroovyLexerGStringEnd)
	} else {
		l.SetChannel(antlr.TokenHiddenChannel)
	}
	l.PopMode()
}

func (b *GroovyLexer) Recover(re antlr.RecognitionException) {
	if _, ok := re.(*antlr.LexerNoViableAltException); ok {
		panic(re)
	} else {
		b.BaseLexer.Recover(re)
	}
}

type PositionAdjustingLexerATNSimulator struct {
	*antlr.LexerATNSimulator
}

func NewPositionAdjustingLexerATNSimulator(recog antlr.Lexer, atn *antlr.ATN, decisionToDFA []*antlr.DFA, sharedContextCache *antlr.PredictionContextCache) *PositionAdjustingLexerATNSimulator {
	return &PositionAdjustingLexerATNSimulator{
		LexerATNSimulator: antlr.NewLexerATNSimulator(recog, atn, decisionToDFA, sharedContextCache),
	}
}

func resetAcceptPosition(sim *antlr.LexerATNSimulator, input antlr.CharStream, index, line, charPositionInLine int) {
	input.Seek(index)
	sim.Line = line
	sim.CharPositionInLine = charPositionInLine
	sim.Consume(input)
}

// isJavaIdentifierStart checks if a given code point is a valid start character for a Java identifier.
// https://docs.oracle.com/javase/8/docs/api/java/lang/Character.html#isJavaIdentifierStart-char-
func isJavaIdentifierStart(codePoint rune) bool {
	return unicode.IsLetter(codePoint) || unicode.Is(unicode.Lm, codePoint) || unicode.Is(unicode.Nl, codePoint) || unicode.Is(unicode.Pc, codePoint)
}

// isIdentifierIgnorable checks if a given rune is an ignorable character in a Java identifier or a Unicode identifier.
// https://docs.oracle.com/javase/8/docs/api/java/lang/Character.html#isIdentifierIgnorable-char-
func isIdentifierIgnorable(ch rune) bool {
	// Check if the character is an ISO control character that is not whitespace
	if (ch >= '\u0000' && ch <= '\u0008') || (ch >= '\u000E' && ch <= '\u001B') || (ch >= '\u007F' && ch <= '\u009F') {
		return true
	}
	// Check if the character has the FORMAT general category value
	return unicode.Is(unicode.Cf, ch)
}

// isJavaIdentifierStartAndNotIdentifierIgnorable checks if a given rune is a valid start character for a Java identifier and not ignorable.
func isJavaIdentifierStartAndNotIdentifierIgnorable(ch int) bool {
	return isJavaIdentifierStart(rune(ch)) && !isIdentifierIgnorable(rune(ch))
}

func isJavaIdentifierPartAndNotIdentifierIgnorable(ch int) bool {
	return isJavaIdentifierPart(rune(ch)) && !isIdentifierIgnorable(rune(ch))
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
// https://docs.oracle.com/javase/8/docs/api/java/lang/Character.html#isJavaIdentifierPart-char-
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

// searchArray searches for a target integer in the given array.
// It returns true if the target is found, otherwise false.
func searchArray(arr []int, target int) bool {
	for _, value := range arr {
		if value == target {
			return true
		}
	}
	return false
}

// allowed iff not in REGEX_CHECK_ARRAY
func (l *GroovyLexer) isRegexAllowed() bool {
	return !searchArray(REGEX_CHECK_ARRAY, l.lastTokenType)
}

// matches is a helper function to check if a string matches a given pattern.
func matches(str string, pattern *unicode.RangeTable) bool {
	for _, r := range str {
		if unicode.Is(pattern, r) {
			return true
		}
	}
	return false
}

// Define the patterns used in the function
var (
	LETTER_AND_LEFTCURLY_PATTERN = &unicode.RangeTable{
		R16: []unicode.Range16{
			{Lo: 'a', Hi: 'z', Stride: 1},
			{Lo: 'A', Hi: 'Z', Stride: 1},
			{Lo: '_', Hi: '_', Stride: 1},
			{Lo: '{', Hi: '{', Stride: 1},
		},
	}
	NONSURROGATE_PATTERN = &unicode.RangeTable{
		R16: []unicode.Range16{
			{Lo: 0x0000, Hi: 0x007F, Stride: 1},
			{Lo: 0xE000, Hi: 0xFFFF, Stride: 1},
		},
	}
	SURROGATE_PAIR1_PATTERN = &unicode.RangeTable{
		R16: []unicode.Range16{
			{Lo: 0xD800, Hi: 0xDBFF, Stride: 1},
		},
	}
	SURROGATE_PAIR2_PATTERN = &unicode.RangeTable{
		R16: []unicode.Range16{
			{Lo: 0xDC00, Hi: 0xDFFF, Stride: 1},
		},
	}
)

// isFollowedByJavaLetterInGString checks if the character following the current position in the CharStream is a valid Java identifier part.
func isFollowedByJavaLetterInGString(cs antlr.CharStream) bool {
	c1 := cs.LA(1)

	if c1 == '$' { // single $ is not a valid identifier
		return false
	}

	str1 := string(rune(c1))

	if matches(str1, LETTER_AND_LEFTCURLY_PATTERN) {
		return true
	}

	if matches(str1, NONSURROGATE_PATTERN) && isJavaIdentifierPart(rune(c1)) {
		return true
	}

	c2 := cs.LA(2)
	str2 := string(rune(c2))

	if matches(str1, SURROGATE_PAIR1_PATTERN) && matches(str2, SURROGATE_PAIR2_PATTERN) {
		codePoint := utf16.DecodeRune(rune(c1), rune(c2))
		if isJavaIdentifierPart(codePoint) {
			return true
		}
	}

	return false
}

func escapeSingleQuotes(input string) string {
	return strings.Replace(input, "'", "\\'", -1)
}

// isFollowedBy checks if the character following the current position in the CharStream matches any of the provided characters.
func isFollowedBy(cs antlr.CharStream, chars ...rune) bool {
	c1 := cs.LA(1)

	for _, c := range chars {
		if c1 == int(c) {
			return true
		}
	}

	return false
}

// isUpperCase checks if a given rune is an uppercase letter.
func isUpperCase(ch int) bool {
	return unicode.IsUpper(rune(ch))
}

type CustomErrorListener struct {
	*antlr.DefaultErrorListener
	filename string
	hasError bool
}

func NewCustomErrorListener(filename string) *CustomErrorListener {
	return &CustomErrorListener{filename: filename}
}

func (c *CustomErrorListener) SyntaxError(_ antlr.Recognizer, _ interface{}, line, column int, msg string, _ antlr.RecognitionException) {
	c.hasError = true
	fmt.Printf("File: %s - line %d:%d %s\n", c.filename, line, column, msg)
}

func (c *CustomErrorListener) ReportAmbiguity(_ antlr.Parser, _ *antlr.DFA, _, _ int, _ bool, _ *antlr.BitSet, _ *antlr.ATNConfigSet) {
	c.hasError = true
	fmt.Printf("File: %s - Ambiguity detected\n", c.filename)
}

func (c *CustomErrorListener) ReportAttemptingFullContext(_ antlr.Parser, _ *antlr.DFA, _, _ int, _ *antlr.BitSet, _ *antlr.ATNConfigSet) {
	c.hasError = true
	fmt.Printf("File: %s - Attempting full context\n", c.filename)
}

func (c *CustomErrorListener) ReportContextSensitivity(_ antlr.Parser, _ *antlr.DFA, _, _, _ int, _ *antlr.ATNConfigSet) {
	c.hasError = true
	fmt.Printf("File: %s - Context sensitivity\n", c.filename)
}

func (c *CustomErrorListener) HasError() bool {
	return c.hasError
}
