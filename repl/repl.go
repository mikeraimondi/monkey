package repl

import (
	"bufio"
	"fmt"
	"io"
	"log"

	"github.com/mikeraimondi/monkey/compiler"
	"github.com/mikeraimondi/monkey/object"
	"github.com/mikeraimondi/monkey/vm"

	"github.com/mikeraimondi/monkey/evaluator"
	"github.com/mikeraimondi/monkey/lexer"
	"github.com/mikeraimondi/monkey/parser"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	macroEnv := object.NewEnvironment()

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

		comp := compiler.New()
		err := comp.Compile(expanded)
		if err != nil {
			fmt.Fprintf(out, "compilation failure:\n %s\n", err)
			continue
		}

		machine := vm.New(comp.Bytecode())
		err = machine.Run()
		if err != nil {
			fmt.Fprintf(out, "bytecode execution failure:\n %s\n", err)
			continue
		}

		stackTop := machine.StackTop()
		fmt.Fprintf(out, "%s\n", stackTop.Inspect())
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
