  .section .text
  .globl _start

_start:

    mov $42, %rax
    sub $222,  %rax        
    imul $11, %rax

    mov $19,  %rbx
    imul $88, %rbx
    add %rbx, %rax  

    
  call imprime_num
  call sair

  .include "runtime.s"
  
