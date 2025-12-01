  .section .text
  .globl _start

_start:
  
    mov $19, %rax
    imul $15, %rax        

    mov $7, %rbx
    imul $10, %rbx 

    sub %rbx, %rax
    
    mov $117, %rbx
    sub $33, %rbx

    add %rbx, %rax
    
  call imprime_num
  call sair

  .include "runtime.s"
  
