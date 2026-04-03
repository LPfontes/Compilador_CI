# Compilador e Interpretador de Linguagem Fun (CI)

Este projeto é a implementação prática do compilador para a linguagem procedimental **Fun** (derivada das antigas linguagens Expr e Cmd), construído para a disciplina de **Construção de Compiladores I**. Ele é um compilador Turing-Completo nativo, capaz de gerar código diretamente para a arquitetura Assembly x86-64, lidando com controle de rotinas, ponteiros de pilha, recursão, funções e vetores.

## 🚀 Funcionalidades

O projeto executa o pipeline completo de compilação:

1.  **Análise Léxica (Lexer):** Tokenização de Identificadores, Keywords (`var`, `fun`, `main`, `if`, `else`, `for`, `while`, `return`), Operadores Aritméticos (`+`, `-`, `*`, `/`), Relacionais (`<`, `>`, `==`, `!=`, `<=`, `>=`), Lógicos (`and`, `or`, `not`) e delimitadores (`{ }`, `[ ]`, `( )`).
2.  **Análise Sintática (Parser):** Parser LL recursivo com lookahead para distinguir entre acesso a variáveis, chamadas de função (`nome(args)`) e acesso a vetores (`nome[indice]`).
3.  **Analisador Semântico:** Validação de escopo cruzando tabelas Globais e Locais, verificação de aridade em chamadas de função e validação de declarações de vetores.
4.  **Interpretador:** Execução simulada via tree-walk com frames isolados por chamada de função, suporte a vetores via chaves compostas no mapa de memória.
5.  **Compilador Assembly x86-64:** Geração de código AT&T com:
    - Stack Frames com prólogo/epílogo (`%rbp`/`%rsp`)
    - Endereçamento de parâmetros via offsets positivos de `%rbp`
    - Variáveis locais e vetores alocados na pilha
    - Vetores globais na seção `.bss` com `.lcomm`
    - Acesso indexado a vetores via endereçamento base+índice (`leaq` + `(%rdx, %rcx, 8)`)
6.  **Otimizador Peephole:** Pós-processamento do assembly gerado que elimina pares redundantes de `pushq`/`popq`, substituindo por `movq` diretos ou removendo completamente quando desnecessários.
7.  **Operações Vetoriais AVX2 (SIMD):** Builtins que mapeiam diretamente para instruções AVX2, processando **4 inteiros de 64 bits em paralelo** usando registradores `%ymm` de 256 bits:
    - `vadd(dest, src1, src2, n)` → `vpaddq` (soma paralela)
    - `vsub(dest, src1, src2, n)` → `vpsubq` (subtração paralela)
    - `vset(dest, valor, n)` → `vpbroadcastq` (preenchimento broadcast)
    - Loop automático para vetores maiores que 4 elementos
    - Alinhamento `.balign 32` para vetores globais
8.  **Integração (Montagem/Ligação):** Execução automática de `as` e `ld` para gerar o binário.

## 📝 Sintaxe da Linguagem

### Variáveis e Vetores
```
var x = 10;          // variável escalar global
var arr[5];          // vetor de 5 inteiros (inicializado com zeros)
```

### Funções
```
fun soma(a, b) {
    return a + b;
}
```

### Vetores em Funções
```
fun preenche(n) {
    var v[10];
    v[0] = n;
    v[1] = n * 2;
    return v[0] + v[1];
}
```

### Estruturas de Controle
```
if x > 0 {
    y = 1;
} else {
    y = 0;
}

while i < 10 {
    v[i] = i * 2;
    i = i + 1;
}
```

### Programa Completo
```
var v[5];
var i = 0;

fun soma_vetor(n) {
    var total = 0;
    var j = 0;
    while j < n {
        total = total + v[j];
        j = j + 1;
    }
    return total;
}

main {
    while i < 5 {
        v[i] = i * 10;
        i = i + 1;
    }
    return soma_vetor(5);
}
```

### Operações Vetoriais (AVX2)
```
var a[4];
var b[4];
var c[4];

main {
    a[0] = 1; a[1] = 2; a[2] = 3; a[3] = 4;
    b[0] = 10; b[1] = 20; b[2] = 30; b[3] = 40;

    vadd(c, a, b, 4);    // c[i] = a[i] + b[i] usando AVX2 (vpaddq)
    vsub(c, c, b, 4);    // c[i] = c[i] - b[i] usando AVX2 (vpsubq)
    vset(c, 42, 4);      // c[i] = 42 para todos usando AVX2 (vpbroadcastq)

    return c[0];
}
```

> **Nota:** Os builtins vetoriais processam 4 inteiros de 64 bits por instrução. O tamanho `n` deve ser múltiplo de 4.

## Pré-requisitos

- **Go** (1.23+)
- **Graphviz** (para visualização da AST)
- **GCC/Binutils** (`as` e `ld`)
- **CPU com AVX2** (para executar binários com operações vetoriais — `grep avx2 /proc/cpuinfo`)

### Instalação (Linux)
```bash
sudo apt update
sudo apt install golang graphviz build-essential
```

## 📂 Estrutura do Projeto

```text
.
├── src/
│   ├── ast.go              # Estruturas da AST (Programa, Decl, FunDecl, AcessoVetor, etc.)
│   ├── compiler.go         # Gerador de Assembly x86-64 + Otimizador Peephole
│   ├── semantic.go         # Verificador semântico de escopos e declarações
│   ├── interpreter.go      # Interpretador tree-walk com suporte a vetores
│   ├── lexer.go            # Scanner léxico
│   ├── main.go             # Ponto de entrada
│   ├── parser.go           # Parser LL recursivo
│   ├── visualizer.go       # Gerador de grafos Graphviz
│   ├── main_test.go        # Testes do compilador x86-64
│   ├── main_func_test.go   # Testes de funções e vetores no compilador
│   └── parser_test.go      # Testes do interpretador (precedência, variáveis, vetores)
├── assembly/
│   └── runtime.s           # Runtime: imprime_num e sair (syscalls Linux)
├── arvore/                 # Saída: imagens .png da AST
├── output/                 # Saída: assembly .s e binários
```

## 🛠️ Como Executar

1. **Acesse a pasta src:**
    ```bash
    cd src
    ```

2.  **Execute com a entrada desejada:**
    ```bash
    go run ./ "fun soma(a, b) { return a + b; } main { return soma(2, 3); }"
    ```

    Exemplo com vetores:
    ```bash
    go run ./ "var v[5]; main { v[0] = 10; v[1] = 20; return v[0] + v[1]; }"
    ```

### Saída Esperada
- Resultado numérico calculado pelo interpretador.
- Imagem da AST em `arvore/arvore.png`.
- Assembly otimizado em `output/output.s`.
- Execução do binário nativo gerado.

## 🧪 Testes

A suíte contém testes de precedência, variáveis, operadores relacionais/lógicos, funções (incluindo Fibonacci recursivo), vetores globais/locais e compilação x86-64.

```bash
cd src
go test -v ./
```
