package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// lexIntegerToken("foo 123", 4) => "123"
// lexIntegerToken("foo 12 3", 4) => "12"
// lexIntegerToken("foo 12a 3", 4) => "12" <-- Ignoring this situation
func Test_lexIntegerToken(t *testing.T) {
	tests := []struct {
		source         string
		cursor         int
		expectedValue  string
		expectedCursor int
	}{
		{
			"foo 123",
			4,
			"123",
			7,
		},
		{
			"foo 12 3",
			4,
			"12",
			6,
		},
		{
			"foo 12a 3",
			4,
			"12",
			6,
		},
	}
	for _, test := range tests {
		lc := lexingContext{
			source:         []rune(test.source),
			sourceFileName: "<memory>",
		}
		cursor, token := lc.lexIntegerToken(test.cursor)
		assert.Equal(t, cursor, test.expectedCursor)
		assert.Equal(t, token.value, test.expectedValue)
		assert.Equal(t, token.kind, integerToken)
	}
}

// lexIdentifierToken("123 ab + ", 4) => "ab"
// lexIdentifierToken("123 ab123 + ", 4) => "ab123"
func Test_lexIdentifierToken(t *testing.T) {
	tests := []struct {
		source         string
		cursor         int
		expectedValue  string
		expectedCursor int
	}{
		{
			"123 ab + ",
			4,
			"ab",
			6,
		},
		{
			"123 ab123 + ",
			4,
			"ab123",
			9,
		},
	}
	for _, test := range tests {
		lc := lexingContext{
			source:         []rune(test.source),
			sourceFileName: "<memory>",
		}
		cursor, token := lc.lexIdentifierToken(test.cursor)
		assert.Equal(t, cursor, test.expectedCursor)
		assert.Equal(t, token.value, test.expectedValue)
		assert.Equal(t, token.kind, identifierToken)
	}
}

// lex(" ( + 13 2  )") should produce: ["(", "+", "13", "2", ")"]
func Test_lex(t *testing.T) {
	tests := []struct {
		source string
		tokens []token
	}{
		{
			" ( + 13 2  )",
			[]token{
				{
					value:    "(",
					kind:     syntaxToken,
					location: 1,
				},
				{
					value:    "+",
					kind:     identifierToken,
					location: 3,
				},
				{
					value:    "13",
					kind:     integerToken,
					location: 5,
				},
				{
					value:    "2",
					kind:     integerToken,
					location: 8,
				},
				{
					value:    ")",
					kind:     syntaxToken,
					location: 11,
				},
			},
		},
	}

	for _, test := range tests {
		lc := lexingContext{
			source:         []rune(test.source),
			sourceFileName: "<memory>",
		}
		tokens := lc.lex()
		for i, token := range tokens {
			// Cheating by setting the received token's
			// lexingContext to the expected lexingContext
			token.lc = test.tokens[i].lc
			assert.Equal(t, token, test.tokens[i])
		}
	}
}
