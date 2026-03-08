# Compilador e Interpretador de Expressões Aritméticas (CI)

Este projeto é uma implementação prática para a disciplina de **Construção de Compiladores I**. Ele consiste em um compilador e interpretador capaz de processar expressões aritméticas, gerar a Árvore Sintática Abstrata (AST), visualizar essa árvore graficamente e compilar o código para Assembly x86-64.

## 🚀 Funcionalidades

O projeto executa o pipeline completo de um compilador:

1.  **Análise Léxica (Lexer):** Tokeniza a entrada, identificando números, parênteses e operadores (`+`, `-`, `*`, `/`).
2.  **Análise Sintática (Parser):** Constrói uma Árvore Sintática Abstrata (AST) respeitando a precedência definida por parênteses.
3.  **Interpretador:** Percorre a árvore recursivamente para calcular o resultado imediatamente.
4.  **Visualização:** Gera uma imagem `.png` da árvore sintática utilizando a biblioteca `graph` e o software Graphviz.
5.  **Geração de Código:** Traduz a AST para Assembly x86-64 (sintaxe AT&T).
6.  **Montagem e Ligação:** Automatiza o uso do `as` (assembler) e `ld` (linker) para gerar um binário executável.

## Pré-requisitos

Para executar este projeto, você precisará das seguintes ferramentas instaladas no seu ambiente Linux:

- **Go** (versão 1.23 ou superior)
- **Graphviz** (para o comando `dot` gerar a imagem da árvore)
- **GCC/Binutils** (para os comandos `as` e `ld`)

### Instalação das dependências (Ubuntu/Debian)

```bash
sudo apt update
sudo apt install golang graphviz build-essential
```

## 📂 Estrutura do Projeto

```text
.
├── src/
│   ├── ast.go          # Definições da Árvore Sintática (AST)
│   ├── compiler.go     # Geração de Assembly e execução
│   ├── interpreter.go  # Lógica do interpretador
│   ├── lexer.go        # Análise Léxica
│   ├── main.go         # Principal
│   ├── parser.go       # Análise Sintática
│   ├── visualizer.go   # Geração da visualização gráfica
│   └── main_test.go    # Testes automatizados
├── assembly/
│   └── runtime.s       # Rotinas de suporte 
├── arvore/             # Saída: Arquivos .dot e .png da árvore gerada
├── output/             # Saída: Arquivos .s, .o e o binário final
├── go.mod              # Gerenciador de dependências Go
└── README.md           # Documentação
```

## 🛠️ Como Executar

1.  **Baixe as dependências do Go:**

    ```bash
    go mod tidy
    ```

2.  **Execute o compilador passando a expressão entre aspas:**

    ```bash
    go run ./src "(3 * (2 + 5))"
    ```

### Saída Esperada

O programa irá imprimir no terminal:

- A estrutura da árvore em texto.
- O resultado calculado pelo interpretador.
- Mensagens de sucesso da geração da imagem e do binário.
- A execução do binário gerado.

Além disso, verifique:

- **Imagem da Árvore:** `arvore/arvore.png`
- **Código Assembly:** `output/output.s`
- **Executável:** `output/output`

## 🧪 Testes

O projeto possui testes automatizados que verificam o interpretador e a compilação de diversos casos (soma, subtração, aninhamento, divisão por zero, etc).

Para rodar os testes:

```bash
go test -v ./src
```
