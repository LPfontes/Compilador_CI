package main

import (
	"errors"
	"fmt"
)

func interpretarCmd(cmd Cmd, globais map[string]int, locais map[string]int, funcoes map[string]FunDecl) error {
	switch v := cmd.(type) {
	case IfCmd:
		cond, err := interpretar(v.Condicao, globais, locais, funcoes)
		if err != nil {
			return err
		}
		if cond != 0 {
			for _, c := range v.CorpoIf {
				if err := interpretarCmd(c, globais, locais, funcoes); err != nil {
					return err
				}
			}
		} else {
			for _, c := range v.CorpoElse {
				if err := interpretarCmd(c, globais, locais, funcoes); err != nil {
					return err
				}
			}
		}
	case WhileCmd:
		for {
			cond, err := interpretar(v.Condicao, globais, locais, funcoes)
			if err != nil {
				return err
			}
			if cond == 0 {
				break
			}
			for _, c := range v.Corpo {
				if err := interpretarCmd(c, globais, locais, funcoes); err != nil {
					return err
				}
			}
		}
	case ForCmd:
		if err := interpretarCmd(v.Init, globais, locais, funcoes); err != nil {
			return err
		}
		for {
			cond, err := interpretar(v.Condicao, globais, locais, funcoes)
			if err != nil {
				return err
			}
			if cond == 0 {
				break
			}
			for _, c := range v.Corpo {
				if err := interpretarCmd(c, globais, locais, funcoes); err != nil {
					return err
				}
			}
			if err := interpretarCmd(v.Passo, globais, locais, funcoes); err != nil {
				return err
			}
		}
	case AtribCmd:
		val, err := interpretar(v.Expressao, globais, locais, funcoes)
		if err != nil {
			return err
		}
		if locais != nil {
			if _, ok := locais[v.Nome]; ok {
				locais[v.Nome] = val
				return nil
			}
		}
		globais[v.Nome] = val
	case AtribVetorCmd:
		idx, err := interpretar(v.Indice, globais, locais, funcoes)
		if err != nil {
			return err
		}
		val, err := interpretar(v.Expressao, globais, locais, funcoes)
		if err != nil {
			return err
		}
		chave := fmt.Sprintf("%s__%d", v.Nome, idx)
		if locais != nil {
			if _, ok := locais[v.Nome]; ok {
				locais[chave] = val
				return nil
			}
		}
		globais[chave] = val
	case CmdBuiltinVec:
		// Resolver o nome do destino
		destVar, ok := v.Args[0].(Var)
		if !ok {
			return fmt.Errorf("primeiro argumento de %s deve ser nome de vetor", v.Operacao)
		}
		destNome := destVar.Nome

		if v.Operacao == "vset" {
			valor, err := interpretar(v.Args[1], globais, locais, funcoes)
			if err != nil {
				return err
			}
			n, err := interpretar(v.Args[2], globais, locais, funcoes)
			if err != nil {
				return err
			}
			for i := 0; i < n; i++ {
				chave := fmt.Sprintf("%s__%d", destNome, i)
				if locais != nil {
					if _, ok := locais[destNome]; ok {
						locais[chave] = valor
						continue
					}
				}
				globais[chave] = valor
			}
		} else {
			// vadd / vsub
			src1Var, _ := v.Args[1].(Var)
			src2Var, _ := v.Args[2].(Var)
			n, err := interpretar(v.Args[3], globais, locais, funcoes)
			if err != nil {
				return err
			}
			for i := 0; i < n; i++ {
				chave1 := fmt.Sprintf("%s__%d", src1Var.Nome, i)
				chave2 := fmt.Sprintf("%s__%d", src2Var.Nome, i)
				chaveDest := fmt.Sprintf("%s__%d", destNome, i)
				var v1, v2 int
				if locais != nil {
					if val, ok := locais[chave1]; ok {
						v1 = val
					} else {
						v1 = globais[chave1]
					}
					if val, ok := locais[chave2]; ok {
						v2 = val
					} else {
						v2 = globais[chave2]
					}
				} else {
					v1 = globais[chave1]
					v2 = globais[chave2]
				}
				var resultado int
				if v.Operacao == "vadd" {
					resultado = v1 + v2
				} else {
					resultado = v1 - v2
				}
				if locais != nil {
					if _, ok := locais[destNome]; ok {
						locais[chaveDest] = resultado
						continue
					}
				}
				globais[chaveDest] = resultado
			}
		}
	}
	return nil
}

func interpretar(exp Exp, globais map[string]int, locais map[string]int, funcoes map[string]FunDecl) (int, error) {
	switch v := exp.(type) {
	case Const:
		return v.Valor, nil
	case Var:
		if locais != nil {
			if val, ok := locais[v.Nome]; ok {
				return val, nil
			}
		}
		return globais[v.Nome], nil
	case AcessoVetor:
		idx, err := interpretar(v.Indice, globais, locais, funcoes)
		if err != nil {
			return 0, err
		}
		chave := fmt.Sprintf("%s__%d", v.Nome, idx)
		if locais != nil {
			if val, ok := locais[chave]; ok {
				return val, nil
			}
		}
		return globais[chave], nil
	case ChamadaFun:
		f := funcoes[v.Nome]
		novoLocal := make(map[string]int)
		for i, p := range f.Params {
			val, err := interpretar(v.Args[i], globais, locais, funcoes)
			if err != nil {
				return 0, err
			}
			novoLocal[p] = val
		}
		for _, decl := range f.Vars {
			if decl.Tamanho > 0 {
				novoLocal[decl.Nome] = 0
				for j := 0; j < decl.Tamanho; j++ {
					novoLocal[fmt.Sprintf("%s__%d", decl.Nome, j)] = 0
				}
			} else {
				val, err := interpretar(decl.Expressao, globais, novoLocal, funcoes)
				if err != nil {
					return 0, err
				}
				novoLocal[decl.Nome] = val
			}
		}
		for _, cmd := range f.Comandos {
			if err := interpretarCmd(cmd, globais, novoLocal, funcoes); err != nil {
				return 0, err
			}
		}
		return interpretar(f.Resultado, globais, novoLocal, funcoes)
	case Programa:
		fmap := make(map[string]FunDecl)
		for _, f := range v.Funcoes {
			fmap[f.Nome] = f
		}
		for _, decl := range v.Globais {
			if decl.Tamanho > 0 {
				globais[decl.Nome] = 0
				for j := 0; j < decl.Tamanho; j++ {
					globais[fmt.Sprintf("%s__%d", decl.Nome, j)] = 0
				}
			} else {
				val, err := interpretar(decl.Expressao, globais, locais, fmap)
				if err != nil {
					return 0, err
				}
				globais[decl.Nome] = val
			}
		}
		for _, cmd := range v.CmdsMain {
			if err := interpretarCmd(cmd, globais, locais, fmap); err != nil {
				return 0, err
			}
		}
		return interpretar(v.Resultado, globais, locais, fmap)
	case OpUnario:
		val, err := interpretar(v.Expressao, globais, locais, funcoes)
		if err != nil {
			return 0, err
		}
		if v.Operador == "not" {
			if val == 0 {
				return 1, nil
			}
			return 0, nil
		}
		return 0, fmt.Errorf("operador unário desconhecido: %s", v.Operador)
	case OpBin:
		esq, err := interpretar(v.OpEsq, globais, locais, funcoes)
		if err != nil {
			return 0, err
		}
		dir, err := interpretar(v.OpDir, globais, locais, funcoes)
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
		case "and":
			if esq != 0 && dir != 0 {
				return 1, nil
			}
			return 0, nil
		case "or":
			if esq != 0 || dir != 0 {
				return 1, nil
			}
			return 0, nil
		case "<":
			if esq < dir {
				return 1, nil
			}
			return 0, nil
		case ">":
			if esq > dir {
				return 1, nil
			}
			return 0, nil
		case "==":
			if esq == dir {
				return 1, nil
			}
			return 0, nil
		case "!=":
			if esq != dir {
				return 1, nil
			}
			return 0, nil
		case "<=":
			if esq <= dir {
				return 1, nil
			}
			return 0, nil
		case ">=":
			if esq >= dir {
				return 1, nil
			}
			return 0, nil
		default:
			return 0, fmt.Errorf("operador desconhecido: %s", v.Operador)
		}
	}
	return 0, errors.New("tipo de expressao desconhecido")
}

