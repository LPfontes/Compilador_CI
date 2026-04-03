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
		case Var:
			file.WriteString(fmt.Sprintf("    %d [label=\"Var(%s)\", shape=ellipse];\n", id, v.Nome))
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
		case Programa:
			file.WriteString(fmt.Sprintf("    %d [label=\"Programa\", shape=invhouse];\n", id))
			for i, decl := range v.Globais {
				declID := idCounter
				idCounter++
				file.WriteString(fmt.Sprintf("    %d [label=\"Decl(%s)\", shape=rectangle];\n", declID, decl.Nome))
				file.WriteString(fmt.Sprintf("    %d -> %d [label=\"decl %d\"];\n", id, declID, i))
				expID := writeNode(decl.Expressao)
				file.WriteString(fmt.Sprintf("    %d -> %d [label=\"exp\"];\n", declID, expID))
			}
			for i, fun := range v.Funcoes {
				funID := idCounter
				idCounter++
				file.WriteString(fmt.Sprintf("    %d [label=\"Fun(%s)\", shape=hexagon];\n", funID, fun.Nome))
				file.WriteString(fmt.Sprintf("    %d -> %d [label=\"func %d\"];\n", id, funID, i))
				resID := writeNode(fun.Resultado)
				file.WriteString(fmt.Sprintf("    %d -> %d [label=\"ret\"];\n", funID, resID))
			}
			resID := writeNode(v.Resultado)
			file.WriteString(fmt.Sprintf("    %d -> %d [label=\"resultado\"];\n", id, resID))
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
