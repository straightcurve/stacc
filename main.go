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

    compile(lex(program));
}

const (
    OP_IF = iota;
    OP_BLOCK_START
    OP_BLOCK_END
    OP_SWAP
    OP_INCREMENT
    OP_LESS_THAN
    OP_CLONE
    OP_ADD
    OP_EXIT
    OP_PRINT
    OP_PUSH
    OP_WHILE
)

type Operation struct {
    op_type int
    args []int64
}

func lex(input []string) []Operation {
    var program []Operation;
    for i := 0; i < len(input); i++ {
        var op Operation;
        switch input[i] {
        case "swap": {
            op = Operation{};
            op.op_type = OP_SWAP;
            program = append(program, op);
        }
        case "inc": {
            op = Operation{};
            op.op_type = OP_INCREMENT;
            program = append(program, op);
        }
        case "<": {
            op = Operation{};
            op.op_type = OP_LESS_THAN;
            program = append(program, op);
        }
        case "if": {
            op = Operation{};
            op.op_type = OP_IF;
            op.args = []int64{};
            op.args = append(op.args, int64(i));
            program = append(program, op);
        }
        case "{": {
            op = Operation{};
            op.op_type = OP_BLOCK_START;
            program = append(program, op);
        }
        case "}": {
            op = Operation{};
            op.op_type = OP_BLOCK_END;
            program = append(program, op);
        }
        case "clone": {
            op = Operation{};
            op.op_type = OP_CLONE;
            program = append(program, op);
        }
        case "+": {
            op = Operation{};
            op.op_type = OP_ADD;
            program = append(program, op);
        }
        case "exit": {
            op = Operation{};
            op.op_type = OP_EXIT;
            program = append(program, op);
        }
        case "print": {
            op = Operation{};
            op.op_type = OP_PRINT;
            program = append(program, op);
        }
        default: {
            value, err := strconv.ParseInt(input[i], 0, 64);
            check(err);

            op = Operation{};
            op.op_type = OP_PUSH;
            op.args = []int64{};
            op.args = append(op.args, value);
            program = append(program, op);
        }
        }
    }

    return program;
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

func compile(program []Operation) {
    var output = header();
    var scope_i = -1;

    for i := 0; i < len(program); i++ {
        var op = program[i];
        switch op.op_type {
        case OP_SWAP:
            output += compile_swap();
        case OP_INCREMENT:
            output += compile_increment();
        case OP_WHILE:
            output += compile_while();
        case OP_LESS_THAN:
            output += compile_less_than(i);
        case OP_IF:
            output += compile_if();
        case OP_BLOCK_START:
            scope_i = i - 1;
            output += compile_scope_begin();
        case OP_BLOCK_END:
            output += compile_scope_end(scope_i);
        case OP_CLONE:
            output += compile_clone();
        case OP_ADD:
            output += compile_add();
        case OP_EXIT:
            output += compile_exit();
        case OP_PRINT:
            output += compile_print();
        case OP_PUSH:
            output += compile_push(op.args[0]);
        default: {
            fmt.Println(op.op_type, " << not sure what is");
        }
        }
    }

    fmt.Println(output);
}

func compile_less_than(i int) string {
    return fmt.Sprint(`
    pop w0
    pop w1
    cmp w1, w0
    b.gt .skip`, i);
}

func compile_swap() string {
    return `
    pop w0
    pop w1
    push w1
    push w0
    `;
}

func compile_print() string {
    return `
    pop w0
    bl print
    `;
}

func compile_if() string {
    return `
    `;
}

func compile_while() string {
    return `
.loop:
    `;
}

func compile_scope_end(i int) string {
    return fmt.Sprint(`
.skip`, i, `:`);
}

func compile_scope_begin() string {
    return `
    `;
}

func compile_increment() string {
    return `
    pop w0
    add w0, w0, #1
    push w0
    `;
}

func compile_clone() string {
    return `
    pop w0
    push w0
    push w0
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
