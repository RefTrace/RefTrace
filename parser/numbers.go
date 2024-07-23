package parser

import (
	"math"
	"math/big"
	"strings"
	"unicode"
)

// IsDigit returns true if the specified character is a base-10 digit.
func IsDigit(c rune) bool {
	return c >= '0' && c <= '9'
}

// IsOctalDigit returns true if the specific character is a base-8 digit.
func IsOctalDigit(c rune) bool {
	return c >= '0' && c <= '7'
}

// IsHexDigit returns true if the specified character is a base-16 digit.
func IsHexDigit(c rune) bool {
	return IsDigit(c) || (c >= 'A' && c <= 'F') || (c >= 'a' && c <= 'f')
}

// IsNumericTypeSpecifier returns true if the specified character is a valid type specifier
// for a numeric value.
func IsNumericTypeSpecifier(c rune, isDecimal bool) bool {
	if isDecimal {
		switch c {
		case 'G', 'g', 'D', 'd', 'F', 'f':
			return true
		}
	} else {
		switch c {
		case 'G', 'g', 'I', 'i', 'L', 'l':
			return true
		}
	}
	return false
}

// ParseInteger builds a Number from the given integer descriptor.
func ParseInteger(text string) interface{} {
	text = strings.ReplaceAll(text, "_", "")

	negative := false
	if text[0] == '-' || text[0] == '+' {
		negative = text[0] == '-'
		text = text[1:]
	}

	radix := 10
	if len(text) > 1 && text[0] == '0' {
		switch text[1] {
		case 'X', 'x':
			radix = 16
			text = text[2:]
		case 'B', 'b':
			radix = 2
			text = text[2:]
		default:
			radix = 8
		}
	}

	typeSpecifier := 'x'
	if IsNumericTypeSpecifier(rune(text[len(text)-1]), false) {
		typeSpecifier = unicode.ToLower(rune(text[len(text)-1]))
		text = text[:len(text)-1]
	}

	if negative {
		text = "-" + text
	}

	value, _ := new(big.Int).SetString(text, radix)

	switch typeSpecifier {
	case 'i':
		if radix == 10 && (value.Cmp(big.NewInt(math.MaxInt64)) > 0 || value.Cmp(big.NewInt(math.MinInt64)) < 0) {
			panic("Number out of int range")
		}
		return int(value.Int64())
	case 'l':
		if radix == 10 && (value.Cmp(big.NewInt(math.MaxInt64)) > 0 || value.Cmp(big.NewInt(math.MinInt64)) < 0) {
			panic("Number out of int64 range")
		}
		return value.Int64()
	case 'g':
		return value
	default:
		if value.IsInt64() {
			if value.Int64() >= math.MinInt32 && value.Int64() <= math.MaxInt32 {
				return int(value.Int64())
			}
			return value.Int64()
		}
		return value
	}
}

// ParseDecimal builds a Number from the given decimal descriptor.
func ParseDecimal(text string) interface{} {
	text = strings.ReplaceAll(text, "_", "")

	typeSpecifier := 'x'
	if IsNumericTypeSpecifier(rune(text[len(text)-1]), true) {
		typeSpecifier = unicode.ToLower(rune(text[len(text)-1]))
		text = text[:len(text)-1]
	}

	value, _ := new(big.Float).SetString(text)

	switch typeSpecifier {
	case 'f':
		f, _ := value.Float64()
		if f >= -math.MaxFloat32 && f <= math.MaxFloat32 {
			return float32(f)
		}
		panic("Number out of float32 range")
	case 'd':
		f, _ := value.Float64()
		return f
	case 'g':
		fallthrough
	default:
		return value
	}
}
