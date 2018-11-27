package repl

import (
	"bufio"
	"fmt"
	"io"
	"log"

	"github.com/mikeraimondi/monkey/compiler"
	"github.com/mikeraimondi/monkey/evaluator"
	"github.com/mikeraimondi/monkey/lexer"
	"github.com/mikeraimondi/monkey/object"
	"github.com/mikeraimondi/monkey/parser"
	"github.com/mikeraimondi/monkey/vm"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	macroEnv := object.NewEnvironment()
	constants := []object.Object{}
	globals := make([]object.Object, vm.GlobalsSize)
	symbolTable := compiler.NewSymbolTable()
	for i, v := range object.Builtins {
		symbolTable.DefineBuiltin(i, v.Name)
	}

	for {
		fmt.Printf(PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)

		program := p.ParseProgram()
		if errs := p.Errors(); len(errs) != 0 {
			printParserErrors(out, errs)
			continue
		}

		evaluator.DefineMacros(program, macroEnv)
		expanded := evaluator.ExpandMacros(program, macroEnv)

		comp := compiler.NewWithState(symbolTable, constants)
		err := comp.Compile(expanded)
		if err != nil {
			fmt.Fprintf(out, "compilation failure:\n %s\n", err)
			continue
		}

		code := comp.Bytecode()
		constants = code.Constants

		machine := vm.NewWithGlobalsStore(code, globals)
		err = machine.Run()
		if err != nil {
			fmt.Fprintf(out, "bytecode execution failure:\n %s\n", err)
			continue
		}

		lastPopped := machine.LastPoppedStackElem()
		fmt.Fprintf(out, "%s\n", lastPopped.Inspect())
	}
}

func printParserErrors(out io.Writer, errors []string) {
	if _, err := io.WriteString(out, "parser errors:\n"); err != nil {
		log.Fatalln(err)
	}
	for _, msg := range errors {
		if _, err := io.WriteString(out, "\t"+msg+"\n"); err != nil {
			log.Fatalln(err)
		}
	}
}
