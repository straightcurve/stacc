//    makes the symbol visible to ld
.global main

// this is used for constants
.data

// this is used for variables
.bss

// this is used for instructions
.text
.macro push Rn:req
      str     \Rn, [sp, #-16]!
.endm
  .macro pop Rn:req
      ldr     \Rn, [sp], #16
  .endm

//.type print, @function

.global print
.type print, @function
print:                              // @print(unsigned long)
        sub     sp, sp, #96
        stp     x29, x30, [sp, #80]             // 16-byte Folded Spill
        add     x29, sp, #80
        mov     x10, #-3689348814741910324
        mov     w9, #10
        mov     w8, #-2
        movk    x10, #52429
        add     x11, sp, #12
        strb    w9, [sp, #76]
.LBB0_1:                                // =>This Inner Loop Header: Depth=1
        umulh   x12, x0, x10
        lsr     x12, x12, #3
        msub    w14, w12, w9, w0
        add     w13, w8, #65
        cmp     x0, #9
        orr     w14, w14, #0x30
        sub     w8, w8, #1
        mov     x0, x12
        strb    w14, [x11, w13, uxtw]
        b.hi    .LBB0_1
        neg     w2, w8
        add     w8, w8, #65
        add     x9, sp, #12
        add     x1, x9, x8
        mov     w0, #1
        strb    wzr, [x1]
        bl      write
        ldp     x29, x30, [sp, #80]             // 16-byte Folded Reload
        add     sp, sp, #96
        ret
main:
    mov w0, 200
    push w0

    mov w0, 220 
    push w0

    pop w1
    pop w0
    add w0, w1, w0

    bl print

    mov x8, 93
    svc 0
