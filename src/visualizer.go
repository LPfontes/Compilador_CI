package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"

	"github.com/dominikbraun/graph"
	"github.com/dominikbraun/graph/draw"
)

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