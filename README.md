# Compilador Simples (Construção de Compiladores I)

Este projeto é uma atividade desenvolvida para a disciplina de **Construção de Compiladores I**.
O projeto consiste em um compilador completo para uma linguagem extremamente simples. O propósito é demonstrar as etapas fundamentais de um compilador no contexto mais simples possível.

O programa já realiza a montagem e linkagem do arquivo. 

A linguagem reconhecida pelo compilador é definida pela seguinte gramática (EBNF):

```
<programa> ::= <literal-inteiro>
<literal-inteiro> ::= <digito>+
<digito> ::= 0 | 1 | 2 | 3 | 4 | 5 | 6 | 7 | 8 | 9
```

(Basicamente, a linguagem aceita apenas um único número inteiro).

## Pré-requisitos

Para executar o compilador, é necessário ter o **compilador de Go (Golang)** instalado em seu sistema.

## Execução

Para compilar e executar o programa, utilize o comando `go run` no seu terminal:

```bash
go run CI <arquivo de entrada>
```

