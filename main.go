package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func getInput(prompt string, r *bufio.Reader) (string, error) {
	fmt.Print(prompt)
	input, err := r.ReadString('\n')

	return strings.TrimSpace(input), err
}

var globalSymbolTable = NewSymbolTable()

func run(t string) {
	globalSymbolTable.Set("true", NewNumber(1))
	globalSymbolTable.Set("false", NewNumber(1))

	lexer := NewLexer(t, "<stdin>", "")
	tokens, err := lexer.MakeTokens()

	if err != nil {
		fmt.Println(err)
		return
	}

	parser := NewParser(tokens, -1)
	ast := parser.Parse()

	if ast.Error != nil {
		fmt.Println(ast.Error)
		return
	}

	interp := NewInterpretor()
	ctx := NewContext("<root>")
	ctx.SymbolTable = globalSymbolTable

	res := interp.Visit(ast.Node, ctx)
	if res != nil {
		fmt.Println(res)
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	for {
		text, err := getInput("\033[33mluminary %\033[37m ", reader)
		if err != nil {
			fmt.Println(err)
			break
		}
		if (text == "exit") {
			break
		}
		run(text)
	}
}
