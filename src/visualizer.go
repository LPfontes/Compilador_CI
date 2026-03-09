package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
)

func gerarVisualizacao(root Exp, nomeBase string) {
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

	file.WriteString("digraph G {\n")

	var idCounter int
	var writeNode func(e Exp) int

	writeNode = func(e Exp) int {
		id := idCounter
		idCounter++

		switch v := e.(type) {
		case Const:
			label := strconv.Itoa(v.Valor)
			file.WriteString(fmt.Sprintf("    %d [label=\"%s\", shape=circle];\n", id, label))
		case OpBin:
			label := v.Operador
			file.WriteString(fmt.Sprintf("    %d [label=\"%s\", shape=box];\n", id, label))

			leftID := writeNode(v.OpEsq)
			rightID := writeNode(v.OpDir)

			file.WriteString(fmt.Sprintf("    %d -> %d;\n", id, leftID))
			file.WriteString(fmt.Sprintf("    %d -> %d;\n", id, rightID))
		}
		return id
	}

	writeNode(root)
	file.WriteString("}\n")

	cmd := exec.Command("dot", "-Tpng", dotPath, "-o", pngPath)
	if err := cmd.Run(); err != nil {
		fmt.Printf("Arquivo '%s' gerado. Para visualizar, instale o Graphviz e execute: dot -Tpng %s -o %s\n", dotPath, dotPath, pngPath)
	} else {
		fmt.Printf("Imagem '%s' gerada com sucesso!\n", pngPath)
	}
}
