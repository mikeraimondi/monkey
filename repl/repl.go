package repl

import (
	"bufio"
	"fmt"
	"io"
	"log"

	"github.com/mikeraimondi/monkey/object"

	"github.com/mikeraimondi/monkey/evaluator"
	"github.com/mikeraimondi/monkey/lexer"
	"github.com/mikeraimondi/monkey/parser"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()
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

		evaluated := evaluator.Eval(expanded, env)
		if evaluated != nil {
			if _, err := io.WriteString(out, evaluated.Inspect()+"\n"); err != nil {
				log.Fatalln(err)
			}
		}
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
