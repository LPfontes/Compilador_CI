Documentação de Extensões - Compilador Fun (CI)
Este documento detalha as extensões implementadas no compilador da linguagem Fun para o Projeto Final da disciplina de Construção de Compiladores I. O projeto foi estruturado para elevar o patamar de desempenho e expressividade do compilador através de múltiplas camadas de funcionalidade.

1. Extensões Simples (Base de Expressividade)
Estas extensões focam na melhoria da sintaxe e das capacidades lógicas básicas da linguagem:

Novos Operadores de Comparação: Adição de suporte completo (Lexer, Parser, Interpretador e Compilador) para os operadores relacionais <=, >= e !=.

Operadores Lógicos/Booleanos: Implementação das palavras-chave and, or e not, permitindo a criação de condições compostas complexas.

Instrução de Loop for: Adição da estrutura de controle for com suporte a inicialização, condição e passo. Foi implementada via desugaring interno para a estrutura de while, garantindo compatibilidade total com o gerador de código.

2. Extensões de Complexidade Média (Estruturas de Dados)
Suporte para Arrays de Inteiros (Vetores): Implementamos o suporte completo para vetores estáticos.

Escopo: Disponível tanto em escopo global (seção .bss) quanto local (alocado na pilha/stack frame).

Funcionalidades: Permite acesso e atribuição indexada dinâmica (ex: arr[i] = x;) e alocação correta para suporte a funções recursivas.

3. Extensões de Complexidade Alta (Performance e Hardware)
Estas funcionalidades representam o topo da complexidade técnica do projeto, focando em otimização de baixo nível:

Otimização de Código (Peephole Optimization): Uma etapa de pós-processamento que analisa o Assembly gerado para eliminar redundâncias.

Mecânica: Transforma sequências custosas como pushq seguido de popq em instruções movq diretas ou as remove quando desnecessárias.

Impacto: Melhora significativamente a velocidade de execução ao reduzir acessos à memória RAM.

Suporte a Operações com Vetores (AVX2 SIMD): Implementação de suporte nativo para instruções Single Instruction, Multiple Data da arquitetura x86-64.

Paralelismo: Permite processar até 4 inteiros de 64 bits em uma única instrução através de registradores de 256 bits (%ymm).

Builtins: Disponibiliza as funções vadd, vsub e vset.

Gerenciamento: O compilador lida automaticamente com o alinhamento de memória em 32 bytes (.balign 32) e gera loops otimizados para blocos de dados.
