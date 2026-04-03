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
			tokens, _ := tokenize("main { return " + caso.entrada + " ; }")
			parser := &Parser{tokens: tokens}
			arvore, _ := parser.analisaPrograma()

			nomeArquivo := strings.ReplaceAll(caso.nome, " ", "_")
			gerarVisualizacao(arvore, nomeArquivo)

			resultado, _ := interpretar(arvore, make(map[string]int), nil, make(map[string]FunDecl))

			if resultado != caso.esperado {
				t.Errorf("Para '%s': esperado %d, obtido %d", caso.entrada, caso.esperado, resultado)
			}
		})
	}
}

func TestVariaveis(t *testing.T) {
	casos := []struct {
		nome     string
		entrada  string
		esperado int
	}{
		{
			nome:     "Atribuicao Simples",
			entrada:  "var x = 10; main { return x; }",
			esperado: 10,
		},
		{
			nome:     "Atribuicao Multipla",
			entrada:  "var x = 10; var y = 20; main { return x + y; }",
			esperado: 30,
		},
		{
			nome:     "Reutilizacao de Variaveis",
			entrada:  "var x = 5; var y = x * 2; var z = x + y; main { return z; }",
			esperado: 15,
		},
		{
			nome:     "Precedencia com Variaveis",
			entrada:  "var a = 2; var b = 3; var c = 4; main { return a + b * c; }",
			esperado: 14,
		},
		{
			nome:     "Uso de Parenteses com Variaveis",
			entrada:  "var a = 2; var b = 3; main { return (a + b) * 2; }",
			esperado: 10,
		},
	}

	for _, caso := range casos {
		t.Run(caso.nome, func(t *testing.T) {
			tokens, _ := tokenize(caso.entrada)
			parser := &Parser{tokens: tokens}
			arvore, _ := parser.analisaPrograma()

			nomeArquivo := strings.ReplaceAll(caso.nome, " ", "_")
			gerarVisualizacao(arvore, nomeArquivo)

			resultado, _ := interpretar(arvore, make(map[string]int), nil, make(map[string]FunDecl))

			if resultado != caso.esperado {
				t.Errorf("Para '%s': esperado %d, obtido %d", caso.entrada, caso.esperado, resultado)
			}
		})
	}
}

func TestOperadoresRelacionais(t *testing.T) {
	casos := []struct {
		nome     string
		entrada  string
		esperado int
	}{
		{
			nome:     "Maior ou Igual Verdadeiro",
			entrada:  "var x = 10; main { return x >= 10; }",
			esperado: 1,
		},
		{
			nome:     "Maior ou Igual Falso",
			entrada:  "var x = 10; main { return x >= 11; }",
			esperado: 0,
		},
		{
			nome:     "Menor ou Igual Verdadeiro",
			entrada:  "var x = 5; main { return x <= 10; }",
			esperado: 1,
		},
		{
			nome:     "Menor ou Igual Falso",
			entrada:  "var x = 5; main { return x <= 4; }",
			esperado: 0,
		},
		{
			nome:     "Maior Verdadeiro",
			entrada:  "var x = 10; main { return x > 5; }",
			esperado: 1,
		},
		{
			nome:     "Menor Verdadeiro",
			entrada:  "var x = 10; main { return x < 20; }",
			esperado: 1,
		},
		{
			nome:     "Igual Verdadeiro",
			entrada:  "var x = 10; main { return x == 10; }",
			esperado: 1,
		},
		{
			nome:     "Diferente Verdadeiro",
			entrada:  "var x = 10; main { return x != 5; }",
			esperado: 1,
		},
		{
			nome:     "Diferente Falso",
			entrada:  "var x = 10; main { return x != 10; }",
			esperado: 0,
		},
	}

	for _, caso := range casos {
		t.Run(caso.nome, func(t *testing.T) {
			tokens, _ := tokenize(caso.entrada)
			parser := &Parser{tokens: tokens}
			arvore, _ := parser.analisaPrograma()

			nomeArquivo := strings.ReplaceAll(caso.nome, " ", "_")
			gerarVisualizacao(arvore, nomeArquivo)

			resultado, _ := interpretar(arvore, make(map[string]int), nil, make(map[string]FunDecl))

			if resultado != caso.esperado {
				t.Errorf("Para '%s': esperado %d, obtido %d", caso.entrada, caso.esperado, resultado)
			}
		})
	}
}

func TestOperadoresLogicos(t *testing.T) {
	casos := []struct {
		nome     string
		entrada  string
		esperado int
	}{
		{
			nome:     "Not Falso",
			entrada:  "var x = 5; main { return not x; }",
			esperado: 0,
		},
		{
			nome:     "Not Verdadeiro",
			entrada:  "var x = 0; main { return not x; }",
			esperado: 1,
		},
		{
			nome:     "And Verdadeiro",
			entrada:  "var x = 1; var y = 2; main { return x and y; }",
			esperado: 1,
		},
		{
			nome:     "And Falso",
			entrada:  "var x = 1; var y = 0; main { return x and y; }",
			esperado: 0,
		},
		{
			nome:     "Or Verdadeiro",
			entrada:  "var x = 0; var y = 2; main { return x or y; }",
			esperado: 1,
		},
		{
			nome:     "Or Falso",
			entrada:  "var x = 0; var y = 0; main { return x or y; }",
			esperado: 0,
		},
		{
			nome:     "Expressao Cruzada",
			entrada:  "main { return 1 or 0 and 0; }",
			esperado: 1,
		},
	}

	for _, caso := range casos {
		t.Run(caso.nome, func(t *testing.T) {
			tokens, _ := tokenize(caso.entrada)
			parser := &Parser{tokens: tokens}
			arvore, _ := parser.analisaPrograma()
			
			resultado, _ := interpretar(arvore, make(map[string]int), nil, make(map[string]FunDecl))

			if resultado != caso.esperado {
				t.Errorf("Para '%s': esperado %d, obtido %d", caso.entrada, caso.esperado, resultado)
			}
		})
	}
}

func TestFuncoes(t *testing.T) {
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

	for _, caso := range casos {
		t.Run(caso.nome, func(t *testing.T) {
			tokens, _ := tokenize(caso.entrada)
			parser := &Parser{tokens: tokens}
			arvore, _ := parser.analisaPrograma()
			resultado, _ := interpretar(arvore, make(map[string]int), nil, make(map[string]FunDecl))
			if resultado != caso.esperado {
				t.Errorf("Esperado %d, obtido %d", caso.esperado, resultado)
			}
		})
	}
}

func TestVetores(t *testing.T) {
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
			nome:     "Vetor com Loop While",
			entrada:  "var v[5]; var i = 0; main { while i < 5 { v[i] = i * 10; i = i + 1; } return v[3]; }",
			esperado: 30,
		},
		{
			nome:     "Vetor Local com Parametro como Indice",
			entrada:  "fun get(idx) { var arr[3]; arr[0] = 100; arr[1] = 200; arr[2] = 300; return arr[idx]; } main { return get(1); }",
			esperado: 200,
		},
	}

	for _, caso := range casos {
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
			resultado, err := interpretar(arvore, make(map[string]int), nil, make(map[string]FunDecl))
			if err != nil {
				t.Fatalf("Erro no interpretador: %v", err)
			}
			if resultado != caso.esperado {
				t.Errorf("Esperado %d, obtido %d", caso.esperado, resultado)
			}
		})
	}
}

func TestVetoresAVX(t *testing.T) {
	casos := []struct {
		nome     string
		entrada  string
		esperado int
	}{
		{
			nome:     "vadd 4 elementos",
			entrada:  "var a[4]; var b[4]; var c[4]; main { a[0] = 1; a[1] = 2; a[2] = 3; a[3] = 4; b[0] = 10; b[1] = 20; b[2] = 30; b[3] = 40; vadd(c, a, b, 4); return c[2]; }",
			esperado: 33,
		},
		{
			nome:     "vsub 4 elementos",
			entrada:  "var a[4]; var b[4]; var c[4]; main { a[0] = 100; a[1] = 200; a[2] = 300; a[3] = 400; b[0] = 1; b[1] = 2; b[2] = 3; b[3] = 4; vsub(c, a, b, 4); return c[3]; }",
			esperado: 396,
		},
		{
			nome:     "vset 4 elementos",
			entrada:  "var v[4]; main { vset(v, 42, 4); return v[0] + v[3]; }",
			esperado: 84,
		},
		{
			nome:     "vadd soma total",
			entrada:  "var a[4]; var b[4]; var c[4]; main { a[0] = 5; a[1] = 10; a[2] = 15; a[3] = 20; b[0] = 1; b[1] = 1; b[2] = 1; b[3] = 1; vadd(c, a, b, 4); return c[0] + c[1] + c[2] + c[3]; }",
			esperado: 54,
		},
	}

	for _, caso := range casos {
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
			resultado, err := interpretar(arvore, make(map[string]int), nil, make(map[string]FunDecl))
			if err != nil {
				t.Fatalf("Erro no interpretador: %v", err)
			}
			if resultado != caso.esperado {
				t.Errorf("Esperado %d, obtido %d", caso.esperado, resultado)
			}
		})
	}
}

func TestFor(t *testing.T) {
	casos := []struct {
		nome     string
		entrada  string
		esperado int
	}{
		{
			nome:     "For Simples Soma",
			entrada:  "var x = 0; var i = 0; main { for i = 0; i < 5; i = i + 1 { x = x + 1; } return x; }",
			esperado: 5,
		},
		{
			nome:     "For Acumulador",
			entrada:  "var soma = 0; var i = 0; main { for i = 1; i <= 10; i = i + 1 { soma = soma + i; } return soma; }",
			esperado: 55,
		},
		{
			nome:     "For com Vetor",
			entrada:  "var v[4]; var i = 0; main { for i = 0; i < 4; i = i + 1 { v[i] = i * 10; } return v[2]; }",
			esperado: 20,
		},
	}

	for _, caso := range casos {
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
			resultado, err := interpretar(arvore, make(map[string]int), nil, make(map[string]FunDecl))
			if err != nil {
				t.Fatalf("Erro no interpretador: %v", err)
			}
			if resultado != caso.esperado {
				t.Errorf("Esperado %d, obtido %d", caso.esperado, resultado)
			}
		})
	}
}
