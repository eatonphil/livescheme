package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func compareValue(a value, b value) bool {
	if a.kind != b.kind {
		fmt.Println("Value kinds not equal", a.kind, b.kind)
		return false
	}

	if a.kind == literalValue {
		if a.literal.value != b.literal.value {
			fmt.Println("Literals not equal", a.literal, b.literal)
			return false
		}

		return true
	}

	return compareAst(*a.list, *b.list)
}

func compareAst(a ast, b ast) bool {
	if len(a) != len(b) {
		fmt.Println("AST lengths not equal", len(a), len(b))
		return false
	}

	for i := range a {
		aI := a[i]
		bI := b[i]

		if !compareValue(aI, bI) {
			return false
		}
	}

	return true
}

func Test_parse(t *testing.T) {
	tests := []struct {
		input        string
		prettyOutput string
		output       ast
	}{
		{
			"(+ 1 2)",
			"(+ 1 2 )",
			ast{
				value{
					kind:    literalValue,
					literal: &token{value: "+"},
				},
				value{
					kind:    literalValue,
					literal: &token{value: "1"},
				},
				value{
					kind:    literalValue,
					literal: &token{value: "2"},
				},
			},
		},
		{
			"(+ 1 (- 12 9))",
			"(+ 1 (- 12 9 ) )",
			ast{
				value{
					kind:    literalValue,
					literal: &token{value: "+"},
				},
				value{
					kind:    literalValue,
					literal: &token{value: "1"},
				},
				value{
					kind: listValue,
					list: &ast{
						value{
							kind:    literalValue,
							literal: &token{value: "-"},
						},
						value{
							kind:    literalValue,
							literal: &token{value: "12"},
						},
						value{
							kind:    literalValue,
							literal: &token{value: "9"},
						},
					},
				},
			},
		},
		{
			"(+ 1 (- 12 9) 12)",
			"(+ 1 (- 12 9 ) 12 )",
			ast{
				value{
					kind:    literalValue,
					literal: &token{value: "+"},
				},
				value{
					kind:    literalValue,
					literal: &token{value: "1"},
				},
				value{
					kind: listValue,
					list: &ast{
						value{
							kind:    literalValue,
							literal: &token{value: "-"},
						},
						value{
							kind:    literalValue,
							literal: &token{value: "12"},
						},
						value{
							kind:    literalValue,
							literal: &token{value: "9"},
						},
					},
				},
				value{
					kind:    literalValue,
					literal: &token{value: "12"},
				},
			},
		},
		{
			"((+ 1 2) 1 (- 12 9) 12)",
			"((+ 1 2 ) 1 (- 12 9 ) 12 )",
			ast{
				value{
					kind: listValue,
					list: &ast{
						value{
							kind:    literalValue,
							literal: &token{value: "+"},
						},
						value{
							kind:    literalValue,
							literal: &token{value: "1"},
						},
						value{
							kind:    literalValue,
							literal: &token{value: "2"},
						},
					},
				},
				value{
					kind:    literalValue,
					literal: &token{value: "1"},
				},
				value{
					kind: listValue,
					list: &ast{
						value{
							kind:    literalValue,
							literal: &token{value: "-"},
						},
						value{
							kind:    literalValue,
							literal: &token{value: "12"},
						},
						value{
							kind:    literalValue,
							literal: &token{value: "9"},
						},
					},
				},
				value{
					kind:    literalValue,
					literal: &token{value: "12"},
				},
			},
		},
	}

	for _, test := range tests {
		lc := lexingContext{
			source:         []rune(test.input),
			sourceFileName: "<memory>",
		}
		tokens := lc.lex()
		ast, _ := parse(tokens, 0)
		assert.Equal(t, test.prettyOutput, ast.pretty())
		assert.True(t, compareAst(test.output, ast))
	}
}
