package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
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
	case "and":
		return "\tcmpq $0, %rax\n\tsetne %al\n\tcmpq $0, %rbx\n\tsetne %bl\n\tandb %bl, %al\n\tmovzbq %al, %rax\n"
	case "or":
		return "\tcmpq $0, %rax\n\tsetne %al\n\tcmpq $0, %rbx\n\tsetne %bl\n\torb %bl, %al\n\tmovzbq %al, %rax\n"
	case "==":
		return "\txorq %rcx, %rcx\n\tcmpq %rbx, %rax\n\tsetz %cl\n\tmovq %rcx, %rax\n"
	case "<=":
		return "\txorq %rcx, %rcx\n\tcmpq %rbx, %rax\n\tsetle %cl\n\tmovq %rcx, %rax\n"
	case ">=":
		return "\txorq %rcx, %rcx\n\tcmpq %rbx, %rax\n\tsetge %cl\n\tmovq %rcx, %rax\n"
	case "<":
		return "\txorq %rcx, %rcx\n\tcmpq %rbx, %rax\n\tsetl %cl\n\tmovq %rcx, %rax\n"
	case ">":
		return "\txorq %rcx, %rcx\n\tcmpq %rbx, %rax\n\tsetg %cl\n\tmovq %rcx, %rax\n"
	case "!=":
		return "\txorq %rcx, %rcx\n\tcmpq %rbx, %rax\n\tsetne %cl\n\tmovq %rcx, %rax\n"
	}
	return ""
}

var rotuloCount int

func novoRotulo(prefixo string) string {
	rotuloCount++
	return fmt.Sprintf("%s%d", prefixo, rotuloCount)
}

func gerarCmd(cmd Cmd, locais map[string]int) string {
	switch v := cmd.(type) {
	case AtribCmd:
		code := gerarAssembly(v.Expressao, locais)
		code += "\tpopq %rax\n"
		if offset, ok := locais[v.Nome]; ok {
			code += fmt.Sprintf("\tmovq %%rax, %d(%%rbp)\n", offset)
		} else {
			code += fmt.Sprintf("\tmovq %%rax, %s(%%rip)\n", v.Nome)
		}
		return code
	case AtribVetorCmd:
		code := gerarAssembly(v.Indice, locais)
		code += gerarAssembly(v.Expressao, locais)
		code += "\tpopq %rax\n"  // valor
		code += "\tpopq %rcx\n"  // indice
		if offset, ok := locais[v.Nome]; ok {
			code += fmt.Sprintf("\tleaq %d(%%rbp), %%rdx\n", offset)
			code += "\tmovq %rax, (%rdx, %rcx, 8)\n"
		} else {
			code += fmt.Sprintf("\tleaq %s(%%rip), %%rdx\n", v.Nome)
			code += "\tmovq %rax, (%rdx, %rcx, 8)\n"
		}
		return code
	case IfCmd:
		lFalso := novoRotulo("Lfalso")
		lFim := novoRotulo("Lfim")
		code := gerarAssembly(v.Condicao, locais)
		code += "\tpopq %rax\n"
		code += "\tcmpq $0, %rax\n"
		code += fmt.Sprintf("\tjz %s\n", lFalso)
		for _, c := range v.CorpoIf {
			code += gerarCmd(c, locais)
		}
		code += fmt.Sprintf("\tjmp %s\n", lFim)
		code += fmt.Sprintf("%s:\n", lFalso)
		for _, c := range v.CorpoElse {
			code += gerarCmd(c, locais)
		}
		code += fmt.Sprintf("%s:\n", lFim)
		return code
	case WhileCmd:
		lInicio := novoRotulo("Linicio")
		lFim := novoRotulo("Lfim")
		code := fmt.Sprintf("%s:\n", lInicio)
		code += gerarAssembly(v.Condicao, locais)
		code += "\tpopq %rax\n"
		code += "\tcmpq $0, %rax\n"
		code += fmt.Sprintf("\tjz %s\n", lFim)
		for _, c := range v.Corpo {
			code += gerarCmd(c, locais)
		}
		code += fmt.Sprintf("\tjmp %s\n", lInicio)
		code += fmt.Sprintf("%s:\n", lFim)
		return code
	case ForCmd:
		code := gerarCmd(v.Init, locais)
		lInicio := novoRotulo("Lfor")
		lFim := novoRotulo("Lfor_fim")
		code += fmt.Sprintf("%s:\n", lInicio)
		code += gerarAssembly(v.Condicao, locais)
		code += "\tpopq %rax\n"
		code += "\tcmpq $0, %rax\n"
		code += fmt.Sprintf("\tjz %s\n", lFim)
		for _, c := range v.Corpo {
			code += gerarCmd(c, locais)
		}
		code += gerarCmd(v.Passo, locais)
		code += fmt.Sprintf("\tjmp %s\n", lInicio)
		code += fmt.Sprintf("%s:\n", lFim)
		return code
	case CmdBuiltinVec:
		code := ""
		destNome := v.Args[0].(Var).Nome

		if v.Operacao == "vset" {
			// vset(dest, valor, n)
			code += gerarAssembly(v.Args[1], locais) // valor na pilha
			code += "\tpopq %rax\n"
			code += "\tvpbroadcastq %rax, %ymm0\n"
			if offset, ok := locais[destNome]; ok {
				code += fmt.Sprintf("\tleaq %d(%%rbp), %%rdx\n", offset)
			} else {
				code += fmt.Sprintf("\tleaq %s(%%rip), %%rdx\n", destNome)
			}
			nConst, isConst := v.Args[2].(Const)
			if isConst && nConst.Valor == 4 {
				code += "\tvmovdqu %ymm0, (%rdx)\n"
			} else {
				code += gerarAssembly(v.Args[2], locais)
				code += "\tpopq %r8\n"
				lLoop := novoRotulo("Lvset")
				lFim := novoRotulo("Lvset_fim")
				code += "\txorq %rcx, %rcx\n"
				code += fmt.Sprintf("%s:\n", lLoop)
				code += "\tcmpq %r8, %rcx\n"
				code += fmt.Sprintf("\tjge %s\n", lFim)
				code += "\tvmovdqu %ymm0, (%rdx, %rcx, 8)\n"
				code += "\taddq $4, %rcx\n"
				code += fmt.Sprintf("\tjmp %s\n", lLoop)
				code += fmt.Sprintf("%s:\n", lFim)
			}
		} else {
			// vadd / vsub
			src1Nome := v.Args[1].(Var).Nome
			src2Nome := v.Args[2].(Var).Nome
			var avxInstr string
			if v.Operacao == "vadd" {
				avxInstr = "vpaddq"
			} else {
				avxInstr = "vpsubq"
			}

			nConst, isConst := v.Args[3].(Const)
			if isConst && nConst.Valor == 4 {
				// Caso simples: exatamente 4 elementos
				if offset, ok := locais[src1Nome]; ok {
					code += fmt.Sprintf("\tleaq %d(%%rbp), %%rax\n", offset)
				} else {
					code += fmt.Sprintf("\tleaq %s(%%rip), %%rax\n", src1Nome)
				}
				code += "\tvmovdqu (%rax), %ymm0\n"
				if offset, ok := locais[src2Nome]; ok {
					code += fmt.Sprintf("\tleaq %d(%%rbp), %%rax\n", offset)
				} else {
					code += fmt.Sprintf("\tleaq %s(%%rip), %%rax\n", src2Nome)
				}
				code += "\tvmovdqu (%rax), %ymm1\n"
				code += fmt.Sprintf("\t%s %%ymm1, %%ymm0, %%ymm2\n", avxInstr)
				if offset, ok := locais[destNome]; ok {
					code += fmt.Sprintf("\tleaq %d(%%rbp), %%rax\n", offset)
				} else {
					code += fmt.Sprintf("\tleaq %s(%%rip), %%rax\n", destNome)
				}
				code += "\tvmovdqu %ymm2, (%rax)\n"
			} else {
				// Loop para N > 4
				code += gerarAssembly(v.Args[3], locais)
				code += "\tpopq %r8\n"
				lLoop := novoRotulo("Lavx")
				lFim := novoRotulo("Lavx_fim")
				code += "\txorq %rcx, %rcx\n"
				code += fmt.Sprintf("%s:\n", lLoop)
				code += "\tcmpq %r8, %rcx\n"
				code += fmt.Sprintf("\tjge %s\n", lFim)
				if offset, ok := locais[src1Nome]; ok {
					code += fmt.Sprintf("\tleaq %d(%%rbp), %%rax\n", offset)
				} else {
					code += fmt.Sprintf("\tleaq %s(%%rip), %%rax\n", src1Nome)
				}
				code += "\tvmovdqu (%rax, %rcx, 8), %ymm0\n"
				if offset, ok := locais[src2Nome]; ok {
					code += fmt.Sprintf("\tleaq %d(%%rbp), %%rax\n", offset)
				} else {
					code += fmt.Sprintf("\tleaq %s(%%rip), %%rax\n", src2Nome)
				}
				code += "\tvmovdqu (%rax, %rcx, 8), %ymm1\n"
				code += fmt.Sprintf("\t%s %%ymm1, %%ymm0, %%ymm2\n", avxInstr)
				if offset, ok := locais[destNome]; ok {
					code += fmt.Sprintf("\tleaq %d(%%rbp), %%rax\n", offset)
				} else {
					code += fmt.Sprintf("\tleaq %s(%%rip), %%rax\n", destNome)
				}
				code += "\tvmovdqu %ymm2, (%rax, %rcx, 8)\n"
				code += "\taddq $4, %rcx\n"
				code += fmt.Sprintf("\tjmp %s\n", lLoop)
				code += fmt.Sprintf("%s:\n", lFim)
			}
		}
		code += "\tvzeroupper\n"
		return code
	}
	return ""
}

func gerarAssembly(n Exp, locais map[string]int) string {
	switch v := n.(type) {
	case Var:
		if offset, ok := locais[v.Nome]; ok {
			return fmt.Sprintf("\tmovq %d(%%rbp), %%rax\n\tpushq %%rax\n", offset)
		}
		return fmt.Sprintf("\tmovq %s(%%rip), %%rax\n\tpushq %%rax\n", v.Nome)
	case AcessoVetor:
		code := gerarAssembly(v.Indice, locais)
		code += "\tpopq %rcx\n" // indice em rcx
		if offset, ok := locais[v.Nome]; ok {
			code += fmt.Sprintf("\tleaq %d(%%rbp), %%rdx\n", offset)
			code += "\tmovq (%rdx, %rcx, 8), %rax\n"
		} else {
			code += fmt.Sprintf("\tleaq %s(%%rip), %%rdx\n", v.Nome)
			code += "\tmovq (%rdx, %rcx, 8), %rax\n"
		}
		code += "\tpushq %rax\n"
		return code
	case Const:
		return fmt.Sprintf("\tpushq $%d\n", v.Valor)
	case ChamadaFun:
		code := ""
		for i := len(v.Args) - 1; i >= 0; i-- {
			code += gerarAssembly(v.Args[i], locais)
		}
		code += fmt.Sprintf("\tcall %s\n", v.Nome)
		if len(v.Args) > 0 {
			code += fmt.Sprintf("\taddq $%d, %%rsp\n", len(v.Args)*8)
		}
		code += "\tpushq %rax\n"
		return code
	case OpBin:
		if right, ok := v.OpDir.(Const); ok {
			code := gerarAssembly(v.OpEsq, locais)
			code += "\tpopq %rax\n"
			code += fmt.Sprintf("\tmovq $%d, %%rbx\n", right.Valor)

			code += getAsmInstruction(v.Operador)
			code += "\tpushq %rax\n"
			return code
		}

		code := gerarAssembly(v.OpEsq, locais)
		code += gerarAssembly(v.OpDir, locais)

		code += "\tpopq %rbx\n"
		code += "\tpopq %rax\n"

		code += getAsmInstruction(v.Operador)
		code += "\tpushq %rax\n"
		return code

	case OpUnario:
		if v.Operador == "not" {
			code := gerarAssembly(v.Expressao, locais)
			code += "\tpopq %rax\n"
			code += "\tcmpq $0, %rax\n"
			code += "\tsetz %al\n"
			code += "\tmovzbq %al, %rax\n"
			code += "\tpushq %rax\n"
			return code
		}

	}
	return ""
}

func gerarFunDecl(f FunDecl) string {
	totalSlots := 0
	locais := make(map[string]int)

	for _, v := range f.Vars {
		locais[v.Nome] = totalSlots * 8
		if v.Tamanho > 0 {
			totalSlots += v.Tamanho
		} else {
			totalSlots++
		}
	}
	for i, p := range f.Params {
		locais[p] = (totalSlots * 8) + 16 + (i * 8)
	}

	code := fmt.Sprintf("%s:\n", f.Nome)
	code += "\tpushq %rbp\n"
	if totalSlots > 0 {
		code += fmt.Sprintf("\tsubq $%d, %%rsp\n", totalSlots*8)
	}
	code += "\tmovq %rsp, %rbp\n"

	for _, decl := range f.Vars {
		if decl.Tamanho > 0 {
			// Arrays: zerar os slots
			baseOffset := locais[decl.Nome]
			for j := 0; j < decl.Tamanho; j++ {
				code += fmt.Sprintf("\tmovq $0, %d(%%rbp)\n", baseOffset+j*8)
			}
		} else {
			code += gerarAssembly(decl.Expressao, locais)
			code += "\tpopq %rax\n"
			offset := locais[decl.Nome]
			code += fmt.Sprintf("\tmovq %%rax, %d(%%rbp)\n", offset)
		}
	}
	for _, cmd := range f.Comandos {
		code += gerarCmd(cmd, locais)
	}

	code += gerarAssembly(f.Resultado, locais)
	code += "\tpopq %rax\n"

	if totalSlots > 0 {
		code += fmt.Sprintf("\taddq $%d, %%rsp\n", totalSlots*8)
	}
	code += "\tpopq %rbp\n"
	code += "\tret\n"
	return code
}

func compilar(root Programa, path string) {
	bss := ""
	if len(root.Globais) > 0 {
		bss = ".section .bss\n"
		tabela := make(map[string]bool)
		for _, decl := range root.Globais {
			if !tabela[decl.Nome] {
				if decl.Tamanho > 0 {
					bss += "\t.balign 32\n"
					bss += fmt.Sprintf("\t.lcomm %s, %d\n", decl.Nome, decl.Tamanho*8)
				} else {
					bss += fmt.Sprintf("\t.lcomm %s, 8\n", decl.Nome)
				}
				tabela[decl.Nome] = true
			}
		}
	}

	assembly := ""
	if bss != "" {
		assembly += bss + "\n"
	}
	assembly += ".section .text\n.globl _start\n_start:\n"
	
	for _, decl := range root.Globais {
		if decl.Tamanho > 0 {
			// Arrays globais já são zerados pelo .bss
		} else {
			assembly += gerarAssembly(decl.Expressao, nil)
			assembly += "\tpopq %rax\n"
			assembly += fmt.Sprintf("\tmovq %%rax, %s(%%rip)\n", decl.Nome)
		}
	}

	for _, cmd := range root.CmdsMain {
		assembly += gerarCmd(cmd, nil)
	}

	assembly += gerarAssembly(root.Resultado, nil)
	assembly += "\tpopq %rax\n\tcall imprime_num\n\tcall sair\n\n"

	for _, f := range root.Funcoes {
		assembly += gerarFunDecl(f)
		assembly += "\n"
	}

	assembly += "\t.include \"assembly/runtime.s\"\n"

	assembly = otimizarAssembly(assembly)

	err := os.WriteFile(path, []byte(assembly), 0644)
	if err != nil {
		fmt.Println("Erro ao gerar arquivo assembly:", err)
		return
	}
	fmt.Println("Arquivo '" + path + "' gerado com sucesso!")
}

func otimizarAssembly(asm string) string {
	linhas := strings.Split(asm, "\n")
	resultado := make([]string, 0, len(linhas))

	for i := 0; i < len(linhas); i++ {
		a := strings.TrimSpace(linhas[i])

		// Olha o par (linhas[i], linhas[i+1])
		if i+1 < len(linhas) {
			b := strings.TrimSpace(linhas[i+1])

			// pushq X / popq %reg  →  movq X, %reg
			if strings.HasPrefix(a, "pushq ") && strings.HasPrefix(b, "popq ") {
				src := strings.TrimPrefix(a, "pushq ")
				dst := strings.TrimPrefix(b, "popq ")

				// pushq %r / popq %r  →  elimina ambos
				if src == dst {
					i++ // pula a próxima linha
					continue
				}

				// pushq $imm / popq %reg  →  movq $imm, %reg
				// pushq %reg1 / popq %reg2  →  movq %reg1, %reg2
				resultado = append(resultado, "\tmovq "+src+", "+dst)
				i++ // pula a próxima linha
				continue
			}
		}

		resultado = append(resultado, linhas[i])
	}

	return strings.Join(resultado, "\n")
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
