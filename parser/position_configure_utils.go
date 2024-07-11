package parser

import (
	"strings"

	"github.com/antlr4-go/antlr/v4"
)

// configureAST sets the position information for any ASTNode
func configureAST[T ASTNode](astNode T, ctx antlr.ParserRuleContext) T {
	start := ctx.GetStart()
	stop := ctx.GetStop()

	astNode.SetLineNumber(start.GetLine())
	astNode.SetColumnNumber(start.GetColumn() + 1)

	endPosition := endPosition(stop)
	astNode.SetLastLineNumber(endPosition.V1)
	astNode.SetLastColumnNumber(endPosition.V2)

	return astNode
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
