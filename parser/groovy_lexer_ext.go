package parser

import (
	"unicode"
	"unicode/utf16"
)

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
