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

func gerarVisualizacao(root Exp, nomeBase string) {
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

			leftID := buildGraph(v.OpEsq)
			rightID := buildGraph(v.OpDir)

			_ = g.AddEdge(currentID, leftID)
			_ = g.AddEdge(currentID, rightID)
		}
		return currentID
	}

	buildGraph(root)

	if _, err := os.Stat("arvore"); os.IsNotExist(err) {
		_ = os.Mkdir("arvore", 0755)
	}

	dotPath := "arvore/" + nomeBase + ".dot"
	pngPath := "arvore/" + nomeBase + ".png"
	file, err := os.Create(dotPath)
	if err != nil {
		fmt.Println("Erro ao criar arquivo DOT:", err)
		return
	}
	defer file.Close()

	if err := draw.DOT(g, file); err != nil {
		fmt.Println("Erro ao gerar DOT:", err)
		return
	}

	cmd := exec.Command("dot", "-Tpng", dotPath, "-o", pngPath)
	if err := cmd.Run(); err != nil {
		fmt.Printf("Arquivo '%s' gerado. Para visualizar, instale o Graphviz e execute: dot -Tpng %s -o %s\n", dotPath, dotPath, pngPath)
	} else {
		fmt.Printf("Imagem '%s' gerada com sucesso!\n", pngPath)
	}
}
