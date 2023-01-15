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
	for _, token := range tokens {
		fmt.Println(token.value)
	}

	ast, _ := parse(tokens, 0)
	fmt.Print(ast.pretty())

	// TODO: execution via AST walking interpretation
	//value := interpret(ast)
	//fmt.Println(value)

	// TODO: compile the AST to JavaScript? Go? C? Assembly? LLVM?
	//compile(ast)
}
