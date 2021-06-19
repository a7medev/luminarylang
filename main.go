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

func run(t string) {
	lexer := NewLexer(t)
	tokens, err := lexer.MakeTokens()

	if err != nil {
		fmt.Println(err)
	}

	for _, v := range tokens {
		fmt.Println(v.ToString())
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	for {
		text, err := getInput("\033[33mluminary %\033[37m ", reader)
		if err != nil {
			panic(err)
		}
		if (text == "exit") {
			break
		}
		run(text)
	}
}
