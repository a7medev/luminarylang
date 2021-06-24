package main

import (
	"bufio"
	"fmt"
	"os"
)

func GetInput(prompt string) (string, error) {
	r := bufio.NewReader(os.Stdin)

	fmt.Print(prompt)
	input, err := r.ReadString('\n')

	return input[:len(input) - 1], err
}

var globalSymbolTable = NewSymbolTable()

func run(t string) {
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

	interp.Visit(ast.Node, ctx)
}

func main() {
	for {
		text, err := GetInput("\033[33mluminary %\033[37m ")
		if err != nil {
			fmt.Println(err)
			break
		}
		run(text)
	}
}
