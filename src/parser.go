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

func (p *Parser) analisaExp() (Exp, error) {
	tok := p.atual()

	if tok.Tipo == "Numero" {
		p.consumir()
		valor, _ := strconv.Atoi(tok.Literal)
		return Const{Valor: valor}, nil

	} else if tok.Tipo == "ParenEsq" {
		p.consumir() // consome '('

		opEsq, err := p.analisaExp() // descobre o operando esquerdo Descida Recursiva
		if err != nil {
			return nil, err
		}

		// analisaOperador()
		operadorTok := p.consumir()
		if _, ok := tipos[operadorTok.Literal]; !ok || operadorTok.Tipo == "ParenEsq" || operadorTok.Tipo == "ParenDir" { // Verifica se é um operador válido
			return nil, fmt.Errorf("esperado operador na posição %d", operadorTok.posicao)
		}

		opDir, err := p.analisaExp() // descobre o operando direito Descida Recursiva
		if err != nil {
			return nil, err
		}

		// verificaProxToken(FECHA_PARENTESE)
		if p.atual().Tipo != "ParenDir" {
			return nil, fmt.Errorf("esperado ')' na posição %d", p.atual().posicao)
		}
		p.consumir() // consome ')'

		return OpBin{
			Operador: operadorTok.Literal,
			OpEsq:    opEsq,
			OpDir:    opDir,
		}, nil
	}

	return nil, fmt.Errorf("erro sintático: token inesperado %s na posição %d", tok.Literal, tok.posicao)
}
