package parser

import (
	"unicode"

	"github.com/antlr4-go/antlr/v4"
)

type MyGroovyParser struct {
	*antlr.BaseParser
	inSwitchExpressionLevel int
}

func isFollowingArgumentsOrClosure() bool {
	/*
		if postfixExprAltContext, ok := context.(*PostfixExprAltContext); ok {
			peacChildren := postfixExprAltContext.GetChildren()

			defer func() {
				if r := recover(); r != nil {
					panic(fmt.Sprintf("Unexpected structure of expression context: %v", context))
				}
			}()

			peacChild := peacChildren[0]
			pecChildren := peacChild.(*PostfixExpressionContext).GetChildren()

			pecChild := pecChildren[0]
			pec := pecChild.(*PathExpressionContext)

			t := pec.GetT()

			return t == 2 || t == 3
		}
	*/

	return false
}

func isInvalidMethodDeclaration(ts antlr.TokenStream) bool {
	tokenType := ts.LT(1).GetTokenType()

	return (tokenType == GroovyParserIdentifier || tokenType == GroovyParserCapitalizedIdentifier || tokenType == GroovyParserStringLiteral || tokenType == GroovyParserYIELD) &&
		ts.LT(2).GetTokenType() == GroovyParserLPAREN
}

const ANNOTATION_TYPE = -999

var MODIFIER_ARRAY = []int{
	ANNOTATION_TYPE,
	GroovyParserDEF,
	GroovyParserVAR,
	GroovyParserNATIVE,
	GroovyParserSYNCHRONIZED,
	GroovyParserTRANSIENT,
	GroovyParserVOLATILE,
	GroovyParserPUBLIC,
	GroovyParserPROTECTED,
	GroovyParserPRIVATE,
	GroovyParserSTATIC,
	GroovyParserABSTRACT,
	GroovyParserSEALED,
	GroovyParserNON_SEALED,
	GroovyParserFINAL,
	GroovyParserSTRICTFP,
	GroovyParserDEFAULT,
}

func contains(arr []int, item int) bool {
	for _, v := range arr {
		if v == item {
			return true
		}
	}
	return false
}

func isInvalidLocalVariableDeclaration(ts antlr.TokenStream) bool {
	index := 2
	tokenType2 := ts.LT(index).GetTokenType()
	var tokenType3 int

	if tokenType2 == GroovyParserDOT {
		tokeTypeN := tokenType2

		for {
			index += 2
			tokeTypeN = ts.LT(index).GetTokenType()
			if tokeTypeN != GroovyParserDOT {
				break
			}
		}

		if tokeTypeN == GroovyParserLT || tokeTypeN == GroovyParserLBRACK {
			return false
		}

		index--
		tokenType2 = ts.LT(index + 1).GetTokenType()
	} else {
		index = 1
	}

	token := ts.LT(index)
	tokenType := token.GetTokenType()
	tokenType3 = ts.LT(index + 2).GetTokenType()
	nextCodePoint := int([]rune(token.GetText())[0])

	return !(tokenType == GroovyParserBuiltInPrimitiveType || contains(MODIFIER_ARRAY, tokenType)) &&
		!unicode.IsUpper(rune(nextCodePoint)) &&
		nextCodePoint != '@' &&
		!(tokenType3 == GroovyParserASSIGN || (tokenType2 == GroovyParserLT || tokenType2 == GroovyParserLBRACK))
}
