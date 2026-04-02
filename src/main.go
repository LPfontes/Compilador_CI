package main

import (
	"fmt"
	"os"
)

// --- Função Principal ---

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Uso: go run main.go \"(1 + 2)\"")
		return
	}
	if _, err := os.Stat("arvore"); os.IsNotExist(err) {
		if err := os.Mkdir("arvore", 0755); err != nil {
			fmt.Println("Erro ao criar diretório 'arvore':", err)
			return
		}
	}
	if _, err := os.Stat("output"); os.IsNotExist(err) {
		if err := os.Mkdir("output", 0755); err != nil {
			fmt.Println("Erro ao criar diretório 'output':", err)
			return
		}
	}
	input := os.Args[1]
	// LEXER
	tokens, err := tokenize(input)
	if err != nil {
		fmt.Println("Erro Léxico:", err)
		return
	}
	// PARSER
	parser := &Parser{tokens: tokens}
	arvore, err := parser.analisaPrograma()
	if err != nil {
		fmt.Println("Erro Sintático:", err)
		return
	}

	err = AnalisarSemantica(arvore)
	if err != nil {
		fmt.Println("Erro Semântico:", err)
		return
	}

	fmt.Printf("Árvore gerada %+v\n", arvore)
	gerarVisualizacao(arvore, "arvore")
	
	amb := make(map[string]int)
	resultado, err := interpretar(arvore, amb)
	if err != nil {
		fmt.Println("Erro de Execução:", err)
		return
	} else {
		fmt.Printf("Resultado: %d\n", resultado)
	}

	compilar(arvore, "output/output.s")
	executarBinario("output/output.s")
}
