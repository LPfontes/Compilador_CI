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
	esq, err := p.analisaExpAnd()
	if err != nil {
		return nil, err
	}

	for p.atual().Tipo == "Or" {
		op := p.consumir()
		dir, err := p.analisaExpAnd()
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

// analisaExpAnd lida com And
func (p *Parser) analisaExpAnd() (Exp, error) {
	esq, err := p.analisaExpRelacional()
	if err != nil {
		return nil, err
	}

	for p.atual().Tipo == "And" {
		op := p.consumir()
		dir, err := p.analisaExpRelacional()
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

// analisaExpRelacional lida com <, >, ==, <=, >=
func (p *Parser) analisaExpRelacional() (Exp, error) {
	esq, err := p.analisaExpAritmetica()
	if err != nil {
		return nil, err
	}

	for p.atual().Tipo == "MenorQue" || p.atual().Tipo == "MaiorQue" || p.atual().Tipo == "IgualIgual" || p.atual().Tipo == "MenorIgual" || p.atual().Tipo == "MaiorIgual" || p.atual().Tipo == "Diferente" {
		op := p.consumir()
		dir, err := p.analisaExpAritmetica()
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

// analisaExpAritmetica (Expression): Lida com Soma e Subtração exp_a
func (p *Parser) analisaExpAritmetica() (Exp, error) {
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

	if tok.Tipo == "Not" {
		op := p.consumir() // consome not
		exp, err := p.analisaFator()
		if err != nil {
			return nil, err
		}
		return OpUnario{Operador: op.Literal, Expressao: exp}, nil
		
	} else if tok.Tipo == "Numero" {
		p.consumir()
		valor, _ := strconv.Atoi(tok.Literal)
		return Const{Valor: valor}, nil

	} else if tok.Tipo == "Identificador" {
		nome := p.consumir().Literal
		if p.pos < len(p.tokens) && p.atual().Tipo == "ParenEsq" {
			p.consumir() // '('
			var args []Exp
			for p.atual().Tipo != "ParenDir" && p.atual().Tipo != "EOF" {
				exp, err := p.analisaExp()
				if err != nil {
					return nil, err
				}
				args = append(args, exp)
				if p.atual().Tipo == "Virgula" {
					p.consumir()
				}
			}
			if p.atual().Tipo != "ParenDir" {
				return nil, fmt.Errorf("esperado ')' encerrando chamada de func")
			}
			p.consumir()
			return ChamadaFun{Nome: nome, Args: args}, nil
		} else if p.pos < len(p.tokens) && p.atual().Tipo == "ColcheteEsq" {
			p.consumir() // '['
			indiceExp, err := p.analisaExp()
			if err != nil {
				return nil, err
			}
			if p.atual().Tipo != "ColcheteDir" {
				return nil, fmt.Errorf("esperado ']' encerrando acesso a vetor")
			}
			p.consumir() // ']'
			return AcessoVetor{Nome: nome, Indice: indiceExp}, nil
		}
		return Var{Nome: nome}, nil

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
	if p.atual().Tipo != "Var" {
		return Decl{}, fmt.Errorf("esperado 'var'")
	}
	p.consumir() // 'var'
	nomeTok := p.consumir() // Identificador
	if nomeTok.Tipo != "Identificador" {
		return Decl{}, fmt.Errorf("esperado identificador apos var")
	}

	if p.atual().Tipo == "ColcheteEsq" {
		p.consumir() // '['
		tokTamanho := p.atual()
		if tokTamanho.Tipo != "Numero" {
			return Decl{}, fmt.Errorf("esperado tamanho numerico constante para vetor")
		}
		tamanho, _ := strconv.Atoi(p.consumir().Literal)
		if p.atual().Tipo != "ColcheteDir" {
			return Decl{}, fmt.Errorf("esperado ']' apos tamanho")
		}
		p.consumir() // ']'
		if p.atual().Tipo != "PontoVirgula" {
			return Decl{}, fmt.Errorf("esperado ';' apos declaracao de vetor")
		}
		p.consumir() // ';'
		return Decl{Nome: nomeTok.Literal, Tamanho: tamanho, Expressao: Const{Valor: 0}}, nil
	}

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

func (p *Parser) analisaFunDecl() (FunDecl, error) {
	p.consumir() // 'fun'
	nomeTok := p.consumir()
	if nomeTok.Tipo != "Identificador" {
		return FunDecl{}, fmt.Errorf("esperado identificador para a funcao")
	}
	if p.atual().Tipo != "ParenEsq" {
		return FunDecl{}, fmt.Errorf("esperado '(' na funcao")
	}
	p.consumir()
	var params []string
	for p.atual().Tipo != "ParenDir" && p.atual().Tipo != "EOF" {
		tok := p.consumir()
		if tok.Tipo == "Identificador" {
			params = append(params, tok.Literal)
		} else {
			return FunDecl{}, fmt.Errorf("esperado parametro identificador")
		}
		if p.atual().Tipo == "Virgula" {
			p.consumir()
		}
	}
	if p.atual().Tipo != "ParenDir" {
		return FunDecl{}, fmt.Errorf("esperado ')' apos parametros")
	}
	p.consumir()
	if p.atual().Tipo != "ChaveEsq" {
		return FunDecl{}, fmt.Errorf("esperado '{' pos fun")
	}
	p.consumir()

	var vars []Decl
	for p.atual().Tipo == "Var" {
		dec, err := p.analisaDecl()
		if err != nil {
			return FunDecl{}, err
		}
		vars = append(vars, dec)
	}

	var cmds []Cmd
	for p.atual().Tipo != "Return" && p.atual().Tipo != "EOF" {
		cmd, err := p.analisaComando()
		if err != nil {
			return FunDecl{}, err
		}
		cmds = append(cmds, cmd)
	}
	if p.atual().Tipo != "Return" {
		return FunDecl{}, fmt.Errorf("esperado 'return' funcao")
	}
	p.consumir()
	res, err := p.analisaExp()
	if err != nil {
		return FunDecl{}, err
	}
	if p.atual().Tipo != "PontoVirgula" {
		return FunDecl{}, fmt.Errorf("esperado ';' apos return da funcao")
	}
	p.consumir()
	if p.atual().Tipo != "ChaveDir" {
		return FunDecl{}, fmt.Errorf("esperado '}'")
	}
	p.consumir()
	return FunDecl{Nome: nomeTok.Literal, Params: params, Vars: vars, Comandos: cmds, Resultado: res}, nil
}


func (p *Parser) analisaComando() (Cmd, error) {
	tok := p.atual()
	if tok.Tipo == "If" {
		p.consumir()
		exp, err := p.analisaExp()
		if err != nil {
			return nil, err
		}
		if p.atual().Tipo != "ChaveEsq" {
			return nil, fmt.Errorf("esperado '{' após if")
		}
		p.consumir()
		var corpoIf []Cmd
		for p.atual().Tipo != "ChaveDir" && p.atual().Tipo != "EOF" {
			cmd, err := p.analisaComando()
			if err != nil {
				return nil, err
			}
			corpoIf = append(corpoIf, cmd)
		}
		if p.atual().Tipo != "ChaveDir" {
			return nil, fmt.Errorf("esperado '}'")
		}
		p.consumir()

		var corpoElse []Cmd
		if p.atual().Tipo == "Else" {
			p.consumir()
			if p.atual().Tipo != "ChaveEsq" {
				return nil, fmt.Errorf("esperado '{' após else")
			}
			p.consumir()
			for p.atual().Tipo != "ChaveDir" && p.atual().Tipo != "EOF" {
				cmd, err := p.analisaComando()
				if err != nil {
					return nil, err
				}
				corpoElse = append(corpoElse, cmd)
			}
			if p.atual().Tipo != "ChaveDir" {
				return nil, fmt.Errorf("esperado '}'")
			}
			p.consumir()
		} else {
			return nil, fmt.Errorf("comando if na linguagem Cmd exige branch else obrigatorio")
		}

		return IfCmd{Condicao: exp, CorpoIf: corpoIf, CorpoElse: corpoElse}, nil

	} else if tok.Tipo == "While" {
		p.consumir()
		exp, err := p.analisaExp()
		if err != nil {
			return nil, err
		}
		if p.atual().Tipo != "ChaveEsq" {
			return nil, fmt.Errorf("esperado '{' após while")
		}
		p.consumir()
		var corpo []Cmd
		for p.atual().Tipo != "ChaveDir" && p.atual().Tipo != "EOF" {
			cmd, err := p.analisaComando()
			if err != nil {
				return nil, err
			}
			corpo = append(corpo, cmd)
		}
		if p.atual().Tipo != "ChaveDir" {
			return nil, fmt.Errorf("esperado '}' no while")
		}
		p.consumir()
		return WhileCmd{Condicao: exp, Corpo: corpo}, nil

	} else if tok.Tipo == "For" {
		p.consumir() // 'for'
		// Init: var = expr;
		initCmd, err := p.analisaComando()
		if err != nil {
			return nil, err
		}
		// Condicao: expr;
		cond, err := p.analisaExp()
		if err != nil {
			return nil, err
		}
		if p.atual().Tipo != "PontoVirgula" {
			return nil, fmt.Errorf("esperado ';' após condição do for")
		}
		p.consumir()
		// Passo: var = expr (sem ';' antes do '{')
		passoNome := p.consumir()
		if passoNome.Tipo != "Identificador" {
			return nil, fmt.Errorf("esperado identificador no passo do for")
		}
		if p.atual().Tipo != "Atribuicao" {
			return nil, fmt.Errorf("esperado '=' no passo do for")
		}
		p.consumir()
		passoExp, err := p.analisaExp()
		if err != nil {
			return nil, err
		}
		passo := AtribCmd{Nome: passoNome.Literal, Expressao: passoExp}

		if p.atual().Tipo != "ChaveEsq" {
			return nil, fmt.Errorf("esperado '{' após for")
		}
		p.consumir()
		var corpo []Cmd
		for p.atual().Tipo != "ChaveDir" && p.atual().Tipo != "EOF" {
			cmd, err := p.analisaComando()
			if err != nil {
				return nil, err
			}
			corpo = append(corpo, cmd)
		}
		if p.atual().Tipo != "ChaveDir" {
			return nil, fmt.Errorf("esperado '}' no for")
		}
		p.consumir()
		return ForCmd{Init: initCmd, Condicao: cond, Passo: passo, Corpo: corpo}, nil

	} else if tok.Tipo == "Identificador" {
		nomeTok := p.consumir()

		// Builtins vetoriais AVX2
		if (nomeTok.Literal == "vadd" || nomeTok.Literal == "vsub" || nomeTok.Literal == "vset") && p.atual().Tipo == "ParenEsq" {
			p.consumir() // '('
			var args []Exp
			for p.atual().Tipo != "ParenDir" && p.atual().Tipo != "EOF" {
				exp, err := p.analisaExp()
				if err != nil {
					return nil, err
				}
				args = append(args, exp)
				if p.atual().Tipo == "Virgula" {
					p.consumir()
				}
			}
			if p.atual().Tipo != "ParenDir" {
				return nil, fmt.Errorf("esperado ')' apos argumentos de %s", nomeTok.Literal)
			}
			p.consumir() // ')'
			if p.atual().Tipo != "PontoVirgula" {
				return nil, fmt.Errorf("esperado ';' apos %s", nomeTok.Literal)
			}
			p.consumir() // ';'
			return CmdBuiltinVec{Operacao: nomeTok.Literal, Args: args}, nil
		}
		
		if p.atual().Tipo == "ColcheteEsq" {
			p.consumir() // '['
			indExp, err := p.analisaExp()
			if err != nil {
				return nil, err
			}
			if p.atual().Tipo != "ColcheteDir" {
				return nil, fmt.Errorf("esperado ']' em atribuicao vetor")
			}
			p.consumir() // ']'
			if p.atual().Tipo != "Atribuicao" {
				return nil, fmt.Errorf("esperado '=' apos indice limitador")
			}
			p.consumir() // '='
			valExp, err := p.analisaExp()
			if err != nil {
				return nil, err
			}
			if p.atual().Tipo != "PontoVirgula" {
				return nil, fmt.Errorf("esperado ';' apos vetor cmd")
			}
			p.consumir() // ';'
			return AtribVetorCmd{Nome: nomeTok.Literal, Indice: indExp, Expressao: valExp}, nil
		}

		if p.atual().Tipo != "Atribuicao" {
			return nil, fmt.Errorf("esperava '=' após variável '%s' num comando", nomeTok.Literal)
		}
		p.consumir()
		exp, err := p.analisaExp()
		if err != nil {
			return nil, err
		}
		if p.atual().Tipo != "PontoVirgula" {
			return nil, fmt.Errorf("esperado ';' apos atribuicao")
		}
		p.consumir()
		return AtribCmd{Nome: nomeTok.Literal, Expressao: exp}, nil
	}

	return nil, fmt.Errorf("comando invalido ou token inesperado: %s", tok.Literal)
}

func (p *Parser) analisaPrograma() (Programa, error) {
	var globais []Decl
	var funcoes []FunDecl

	for p.atual().Tipo == "Var" || p.atual().Tipo == "Fun" {
		if p.atual().Tipo == "Var" {
			dec, err := p.analisaDecl()
			if err != nil {
				return Programa{}, err
			}
			globais = append(globais, dec)
		} else {
			f, err := p.analisaFunDecl()
			if err != nil {
				return Programa{}, err
			}
			funcoes = append(funcoes, f)
		}
	}

	if p.atual().Tipo != "Main" {
		return Programa{}, fmt.Errorf("esperado 'main'")
	}
	p.consumir()

	if p.atual().Tipo != "ChaveEsq" {
		return Programa{}, fmt.Errorf("esperado '{' main")
	}
	p.consumir()

	var cmds []Cmd
	for p.atual().Tipo != "Return" && p.atual().Tipo != "EOF" {
		cmd, err := p.analisaComando()
		if err != nil {
			return Programa{}, err
		}
		cmds = append(cmds, cmd)
	}

	if p.atual().Tipo != "Return" {
		return Programa{}, fmt.Errorf("esperado 'return' main")
	}
	p.consumir()
	res, err := p.analisaExp()
	if err != nil {
		return Programa{}, err
	}
	if p.atual().Tipo != "PontoVirgula" {
		return Programa{}, fmt.Errorf("esperado ';' apos main")
	}
	p.consumir()
	if p.atual().Tipo != "ChaveDir" {
		return Programa{}, fmt.Errorf("esperado '}' main")
	}
	p.consumir()
	return Programa{Globais: globais, Funcoes: funcoes, CmdsMain: cmds, Resultado: res}, nil
}
