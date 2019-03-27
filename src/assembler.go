package src

import (
	"bytes"
	"fmt"
	"github.com/lunixbochs/struc"
	"os"
	"strconv"
	"strings"
)

type Assembler struct {
	Source string
	mem    Memory
}

func (asm *Assembler) Parse() error {
	commands := strings.Split(asm.Source, "\n")
	asm.mem.Data = make([]int64, 10)
	code := make([]int64, 0)
	varMap := make(map[string]int64)
	strMap := make(map[string]int64)
	labelMap := make(map[string]int64)
	goMap := make(map[int64]string)
	for _, cmd := range commands {
		cmd := strings.Trim(cmd, " \t")
		args := strings.Split(cmd, " ")

		if strings.HasPrefix(args[0], "!") {
			if args[1] == "int" {
				if len(args) != 3 {
					return fmt.Errorf("Assembler.parse: [!] expected 3 arguments, actual: %d", len(args))
				}

				num, _ := strconv.Atoi(args[2])
				asm.mem.Data = append(asm.mem.Data, int64(num))
				varMap[args[0][1:]] = int64(len(asm.mem.Data) - 1)
			} else if args[1] == "str" {
				str := strings.Join(args[2:], " ")
				strMap[args[0][1:]] = int64(len(asm.mem.Text))
				raw := []byte(str)
				raw = append(raw, 0)
				asm.mem.Text = append(asm.mem.Text, raw...)
			}
		} else if strings.HasPrefix(args[0], ".") {
			if len(args) != 1 {
				return fmt.Errorf("Assembler.parse: [.] expected 1 argument, actual: %d", len(args))
			}

			labelMap[args[0][1:]] = int64(len(code))
		} else if len(args) == 3 {
			var arg Command
			switch args[0] {
			case "mov":
				arg = MOV
			case "add":
				arg = ADD
			case "sub":
				arg = SUB
			case "if<":
				arg = IF_LESS
			case "if>":
				arg = IF_MORE
			case "if=":
				arg = IF_EQ
			}
			code = append(code, int64(arg))
			code = asm.parseOperand(args[1], code, varMap, strMap)
			code = asm.parseOperand(args[2], code, varMap, strMap)
		} else if len(args) == 2 {
			var arg Command
			switch args[0] {
			case "in":
				arg = IN
			case "out":
				arg = OUT
			case "push":
				arg = PUSH
			case "pop":
				arg = POP
			case "go":
				arg = GO
				goMap[int64(len(code) + 2)] = args[1]
			}

			code = append(code, int64(arg))
			code = asm.parseOperand(args[1], code, varMap, strMap)

		} else if len(args) == 1 && args[0] != "" {
			var arg Command
			switch args[0] {
			case "ret":
				arg = RET
			case "fi":
				arg = FI
			}
			code = append(code, int64(arg))
		}
	}

	asm.mem.TextSize = int64(len(asm.mem.Text))

	asm.mem.CodeOffset = int64(len(asm.mem.Data))
	for pos, label := range goMap {
		code[pos] = asm.mem.CodeOffset + labelMap[label]
	}

	asm.mem.Data = append(asm.mem.Data, code...)
	asm.mem.StackOffset = int64(len(asm.mem.Data))
	asm.mem.StackSize = 4096
	stack := make([]int64, asm.mem.StackSize)
	asm.mem.Data = append(asm.mem.Data, stack...)
	asm.mem.Size = int64(len(asm.mem.Data))
	return nil
}

func (asm *Assembler) parseOperand(arg string, code []int64, varMap map[string]int64, strMap map[string]int64) []int64 {
	if strings.HasPrefix(arg, "%") {
		var reg Register
		switch arg[1:] {
		case "p1":
			reg = P1
		case "p2":
			reg = P2
		case "p3":
			reg = P3
		case "r1":
			reg = R1
		case "r2":
			reg = R2
		case "r3":
			reg = R3
		case "tmp":
			reg = TMP
		case "sp":
			reg = SP
		case "bp":
			reg = BP
		case "ip":
			reg = IP
		}
		code = append(code, 1)
		code = append(code, int64(reg))
	} else if strings.HasPrefix(arg, "!") {
		if _, ok := strMap[arg[1:]]; ok {
			code = append(code, 0)
			code = append(code, strMap[arg[1:]])
		} else {
			code = append(code, 1)
			code = append(code, varMap[arg[1:]])
		}
	} else {
		code = append(code, 1)
		code = append(code, -1)
	}
	return code
}

func (asm *Assembler) WriteBinary(fileName string) error {
	file, err := os.Create(fileName)
	defer file.Close()
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	err = struc.Pack(&buf, &asm.mem)
	if err != nil {
		return err
	}
	_, err = file.Write(buf.Bytes())
	return err
}
