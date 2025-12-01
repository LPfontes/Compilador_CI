  .section .text
  .globl _start

_start:

    mov $9,  %rax
    imul $8, %rax        
    imul $7, %rax

    mov $6,  %rbx
    imul $5, %rbx
    imul $4, %rbx
    imul $3, %rbx
    imul $2, %rbx

    idiv %rbx
    
  call imprime_num
  call sair

  .include "runtime.s"
  
