package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func GetInput(prompt string) (string, error) {
	r := bufio.NewReader(os.Stdin)

	fmt.Print(prompt)
	input, err := r.ReadString('\n')

	return input[:len(input) - 1], err
}

var globalSymbolTable = NewSymbolTable()

func Run(t string) []interface{} {
	if strings.TrimSpace(t) == "" {
		return []interface{}{}
	}

	lexer := NewLexer(t, "<stdin>", "")
	tokens, err := lexer.MakeTokens()

	if err != nil {
		fmt.Println(err)
		return []interface{}{}
	}

	parser := NewParser(tokens, -1)
	ast := parser.Parse()
	
	if ast.Error != nil {
		fmt.Println(ast.Error)
		return []interface{}{}
	}
	
	interp := NewInterpretor()
	ctx := NewContext("<root>")
	ctx.SymbolTable = globalSymbolTable

	res := interp.Visit(ast.Node, ctx)

	if res.Error != nil {
		fmt.Println(res.Error)
		return []interface{}{}
	}

	return res.Value.GetVal().([]interface{})
}

func main() {
	if len(os.Args) < 2 {
		for {
			text, err := GetInput("\033[33mLuminary %\033[37m ")

			if err != nil {
				fmt.Println(err)
				break
			}

			res := Run(text)
			for _, val := range res {
				fmt.Println(val)
			}
		}
		return
	}

	file := os.Args[1]
	content, err := os.ReadFile(file)
	if err != nil {
		fmt.Println("Failed to load file")
		return
	}
	Run(string(content))
}
