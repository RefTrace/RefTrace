package parser

// Token represents a CSTNode produced by the Lexer.
type Token struct {
	tokenType   int
	meaning     int
	text        string
	startLine   int
	startColumn int
}

var (
	TokenEOF  = NewToken(EOF, "", -1, -1)
	TokenNULL = NewToken(UNKNOWN, "", -1, -1)
)

// NewToken initializes a Token with the specified information.
func NewToken(tokenType int, text string, startLine, startColumn int) *Token {
	return &Token{
		tokenType:   tokenType,
		meaning:     tokenType,
		text:        text,
		startLine:   startLine,
		startColumn: startColumn,
	}
}

// Dup returns a copy of this Token.
/*
func (t *Token) Dup() *Token {
	token := NewToken(t.tokenType, t.text, t.startLine, t.startColumn)
	token.SetMeaning(t.meaning)
	return token
}
*/

// GetMeaning returns the meaning of this node.
func (t *Token) GetMeaning() int {
	return t.meaning
}

// SetMeaning sets the meaning for this node.
func (t *Token) SetMeaning(meaning int) *Token {
	if t != TokenEOF && t != TokenNULL {
		t.meaning = meaning
	}
	return t
}

// GetType returns the actual type of the node.
func (t *Token) GetType() int {
	return t.tokenType
}

// Size returns the number of elements in the node (always 1 for Token).
func (t *Token) Size() int {
	return 1
}

// Get returns the specified element, or nil.
/*
func (t *Token) Get(index int) CSTNode {
	if index > 0 {
		panic("attempt to access Token element other than root")
	}
	return t
}
*/

// GetRoot returns the root of the node (self for Token).
func (t *Token) GetRoot() *Token {
	return t
}

// GetRootText returns the text of the root node.
func (t *Token) GetRootText() string {
	return t.text
}

// GetText returns the text of the token.
func (t *Token) GetText() string {
	return t.text
}

// SetText sets the text of the token.
func (t *Token) SetText(text string) {
	if t != TokenEOF && t != TokenNULL {
		t.text = text
	}
}

// GetStartLine returns the starting line of the node.
func (t *Token) GetStartLine() int {
	return t.startLine
}

// GetStartColumn returns the starting column of the node.
func (t *Token) GetStartColumn() int {
	return t.startColumn
}

// TODO: handle reductions

// NewKeyword creates a token that represents a keyword.
func NewKeyword(text string, startLine, startColumn int) *Token {
	tokenType := LookupKeyword(text)
	if tokenType != UNKNOWN {
		return NewToken(tokenType, text, startLine, startColumn)
	}
	return nil
}

// NewString creates a token that represents a double-quoted string.
func NewString(text string, startLine, startColumn int) *Token {
	return NewToken(STRING, text, startLine, startColumn)
}

// NewIdentifier creates a token that represents an identifier.
func NewIdentifier(text string, startLine, startColumn int) *Token {
	return NewToken(IDENTIFIER, text, startLine, startColumn)
}

// NewInteger creates a token that represents an integer.
func NewInteger(text string, startLine, startColumn int) *Token {
	return NewToken(INTEGER_NUMBER, text, startLine, startColumn)
}

// NewDecimal creates a token that represents a decimal number.
func NewDecimal(text string, startLine, startColumn int) *Token {
	return NewToken(DECIMAL_NUMBER, text, startLine, startColumn)
}

// NewSymbol creates a token that represents a symbol.
func NewSymbol(tokenType int, startLine, startColumn int) *Token {
	return NewToken(tokenType, GetText(tokenType), startLine, startColumn)
}

// NewSymbolFromText creates a token that represents a symbol, using a library for the type.
func NewSymbolFromText(text string, startLine, startColumn int) *Token {
	return NewToken(LookupSymbol(text), text, startLine, startColumn)
}

// NewPlaceholder creates a token with the specified meaning.
func NewPlaceholder(meaning int) *Token {
	token := NewToken(UNKNOWN, "", -1, -1)
	token.SetMeaning(meaning)
	return token
}
