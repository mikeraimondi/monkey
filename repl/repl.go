package repl

import (
	"bufio"
	"fmt"
	"io"
	"log"

	"github.com/mikeraimondi/monkey/lexer"
	"github.com/mikeraimondi/monkey/parser"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

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

		if _, err := io.WriteString(out, program.String()+"\n"); err != nil {
			log.Fatalln(err)
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
