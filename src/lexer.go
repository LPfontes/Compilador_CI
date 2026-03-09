package main

import "errors"

var tipos = map[string]string{
	"(": "ParenEsq",
	")": "ParenDir",
	"+": "Soma",
	"-": "Sub",
	"*": "Mult",
	"/": "Div",
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
	if val, ok := tipos[char]; ok { // Verifica se o caractere é um símbolo conhecido
		*index++
		return Token{Tipo: val, Literal: char, posicao: start}
	}

	if input[*index] >= '0' && input[*index] <= '9' { // Verifica se é um número
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
			return nil, errors.New("Token inválido: " + token.Literal)
		}
		tokens = append(tokens, token)
	}
	tokens = append(tokens, Token{Tipo: "EOF", Literal: ""})
	return tokens, nil
}
