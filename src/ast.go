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

type OpUnario struct {
	Operador string
	Expressao Exp
}

func (u OpUnario) isExp() {}

type Var struct {
	Nome string
}

func (v Var) isExp() {}

type Decl struct {
	Nome      string
	Tamanho   int
	Expressao Exp
}

type Cmd interface {
	isCmd()
}

type IfCmd struct {
	Condicao  Exp
	CorpoIf   []Cmd
	CorpoElse []Cmd
}

func (i IfCmd) isCmd() {}

type WhileCmd struct {
	Condicao Exp
	Corpo    []Cmd
}

func (w WhileCmd) isCmd() {}

type ForCmd struct {
	Init     Cmd
	Condicao Exp
	Passo    Cmd
	Corpo    []Cmd
}

func (f ForCmd) isCmd() {}

type AtribCmd struct {
	Nome      string
	Expressao Exp
}

func (a AtribCmd) isCmd() {}

type ChamadaFun struct {
	Nome string
	Args []Exp
}

func (c ChamadaFun) isExp() {}

type AcessoVetor struct {
	Nome   string
	Indice Exp
}

func (a AcessoVetor) isExp() {}

type AtribVetorCmd struct {
	Nome      string
	Indice    Exp
	Expressao Exp
}

func (a AtribVetorCmd) isCmd() {}

type CmdBuiltinVec struct {
	Operacao string // "vadd", "vsub", "vset"
	Args     []Exp
}

func (c CmdBuiltinVec) isCmd() {}

type FunDecl struct {
	Nome      string
	Params    []string
	Vars      []Decl
	Comandos  []Cmd
	Resultado Exp
}

type Programa struct {
	Globais   []Decl
	Funcoes   []FunDecl
	CmdsMain  []Cmd
	Resultado Exp
}

func (p Programa) isExp() {}