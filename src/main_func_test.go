package main

import (
	"os/exec"
	"strconv"
	"strings"
	"testing"
)

func TestCompiladorFuncoes(t *testing.T) {
	casos := []struct {
		nome     string
		entrada  string
		esperado int
	}{
		{
			nome:     "Funcao Simples",
			entrada:  "fun soma(a, b) { return a + b; } main { return soma(10, 20); }",
			esperado: 30,
		},
		{
			nome:     "Funcao Com Variavel Local",
			entrada:  "fun sub(a, b) { var x = a - b; return x; } main { return sub(50, 20); }",
			esperado: 30,
		},
		{
			nome:     "Funcao Shadowing Global",
			entrada:  "var x = 100; fun get(x) { return x; } main { return get(50); }",
			esperado: 50,
		},
		{
			nome:     "Funcao Local Modificada",
			entrada:  "fun teste(x) { var y = 5; if x > 0 { y = 10; } else { y = 0; } return y; } main { return teste(1); }",
			esperado: 10,
		},
		{
			nome:     "Funcao Chamando Funcao",
			entrada:  "fun duplo(x) { return x * 2; } fun quadruplo(x) { return duplo(duplo(x)); } main { return quadruplo(5); }",
			esperado: 20,
		},
		{
			nome: "Fibonacci Recursao",
			entrada: "fun fib(n) { var res = 0; if n < 2 { res = 1; } else { res = fib(n - 1) + fib(n - 2); } return res; } main { return fib(5); }",
			esperado: 8,
		},
	}

	for i, caso := range casos {
		t.Run(caso.nome, func(t *testing.T) {
			tokens, _ := tokenize(caso.entrada)
			parser := &Parser{tokens: tokens}
			arvore, err := parser.analisaPrograma()
			if err != nil {
				t.Fatalf("Erro no parser: %v", err)
			}
			
			err = AnalisarSemantica(arvore)
			if err != nil {
				t.Fatalf("Erro Semantico: %v", err)
			}

			path := "output/output_func" + strconv.Itoa(i) + ".s"
			compilar(arvore, path)

			if out, err := exec.Command("as", path, "-o", "output/output_func.o").CombinedOutput(); err != nil {
				t.Fatalf("Erro na montagem: %s", out)
			}
			if out, err := exec.Command("ld", "output/output_func.o", "-o", "output/output_func").CombinedOutput(); err != nil {
				t.Fatalf("Erro na ligacao: %s", out)
			}

			out, err := exec.Command("./output/output_func").CombinedOutput()
			if err != nil {
				t.Fatalf("Erro na execucao: %v (Saida: %s)", err, out)
			}
			
			saidaLimpa := strings.TrimSpace(string(out))
			esperadoStr := strconv.Itoa(caso.esperado)
			if saidaLimpa != esperadoStr {
				t.Errorf("Esperado %s, obtido %s", esperadoStr, saidaLimpa)
			}
		})
	}
}

func TestCompiladorVetores(t *testing.T) {
	casos := []struct {
		nome     string
		entrada  string
		esperado int
	}{
		{
			nome:     "Vetor Global Escrita e Leitura",
			entrada:  "var v[5]; main { v[0] = 10; v[1] = 20; return v[0] + v[1]; }",
			esperado: 30,
		},
		{
			nome:     "Vetor Global Indice Dinamico",
			entrada:  "var v[3]; var i = 2; main { v[i] = 99; return v[2]; }",
			esperado: 99,
		},
		{
			nome:     "Vetor Local em Funcao",
			entrada:  "fun soma3() { var a[3]; a[0] = 10; a[1] = 20; a[2] = 30; return a[0] + a[1] + a[2]; } main { return soma3(); }",
			esperado: 60,
		},
		{
			nome:     "Vetor Local com Parametro",
			entrada:  "fun get(idx) { var arr[3]; arr[0] = 100; arr[1] = 200; arr[2] = 300; return arr[idx]; } main { return get(1); }",
			esperado: 200,
		},
	}

	for i, caso := range casos {
		t.Run(caso.nome, func(t *testing.T) {
			tokens, err := tokenize(caso.entrada)
			if err != nil {
				t.Fatalf("Erro léxico: %v", err)
			}
			parser := &Parser{tokens: tokens}
			arvore, err := parser.analisaPrograma()
			if err != nil {
				t.Fatalf("Erro no parser: %v", err)
			}

			err = AnalisarSemantica(arvore)
			if err != nil {
				t.Fatalf("Erro Semantico: %v", err)
			}

			path := "output/output_vec" + strconv.Itoa(i) + ".s"
			compilar(arvore, path)

			if out, err := exec.Command("as", path, "-o", "output/output_vec.o").CombinedOutput(); err != nil {
				t.Fatalf("Erro na montagem: %s\nASM: %s", out, path)
			}
			if out, err := exec.Command("ld", "output/output_vec.o", "-o", "output/output_vec").CombinedOutput(); err != nil {
				t.Fatalf("Erro na ligacao: %s", out)
			}

			out, err := exec.Command("./output/output_vec").CombinedOutput()
			if err != nil {
				t.Fatalf("Erro na execucao: %v (Saida: %s)", err, out)
			}

			saidaLimpa := strings.TrimSpace(string(out))
			esperadoStr := strconv.Itoa(caso.esperado)
			if saidaLimpa != esperadoStr {
				t.Errorf("Esperado %s, obtido %s", esperadoStr, saidaLimpa)
			}
		})
	}
}

