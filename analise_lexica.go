package main

import (
	"errors"
	"fmt"
	"os"
)

var tipos = map[string]string{
	"(": "ParenEsq",
	")": "ParenDir",
	"+": "Soma",
	"-": "Sub",
	"*": "Mult",
	"/": "Div",
}

type Token struct {
	Tipo    string
	Literal string
	posicao int
}

func getProximoToken(input string, index *int) Token {
	for *index < len(input) && (input[*index] == ' ' || input[*index] == '\n' || input[*index] == '\t') {
		*index++
	}

	if *index >= len(input) {
		return Token{Tipo: "EOF", Literal: ""}
	}

	start := *index
	char := string(input[*index])

	if val, ok := tipos[char]; ok {
		*index++
		return Token{Tipo: val, Literal: char, posicao: start}
	}

	if input[*index] >= '0' && input[*index] <= '9' {
		for *index < len(input) && input[*index] >= '0' && input[*index] <= '9' {
			*index++
		}
		return Token{Tipo: "Numero", Literal: input[start:*index], posicao: start}
	}

	*index++
	return Token{Tipo: "INVALIDO", Literal: char, posicao: start}
}

func tokenize(input string) ([]Token, error) {
	var tokens []Token
	index := 0

	for index < len(input) {
		token := getProximoToken(input, &index)
		if token.Tipo == "EOF" {
			break
		}
		if token.Tipo == "INVALIDO" {
			return nil, errors.New("Token inválido encontrado: " + token.Literal + " na posição " + fmt.Sprint(token.posicao))
		}
		tokens = append(tokens, token)
	}
	return tokens, nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Uso: go run analise_lexica.go <arquivo>")
		return
	}

	content, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Println("Erro ao ler arquivo:", err)
		return
	}

	input := string(content)
	tokens, err := tokenize(input)
	if err != nil {
		fmt.Println("Erro na análise léxica:", err)
	}
	for _, token := range tokens {
		fmt.Printf("<%s,\"%s\", %d>\n", token.Tipo, token.Literal, token.posicao)
	}

}
