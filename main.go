package main

import "os"
import "fmt"

func main() {
	// accept program
	program, err := os.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}

	tokens := lex(string(program))
	debug := false
	if debug {
		for _, token := range tokens {
			fmt.Println(token.value)
		}
	}

	var parseIndex int
	var a = ast{
		value{
			kind: literalValue,
			literal: &token{
				value: "begin",
				kind:  identifierToken,
			},
		},
	}

	// Need to keep parsing until end of ALL tokens
	for parseIndex < len(tokens) {
		childAst, nextIndex := parse(tokens, parseIndex)
		a = append(a, value{
			kind: listValue,
			list: &childAst,
		})
		parseIndex = nextIndex
	}

	if debug {
		fmt.Println(a.pretty())
	}

	// Other potential steps:
	// 1. static type checking?
	// not in our language

	// 2. other optimization steps: constant propagation? (+ 5 2) => 7
	// not for now

	initializeBuiltins()
	ctx := map[string]any{}
	value := astWalk(a, ctx)
	fmt.Println(value)

	// TODO: compile the AST to JavaScript? Go? C? Assembly? LLVM?
	//compile(ast)
}
