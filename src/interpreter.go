package main

import (
	"errors"
	"fmt"
)

func interpretar(n Exp) (int, error) {
	switch v := n.(type) {
	case Const:
		return v.Valor, nil
	case OpBin:
		esq, err := interpretar(v.OpEsq)
		if err != nil {
			return 0, err
		}
		dir, err := interpretar(v.OpDir)
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
