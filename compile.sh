
cd target
go run ../main.go | llc -O0 --x86-asm-syntax=intel > a.asm

clang -c -O3 -march=native -ffast-math a.asm