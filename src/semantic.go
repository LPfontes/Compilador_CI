package main

import (
	"fmt"
)

func validaVariaveisExp(exp Exp, tabela map[string]bool) error {
	switch v := exp.(type) {
	case Const:
		return nil
	case Var:
		if !tabela[v.Nome] {
			return fmt.Errorf("erro semântico: variável '%s' não declarada", v.Nome)
		}
		return nil
	case OpBin:
		if err := validaVariaveisExp(v.OpEsq, tabela); err != nil {
			return err
		}
		return validaVariaveisExp(v.OpDir, tabela)
	case Programa:
		return AnalisarSemantica(v)
	}
	return nil
}

func AnalisarSemantica(prog Programa) error {
	tabela := make(map[string]bool)

	for _, decl := range prog.Declaracoes {
		err := validaVariaveisExp(decl.Expressao, tabela)
		if err != nil {
			return err
		}
		tabela[decl.Nome] = true
	}

	err := validaVariaveisExp(prog.Resultado, tabela)
	if err != nil {
		return err
	}

	return nil
}
