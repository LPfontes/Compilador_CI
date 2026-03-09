package main

import (
	"strings"
	"testing"
)

func TestPrecedenciaOperadores(t *testing.T) {
	casos := []struct {
		nome     string
		entrada  string
		esperado int
	}{
		{
			nome:     "Soma e Multiplicacao",
			entrada:  "1 + 2 * 3",
			esperado: 7, // Deve ser 1 + (2 * 3)
		},
		{
			nome:     "Subtracao e Divisao",
			entrada:  "10 - 8 / 2",
			esperado: 6, // Deve ser 10 - (8 / 2)
		},
		{
			nome:     "Multiplicacao e Soma",
			entrada:  "2 * 3 + 4",
			esperado: 10, // Deve ser (2 * 3) + 4
		},
		{
			nome:     "Divisao e Subtracao",
			entrada:  "10 / 2 - 1",
			esperado: 4, // Deve ser (10 / 2) - 1
		},
		{
			nome:     "Ass. Esquerda Soma Sub",
			entrada:  "10 - 5 + 2",
			esperado: 7, // Deve ser (10 - 5) + 2
		},
		{
			nome:     "Ass. Esquerda Mult Div",
			entrada:  "8 / 4 * 2",
			esperado: 4, // Deve ser (8 / 4) * 2
		},
		{
			nome:     "Parenteses Alterando Precedencia",
			entrada:  "(1 + 2) * 3",
			esperado: 9,
		},
		{
			nome:     "Expressao Longa Mista",
			entrada:  "1 + 2 * 3 - 4 / 2",
			esperado: 5, // 1 + 6 - 2
		},
		{
			nome:     "Parenteses Aninhados",
			entrada:  "((1 + 2) * (3 + 4)) / 7",
			esperado: 3, // (3 * 7) / 7
		},
		{
			nome:     "Complexa com Precedencia",
			entrada:  "10 + 3 * 5 - 20 / 4",
			esperado: 20, // 10 + 15 - 5
		},
		{
			nome:     "Muitas Operacoes",
			entrada:  "2 + 3 * 4 / 2 - 1",
			esperado: 7, // 2 + (12/2) - 1 = 2 + 6 - 1
		},
		{
			nome:     "Associatividade Longa",
			entrada:  "20 - 5 - 2 - 1",
			esperado: 12, // ((20 - 5) - 2) - 1
		},
	}

	for _, caso := range casos {
		t.Run(caso.nome, func(t *testing.T) {
			tokens, _ := tokenize(caso.entrada)
			parser := &Parser{tokens: tokens}
			arvore, _ := parser.analisaExp()

			nomeArquivo := strings.ReplaceAll(caso.nome, " ", "_")
			gerarVisualizacao(arvore, nomeArquivo)

			resultado, _ := interpretar(arvore)

			if resultado != caso.esperado {
				t.Errorf("Para '%s': esperado %d, obtido %d", caso.entrada, caso.esperado, resultado)
			}
		})
	}
}
