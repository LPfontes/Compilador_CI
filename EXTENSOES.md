# Documentação de Extensões - Compilador Fun (CI)

Este documento detalha as extensões implementadas no compilador da linguagem **Fun** para o Projeto Final da disciplina de **Construção de Compiladores I**. Seguindo os critérios de complexidade definidos, o projeto implementou múltiplas extensões que elevam o compilador a um patamar de alto desempenho e expressividade.

---

## 1. Extensões de Complexidade Alta

### 1.1. Suporte a Operações com Vetores (AVX2 SIMD)
Implementamos suporte nativo para instruções SIMD (Single Instruction, Multiple Data) utilizando as extensões **AVX2** da arquitetura x86-64. Isso permite o processamento paralelo de 4 inteiros de 64 bits em uma única instrução.

- **Instruções Mapeadas**: `vpaddq` (soma), `vpsubq` (subtração), `vpbroadcastq` (broadcast).
- **Builtins**: `vadd(dest, src1, src2, n)`, `vsub(dest, src1, src2, n)`, `vset(dest, valor, n)`.
- **Destaque**: O compilador gerencia automaticamente o alinhamento de memória (`.balign 32`) e gera loops otimizados para processar blocos de 4 elementos.

### 1.2. Otimização de Código (Peephole Optimization)
Foi adicionada uma etapa de otimização de pós-processamento do código Assembly gerado.

- **Funcionalidade**: Analisa as sequências de instruções para eliminar redundâncias.
- **Exemplo**: Transforma pares `pushq $X` + `popq %rax` em um único `movq $X, %rax`, e elimina `pushq %r` + `popq %r` redundantes.
- **Impacto**: Redução significativa no número de acessos à memória (pilha) e melhoria na velocidade de execução.

---

## 2. Extensões de Complexidade Média

### 2.1. Suporte para Arrays de Inteiros (Vetores)
Implementamos suporte completo para vetores estáticos, tanto em escopo global quanto local (pilha).

- **Sintaxe**: `var arr[10];`
- **Funcionalidades**: Acesso e atribuição indexada (`arr[i] = x;`), suporte a indexação dinâmica (usando expressões/variáveis como índice) e alocação correta em Frames de Pilha para suporte a recursão.

---

## 3. Extensões Simples

### 3.1. Novos Operadores de Comparação
Adicionamos suporte completo (Lexer, Parser, Interpretador e Compilador) para os operadores relacionais faltantes:
- **Operadores**: `<=`, `>=`, `!=` (diferente).

### 3.2. Operadores Lógicos/Booleanos
Implementamos suporte para lógica booleana, permitindo composições complexas em condições de `if`, `while` e `for`.
- **Palavras-chave**: `and`, `or`, `not`.

### 3.3. Instrução de Loop `for`
Embora não estivesse explicitamente na lista de sugestões, adicionamos a estrutura de controle `for` para melhorar a expressividade da linguagem.
- **Sintaxe**: `for i = 0; i < n; i = i + 1 { ... }`
- **Implementação**: Desugaring interno para a estrutura de `while`, garantindo compatibilidade total com o gerador de código x86-64.

---

## 4. Resumo de Conformidade

De acordo com o documento `projeto.txt`, o grupo deveria entregar uma extensão de complexidade média/alta ou três simples. Este projeto entrega:

- **2 Extensões de Complexidade Alta** (AVX2 e Otimização)
- **1 Extensão de Complexidade Média** (Arrays)
- **3 Extensões Simples** (Operadores e Loop For)
---

## Como Usar as Extensões

Consulte o [README.md](README.md) para exemplos detalhados de sintaxe e comandos para compilação e testes.
