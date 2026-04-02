package main

import (
	"errors"
	"fmt"
)

func interpretar(n Exp, amb map[string]int) (int, error) {
	switch v := n.(type) {
	case Const:
		return v.Valor, nil
	case Var:
		return amb[v.Nome], nil
	case Programa:
		for _, decl := range v.Declaracoes {
			val, err := interpretar(decl.Expressao, amb)
			if err != nil {
				return 0, err
			}
			amb[decl.Nome] = val
		}
		return interpretar(v.Resultado, amb)
	case OpBin:
		esq, err := interpretar(v.OpEsq, amb)
		if err != nil {
			return 0, err
		}
		dir, err := interpretar(v.OpDir, amb)
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
		default:
			return 0, fmt.Errorf("operador desconhecido: %s", v.Operador)
		}
	}
	return 0, errors.New("tipo de nó desconhecido na árvore")
}
