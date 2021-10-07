set shell := ["bash", "-c"]

default:
	@just --list --unsorted

compile FILE:
    [[ ! -d build ]] && mkdir -p build ; go run main.go -p {{FILE}} > build/main.asm && as -D -g -o build/main.o build/main.asm && gcc -g -fverbose-asm -o build/main build/main.o

run FILE:
    [[ ! -d build ]] && mkdir -p build ; go run main.go -p {{FILE}} > build/main.asm && as -D -g -o build/main.o build/main.asm && gcc -g -fverbose-asm -o build/main build/main.o && ./build/main

