package main

import (
	"os"
	"os/exec"
	"strconv"
	"strings"
	"testing"
)

func TestInterpretador(t *testing.T) {
	casos := []struct {
		nome     string
		entrada  string
		esperado int
		temErro  bool
	}{
		{"Numero Simples", "42", 42, false},
		{"Soma", "(10 + 5)", 15, false},
		{"Subtracao", "(20 - 8)", 12, false},
		{"Multiplicacao", "(6 * 7)", 42, false},
		{"Divisao", "(100 / 2)", 50, false},
		{"Expressao Aninhada Esq", "((2 + 3) * 4)", 20, false},
		{"Expressao Aninhada Dir", "(2 + (3 * 4))", 14, false},
		{"Expressao Complexa", "((10 * 2) + (30 / 3))", 30, false},
		{"Subtracao Negativa", "(5 - 10)", -5, false},
		{"Divisao por Zero", "(10 / 0)", 0, true},
	}

	for _, caso := range casos {
		t.Run(caso.nome, func(t *testing.T) {
			// 1. Tokenização
			tokens, err := tokenize(caso.entrada)
			if err != nil {
				t.Fatalf("Erro na tokenização: %v", err)
			}

			// 2. Parsing
			parser := &Parser{tokens: tokens}
			arvore, err := parser.analisaExp()
			if err != nil {
				t.Fatalf("Erro no parser: %v", err)
			}

			// 3. Interpretação
			resultado, err := interpretar(arvore)

			if caso.temErro {
				if err == nil {
					t.Errorf("Esperava erro para entrada '%s', mas obteve sucesso: %d", caso.entrada, resultado)
				}
			} else if resultado != caso.esperado {
				t.Errorf("Para '%s': esperado %d, obtido %d", caso.entrada, caso.esperado, resultado)
			}
		})
	}
}

func TestCompilador(t *testing.T) {
	casos := []struct {
		nome     string
		entrada  string
		esperado int
		temErro  bool
	}{
		{"Numero Simples", "42", 42, false},
		{"Soma", "(10 + 5)", 15, false},
		{"Subtracao", "(20 - 8)", 12, false},
		{"Multiplicacao", "(6 * 7)", 42, false},
		{"Divisao", "(100 / 2)", 50, false},
		{"Expressao Aninhada Esq", "((2 + 3) * 4)", 20, false},
		{"Expressao Aninhada Dir", "(2 + (3 * 4))", 14, false},
		{"Expressao Complexa", "((10 * 2) + (30 / 3))", 30, false},
		{"Subtracao Negativa", "(5 - 10)", -5, false},
		{"Divisao por Zero", "(10 / 0)", 0, true},
	}

	// Garante que o diretório de saída existe
	if _, err := os.Stat("output"); os.IsNotExist(err) {
		_ = os.Mkdir("output", 0755)
	}
	index := 0
	for _, caso := range casos {
		t.Run(caso.nome, func(t *testing.T) {
			tokens, _ := tokenize(caso.entrada)
			parser := &Parser{tokens: tokens}
			arvore, _ := parser.analisaExp()

			// 1. Gera o Assembly (output/output.s)
			compilar(arvore, "output/output"+strconv.Itoa(index)+".s")

			// 2. Montagem (as)
			if out, err := exec.Command("as", "output/output"+strconv.Itoa(index)+".s", "-o", "output/output.o").CombinedOutput(); err != nil {
				t.Fatalf("Erro na montagem: %s", out)
			}
			index++
			// 3. Ligação (ld)
			if out, err := exec.Command("ld", "output/output.o", "-o", "output/output").CombinedOutput(); err != nil {
				t.Fatalf("Erro na ligação: %s", out)
			}

			// 4. Execução do Binário
			out, err := exec.Command("./output/output").CombinedOutput()

			if caso.temErro {
				if err == nil {
					t.Errorf("Esperava erro de execução (ex: divisão por zero), mas obteve sucesso. Saída: %s", out)
				}
			} else {
				if err != nil {
					t.Fatalf("Erro na execução do binário: %v. Saída: %s", err, out)
				}
				saidaLimpa := strings.TrimSpace(string(out))
				esperadoStr := strconv.Itoa(caso.esperado)
				if saidaLimpa != esperadoStr {
					t.Errorf("Esperado %s, obtido %s", esperadoStr, saidaLimpa)
				}
			}
		})
	}
}
