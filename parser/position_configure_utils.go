package parser

import (
	"strings"

	"github.com/antlr4-go/antlr/v4"
)

type sourcePosition struct {
	startLine, startColumn, stopLine, stopColumn int
}

func (sp sourcePosition) GetStartLine() int   { return sp.startLine }
func (sp sourcePosition) GetStartColumn() int { return sp.startColumn }
func (sp sourcePosition) GetStopLine() int    { return sp.stopLine }
func (sp sourcePosition) GetStopColumn() int  { return sp.stopColumn }

// configureAST sets the position information for any ASTNode
func configureAST[T ASTNode](astNode T, ctx antlr.ParserRuleContext) T {
	start := ctx.GetStart()
	stop := ctx.GetStop()

	return configureASTWithTokens(astNode, start, stop)
}

func configureASTWithInitialStop[T ASTNode](astNode T, ctx antlr.ParserRuleContext, initialStop ASTNode) T {
	start := ctx.GetStart()
	astNode.SetLineNumber(start.GetLine())
	astNode.SetColumnNumber(start.GetColumn() + 1)

	if initialStop != nil {
		astNode.SetLastLineNumber(initialStop.GetLastLineNumber())
		astNode.SetLastColumnNumber(initialStop.GetLastColumnNumber())
	} else {
		stop := ctx.GetStop()
		configureEndPosition(astNode, stop)
	}

	return astNode
}

// configureASTWithToken sets the position information for any ASTNode using a single token
func configureASTWithToken[T ASTNode](astNode T, token antlr.Token) T {
	return configureASTWithTokens(astNode, token, token)
}

// configureASTWithTokens sets the position information for any ASTNode using start and stop tokens
func configureASTWithTokens[T ASTNode](astNode T, start, stop antlr.Token) T {
	astNode.SetLineNumber(start.GetLine())
	astNode.SetColumnNumber(start.GetColumn() + 1)

	astNode.SetLastLineNumber(stop.GetLine())
	astNode.SetLastColumnNumber(stop.GetColumn() + stop.GetStop() - stop.GetStart() + 1)

	return astNode
}

func configureEndPosition[T ASTNode](astNode T, token antlr.Token) {
	endPos := endPosition(token)
	astNode.SetLastLineNumber(endPos.V1)
	astNode.SetLastColumnNumber(endPos.V2)
}

func configureASTFromSource[T ASTNode](astNode T, source ASTNode) T {
	astNode.SetLineNumber(source.GetLineNumber())
	astNode.SetColumnNumber(source.GetColumnNumber())
	astNode.SetLastLineNumber(source.GetLastLineNumber())
	astNode.SetLastColumnNumber(source.GetLastColumnNumber())

	return astNode
}

// endPosition calculates the end position of a token
func endPosition(token antlr.Token) Tuple2[int, int] {
	stopText := token.GetText()
	stopTextLength := len(stopText)
	newLineCnt := strings.Count(stopText, "\n")

	if newLineCnt == 0 {
		return Tuple2[int, int]{
			V1: token.GetLine(),
			V2: token.GetColumn() + 1 + stopTextLength,
		}
	} else {
		return Tuple2[int, int]{
			V1: token.GetLine() + newLineCnt,
			V2: stopTextLength - strings.LastIndex(stopText, "\n"),
		}
	}
}

// Tuple2 is a simple tuple type for returning two values
type Tuple2[T1, T2 any] struct {
	V1 T1
	V2 T2
}

func NewTuple2[T1, T2 any](v1 T1, v2 T2) Tuple2[T1, T2] {
	return Tuple2[T1, T2]{V1: v1, V2: v2}
}

type Tuple3[T1, T2, T3 any] struct {
	V1 T1
	V2 T2
	V3 T3
}

func NewTuple3[T1, T2, T3 any](v1 T1, v2 T2, v3 T3) Tuple3[T1, T2, T3] {
	return Tuple3[T1, T2, T3]{V1: v1, V2: v2, V3: v3}
}
