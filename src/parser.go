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

	} else if tok.Tipo == "Identificador" {
		p.consumir()
		return Var{Nome: tok.Literal}, nil

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

func (p *Parser) analisaDecl() (Decl, error) {
	nomeTok := p.consumir() // Identificador
	if p.atual().Tipo != "Atribuicao" {
		return Decl{}, fmt.Errorf("esperava '=' após variável '%s', mas encontrou '%s'", nomeTok.Literal, p.atual().Literal)
	}
	p.consumir() // consome '='

	exp, err := p.analisaExp()
	if err != nil {
		return Decl{}, err
	}

	if p.atual().Tipo != "PontoVirgula" {
		return Decl{}, fmt.Errorf("esperava ';' no final da declaração de '%s', mas encontrou '%s'", nomeTok.Literal, p.atual().Literal)
	}
	p.consumir() // consome ';'

	return Decl{Nome: nomeTok.Literal, Expressao: exp}, nil
}

func (p *Parser) analisaPrograma() (Programa, error) {
	var decs []Decl

	for p.atual().Tipo == "Identificador" {
		dec, err := p.analisaDecl()
		if err != nil {
			return Programa{}, err
		}
		decs = append(decs, dec)
	}

	if p.atual().Tipo != "Atribuicao" {
		return Programa{}, fmt.Errorf("esperava '=' para a expressão de resultado, mas encontrou '%s'", p.atual().Literal)
	}
	p.consumir() // '='

	result, err := p.analisaExp()
	if err != nil {
		return Programa{}, err
	}

	return Programa{Declaracoes: decs, Resultado: result}, nil
}
