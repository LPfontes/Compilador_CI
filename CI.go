package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Uso: go run CI.go <arquivo_de_entrada>")
		fmt.Println("Por favor, forneça o nome do arquivo de entrada.")
		return
	}

	// 1 Etapa de análise: Ler todo o conteúdo do arquivo.
	inputFilename := os.Args[1]
	content, err := os.ReadFile(inputFilename)
	if err != nil {
		fmt.Printf("Erro ao ler o arquivo '%s': %v\n", inputFilename, err)
		return
	}

	// Limpa espaços em branco e quebras de linha antes de converter.
	numberStr := strings.TrimSpace(string(content))
	literal, err := strconv.Atoi(numberStr)
	if err != nil {
		fmt.Printf("Erro de conversão: o conteúdo do arquivo ('%s') não é um inteiro válido.\n", numberStr)
		return
	}

	// 2 Etapa de síntese: Gerar o arquivo .s
	fileOutput, err := os.Create("output.s")
	if err != nil {
		fmt.Println(err)
		return
	}

	defer fileOutput.Close()

	assemblyTemplate := `
.section .text
.globl _start
_start:
	mov $%d, %%rax
	call imprime_num
	call sair
	.include "runtime.s"`
	data := []byte(fmt.Sprintf(assemblyTemplate, literal))

	_, err = fileOutput.Write(data)
	if err != nil {
		fmt.Println("Erro ao escrever no arquivo:", err)
		return
	}

	fileOutput.Close()
	fmt.Println("Arquivo output.s gerado com sucesso!")

	// 3 Etapa de montagem (as)
	fmt.Println("Montando o código com 'as'...")
	cmdAs := exec.Command("as", "-o", "output.o", "output.s")
	output, err := cmdAs.CombinedOutput()
	if err != nil {
		fmt.Printf("Erro ao montar o código assembly:\n%s\n", string(output))
		return
	}

	// 4 Etapa de ligação (ld)
	fmt.Println("Ligando o objeto com 'ld'...")
	cmdLd := exec.Command("ld", "-o", "output", "output.o")
	output, err = cmdLd.CombinedOutput()
	if err != nil {
		fmt.Printf("Erro ao ligar o arquivo objeto:\n%s\n", string(output))
		return
	}

	fmt.Println("Executável 'output' criado com sucesso! Você pode executá-lo com './output'")
}