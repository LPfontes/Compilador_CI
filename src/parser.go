package main

import (
	"fmt"
	"strconv"
)

type Parser struct {
	tokens []Token
	pos    int
}

func (p *Parser) atual() Token {
	return p.tokens[p.pos]
}

func (p *Parser) consumir() Token {
	t := p.atual()
	if t.Tipo != "EOF" {
		p.pos++
	}
	return t
}

// analisaExp (Expression): Lida com Soma e Subtração exp_a
func (p *Parser) analisaExp() (Exp, error) {
	esq, err := p.analisaTermo()
	if err != nil {
		return nil, err
	}

	for p.atual().Tipo == "Soma" || p.atual().Tipo == "Sub" {
		op := p.consumir()
		dir, err := p.analisaTermo()
		if err != nil {
			return nil, err
		}
		esq = OpBin{
			Operador: op.Literal,
			OpEsq:    esq,
			OpDir:    dir,
		}
	}
	return esq, nil
}

// analisaTermo (Term): Lida com Multiplicação e Divisão <exp_m>
func (p *Parser) analisaTermo() (Exp, error) {
	esq, err := p.analisaFator()
	if err != nil {
		return nil, err
	}

	for p.atual().Tipo == "Mult" || p.atual().Tipo == "Div" {
		op := p.consumir()
		dir, err := p.analisaFator()
		if err != nil {
			return nil, err
		}
		esq = OpBin{
			Operador: op.Literal,
			OpEsq:    esq,
			OpDir:    dir,
		}
	}
	return esq, nil
}

// analisaFator (Factor): Lida com números e parênteses (a unidade básica) <prim>
func (p *Parser) analisaFator() (Exp, error) {
	tok := p.atual()

	if tok.Tipo == "Numero" {
		p.consumir()
		valor, _ := strconv.Atoi(tok.Literal)
		return Const{Valor: valor}, nil

	} else if tok.Tipo == "ParenEsq" {
		p.consumir() // consome '('

		exp, err := p.analisaExp()
		if err != nil {
			return nil, err
		}

		if p.atual().Tipo != "ParenDir" {
			return nil, fmt.Errorf("esperado ')' na posição %d", p.atual().posicao)
		}
		p.consumir() // consome ')'

		return exp, nil
	}

	return nil, fmt.Errorf("erro sintático: token inesperado %s na posição %d", tok.Literal, tok.posicao)
}
