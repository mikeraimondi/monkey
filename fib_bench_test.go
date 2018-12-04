package main

import (
	"fmt"
	"testing"

	"github.com/mikeraimondi/monkey/ast"
	"github.com/mikeraimondi/monkey/compiler"
	"github.com/mikeraimondi/monkey/evaluator"
	"github.com/mikeraimondi/monkey/lexer"
	"github.com/mikeraimondi/monkey/object"
	"github.com/mikeraimondi/monkey/parser"
	"github.com/mikeraimondi/monkey/vm"
)

var input = `
let fibonacci = fn(x) {
	if (x == 0) {
		0
	} else {
		if (x == 1) {
			return 1;
		} else {
			fibonacci(x - 1) + fibonacci(x - 2);
		}
	}
};
fibonacci(15);
`
var expected = fmt.Sprintf("%d", 610)

func setupBenchmark(b *testing.B) *ast.Program {
	b.Helper()

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	return program
}

func BenchmarkCompiledExecutionFib(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		program := setupBenchmark(b)
		comp := compiler.New()
		err := comp.Compile(program)
		if err != nil {
			b.Fatalf("compiler error: %s", err)
		}

		machine := vm.New(comp.Bytecode())

		b.StartTimer()
		err = machine.Run()
		if err != nil {
			b.Fatalf("vm error: %s", err)
		}

		result := machine.LastPoppedStackElem()
		if actual := result.Inspect(); actual != expected {
			b.Fatalf("wrong result. expected %s, got %s", expected, actual)
		}
	}
}

func BenchmarkInterpretedExecutionFib(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		program := setupBenchmark(b)
		env := object.NewEnvironment()

		b.StartTimer()
		result := evaluator.Eval(program, env)
		if actual := result.Inspect(); actual != expected {
			b.Fatalf("wrong result. expected %s, got %s", expected, actual)
		}
	}
}
