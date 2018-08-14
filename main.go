package main

import (
	"os"

	"github.com/mikeraimondi/monkey/repl"
)

func main() {
	repl.Start(os.Stdin, os.Stdout)
}
