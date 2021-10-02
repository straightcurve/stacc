package main

import "flag"
import "fmt"
import "strconv"
import "os"
import strings "strings"

func check(err error) {
    if err != nil {
        panic(err);
    }
}

func main() {
    var file = flag.String("p", "", "path to source file");

    flag.Parse();

    if *file == "" {
        flag.PrintDefaults();
        os.Exit(1);
    }

    dat, err := os.ReadFile(*file);
    check(err);

    var str_data = string(dat);
    var program []string;
    var initial = strings.Split(str_data, "\n");
    for i := 0; i < len(initial); i++ {
        var split_by_space = strings.Split(initial[i], " ");
        for j := 0; j < len(split_by_space); j++ {
            if len(split_by_space[j]) == 0 {
                continue;
            }

            program = append(program, split_by_space[j]);
        }
    }

    compile(program);
}

func header() string {
    return `
.global main
.type main, @function
.global print
.type print, @function

.data
.bss
.text
.macro push Rn:req
    str     \Rn, [sp, #-16]!
.endm
.macro pop Rn:req
    ldr     \Rn, [sp], #16
.endm

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
    `;
}

func compile(program []string) {
    var output = header();

    for i := 0; i < len(program); i++ {
        switch program[i] {
            case "+":
                output += compile_add();
            case "exit":
                output += compile_exit();
            case "print":
                output += compile_print();
            default: {
                value, err := strconv.ParseInt(program[i], 0, 64);
                check(err);

                output += compile_push(value);
            }
        }
    }

    fmt.Println(output);
}

func compile_print() string {
    return `
    pop w0
    bl print
    `;
}

func compile_push(value int64) string {
    return fmt.Sprintf(`
    mov w0, %d  // #%d
    push w0
    `, value, value);
}

func compile_add() string {
    return `
    pop w1
    pop w0
    add w0, w1, w0
    push w0
    `;
}

func compile_exit() string {
    return `
    mov x8, 93
    svc 0
    `;
}
