package models

import (
	"fmt"
	"strings"
)

type TokenType int

type Token struct {
	Value string
	Type  TokenType
}

const (
	UNKNOWN_TOKEN TokenType = iota
	KEYWORD
	SYMBOL
	IDENTIFIER
	INT_CONST
	STRING_CONST
)

var (
	TokenTypeString = map[TokenType]string{
		KEYWORD:      "keyword",
		SYMBOL:       "symbol",
		IDENTIFIER:   "identifier",
		INT_CONST:    "integerConstant",
		STRING_CONST: "stringConstant",
	}

	TokenTags = map[TokenType]string{
		KEYWORD:      `<keyword> %s </keyword>`,
		SYMBOL:       `<symbol> %s </symbol>`,
		IDENTIFIER:   `<identifier> %s </identifier>`,
		INT_CONST:    `<integerConstant> %s </integerConstant>`,
		STRING_CONST: `<stringConstant> %s </stringConstant>`,
	}

	SpecialTokens = map[rune]string{
		'<': "&lt;",
		'>': "&gt;",
		'&': "&amp;",
	}
)

func WrapTokenTag(t *Token) string {
	return fmt.Sprintf(TokenTags[t.Type], escapeToken(t.Value))
}

func NewToken(v string, t TokenType) *Token {
	return &Token{
		Value: v,
		Type:  t,
	}
}

func escapeToken(s string) string {
	var b strings.Builder
	for _, r := range s {
		if c, ok := SpecialTokens[r]; ok {
			fmt.Fprint(&b, c)
			continue
		}

		fmt.Fprint(&b, string(r))
	}

	return b.String()
}

func (t TokenType) IsKeyword() bool {
	return t == KEYWORD
}

func (t *Token) IsPrimitiveType() bool {
	return t.Type.IsKeyword() && Keywords[t.Value].IsPrimitiveType()
}
