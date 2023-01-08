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
	fmt.Println(tokens)

	// TODO: parsing
	//ast := parse(tokens)

	// TODO: execution via AST walking interpretation
	//value := interpret(ast)
	//fmt.Println(value)

	// TODO: compile the AST to JavaScript? Go? C? Assembly? LLVM?
	//compile(ast)
}
