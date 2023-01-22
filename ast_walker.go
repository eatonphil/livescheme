package main

import (
	"fmt"
	"strconv"
)

var builtins = map[string]func([]value, map[string]any) any{}

func copyContext(in map[string]any) map[string]any {
	out := map[string]any{}
	for key, val := range in {
		out[key] = val
	}

	return out
}

// Eventually with all operations implemented we can do stuff like fib:
// (def fib (a)
//   (if (a < 1)
//     a
//     (fib (- a 1) (- a 2))

func initializeBuiltins() {
	// (if (< a 2)
	//     then
	//     elsecase)
	builtins["if"] = func(args []value, ctx map[string]any) any {
		condition := astWalk2(args[0], ctx)
		then := args[1]
		_else := args[2]

		if condition.(bool) == true {
			return astWalk2(then, ctx)
		}

		return astWalk2(_else, ctx)
	}

	builtins["+"] = func(args []value, ctx map[string]any) any {
		var i int64
		for _, arg := range args {
			i += astWalk2(arg, ctx).(int64)
		}
		return i
	}

	builtins["-"] = func(args []value, ctx map[string]any) any {
		i := astWalk2(args[0], ctx).(int64)
		for _, arg := range args[1:] {
			i -= astWalk2(arg, ctx).(int64)
		}
		return i
	}

	builtins["begin"] = func(args []value, ctx map[string]any) any {
		var last any
		for _, arg := range args {
			last = astWalk2(arg, ctx)
		}

		return last
	}

	// (var a 2)
	// (func plus (a b) (+ a b))
	builtins["func"] = func(args []value, ctx map[string]any) any {
		// e.g. `plus`
		functionName := (*args[0].literal).value

		// e.g. `(a b)`
		params := *args[1].list

		// e.g. (+ a b)
		body := *args[2].list

		// e.g. `(plus 1 2)`
		ctx[functionName] = func(args []any, ctx map[string]any) any {
			childCtx := copyContext(ctx)
			if len(params) != len(args) {
				panic(fmt.Sprintf("Expected %d args to `%s`, got %d", len(params), functionName, len(args)))
			}
			for i, param := range params {
				childCtx[(*param.literal).value] = args[i]
			}

			return astWalk(body, childCtx)
		}

		return ctx[functionName]
	}
}

// Later: user defined functions
// (func plus (a b) (+ a b))
// And also user defined variables
// (var a 12)

// Example of evaluation
// ( + 13 ( - 12 1) )
//   + 13 11
//   24

// Example file that is ok:
// (+ 12 1)
// (- 134 9)
// All expression get evaluated. Only the last one is returned.
// That file is transformed into:
// (begin
//
//	(+ 12 1)
//	(- 134 9))
//
// Before the astWalk stage
func astWalk(ast []value, ctx map[string]any) any {
	// Default case: we've got a list
	// Example: (+ 1 2)
	// Example: `if`, `+`
	functionName := (*ast[0].literal).value

	if builtinFunction, ok := builtins[functionName]; ok {
		return builtinFunction(ast[1:], ctx)
	}

	// Case: calling a function that is not built in
	userDefinedFunction := ctx[functionName].(func([]any, map[string]any) any)

	// Do we evaluate args here?
	// If so, special functions like `if` must be handled separately
	var args []any
	for _, unevaluatedArg := range ast[1:] {
		args = append(args, astWalk2(unevaluatedArg, ctx))
	}

	return userDefinedFunction(args, ctx)
}

func astWalk2(v value, ctx map[string]any) any {
	if v.kind == literalValue {
		t := *v.literal
		switch t.kind {
		// `12`, `1`
		case integerToken:
			// parseInt here
			i, err := strconv.ParseInt(t.value, 10, 64)
			if err != nil {
				fmt.Println("Expected an integer, got: ", t.value)
				panic(err)
			}

			return i
		// (var a 3)
		// (+ a 1) => result should be `4`
		// `+`, or `if`
		case identifierToken:
			return ctx[t.value]
		}
	}

	return astWalk(*v.list, ctx)
}
