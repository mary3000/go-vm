# Golang Virtual Machine

It is simple implementation of Virtual Machine with its own asm language and Von Neumann architecture.


## Commands and marks of `gasm` language:
- `!var int 42` create int variable with name `var` and value `42`
- `.lbl` defines label, related to the code above
- `%` prefix means register
- `!` prefix means variable
- `in, out <smth>` read from / write to stdin / stdout
- `go <smth>` jump to specified label
- `if<, if>, if= <a> <b>` compares `a` and `b`, executes code right after compare if true, otherwise jumps to `fi`
- `ret` return to the caller
- `push, pop <smth>` push or pop value to or from stack
- `mov <dest> <src>` move `src` to `dest`
- `add, sub <dest> <src>` add or sub `src` to `dest`

## Usage

`./main <path-to-file>` or `go main <path-to-file>`, where file is written in `gasm` language and has `.gasm` extension.

Or you can run demo:
`./main` or `go main`  
It will execute `asm_data/fibonacci.gasm` file, which is a recursive version of fibonacci calculation.