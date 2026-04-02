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

type Var struct {
	Nome string
}

func (v Var) isExp() {}

type Decl struct {
	Nome      string
	Expressao Exp
}

type Programa struct {
	Declaracoes []Decl
	Resultado   Exp
}

func (p Programa) isExp() {}