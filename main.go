package main

import (
	"fmt"
	"github.com/apex/log"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
	"virtualMachine/src"
)

func main() {
	if len(os.Args) > 2 {
		log.Errorf("Wrong number of arguments. Expected: 0 or 1, given: %d", len(os.Args)-1)
		log.Errorf(src.FailEmoji)
		return
	}

	fileName := "asm_data/fibonacci.gasm"
	if len(os.Args) > 1 {
		fileName = os.Args[1]
	}

	ext := filepath.Ext(fileName)
	if ext != src.AsmExt {
		log.Errorf("Wrong extention. Expected: %s, given: %s", src.AsmExt, ext)
		log.Errorf(src.FailEmoji)
		return
	}

	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Errorf(err.Error())
		log.Errorf(src.FailEmoji)
		return
	}

	str := string(b)
	asm := src.Assembler{Source: str}

	err = asm.Parse()
	if err != nil {
		log.Errorf(err.Error())
		log.Errorf(src.FailEmoji)
		return
	}

	name := strings.TrimSuffix(path.Base(fileName), ext)
	err = asm.WriteBinary(name + src.BinExt)
	if err != nil {
		log.Errorf(err.Error())
		log.Errorf(src.FailEmoji)
		return
	}

	fmt.Println("\\ (•◡•) /  Welcome aboard!  \\ (•◡•) /")
	time.Sleep(500 * time.Millisecond)
	fmt.Println("I am a tiny Virtual Machine, written on Go.")

	if len(os.Args) == 1 {
		fmt.Println("...But you are actually weaker than me, bag of bones! (¬‿¬)")
		fmt.Println("And you know why? I can calculate Fibonacci sequences 100500x faster than you! Lol")
		fmt.Println("☜(˚▽˚)☞")
		fmt.Println("And I can prove it. See!")
	}

	vm := src.VirtualMachine{}
	err = vm.Exec(name + src.BinExt)
	if err != nil {
		log.Errorf(err.Error())
		log.Errorf(src.FailEmoji)
		return
	}

	time.Sleep(500 * time.Millisecond)
	fmt.Println("What a power! I'm so strooong ᕙ(⇀‸↼‶)ᕗ")
	fmt.Println("See ya next time!")

	disasm := src.Disassembler{}
	err = disasm.Parse(name + src.BinExt)
	if err != nil {
		log.Errorf(err.Error())
		log.Errorf(src.FailEmoji)
		return
	}

	err = disasm.WriteAsm(name + src.AsmExt)
	if err != nil {
		log.Errorf(err.Error())
		log.Errorf(src.FailEmoji)
	}
}
