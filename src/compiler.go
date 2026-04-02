package main

import (
	"fmt"
	"os"
	"os/exec"
)

func getAsmInstruction(op string) string {
	switch op {
	case "+":
		return "\taddq %rbx, %rax\n"
	case "-":
		return "\tsubq %rbx, %rax\n"
	case "*":
		return "\timulq %rbx, %rax\n"
	case "/":
		return "\tcqo\n\tidivq %rbx\n"
	}
	return ""
}

func gerarAssembly(n Exp) string {
	switch v := n.(type) {
	case Var:
		return fmt.Sprintf("\tmovq %s(%%rip), %%rax\n\tpushq %%rax\n", v.Nome)
	case Const:
		return fmt.Sprintf("\tpushq $%d\n", v.Valor)
	case OpBin:
		if right, ok := v.OpDir.(Const); ok {
			code := gerarAssembly(v.OpEsq)
			code += "\tpopq %rax\n"
			code += fmt.Sprintf("\tmovq $%d, %%rbx\n", right.Valor)

			code += getAsmInstruction(v.Operador)
			code += "\tpushq %rax\n"
			return code
		}

		code := gerarAssembly(v.OpEsq)
		code += gerarAssembly(v.OpDir)

		code += "\tpopq %rbx\n"
		code += "\tpopq %rax\n"

		code += getAsmInstruction(v.Operador)
		code += "\tpushq %rax\n"
		return code

	}
	return ""
}

func compilar(root Programa, path string) {
	bss := ""
	if len(root.Declaracoes) > 0 {
		bss = ".section .bss\n"
		tabela := make(map[string]bool)
		for _, decl := range root.Declaracoes {
			if !tabela[decl.Nome] {
				bss += fmt.Sprintf("\t.lcomm %s, 8\n", decl.Nome)
				tabela[decl.Nome] = true
			}
		}
	}

	assembly := ""
	if bss != "" {
		assembly += bss + "\n"
	}
	assembly += ".section .text\n.globl _start\n_start:\n"
	
	for _, decl := range root.Declaracoes {
		assembly += gerarAssembly(decl.Expressao)
		assembly += "\tpopq %rax\n"
		assembly += fmt.Sprintf("\tmovq %%rax, %s(%%rip)\n", decl.Nome)
	}

	assembly += gerarAssembly(root.Resultado)
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
