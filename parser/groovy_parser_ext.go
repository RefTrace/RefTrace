package parser

import (
	"errors"
	"fmt"
	"sort"
	"unicode"

	"github.com/antlr4-go/antlr/v4"
)

var MODIFIER_ARRAY []int

func init() {
	// Initialize MODIFIER_ARRAY based on ModifierOpcodeMap keys
	for k := range ModifierOpcodeMap {
		MODIFIER_ARRAY = append(MODIFIER_ARRAY, k)
	}

	// Sort the array
	sort.Ints(MODIFIER_ARRAY)
}

type GroovyParserRuleContext struct {
	antlr.BaseParserRuleContext
	metaDataMap map[interface{}]interface{}
}

// NewGroovyParserRuleContext creates a new GroovyParserRuleContext
func NewGroovyParserRuleContext(parent antlr.ParserRuleContext, invokingStateNumber int) *GroovyParserRuleContext {
	return &GroovyParserRuleContext{
		BaseParserRuleContext: *antlr.NewBaseParserRuleContext(parent, invokingStateNumber),
		metaDataMap:           make(map[interface{}]interface{}),
	}
}

func (g *GroovyParserRuleContext) GetMetaDataMap() map[interface{}]interface{} {
	return g.metaDataMap
}

func (g *GroovyParserRuleContext) SetMetaDataMap(metaDataMap map[interface{}]interface{}) {
	g.metaDataMap = metaDataMap
}

func (g *GroovyParserRuleContext) NewMetaDataMap() map[interface{}]interface{} {
	return make(map[interface{}]interface{})
}

func (g *GroovyParserRuleContext) GetNodeMetaData(key interface{}) interface{} {
	if g.metaDataMap == nil {
		return nil
	}
	return g.metaDataMap[key]
}

func (g *GroovyParserRuleContext) SetNodeMetaData(key, value interface{}) {
	if old := g.PutNodeMetaData(key, value); old != nil {
		panic(errors.New("Tried to overwrite existing meta data"))
	}
}

func (g *GroovyParserRuleContext) GetNodeMetaDataWithFunc(key interface{}, valFn func() interface{}) interface{} {
	if key == nil {
		panic(errors.New("Tried to get/set meta data with null key"))
	}

	if g.metaDataMap == nil {
		g.metaDataMap = g.NewMetaDataMap()
		g.SetMetaDataMap(g.metaDataMap)
	}
	if val, ok := g.metaDataMap[key]; ok {
		return val
	}
	val := valFn()
	g.metaDataMap[key] = val
	return val
}

func (g *GroovyParserRuleContext) CopyNodeMetaData(other NodeMetaDataHandler) {
	otherMetaDataMap := other.GetMetaDataMap()
	if otherMetaDataMap == nil {
		return
	}
	if g.metaDataMap == nil {
		g.metaDataMap = g.NewMetaDataMap()
		g.SetMetaDataMap(g.metaDataMap)
	}
	for k, v := range otherMetaDataMap {
		g.metaDataMap[k] = v
	}
}

func (g *GroovyParserRuleContext) PutNodeMetaData(key, value interface{}) interface{} {
	if key == nil {
		panic(errors.New("Tried to set meta data with null key"))
	}

	if g.metaDataMap == nil {
		if value == nil {
			return nil
		}
		g.metaDataMap = g.NewMetaDataMap()
		g.SetMetaDataMap(g.metaDataMap)
	} else if value == nil {
		return g.metaDataMap[key]
	}
	oldValue := g.metaDataMap[key]
	g.metaDataMap[key] = value
	return oldValue
}

func (g *GroovyParserRuleContext) RemoveNodeMetaData(key interface{}) {
	if key == nil {
		panic(errors.New("Tried to remove meta data with null key"))
	}

	if g.metaDataMap != nil {
		delete(g.metaDataMap, key)
	}
}

func (g *GroovyParserRuleContext) GetNodeMetaDataMap() map[interface{}]interface{} {
	if g.metaDataMap == nil {
		return map[interface{}]interface{}{}
	}
	return g.metaDataMap
}

type MyGroovyParser struct {
	*antlr.BaseParser
	inSwitchExpressionLevel int
}

func isFollowingArgumentsOrClosure(ctx interface{}) bool {

	if postfixExprAltContext, ok := ctx.(*PostfixExprAltContext); ok {
		peacChildren := postfixExprAltContext.GetChildren()

		defer func() {
			if r := recover(); r != nil {
				panic(fmt.Sprintf("Unexpected structure of expression context: %v", ctx))
			}
		}()

		peacChild := peacChildren[0]
		pecChildren := peacChild.(*PostfixExpressionContext).GetChildren()

		pecChild := pecChildren[0]
		pec := pecChild.(*PathExpressionContext)

		t := pec.GetT()

		return t == 2 || t == 3
	}
	/*
		if pathExpressionContext, ok := context.(*PathExpressionContext); ok {
			t := pathExpressionContext.GetT()
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
