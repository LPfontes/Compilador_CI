package main

import "errors"

var tipos = map[string]string{
	"(": "ParenEsq",
	")": "ParenDir",
	"+": "Soma",
	"-": "Sub",
	"*": "Mult",
	"/": "Div",
	"=": "Atribuicao",
	";": "PontoVirgula",
	",": "Virgula",
	"{": "ChaveEsq",
	"}": "ChaveDir",
	"<": "MenorQue",
	">": "MaiorQue",
	"[": "ColcheteEsq",
	"]": "ColcheteDir",
}

func getProximoToken(input string, index *int) Token {
	for *index < len(input) && (input[*index] == ' ' || input[*index] == '\n' || input[*index] == '\t') {
		*index++
	}
	if *index >= len(input) {
		return Token{Tipo: "EOF", Literal: ""}
	}
	start := *index

	if *index+1 < len(input) {
		chars := input[*index : *index+2]
		if chars == "==" {
			*index += 2
			return Token{Tipo: "IgualIgual", Literal: "==", posicao: start}
		} else if chars == "<=" {
			*index += 2
			return Token{Tipo: "MenorIgual", Literal: "<=", posicao: start}
		} else if chars == "!=" {
			*index += 2
			return Token{Tipo: "Diferente", Literal: "!=", posicao: start}
		} else if chars == ">=" {
			*index += 2
			return Token{Tipo: "MaiorIgual", Literal: ">=", posicao: start}
		}
	}

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

	if (input[*index] >= 'a' && input[*index] <= 'z') || (input[*index] >= 'A' && input[*index] <= 'Z') { // Verifica se é um identificador
		for *index < len(input) && ((input[*index] >= 'a' && input[*index] <= 'z') || (input[*index] >= 'A' && input[*index] <= 'Z') || (input[*index] >= '0' && input[*index] <= '9')) {
			*index++
		}
		literal := input[start:*index]
		switch literal {
		case "if":
			return Token{Tipo: "If", Literal: literal, posicao: start}
		case "else":
			return Token{Tipo: "Else", Literal: literal, posicao: start}
		case "while":
			return Token{Tipo: "While", Literal: literal, posicao: start}
		case "for":
			return Token{Tipo: "For", Literal: literal, posicao: start}
		case "return":
			return Token{Tipo: "Return", Literal: literal, posicao: start}
		case "and":
			return Token{Tipo: "And", Literal: literal, posicao: start}
		case "or":
			return Token{Tipo: "Or", Literal: literal, posicao: start}
		case "not":
			return Token{Tipo: "Not", Literal: literal, posicao: start}
		case "fun":
			return Token{Tipo: "Fun", Literal: literal, posicao: start}
		case "var":
			return Token{Tipo: "Var", Literal: literal, posicao: start}
		case "main":
			return Token{Tipo: "Main", Literal: literal, posicao: start}
		default:
			return Token{Tipo: "Identificador", Literal: literal, posicao: start}
		}
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
