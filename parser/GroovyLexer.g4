/*
 * This file is adapted from the Antlr4 Java grammar which has the following license
 *
 *  Copyright (c) 2013 Terence Parr, Sam Harwell
 *  All rights reserved.
 *  [The "BSD licence"]
 *
 *    http://www.opensource.org/licenses/bsd-license.php
 *
 * Subsequent modifications by the Groovy community have been done under the Apache License v2:
 *
 *  Licensed to the Apache Software Foundation (ASF) under one
 *  or more contributor license agreements.  See the NOTICE file
 *  distributed with this work for additional information
 *  regarding copyright ownership.  The ASF licenses this file
 *  to you under the Apache License, Version 2.0 (the
 *  "License"); you may not use this file except in compliance
 *  with the License.  You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing,
 *  software distributed under the License is distributed on an
 *  "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 *  KIND, either express or implied.  See the License for the
 *  specific language governing permissions and limitations
 *  under the License.
 */

/**
 * The Groovy grammar is based on the official grammar for Java:
 * https://github.com/antlr/grammars-v4/blob/master/java/Java.g4
 */
lexer grammar GroovyLexer;

options {
    language = Go;
    superClass = MyGroovyLexer;
}

/*
@members {
    package parser

    import (
        "github.com/antlr/antlr4/runtime/Go/antlr"
        "sort"
        "unicode"
    )

    type GroovyLexer struct {
        *antlr.BaseLexer
        errorIgnored      bool
        tokenIndex        int64
        lastTokenType     int
        invalidDigitCount int
        parenStack        []Paren
    }

    type Paren struct {
        text          string
        lastTokenType int
        line          int
        column        int
    }

    func NewGroovyLexer(input antlr.CharStream) *GroovyLexer {
        l := new(GroovyLexer)
        l.BaseLexer = antlr.NewBaseLexer(input)
        l.parenStack = make([]Paren, 0, 32)
        return l
    }

    func (l *GroovyLexer) Emit() antlr.Token {
        l.tokenIndex++
        token := l.BaseLexer.Emit()
        tokenType := token.GetTokenType()
        if token.GetChannel() == antlr.TokenDefaultChannel {
            l.lastTokenType = tokenType
        }
        if tokenType == RollBackOne {
            l.rollbackOneChar()
        }
        return token
    }

    var REGEX_CHECK_ARRAY = []int{
        DEC, INC, THIS, RBRACE, RBRACK, RPAREN, GStringEnd, NullLiteral,
        StringLiteral, BooleanLiteral, IntegerLiteral, FloatingPointLiteral,
        Identifier, CapitalizedIdentifier,
    }

    func init() {
        sort.Ints(REGEX_CHECK_ARRAY)
    }

    func (l *GroovyLexer) isRegexAllowed() bool {
        return sort.SearchInts(REGEX_CHECK_ARRAY, l.lastTokenType) < 0
    }

    func (l *GroovyLexer) rollbackOneChar() {
        // This method is intended to be overridden
    }

    func (l *GroovyLexer) enterParen() {
        text := l.GetText()
        l.enterParenCallback(text)
        l.parenStack = append(l.parenStack, Paren{text, l.lastTokenType, l.GetLine(), l.GetCharPositionInLine()})
    }

    func (l *GroovyLexer) exitParen() {
        text := l.GetText()
        l.exitParenCallback(text)
        if len(l.parenStack) > 0 {
            l.parenStack = l.parenStack[:len(l.parenStack)-1]
        }
    }

    func (l *GroovyLexer) isInsideParens() bool {
        if len(l.parenStack) == 0 {
            return false
        }
        paren := l.parenStack[len(l.parenStack)-1]
        text := paren.text
        return (text == "(" && paren.lastTokenType != TRY) || text == "[" || text == "?["
    }

    func (l *GroovyLexer) ignoreTokenInsideParens() {
        if !l.isInsideParens() {
            return
        }
        l.SetChannel(antlr.TokenHiddenChannel)
    }

    func (l *GroovyLexer) ignoreMultiLineCommentConditionally() {
        if !l.isInsideParens() && l.isFollowedByWhiteSpaces() {
            return
        }
        l.SetChannel(antlr.TokenHiddenChannel)
    }

    func (l *GroovyLexer) GetSyntaxErrorSource() int {
        return GroovySyntaxError_LEXER
    }

    func (l *GroovyLexer) GetErrorLine() int {
        return l.GetLine()
    }

    func (l *GroovyLexer) GetErrorColumn() int {
        return l.GetCharPositionInLine() + 1
    }

    func (l *GroovyLexer) PopMode() int {
        defer func() {
            if r := recover(); r != nil {
                // Handle EmptyStackException
            }
        }()
        return l.BaseLexer.PopMode()
    }

    func (l *GroovyLexer) addComment(_type int) {
        text := l.GetInputStream().GetText(antlr.NewInterval(l.GetTokenStartCharIndex(), l.GetCharIndex()-1))
        // Handle the comment text as needed
    }

    func isJavaIdentifierStartAndNotIdentifierIgnorable(codePoint rune) bool {
        return unicode.IsLetter(codePoint) && !unicode.Is(unicode.Cf, codePoint)
    }

    func isJavaIdentifierPartAndNotIdentifierIgnorable(codePoint rune) bool {
        return unicode.IsLetter(codePoint) || unicode.IsDigit(codePoint) && !unicode.Is(unicode.Cf, codePoint)
    }

    func (l *GroovyLexer) IsErrorIgnored() bool {
        return l.errorIgnored
    }

    func (l *GroovyLexer) SetErrorIgnored(errorIgnored bool) {
        l.errorIgnored = errorIgnored
    }

    func (l *GroovyLexer) enterParenCallback(text string) {
        // This method is intended to be overridden
    }

    func (l *GroovyLexer) exitParenCallback(text string) {
        // This method is intended to be overridden
    }

    func (l *GroovyLexer) isFollowedByWhiteSpaces() bool {
        // Implement this method based on your requirements
        return false
    }
}
*/


// §3.10.5 String Literals
StringLiteral
    :   GStringQuotationMark  DqStringCharacter*  GStringQuotationMark
    |   SqStringQuotationMark  SqStringCharacter*  SqStringQuotationMark
    |   Slash { this.isRegexAllowed() && _input.LA(1) != '*' }?  SlashyStringCharacter+  Slash

    |   TdqStringQuotationMark  TdqStringCharacter*  TdqStringQuotationMark
    |   TsqStringQuotationMark  TsqStringCharacter*  TsqStringQuotationMark
    |   DollarSlashyGStringQuotationMarkBegin  DollarSlashyStringCharacter+  DollarSlashyGStringQuotationMarkEnd
    ;

GStringBegin
    :   GStringQuotationMark DqStringCharacter* Dollar -> pushMode(DQ_GSTRING_MODE), pushMode(GSTRING_TYPE_SELECTOR_MODE)
    ;
TdqGStringBegin
    :   TdqStringQuotationMark   TdqStringCharacter* Dollar -> type(GStringBegin), pushMode(TDQ_GSTRING_MODE), pushMode(GSTRING_TYPE_SELECTOR_MODE)
    ;
SlashyGStringBegin
    :   Slash { this.isRegexAllowed() && _input.LA(1) != '*' }? SlashyStringCharacter* Dollar { isFollowedByJavaLetterInGString(_input) }? -> type(GStringBegin), pushMode(SLASHY_GSTRING_MODE), pushMode(GSTRING_TYPE_SELECTOR_MODE)
    ;
DollarSlashyGStringBegin
    :   DollarSlashyGStringQuotationMarkBegin DollarSlashyStringCharacter* Dollar { isFollowedByJavaLetterInGString(_input) }? -> type(GStringBegin), pushMode(DOLLAR_SLASHY_GSTRING_MODE), pushMode(GSTRING_TYPE_SELECTOR_MODE)
    ;

mode DQ_GSTRING_MODE;
GStringEnd
    :   GStringQuotationMark     -> popMode
    ;
GStringPart
    :   Dollar  -> pushMode(GSTRING_TYPE_SELECTOR_MODE)
    ;
GStringCharacter
    :   DqStringCharacter -> more
    ;

mode TDQ_GSTRING_MODE;
TdqGStringEnd
    :   TdqStringQuotationMark    -> type(GStringEnd), popMode
    ;
TdqGStringPart
    :   Dollar   -> type(GStringPart), pushMode(GSTRING_TYPE_SELECTOR_MODE)
    ;
TdqGStringCharacter
    :   TdqStringCharacter -> more
    ;

mode SLASHY_GSTRING_MODE;
SlashyGStringEnd
    :   Dollar? Slash  -> type(GStringEnd), popMode
    ;
SlashyGStringPart
    :   Dollar { isFollowedByJavaLetterInGString(_input) }?   -> type(GStringPart), pushMode(GSTRING_TYPE_SELECTOR_MODE)
    ;
SlashyGStringCharacter
    :   SlashyStringCharacter -> more
    ;

mode DOLLAR_SLASHY_GSTRING_MODE;
DollarSlashyGStringEnd
    :   DollarSlashyGStringQuotationMarkEnd      -> type(GStringEnd), popMode
    ;
DollarSlashyGStringPart
    :   Dollar { isFollowedByJavaLetterInGString(_input) }?   -> type(GStringPart), pushMode(GSTRING_TYPE_SELECTOR_MODE)
    ;
DollarSlashyGStringCharacter
    :   DollarSlashyStringCharacter -> more
    ;

mode GSTRING_TYPE_SELECTOR_MODE;
GStringLBrace
    :   '{' { l.enterParen();  } -> type(LBRACE), popMode, pushMode(DEFAULT_MODE)
    ;
GStringIdentifier
    :   IdentifierInGString -> type(Identifier), popMode, pushMode(GSTRING_PATH_MODE)
    ;


mode GSTRING_PATH_MODE;
GStringPathPart
    :   Dot IdentifierInGString
    ;
RollBackOne
    :   . {
            readChar := l.GetInputStream().LA(-1)
            if l.GetInputStream().LA(1) == antlr.TokenEOF && (readChar == '"' || readChar == '/') {
                l.SetType(GroovyLexerGStringEnd)
            } else {
                l.SetChannel(antlr.TokenHiddenChannel)
            }
          } -> popMode
    ;


mode DEFAULT_MODE;
// character in the double quotation string. e.g. "a"
fragment
DqStringCharacter
    :   ~["\r\n\\$]
    |   EscapeSequence
    ;

// character in the single quotation string. e.g. 'a'
fragment
SqStringCharacter
    :   ~['\r\n\\]
    |   EscapeSequence
    ;

// character in the triple double quotation string. e.g. """a"""
fragment TdqStringCharacter
    :   ~["\\$]
    |   GStringQuotationMark { _input.LA(1) != '"' || _input.LA(2) != '"' || _input.LA(3) == '"' && (_input.LA(4) != '"' || _input.LA(5) != '"') }?
    |   EscapeSequence
    ;

// character in the triple single quotation string. e.g. '''a'''
fragment TsqStringCharacter
    :   ~['\\]
    |   SqStringQuotationMark { _input.LA(1) != '\'' || _input.LA(2) != '\'' || _input.LA(3) == '\'' && (_input.LA(4) != '\'' || _input.LA(5) != '\'') }?
    |   EscapeSequence
    ;

// character in the slashy string. e.g. /a/
fragment SlashyStringCharacter
    :   SlashEscape
    |   Dollar { !isFollowedByJavaLetterInGString(_input) }?
    |   ~[/$\u0000]
    ;

// character in the dollar slashy string. e.g. $/a/$
fragment DollarSlashyStringCharacter
    :   DollarDollarEscape
    |   DollarSlashDollarEscape { _input.LA(-4) != '$' }?
    |   DollarSlashEscape { _input.LA(1) != '$' }?
    |   Slash { _input.LA(1) != '$' }?
    |   Dollar { !isFollowedByJavaLetterInGString(_input) }?
    |   ~[/$\u0000]
    ;

// Groovy keywords
AS              : 'as';
DEF             : 'def';
IN              : 'in';
TRAIT           : 'trait';
THREADSAFE      : 'threadsafe'; // reserved keyword

// the reserved type name of Java10
VAR             : 'var';

// §3.9 Keywords
BuiltInPrimitiveType
    :   BOOLEAN
    |   CHAR
    |   BYTE
    |   SHORT
    |   INT
    |   LONG
    |   FLOAT
    |   DOUBLE
    ;

ABSTRACT      : 'abstract';
ASSERT        : 'assert';

fragment
BOOLEAN       : 'boolean';

BREAK         : 'break';
YIELD         : 'yield';

fragment
BYTE          : 'byte';

CASE          : 'case';
CATCH         : 'catch';

fragment
CHAR          : 'char';

CLASS         : 'class';
CONST         : 'const';
CONTINUE      : 'continue';
DEFAULT       : 'default';
DO            : 'do';

fragment
DOUBLE        : 'double';

ELSE          : 'else';
ENUM          : 'enum';
EXTENDS       : 'extends';
FINAL         : 'final';
FINALLY       : 'finally';

fragment
FLOAT         : 'float';


FOR           : 'for';
IF            : 'if';
GOTO          : 'goto';
IMPLEMENTS    : 'implements';
IMPORT        : 'import';
INSTANCEOF    : 'instanceof';

fragment
INT           : 'int';

INTERFACE     : 'interface';

fragment
LONG          : 'long';

NATIVE        : 'native';
NEW           : 'new';
NON_SEALED    : 'non-sealed';

PACKAGE       : 'package';
PERMITS       : 'permits';
PRIVATE       : 'private';
PROTECTED     : 'protected';
PUBLIC        : 'public';

RECORD        : 'record';
RETURN        : 'return';

SEALED        : 'sealed';

fragment
SHORT         : 'short';


STATIC        : 'static';
STRICTFP      : 'strictfp';
SUPER         : 'super';
SWITCH        : 'switch';
SYNCHRONIZED  : 'synchronized';
THIS          : 'this';
THROW         : 'throw';
THROWS        : 'throws';
TRANSIENT     : 'transient';
TRY           : 'try';
VOID          : 'void';
VOLATILE      : 'volatile';
WHILE         : 'while';


// §3.10.1 Integer Literals

IntegerLiteral
    :   (   DecimalIntegerLiteral
        |   HexIntegerLiteral
        |   OctalIntegerLiteral
        |   BinaryIntegerLiteral
        ) (Underscore { require(l.errorIgnored, "Number ending with underscores is invalid", -1, l); })?

    // !!! Error Alternative !!!
    |   Zero ([0-9] { l.invalidDigitCount++; })+ { require(l.errorIgnored, "Invalid octal number", -(l.invalidDigitCount + 1), l); } IntegerTypeSuffix?
    ;

fragment
Zero
    :   '0'
    ;

fragment
DecimalIntegerLiteral
    :   DecimalNumeral IntegerTypeSuffix?
    ;

fragment
HexIntegerLiteral
    :   HexNumeral IntegerTypeSuffix?
    ;

fragment
OctalIntegerLiteral
    :   OctalNumeral IntegerTypeSuffix?
    ;

fragment
BinaryIntegerLiteral
    :   BinaryNumeral IntegerTypeSuffix?
    ;

fragment
IntegerTypeSuffix
    :   [lLiIgG]
    ;

fragment
DecimalNumeral
    :   Zero
    |   NonZeroDigit (Digits? | Underscores Digits)
    ;

fragment
Digits
    :   Digit (DigitOrUnderscore* Digit)?
    ;

fragment
Digit
    :   Zero
    |   NonZeroDigit
    ;

fragment
NonZeroDigit
    :   [1-9]
    ;

fragment
DigitOrUnderscore
    :   Digit
    |   Underscore
    ;

fragment
Underscores
    :   Underscore+
    ;

fragment
Underscore
    :   '_'
    ;

fragment
HexNumeral
    :   Zero [xX] HexDigits
    ;

fragment
HexDigits
    :   HexDigit (HexDigitOrUnderscore* HexDigit)?
    ;

fragment
HexDigit
    :   [0-9a-fA-F]
    ;

fragment
HexDigitOrUnderscore
    :   HexDigit
    |   Underscore
    ;

fragment
OctalNumeral
    :   Zero Underscores? OctalDigits
    ;

fragment
OctalDigits
    :   OctalDigit (OctalDigitOrUnderscore* OctalDigit)?
    ;

fragment
OctalDigit
    :   [0-7]
    ;

fragment
OctalDigitOrUnderscore
    :   OctalDigit
    |   Underscore
    ;

fragment
BinaryNumeral
    :   Zero [bB] BinaryDigits
    ;

fragment
BinaryDigits
    :   BinaryDigit (BinaryDigitOrUnderscore* BinaryDigit)?
    ;

fragment
BinaryDigit
    :   [01]
    ;

fragment
BinaryDigitOrUnderscore
    :   BinaryDigit
    |   Underscore
    ;

// §3.10.2 Floating-Point Literals

FloatingPointLiteral
    :   (   DecimalFloatingPointLiteral
        |   HexadecimalFloatingPointLiteral
        ) (Underscore { require(l.errorIgnored, "Number ending with underscores is invalid", -1, l); })?
    ;

fragment
DecimalFloatingPointLiteral
    :   Digits? Dot Digits ExponentPart? FloatTypeSuffix?
    |   Digits ExponentPart FloatTypeSuffix?
    |   Digits FloatTypeSuffix
    ;

fragment
ExponentPart
    :   ExponentIndicator SignedInteger
    ;

fragment
ExponentIndicator
    :   [eE]
    ;

fragment
SignedInteger
    :   Sign? Digits
    ;

fragment
Sign
    :   [+\-]
    ;

fragment
FloatTypeSuffix
    :   [fFdDgG]
    ;

fragment
HexadecimalFloatingPointLiteral
    :   HexSignificand BinaryExponent FloatTypeSuffix?
    ;

fragment
HexSignificand
    :   HexNumeral Dot?
    |   Zero [xX] HexDigits? Dot HexDigits
    ;

fragment
BinaryExponent
    :   BinaryExponentIndicator SignedInteger
    ;

fragment
BinaryExponentIndicator
    :   [pP]
    ;

fragment
Dot :   '.'
    ;

// §3.10.3 Boolean Literals

BooleanLiteral
    :   'true'
    |   'false'
    ;


// §3.10.6 Escape Sequences for Character and String Literals

fragment
EscapeSequence
    :   Backslash [btnfrs"'\\]
    |   OctalEscape
    |   UnicodeEscape
    |   DollarEscape
    |   LineEscape
    ;


fragment
OctalEscape
    :   Backslash OctalDigit
    |   Backslash OctalDigit OctalDigit
    |   Backslash ZeroToThree OctalDigit OctalDigit
    ;

// Groovy allows 1 or more u's after the backslash
fragment
UnicodeEscape
    :   Backslash 'u' HexDigit HexDigit HexDigit HexDigit
    ;

fragment
ZeroToThree
    :   [0-3]
    ;

// Groovy Escape Sequences

fragment
DollarEscape
    :   Backslash Dollar
    ;

fragment
LineEscape
    :   Backslash LineTerminator
    ;

fragment
LineTerminator
    :   '\r'? '\n' | '\r'
    ;

fragment
SlashEscape
    :   Backslash Slash
    ;

fragment
Backslash
    :   '\\'
    ;

fragment
Slash
    :   '/'
    ;

fragment
Dollar
    :   '$'
    ;

fragment
GStringQuotationMark
    :   '"'
    ;

fragment
SqStringQuotationMark
    :   '\''
    ;

fragment
TdqStringQuotationMark
    :   '"""'
    ;

fragment
TsqStringQuotationMark
    :   '\'\'\''
    ;

fragment
DollarSlashyGStringQuotationMarkBegin
    :   '$/'
    ;

fragment
DollarSlashyGStringQuotationMarkEnd
    :   '/$'
    ;

// escaped forward slash
fragment
DollarSlashEscape
    :   '$/'
    ;

// escaped dollar sign
fragment
DollarDollarEscape
    :   '$$'
    ;

// escaped dollar slashy string delimiter
fragment
DollarSlashDollarEscape
    :   '$/$'
    ;

// §3.10.7 The Null Literal
NullLiteral
    :   'null'
    ;

// Groovy Operators

RANGE_INCLUSIVE         : '..';
RANGE_EXCLUSIVE_LEFT    : '<..';
RANGE_EXCLUSIVE_RIGHT   : '..<';
RANGE_EXCLUSIVE_FULL    : '<..<';
SPREAD_DOT              : '*.';
SAFE_DOT                : '?.';
SAFE_INDEX              : '?[' { l.enterParen();     } -> pushMode(DEFAULT_MODE);
SAFE_CHAIN_DOT          : '??.';
ELVIS                   : '?:';
METHOD_POINTER          : '.&';
METHOD_REFERENCE        : '::';
REGEX_FIND              : '=~';
REGEX_MATCH             : '==~';
POWER                   : '**';
POWER_ASSIGN            : '**=';
SPACESHIP               : '<=>';
IDENTICAL               : '===';
IMPLIES                 : '==>';
NOT_IDENTICAL           : '!==';
ARROW                   : '->';

// !internalPromise will be parsed as !in ternalPromise, so semantic predicates are necessary
NOT_INSTANCEOF      : '!instanceof' { isFollowedBy(_input, ' ', '\t', '\r', '\n') }?;
NOT_IN              : '!in'         { isFollowedBy(_input, ' ', '\t', '\r', '\n', '[', '(', '{') }?;


// §3.11 Separators

LPAREN          : '('  { l.enterParen();     } -> pushMode(DEFAULT_MODE);
RPAREN          : ')'  { l.exitParen();      } -> popMode;

LBRACE          : '{'  { l.enterParen();     } -> pushMode(DEFAULT_MODE);
RBRACE          : '}'  { l.exitParen();      } -> popMode;

LBRACK          : '['  { l.enterParen();     } -> pushMode(DEFAULT_MODE);
RBRACK          : ']'  { l.exitParen();      } -> popMode;

SEMI            : ';';
COMMA           : ',';
DOT             : Dot;

// §3.12 Operators

ASSIGN          : '=';
GT              : '>';
LT              : '<';
NOT             : '!';
BITNOT          : '~';
QUESTION        : '?';
COLON           : ':';
EQUAL           : '==';
LE              : '<=';
GE              : '>=';
NOTEQUAL        : '!=';
AND             : '&&';
OR              : '||';
INC             : '++';
DEC             : '--';
ADD             : '+';
SUB             : '-';
MUL             : '*';
DIV             : Slash;
BITAND          : '&';
BITOR           : '|';
XOR             : '^';
MOD             : '%';


ADD_ASSIGN      : '+=';
SUB_ASSIGN      : '-=';
MUL_ASSIGN      : '*=';
DIV_ASSIGN      : '/=';
AND_ASSIGN      : '&=';
OR_ASSIGN       : '|=';
XOR_ASSIGN      : '^=';
MOD_ASSIGN      : '%=';
LSHIFT_ASSIGN   : '<<=';
RSHIFT_ASSIGN   : '>>=';
URSHIFT_ASSIGN  : '>>>=';
ELVIS_ASSIGN    : '?=';


// §3.8 Identifiers (must appear after all keywords in the grammar)
CapitalizedIdentifier
    :   JavaLetter {Character.isUpperCase(_input.LA(-1))}? JavaLetterOrDigit*
    ;

Identifier
    :   JavaLetter JavaLetterOrDigit*
    ;

fragment
IdentifierInGString
    :   JavaLetterInGString JavaLetterOrDigitInGString*
    ;

fragment
JavaLetter
    :   [a-zA-Z$_] // these are the "java letters" below 0x7F
    |   // covers all characters above 0x7F which are not a surrogate
        ~[\u0000-\u007F\uD800-\uDBFF]
        { isJavaIdentifierStartAndNotIdentifierIgnorable(_input.LA(-1)) }?
    |   // covers UTF-16 surrogate pairs encodings for U+10000 to U+10FFFF
        [\uD800-\uDBFF] [\uDC00-\uDFFF]
        { isJavaIdentifierStartFromSurrogatePair(_input.LA(-2), _input.LA(-1)) }?
    ;

fragment
JavaLetterInGString
    :   JavaLetter { _input.LA(-1) != '$' }?
    ;

fragment
JavaLetterOrDigit
    :   [a-zA-Z0-9$_] // these are the "java letters or digits" below 0x7F
    |   // covers all characters above 0x7F which are not a surrogate
        ~[\u0000-\u007F\uD800-\uDBFF]
        { isJavaIdentifierPartAndNotIdentifierIgnorable(_input.LA(-1)) }?
    |   // covers UTF-16 surrogate pairs encodings for U+10000 to U+10FFFF
        [\uD800-\uDBFF] [\uDC00-\uDFFF]
        { isJavaIdentifierPartFromSurrogatePair(_input.LA(-2), _input.LA(-1)) }?
    ;

fragment
JavaLetterOrDigitInGString
    :   JavaLetterOrDigit  { _input.LA(-1) != '$' }?
    ;

fragment
ShCommand
    :   ~[\r\n\uFFFF]*
    ;

//
// Additional symbols not defined in the lexical specification
//

AT : '@';
ELLIPSIS : '...';

//
// Whitespace, line escape and comments
//
WS  : ([ \t]+ | LineEscape+) -> skip
    ;

// Inside (...) and [...] but not {...}, ignore newlines.
NL  : LineTerminator   { l.ignoreTokenInsideParens(); }
    ;

// Multiple-line comments (including groovydoc comments)
ML_COMMENT
    :   '/*' .*? '*/'       { l.addComment(0); l.ignoreMultiLineCommentConditionally(); } -> type(NL)
    ;

// Single-line comments
SL_COMMENT
    :   '//' ~[\r\n\uFFFF]* { l.addComment(1); l.ignoreTokenInsideParens(); }             -> type(NL)
    ;

// Script-header comments.
// The very first characters of the file may be "#!".  If so, ignore the first line.
SH_COMMENT
    :   '#!' { require(l.errorIgnored || l.tokenIndex == 0, "Shebang comment should appear at the first line", -2, l); } ShCommand (LineTerminator '#!' ShCommand)* -> skip
    ;

// Unexpected characters will be handled by groovy parser later.
UNEXPECTED_CHAR
    :   . { require(l.errorIgnored, "Unexpected character: '" + l.getText().replace("'", "\\'") + "'", -1, l); }
    ;
