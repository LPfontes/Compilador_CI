package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"

	"github.com/dominikbraun/graph"
	"github.com/dominikbraun/graph/draw"
)

// --- Estruturas Originais Mantidas ---

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

type Exp interface {
	isExp()
}

type Const struct {
	Valor int
}

func (c Const) isExp() {}

type OpBin struct {
	Operador string
	OpEsq    Exp
	OpDir    Exp
}

func (o OpBin) isExp() {}

// --- Lógica do Lexer  ---

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
			return nil, errors.New("Token inválido: " + token.Literal)
		}
		tokens = append(tokens, token)
	}
	tokens = append(tokens, Token{Tipo: "EOF", Literal: ""})
	return tokens, nil
}

// --- Lógica do Parser  ---

type Parser struct {
	tokens []Token
	pos    int
}

func (p *Parser) atual() Token {
	return p.tokens[p.pos]
}

func (p *Parser) consumir() Token {
	t := p.atual()
	if t.Tipo != "EOF" {
		p.pos++
	}
	return t
}

func (p *Parser) analisaExp() (Exp, error) {
	tok := p.atual()

	if tok.Tipo == "Numero" {
		p.consumir()
		valor, _ := strconv.Atoi(tok.Literal)
		return Const{Valor: valor}, nil

	} else if tok.Tipo == "ParenEsq" {
		p.consumir() // consome '('

		opEsq, err := p.analisaExp()
		if err != nil {
			return nil, err
		}

		// analisaOperador()
		operadorTok := p.consumir()
		if _, ok := tipos[operadorTok.Literal]; !ok || operadorTok.Tipo == "ParenEsq" || operadorTok.Tipo == "ParenDir" {
			return nil, fmt.Errorf("esperado operador na posição %d", operadorTok.posicao)
		}

		opDir, err := p.analisaExp()
		if err != nil {
			return nil, err
		}

		// verificaProxToken(FECHA_PARENTESE)
		if p.atual().Tipo != "ParenDir" {
			return nil, fmt.Errorf("esperado ')' na posição %d", p.atual().posicao)
		}
		p.consumir() // consome ')'

		return OpBin{
			Operador: operadorTok.Literal,
			OpEsq:    opEsq,
			OpDir:    opDir,
		}, nil
	}

	return nil, fmt.Errorf("erro sintático: token inesperado %s na posição %d", tok.Literal, tok.posicao)
}

// --- Interpretador ---

func interpretar(n Exp) (int, error) {
	switch v := n.(type) {
	case Const:
		return v.Valor, nil
	case OpBin:
		esq, err := interpretar(v.OpEsq)
		if err != nil {
			return 0, err
		}
		dir, err := interpretar(v.OpDir)
		if err != nil {
			return 0, err
		}

		switch v.Operador {
		case "+":
			return esq + dir, nil
		case "-":
			return esq - dir, nil
		case "*":
			return esq * dir, nil
		case "/":
			if dir == 0 {
				return 0, errors.New("divisão por zero")
			}
			return esq / dir, nil
		default:
			return 0, fmt.Errorf("operador desconhecido: %s", v.Operador)
		}
	}
	return 0, errors.New("tipo de nó desconhecido na árvore")
}

type graphNode struct {
	ID    int
	Label string
}

func gerarVisualizacao(root Exp) {
	g := graph.New(func(n *graphNode) int {
		return n.ID
	}, graph.Directed())

	var idCounter int
	var buildGraph func(e Exp) int

	buildGraph = func(e Exp) int {
		currentID := idCounter
		idCounter++

		switch v := e.(type) {
		case Const:
			label := strconv.Itoa(v.Valor)
			node := &graphNode{ID: currentID, Label: label}
			_ = g.AddVertex(node, graph.VertexAttribute("label", label), graph.VertexAttribute("shape", "circle"))
		case OpBin:
			label := v.Operador
			node := &graphNode{ID: currentID, Label: label}
			_ = g.AddVertex(node, graph.VertexAttribute("label", label), graph.VertexAttribute("shape", "box"))

			rightID := buildGraph(v.OpEsq)
			leftID := buildGraph(v.OpDir)

			_ = g.AddEdge(currentID, rightID)
			_ = g.AddEdge(currentID, leftID)
		}
		return currentID
	}

	buildGraph(root)

	file, err := os.Create("arvore/arvore.dot")
	if err != nil {
		fmt.Println("Erro ao criar arquivo DOT:", err)
		return
	}
	defer file.Close()

	if err := draw.DOT(g, file); err != nil {
		fmt.Println("Erro ao gerar DOT:", err)
		return
	}

	cmd := exec.Command("dot", "-Tpng", "arvore/arvore.dot", "-o", "arvore/arvore.png")
	if err := cmd.Run(); err != nil {
		fmt.Println("Arquivo 'arvore.dot' gerado. Para visualizar, instale o Graphviz e execute: dot -Tpng arvore/arvore.dot -o arvore/arvore.png")
	} else {
		fmt.Println("Imagem 'arvore.png' gerada com sucesso!")
	}
}

// --- Geração de Código Assembly ---

func gerarAssembly(n Exp) string {
	switch v := n.(type) {
	case Const:
		return fmt.Sprintf("\tpushq $%d\n", v.Valor)
	case OpBin:
		if right, ok := v.OpDir.(Const); ok {
			code := gerarAssembly(v.OpEsq)
			code += "\tpopq %rax\n"
			code += fmt.Sprintf("\tmovq $%d, %%rbx\n", right.Valor)

			switch v.Operador {
			case "+":
				code += "\taddq %rbx, %rax\n"
			case "-":
				code += "\tsubq %rbx, %rax\n"
			case "*":
				code += "\timulq %rbx, %rax\n"
			case "/":
				code += "\tcqo\n"
				code += "\tidivq %rbx\n"
			}
			code += "\tpushq %rax\n"
			return code
		}

		code := gerarAssembly(v.OpEsq)
		code += gerarAssembly(v.OpDir)

		code += "\tpopq %rbx\n"
		code += "\tpopq %rax\n"

		switch v.Operador {
		case "+":
			code += "\taddq %rbx, %rax\n"
		case "-":
			code += "\tsubq %rbx, %rax\n"
		case "*":
			code += "\timulq %rbx, %rax\n"
		case "/":
			code += "\tcqo\n"
			code += "\tidivq %rbx\n"
		}
		code += "\tpushq %rax\n"
		return code

	}
	return ""
}

func compilar(root Exp, path string) {
	assembly := ".section .text\n.globl _start\n_start:\n"
	assembly += gerarAssembly(root)
	assembly += "\tpopq %rax\n\tcall imprime_num\n\tcall sair\n\t.include \"assembly/runtime.s\"\n"

	err := os.WriteFile(path, []byte(assembly), 0644)
	if err != nil {
		fmt.Println("Erro ao gerar arquivo assembly:", err)
		return
	}
	fmt.Println("Arquivo '" + path + "' gerado com sucesso!")
}

func executarBinario(path string) {
	cmdAs := exec.Command("as", path, "-o", "output/output.o")
	if out, err := cmdAs.CombinedOutput(); err != nil {
		fmt.Printf("Erro na montagem: %s\n", out)
		return
	}

	cmdLd := exec.Command("ld", "output/output.o", "-o", "output/output")
	if out, err := cmdLd.CombinedOutput(); err != nil {
		fmt.Printf("Erro na ligação: %s\n", out)
		return
	}

	fmt.Println("--- Executando binário ---")
	cmdExec := exec.Command("./output/output")
	cmdExec.Stdout = os.Stdout
	cmdExec.Stderr = os.Stderr
	if err := cmdExec.Run(); err != nil {
		fmt.Printf("Erro na execução: %v\n", err)
	}
}

// --- Função Principal ---

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Uso: go run principal.go \"(1 + 2)\"")
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
	tokens, err := tokenize(input)
	if err != nil {
		fmt.Println("Erro Léxico:", err)
		return
	}

	parser := &Parser{tokens: tokens}
	arvore, err := parser.analisaExp()
	if err != nil {
		fmt.Println("Erro Sintático:", err)
		return
	}

	fmt.Printf("Árvore gerada %+v\n", arvore)
	gerarVisualizacao(arvore)
	resultado, err := interpretar(arvore)
	if err != nil {
		fmt.Println("Erro de Execução:", err)
		return
	} else {
		fmt.Printf("Resultado: %d\n", resultado)
	}

	compilar(arvore, "output/output.s")
	executarBinario("output/output.s")
}
