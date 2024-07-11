package parser

import (
	"strings"
	"unicode"
)

var invalidJavaIdentifiers = map[string]bool{
	"abstract": true, "assert": true, "boolean": true, "break": true, "byte": true,
	"case": true, "catch": true, "char": true, "class": true, "const": true,
	"continue": true, "default": true, "do": true, "double": true, "else": true,
	"enum": true, "extends": true, "final": true, "finally": true, "float": true,
	"for": true, "goto": true, "if": true, "implements": true, "import": true,
	"instanceof": true, "int": true, "interface": true, "long": true, "native": true,
	"new": true, "package": true, "private": true, "protected": true, "public": true,
	"short": true, "static": true, "strictfp": true, "super": true, "switch": true,
	"synchronized": true, "this": true, "throw": true, "throws": true,
	"transient": true, "try": true, "void": true, "volatile": true, "while": true,
	"true": true, "false": true, "null": true,
}

// RepeatString returns a string made up of repetitions of the specified string.
func RepeatString(pattern string, repeats int) string {
	return strings.Repeat(pattern, repeats)
}

// EOL returns the end-of-line marker.
func EOL() string {
	return "\n"
}

// IsJavaIdentifier tells if the given string is a valid Java identifier.
func IsJavaIdentifier(name string) bool {
	if len(name) == 0 || invalidJavaIdentifiers[name] {
		return false
	}
	runes := []rune(name)
	if !unicode.IsLetter(runes[0]) && runes[0] != '_' {
		return false
	}
	for _, r := range runes[1:] {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_' {
			return false
		}
	}
	return true
}
