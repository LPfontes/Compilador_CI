package main

type Token struct {
	Tipo    string
	Literal string
	posicao int
}

type Exp interface {
	isExp()
}

type Const struct {
	Valor int
}

func (c Const) isExp() {}

type OpBin struct {
	Operador string
	OpEsq    Exp
	OpDir    Exp
}

func (o OpBin) isExp() {}