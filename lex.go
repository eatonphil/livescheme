package main

import (
	"unicode"
)

type tokenKind uint
const (
	// e.g. "(", ")"
	syntaxToken tokenKind = iota
	// e.g. "1", "12"
	integerToken
	// e.g. "+", "define"
	identifierToken
)

type token struct {
	value string
	kind tokenKind
	location int
}

func (t token) debug(source []rune) {
	// maybe implement where we go through the source and show
	// where the token is in source, for error debugging
}

func eatWhitespace(source []rune, cursor int) int {
	for cursor < len(source) {
		if unicode.IsSpace(source[cursor]) {
			cursor++
			continue
		}

		break
	}

	return cursor
}

func lexSyntaxToken(source []rune, cursor int) (int, *token) {
	if source[cursor] == '(' || source[cursor] == ')' {
		return cursor + 1, &token{
			value: string([]rune{source[cursor]}),
			kind: syntaxToken,
			location: cursor,
		}
	}

	return cursor, nil
}

// lexIntegerToken("foo 123", 4) => "123"
// lexIntegerToken("foo 12 3", 4) => "12"
// lexIntegerToken("foo 12a 3", 4) => "12" <-- Ignoring this situation
func lexIntegerToken(source []rune, cursor int) (int, *token) {
	originalCursor := cursor

	var value []rune
	for cursor < len(source) {
		r := source[cursor]
		if r >= '0' && r <= '9' {
			value = append(value, r)
			cursor++
			continue
		}

		break
	}

	if len(value) == 0 {
		return originalCursor, nil
	}

	return cursor, &token{
		value: string(value),
		kind: integerToken,
		location: originalCursor,
	}
}

// lexIdentifierToken("123 ab + ", 4) => "ab"
// lexIdentifierToken("123 ab123 + ", 4) => "ab123"
func lexIdentifierToken(source []rune, cursor int) (int, *token) {
	originalCursor := cursor
	var value []rune

	for cursor < len(source) {
		r := source[cursor]
		if !unicode.IsSpace(r) {
			value = append(value, r)
			cursor++
			continue
		}

		break
	}

	if len(value) == 0 {
		return originalCursor, nil
	}

	return cursor, &token{
		value: string(value),
		kind: identifierToken,
		location: originalCursor,
	}
}

// for example: "(+ 13 2)"
// lex(" ( + 13 2  )") should produce: ["(", "+", "13", "2", ")"]
func lex(raw string) []token {
	source := []rune(raw)
	var tokens []token
	var t *token

	cursor := 0
	for cursor < len(source) {
		// eat whitespace
		cursor = eatWhitespace(source, cursor)
		if cursor == len(source) {
			break
		}

		cursor, t = lexSyntaxToken(source, cursor)
		if t != nil {
			tokens = append(tokens, *t)
			continue
		}

		cursor, t = lexIntegerToken(source, cursor)
		if t != nil {
			tokens = append(tokens, *t)
			continue
		}

		cursor, t = lexIdentifierToken(source, cursor)
		if t != nil {
			tokens = append(tokens, *t)
			continue
		}

		// Lexed nothing, not good!
		// fmt.Println(tokens[len(tokens)-1].debug()) // line of code
		panic("Could not lex")
	}

	return tokens
}
