package parser

import (
	"fmt"
	"strings"
)

// Constants for token types
const (
	EOF     = -1
	UNKNOWN = 0
	NEWLINE = 5

	LEFT_CURLY_BRACE     = 10
	RIGHT_CURLY_BRACE    = 20
	LEFT_SQUARE_BRACKET  = 30
	RIGHT_SQUARE_BRACKET = 40
	LEFT_PARENTHESIS     = 50
	RIGHT_PARENTHESIS    = 60

	DOT         = 70
	DOT_DOT     = 75
	DOT_DOT_DOT = 77

	NAVIGATE = 80

	FIND_REGEX    = 90
	MATCH_REGEX   = 94
	REGEX_PATTERN = 97
	IMPLIES       = 99

	EQUAL  = 100
	EQUALS = EQUAL
	ASSIGN = EQUAL

	COMPARE_NOT_EQUAL          = 120
	COMPARE_IDENTICAL          = 121
	COMPARE_NOT_IDENTICAL      = 122
	COMPARE_EQUAL              = 123
	COMPARE_LESS_THAN          = 124
	COMPARE_LESS_THAN_EQUAL    = 125
	COMPARE_GREATER_THAN       = 126
	COMPARE_GREATER_THAN_EQUAL = 127
	COMPARE_TO                 = 128
	COMPARE_NOT_IN             = 129
	COMPARE_NOT_INSTANCEOF     = 130

	NOT         = 160
	LOGICAL_OR  = 162
	LOGICAL_AND = 164

	LOGICAL_OR_EQUAL  = 166
	LOGICAL_AND_EQUAL = 168

	PLUS      = 200
	MINUS     = 201
	MULTIPLY  = 202
	DIVIDE    = 203
	INTDIV    = 204
	MOD       = 205
	STAR_STAR = 206
	POWER     = STAR_STAR

	PLUS_EQUAL     = 210
	MINUS_EQUAL    = 211
	MULTIPLY_EQUAL = 212
	DIVIDE_EQUAL   = 213
	INTDIV_EQUAL   = 214
	MOD_EQUAL      = 215
	POWER_EQUAL    = 216
	ELVIS_EQUAL    = 217

	PLUS_PLUS         = 250
	PREFIX_PLUS_PLUS  = 251
	POSTFIX_PLUS_PLUS = 252
	PREFIX_PLUS       = 253

	MINUS_MINUS         = 260
	PREFIX_MINUS_MINUS  = 261
	POSTFIX_MINUS_MINUS = 262
	PREFIX_MINUS        = 263

	LEFT_SHIFT           = 280
	RIGHT_SHIFT          = 281
	RIGHT_SHIFT_UNSIGNED = 282

	LEFT_SHIFT_EQUAL           = 285
	RIGHT_SHIFT_EQUAL          = 286
	RIGHT_SHIFT_UNSIGNED_EQUAL = 287

	STAR = MULTIPLY

	COMMA     = 300
	COLON     = 310
	SEMICOLON = 320
	QUESTION  = 330

	PIPE        = 340
	DOUBLE_PIPE = LOGICAL_OR
	BITWISE_OR  = PIPE
	BITWISE_AND = 341
	BITWISE_XOR = 342

	BITWISE_OR_EQUAL  = 350
	BITWISE_AND_EQUAL = 351
	BITWISE_XOR_EQUAL = 352
	BITWISE_NEGATION  = REGEX_PATTERN
	REMAINDER         = 353
	REMAINDER_EQUAL   = 354

	STRING = 400

	IDENTIFIER = 440

	INTEGER_NUMBER = 450
	DECIMAL_NUMBER = 451

	KEYWORD_PRIVATE   = 500
	KEYWORD_PROTECTED = 501
	KEYWORD_PUBLIC    = 502

	KEYWORD_ABSTRACT  = 510
	KEYWORD_FINAL     = 511
	KEYWORD_NATIVE    = 512
	KEYWORD_TRANSIENT = 513
	KEYWORD_VOLATILE  = 514

	KEYWORD_SYNCHRONIZED = 520
	KEYWORD_STATIC       = 521

	KEYWORD_DEF       = 530
	KEYWORD_DEFMACRO  = 539
	KEYWORD_CLASS     = 531
	KEYWORD_INTERFACE = 532
	KEYWORD_MIXIN     = 533

	KEYWORD_IMPLEMENTS = 540
	KEYWORD_EXTENDS    = 541
	KEYWORD_THIS       = 542
	KEYWORD_SUPER      = 543
	KEYWORD_INSTANCEOF = 544
	KEYWORD_PROPERTY   = 545
	KEYWORD_NEW        = 546

	KEYWORD_PACKAGE = 550
	KEYWORD_IMPORT  = 551
	KEYWORD_AS      = 552

	KEYWORD_RETURN   = 560
	KEYWORD_IF       = 561
	KEYWORD_ELSE     = 562
	KEYWORD_DO       = 570
	KEYWORD_WHILE    = 571
	KEYWORD_FOR      = 572
	KEYWORD_IN       = 573
	KEYWORD_BREAK    = 574
	KEYWORD_CONTINUE = 575
	KEYWORD_SWITCH   = 576
	KEYWORD_CASE     = 577
	KEYWORD_DEFAULT  = 578

	KEYWORD_TRY     = 580
	KEYWORD_CATCH   = 581
	KEYWORD_FINALLY = 582
	KEYWORD_THROW   = 583
	KEYWORD_THROWS  = 584
	KEYWORD_ASSERT  = 585

	KEYWORD_VOID    = 600
	KEYWORD_BOOLEAN = 601
	KEYWORD_BYTE    = 602
	KEYWORD_SHORT   = 603
	KEYWORD_INT     = 604
	KEYWORD_LONG    = 605
	KEYWORD_FLOAT   = 606
	KEYWORD_DOUBLE  = 607
	KEYWORD_CHAR    = 608

	KEYWORD_TRUE  = 610
	KEYWORD_FALSE = 611
	KEYWORD_NULL  = 612

	KEYWORD_CONST = 700
	KEYWORD_GOTO  = 701

	SYNTH_COMPILATION_UNIT = 800

	SYNTH_CLASS                 = 801
	SYNTH_INTERFACE             = 802
	SYNTH_MIXIN                 = 803
	SYNTH_METHOD                = 804
	SYNTH_PROPERTY              = 805
	SYNTH_PARAMETER_DECLARATION = 806

	SYNTH_LIST    = 810
	SYNTH_MAP     = 811
	SYNTH_GSTRING = 812

	SYNTH_METHOD_CALL = 814
	SYNTH_CAST        = 815
	SYNTH_BLOCK       = 816
	SYNTH_CLOSURE     = 817
	SYNTH_LABEL       = 818
	SYNTH_TERNARY     = 819
	SYNTH_TUPLE       = 820

	SYNTH_VARIABLE_DECLARATION = 830

	GSTRING_START            = 901
	GSTRING_END              = 902
	GSTRING_EXPRESSION_START = 903
	GSTRING_EXPRESSION_END   = 904

	ANY                      = 1000
	NOT_EOF                  = 1001
	GENERAL_END_OF_STATEMENT = 1002
	ANY_END_OF_STATEMENT     = 1003

	ASSIGNMENT_OPERATOR       = 1100
	COMPARISON_OPERATOR       = 1101
	MATH_OPERATOR             = 1102
	LOGICAL_OPERATOR          = 1103
	RANGE_OPERATOR            = 1104
	REGEX_COMPARISON_OPERATOR = 1105
	DEREFERENCE_OPERATOR      = 1106
	BITWISE_OPERATOR          = 1107
	INSTANCEOF_OPERATOR       = 1108

	PREFIX_OPERATOR          = 1200
	POSTFIX_OPERATOR         = 1210
	INFIX_OPERATOR           = 1220
	PREFIX_OR_INFIX_OPERATOR = 1230
	PURE_PREFIX_OPERATOR     = 1235

	KEYWORD                  = 1300
	SYMBOL                   = 1301
	LITERAL                  = 1310
	NUMBER                   = 1320
	SIGN                     = 1325
	NAMED_VALUE              = 1330
	TRUTH_VALUE              = 1331
	PRIMITIVE_TYPE           = 1340
	CREATABLE_PRIMITIVE_TYPE = 1341
	LOOP                     = 1350
	RESERVED_KEYWORD         = 1360
	KEYWORD_IDENTIFIER       = 1361
	SYNTHETIC                = 1370

	TYPE_DECLARATION     = 1400
	DECLARATION_MODIFIER = 1410

	TYPE_NAME           = 1420
	CREATABLE_TYPE_NAME = 1430

	MATCHED_CONTAINER          = 1500
	LEFT_OF_MATCHED_CONTAINER  = 1501
	RIGHT_OF_MATCHED_CONTAINER = 1502

	EXPRESSION = 1900

	OPERATOR_EXPRESSION = 1901
	SYNTH_EXPRESSION    = 1902
	KEYWORD_EXPRESSION  = 1903
	LITERAL_EXPRESSION  = 1904
	ARRAY_EXPRESSION    = 1905

	SIMPLE_EXPRESSION  = 1910
	COMPLEX_EXPRESSION = 1911

	PARAMETER_TERMINATORS       = 2000
	ARRAY_ITEM_TERMINATORS      = 2001
	TYPE_LIST_TERMINATORS       = 2002
	OPTIONAL_DATATYPE_FOLLOWERS = 2003

	SWITCH_BLOCK_TERMINATORS = 2004
	SWITCH_ENTRIES           = 2005

	METHOD_CALL_STARTERS = 2006
	UNSAFE_OVER_NEWLINES = 2007

	PRECLUDES_CAST_OPERATOR = 2008
)

var (
	texts        = make(map[int]string)
	lookup       = make(map[string]int)
	keywords     = make(map[string]bool)
	descriptions = make(map[int]string)
)

func init() {
	// Initialize maps with token types, texts, and descriptions
	addTranslation("\n", NEWLINE)
	addTranslation("{", LEFT_CURLY_BRACE)
	addTranslation("}", RIGHT_CURLY_BRACE)
	// ... (other translations)

	addKeyword("abstract", KEYWORD_ABSTRACT)
	addKeyword("as", KEYWORD_AS)
	// ... (other keywords)

	addDescription(NEWLINE, "<newline>")
	addDescription(PREFIX_PLUS_PLUS, "<prefix ++>")
	// ... (other descriptions)
}

func addTranslation(text string, tokenType int) {
	texts[tokenType] = text
	lookup[text] = tokenType
}

func addKeyword(text string, tokenType int) {
	keywords[text] = true
	addTranslation(text, tokenType)
}

func addDescription(tokenType int, description string) {
	if strings.HasPrefix(description, "<") && strings.HasSuffix(description, ">") {
		descriptions[tokenType] = description
	} else {
		descriptions[tokenType] = fmt.Sprintf("\"%s\"", description)
	}
}

func GetKeywords() []string {
	keys := make([]string, 0, len(keywords))
	for k := range keywords {
		keys = append(keys, k)
	}
	return keys
}

func IsKeyword(text string) bool {
	_, ok := keywords[text]
	return ok
}

func Lookup(text string, filter int) int {
	if tokenType, ok := lookup[text]; ok {
		if filter == UNKNOWN || OfType(tokenType, filter) {
			return tokenType
		}
	}
	return UNKNOWN
}

func LookupKeyword(text string) int {
	return Lookup(text, KEYWORD)
}

func LookupSymbol(text string) int {
	return Lookup(text, SYMBOL)
}

func GetText(tokenType int) string {
	if text, ok := texts[tokenType]; ok {
		return text
	}
	return ""
}

func GetDescription(tokenType int) string {
	if desc, ok := descriptions[tokenType]; ok {
		return desc
	}
	return "<>"
}

func OfType(specific, general int) bool {
	if general == specific {
		return true
	}

	switch general {
	case ANY:
		return true

	case NOT_EOF:
		return specific >= UNKNOWN && specific <= SYNTH_VARIABLE_DECLARATION

	case GENERAL_END_OF_STATEMENT:
		switch specific {
		case EOF, NEWLINE, SEMICOLON:
			return true
		}

	case ANY_END_OF_STATEMENT:
		switch specific {
		case EOF, NEWLINE, SEMICOLON, RIGHT_CURLY_BRACE:
			return true
		}

	case ASSIGNMENT_OPERATOR:
		return specific == EQUAL || (specific >= PLUS_EQUAL && specific <= ELVIS_EQUAL) ||
			(specific >= LOGICAL_OR_EQUAL && specific <= LOGICAL_AND_EQUAL) ||
			(specific >= LEFT_SHIFT_EQUAL && specific <= RIGHT_SHIFT_UNSIGNED_EQUAL) ||
			(specific >= BITWISE_OR_EQUAL && specific <= BITWISE_XOR_EQUAL)

	case COMPARISON_OPERATOR:
		return specific >= COMPARE_NOT_EQUAL && specific <= COMPARE_TO

	case INSTANCEOF_OPERATOR:
		return specific == KEYWORD_INSTANCEOF || specific == COMPARE_NOT_INSTANCEOF

	case MATH_OPERATOR:
		return (specific >= PLUS && specific <= RIGHT_SHIFT_UNSIGNED) ||
			(specific >= NOT && specific <= LOGICAL_AND) ||
			(specific >= BITWISE_OR && specific <= BITWISE_XOR)

	case LOGICAL_OPERATOR:
		return specific >= NOT && specific <= LOGICAL_AND

	case BITWISE_OPERATOR:
		return (specific >= BITWISE_OR && specific <= BITWISE_XOR) || specific == BITWISE_NEGATION

	case RANGE_OPERATOR:
		return specific == DOT_DOT || specific == DOT_DOT_DOT

	case REGEX_COMPARISON_OPERATOR:
		return specific == FIND_REGEX || specific == MATCH_REGEX

	case DEREFERENCE_OPERATOR:
		return specific == DOT || specific == NAVIGATE

	case PREFIX_OPERATOR:
		switch specific {
		case MINUS, PLUS_PLUS, MINUS_MINUS:
			return true
		}
		fallthrough

	case PURE_PREFIX_OPERATOR:
		switch specific {
		case REGEX_PATTERN, NOT, PREFIX_PLUS, PREFIX_PLUS_PLUS, PREFIX_MINUS, PREFIX_MINUS_MINUS, SYNTH_CAST:
			return true
		}

	case POSTFIX_OPERATOR:
		switch specific {
		case PLUS_PLUS, POSTFIX_PLUS_PLUS, MINUS_MINUS, POSTFIX_MINUS_MINUS:
			return true
		}

	case INFIX_OPERATOR:
		switch specific {
		case DOT, NAVIGATE, LOGICAL_OR, LOGICAL_AND, BITWISE_OR, BITWISE_AND, BITWISE_XOR,
			LEFT_SHIFT, RIGHT_SHIFT, RIGHT_SHIFT_UNSIGNED, FIND_REGEX, MATCH_REGEX,
			DOT_DOT, DOT_DOT_DOT, KEYWORD_INSTANCEOF:
			return true
		}
		return (specific >= COMPARE_NOT_EQUAL && specific <= COMPARE_TO) ||
			(specific >= PLUS && specific <= MOD_EQUAL) ||
			specific == EQUAL ||
			(specific >= PLUS_EQUAL && specific <= ELVIS_EQUAL) ||
			(specific >= LOGICAL_OR_EQUAL && specific <= LOGICAL_AND_EQUAL) ||
			(specific >= LEFT_SHIFT_EQUAL && specific <= RIGHT_SHIFT_UNSIGNED_EQUAL) ||
			(specific >= BITWISE_OR_EQUAL && specific <= BITWISE_XOR_EQUAL)

	case PREFIX_OR_INFIX_OPERATOR:
		switch specific {
		case POWER, PLUS, MINUS, PREFIX_PLUS, PREFIX_MINUS:
			return true
		}

		// ... (continue with the rest of the cases)

	}

	return false
}

func GetPrecedence(tokenType int, throwIfInvalid bool) int {
	switch tokenType {
	case LEFT_PARENTHESIS:
		return 0

	case EQUAL, PLUS_EQUAL, MINUS_EQUAL, MULTIPLY_EQUAL, DIVIDE_EQUAL,
		INTDIV_EQUAL, MOD_EQUAL, POWER_EQUAL, ELVIS_EQUAL,
		LOGICAL_OR_EQUAL, LOGICAL_AND_EQUAL,
		LEFT_SHIFT_EQUAL, RIGHT_SHIFT_EQUAL, RIGHT_SHIFT_UNSIGNED_EQUAL,
		BITWISE_OR_EQUAL, BITWISE_AND_EQUAL, BITWISE_XOR_EQUAL, REMAINDER_EQUAL:
		return 5

	case QUESTION:
		return 10

	case IMPLIES:
		return 12

	case LOGICAL_OR:
		return 15

	case LOGICAL_AND:
		return 20

	case BITWISE_OR, BITWISE_AND, BITWISE_XOR:
		return 22

	case COMPARE_IDENTICAL, COMPARE_NOT_IDENTICAL:
		return 24

	case COMPARE_NOT_EQUAL, COMPARE_EQUAL, COMPARE_LESS_THAN, COMPARE_LESS_THAN_EQUAL,
		COMPARE_GREATER_THAN, COMPARE_GREATER_THAN_EQUAL, COMPARE_TO,
		FIND_REGEX, MATCH_REGEX, KEYWORD_INSTANCEOF, COMPARE_NOT_INSTANCEOF:
		return 25

	case DOT_DOT, DOT_DOT_DOT:
		return 30

	case LEFT_SHIFT, RIGHT_SHIFT, RIGHT_SHIFT_UNSIGNED:
		return 35

	case PLUS, MINUS:
		return 40

	case MULTIPLY, DIVIDE, INTDIV, MOD, REMAINDER:
		return 45

	case NOT, REGEX_PATTERN:
		return 50

	case SYNTH_CAST:
		return 55

	case PLUS_PLUS, MINUS_MINUS, PREFIX_PLUS_PLUS, PREFIX_MINUS_MINUS,
		POSTFIX_PLUS_PLUS, POSTFIX_MINUS_MINUS:
		return 65

	case PREFIX_PLUS, PREFIX_MINUS:
		return 70

	case POWER:
		return 72

	case SYNTH_METHOD, LEFT_SQUARE_BRACKET:
		return 75

	case DOT, NAVIGATE:
		return 80

	case KEYWORD_NEW:
		return 85
	}

	if throwIfInvalid {
		panic("precedence requested for non-operator")
	}

	return -1
}