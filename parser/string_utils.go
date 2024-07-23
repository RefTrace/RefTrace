package parser

import (
	"regexp"
	"strconv"
	"strings"
)

const (
	BACKSLASH       = "\\"
	NONE_SLASHY     = 0
	SLASHY          = 1
	DOLLAR_SLASHY   = 2
	INDEX_NOT_FOUND = -1
)

var (
	HEX_ESCAPES_PATTERN      = regexp.MustCompile(`(\\*)\\u([0-9a-fA-F]{4})`)
	OCTAL_ESCAPES_PATTERN    = regexp.MustCompile(`(\\*)\\([0-3]?[0-7]?[0-7])`)
	STANDARD_ESCAPES_PATTERN = regexp.MustCompile(`(\\*)\\([btnfrs"'])`)
	LINE_ESCAPE_PATTERN      = regexp.MustCompile(`(\\*)\\\r?\n`)

	STANDARD_ESCAPES = map[rune]rune{
		'b': '\b',
		't': '\t',
		'n': '\n',
		'f': '\f',
		'r': '\r',
		's': ' ',
	}
)

func ReplaceHexEscapes(text string) string {
	if !strings.Contains(text, BACKSLASH) {
		return text
	}

	return HEX_ESCAPES_PATTERN.ReplaceAllStringFunc(text, func(match string) string {
		parts := HEX_ESCAPES_PATTERN.FindStringSubmatch(match)
		if isLengthOdd(parts[1]) {
			return match
		}
		code, _ := strconv.ParseInt(parts[2], 16, 32)
		return parts[1] + string(rune(code))
	})
}

func ReplaceOctalEscapes(text string) string {
	if !strings.Contains(text, BACKSLASH) {
		return text
	}

	return OCTAL_ESCAPES_PATTERN.ReplaceAllStringFunc(text, func(match string) string {
		parts := OCTAL_ESCAPES_PATTERN.FindStringSubmatch(match)
		if isLengthOdd(parts[1]) {
			return match
		}
		code, _ := strconv.ParseInt(parts[2], 8, 32)
		return parts[1] + string(rune(code))
	})
}

func ReplaceStandardEscapes(text string) string {
	if !strings.Contains(text, BACKSLASH) {
		return text
	}

	result := STANDARD_ESCAPES_PATTERN.ReplaceAllStringFunc(text, func(match string) string {
		parts := STANDARD_ESCAPES_PATTERN.FindStringSubmatch(match)
		if isLengthOdd(parts[1]) {
			return match
		}
		char, ok := STANDARD_ESCAPES[rune(parts[2][0])]
		if ok {
			return parts[1] + string(char)
		}
		return parts[1] + parts[2]
	})

	return strings.ReplaceAll(result, "\\\\", "\\")
}

func ReplaceEscapes(text string, slashyType int) string {
	switch slashyType {
	case SLASHY, DOLLAR_SLASHY:
		text = ReplaceHexEscapes(text)
		text = ReplaceLineEscape(text)

		if slashyType == SLASHY {
			text = strings.ReplaceAll(text, "\\/", "/")
		}

		if slashyType == DOLLAR_SLASHY {
			text = strings.ReplaceAll(text, "$/", "/")
			text = strings.ReplaceAll(text, "$$", "$")
		}

	case NONE_SLASHY:
		text = replaceEscapes(text)

	default:
		panic("Invalid slashyType")
	}

	return text
}

func replaceEscapes(text string) string {
	if !strings.Contains(text, BACKSLASH) {
		return text
	}

	text = strings.ReplaceAll(text, "\\$", "$")
	text = ReplaceLineEscape(text)
	return ReplaceStandardEscapes(ReplaceHexEscapes(ReplaceOctalEscapes(text)))
}

func ReplaceLineEscape(text string) string {
	if !strings.Contains(text, BACKSLASH) {
		return text
	}

	return LINE_ESCAPE_PATTERN.ReplaceAllStringFunc(text, func(match string) string {
		parts := LINE_ESCAPE_PATTERN.FindStringSubmatch(match)
		if isLengthOdd(parts[1]) {
			return match
		}
		return parts[1]
	})
}

func isLengthOdd(str string) bool {
	return str != "" && len(str)%2 == 1
}

func RemoveCR(text string) string {
	return strings.ReplaceAll(text, "\r\n", "\n")
}

func CountChar(text string, c rune) int {
	count := 0
	for _, r := range text {
		if r == c {
			count++
		}
	}
	return count
}

func TrimQuotations(text string, quotationLength int) string {
	length := len(text)
	if length == quotationLength<<1 {
		return ""
	}
	return text[quotationLength : length-quotationLength]
}

func Matches(text string, pattern *regexp.Regexp) bool {
	return pattern.MatchString(text)
}

func Replace(text, searchString, replacement string) string {
	if IsEmpty(text) || IsEmpty(searchString) || replacement == "" {
		return text
	}

	var buf strings.Builder
	start := 0
	for {
		end := strings.Index(text[start:], searchString)
		if end == INDEX_NOT_FOUND {
			break
		}
		end += start
		buf.WriteString(text[start:end])
		buf.WriteString(replacement)
		start = end + len(searchString)
	}
	buf.WriteString(text[start:])
	return buf.String()
}

func IsEmpty(cs string) bool {
	return cs == ""
}
