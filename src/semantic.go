package main

import (
	"fmt"
)

func validaVariaveisExp(exp Exp, global map[string]bool, local map[string]bool, funcoes map[string]int) error {
	switch v := exp.(type) {
	case Const:
		return nil
	case Var:
		if local != nil && local[v.Nome] {
			return nil
		}
		if !global[v.Nome] {
			return fmt.Errorf("erro semântico: variável '%s' não declarada", v.Nome)
		}
		return nil
	case ChamadaFun:
		args_expected, ok := funcoes[v.Nome]
		if !ok {
			return fmt.Errorf("função '%s' não declarada", v.Nome)
		}
		if len(v.Args) != args_expected {
			return fmt.Errorf("função '%s' espera %d argumentos, mas recebeu %d", v.Nome, args_expected, len(v.Args))
		}
		for _, arg := range v.Args {
			if err := validaVariaveisExp(arg, global, local, funcoes); err != nil {
				return err
			}
		}
		return nil
	case OpBin:
		if err := validaVariaveisExp(v.OpEsq, global, local, funcoes); err != nil {
			return err
		}
		return validaVariaveisExp(v.OpDir, global, local, funcoes)
	case OpUnario:
		return validaVariaveisExp(v.Expressao, global, local, funcoes)
	case AcessoVetor:
		if local != nil && local[v.Nome] {
			return validaVariaveisExp(v.Indice, global, local, funcoes)
		}
		if !global[v.Nome] {
			return fmt.Errorf("erro semântico: vetor '%s' não declarado", v.Nome)
		}
		return validaVariaveisExp(v.Indice, global, local, funcoes)
	case Programa:
		return AnalisarSemantica(v)
	}
	return nil
}

func validaVariaveisCmd(cmd Cmd, global map[string]bool, local map[string]bool, funcoes map[string]int) error {
	switch v := cmd.(type) {
	case IfCmd:
		if err := validaVariaveisExp(v.Condicao, global, local, funcoes); err != nil {
			return err
		}
		for _, c := range v.CorpoIf {
			if err := validaVariaveisCmd(c, global, local, funcoes); err != nil {
				return err
			}
		}
		for _, c := range v.CorpoElse {
			if err := validaVariaveisCmd(c, global, local, funcoes); err != nil {
				return err
			}
		}
	case WhileCmd:
		if err := validaVariaveisExp(v.Condicao, global, local, funcoes); err != nil {
			return err
		}
		for _, c := range v.Corpo {
			if err := validaVariaveisCmd(c, global, local, funcoes); err != nil {
				return err
			}
		}
	case ForCmd:
		if err := validaVariaveisCmd(v.Init, global, local, funcoes); err != nil {
			return err
		}
		if err := validaVariaveisExp(v.Condicao, global, local, funcoes); err != nil {
			return err
		}
		for _, c := range v.Corpo {
			if err := validaVariaveisCmd(c, global, local, funcoes); err != nil {
				return err
			}
		}
		if err := validaVariaveisCmd(v.Passo, global, local, funcoes); err != nil {
			return err
		}
	case AtribCmd:
		if (local == nil || !local[v.Nome]) && !global[v.Nome] {
			return fmt.Errorf("erro semântico: atribuição a variável não declarada '%s'", v.Nome)
		}
		if err := validaVariaveisExp(v.Expressao, global, local, funcoes); err != nil {
			return err
		}
	case AtribVetorCmd:
		if (local == nil || !local[v.Nome]) && !global[v.Nome] {
			return fmt.Errorf("erro semântico: vetor '%s' não declarado para atribuição", v.Nome)
		}
		if err := validaVariaveisExp(v.Indice, global, local, funcoes); err != nil {
			return err
		}
		if err := validaVariaveisExp(v.Expressao, global, local, funcoes); err != nil {
			return err
		}
	case CmdBuiltinVec:
		expectedArgs := 4
		if v.Operacao == "vset" {
			expectedArgs = 3
		}
		if len(v.Args) != expectedArgs {
			return fmt.Errorf("erro semântico: %s espera %d argumentos, recebeu %d", v.Operacao, expectedArgs, len(v.Args))
		}
		for _, arg := range v.Args {
			if err := validaVariaveisExp(arg, global, local, funcoes); err != nil {
				return err
			}
		}
	}
	return nil
}

func AnalisarSemantica(prog Programa) error {
	globais := make(map[string]bool)
	funcoes := make(map[string]int)

	for _, decl := range prog.Globais {
		err := validaVariaveisExp(decl.Expressao, globais, nil, funcoes)
		if err != nil {
			return err
		}
		globais[decl.Nome] = true
	}

	for _, f := range prog.Funcoes {
		funcoes[f.Nome] = len(f.Params)
	}

	for _, f := range prog.Funcoes {
		local := make(map[string]bool)
		for _, p := range f.Params {
			local[p] = true
		}
		for _, decl := range f.Vars {
			err := validaVariaveisExp(decl.Expressao, globais, local, funcoes)
			if err != nil {
				return err
			}
			local[decl.Nome] = true
		}
		for _, cmd := range f.Comandos {
			err := validaVariaveisCmd(cmd, globais, local, funcoes)
			if err != nil {
				return err
			}
		}
		err := validaVariaveisExp(f.Resultado, globais, local, funcoes)
		if err != nil {
			return err
		}
	}

	for _, cmd := range prog.CmdsMain {
		err := validaVariaveisCmd(cmd, globais, nil, funcoes)
		if err != nil {
			return err
		}
	}

	err := validaVariaveisExp(prog.Resultado, globais, nil, funcoes)
	if err != nil {
		return err
	}

	return nil
}
