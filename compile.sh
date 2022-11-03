
cd target
go run ../main.go | llc -O0 --x86-asm-syntax=intel > a.asm

clang -c -O0 -march=native a.asm
ld a.o