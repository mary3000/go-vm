package src

import (
	"bytes"
	"fmt"
	"github.com/lunixbochs/struc"
	"io/ioutil"
	"path/filepath"
	"strconv"
)

type Disassembler struct {
	src Memory
	asm string
}

func (disasm *Disassembler) Parse(fileName string) error {
	err := disasm.load(fileName)
	if err != nil {
		return err
	}

	varMap := make(map[int]string)
	varName := "var"
	varIdx := 0
	for i := 10; i < int(disasm.src.CodeOffset); i++ {
		varMap[i] = varName + strconv.Itoa(varIdx)
		varIdx++
		disasm.asm += "!" + varMap[i] + " int " + strconv.Itoa(int(disasm.src.Data[i])) + "\n"
	}

	strMap := make(map[int]string)

	for i := 0; i < int(disasm.src.TextSize); {
		data := disasm.src.Text[i:]
		end := bytes.IndexByte(data, 0)
		if end == -1 {
			break
		}
		data = data[:end]
		strMap[i] = varName + strconv.Itoa(varIdx)
		varIdx++
		disasm.asm += "!" + strMap[i] + " str " + string(data) + "\n"
		i += end + 1
	}

	labelName := "label"
	labelIdx := 0
	labelMap := make(map[int]string)
	for i := disasm.src.CodeOffset; i < disasm.src.StackOffset; {
		cmd := Command(disasm.src.Data[i])
		i++
		switch cmd {
		case GO:
			ref := int(disasm.src.Data[i+1])
			if _, ok := labelMap[ref]; !ok {
				labelMap[ref] = labelName + strconv.Itoa(labelIdx)
				labelIdx++
			}
			fallthrough
		case IN, OUT, PUSH, POP:
			i += 2
		case IF_LESS, IF_MORE, IF_EQ, MOV, ADD, SUB:
			i += 4
		}
	}

	cmdName := ""
	for i := disasm.src.CodeOffset; i < disasm.src.StackOffset; {
		cmdName = ""
		cmd := Command(disasm.src.Data[i])

		if label, ok := labelMap[int(i)]; ok {
			disasm.asm += "." + label + "\n"
		}

		i++
		switch cmd {
		case IN:
			disasm.asm += "in "
			i++
			disasm.asm += disasm.operandToString(int(i), varMap, false)
		case OUT:
			disasm.asm += "out "
			isString := disasm.src.Data[i] == 0
			i++
			refMap := varMap
			if isString {
				refMap = strMap
			}
			disasm.asm += disasm.operandToString(int(i), refMap, isString)
		case GO:
			disasm.asm += "go "
			i++
			disasm.asm += labelMap[int(disasm.src.Data[i])]
		case RET:
			disasm.asm += "ret"
			i--
		case FI:
			disasm.asm += "fi"
			i--
		case PUSH:
			disasm.asm += "push "
			i++
			disasm.asm += disasm.operandToString(int(i), varMap, false)
		case POP:
			disasm.asm += "pop "
			i++
			disasm.asm += disasm.operandToString(int(i), varMap, false)
		case IF_LESS:
			cmdName = "if< "
			fallthrough
		case IF_MORE:
			if cmdName == "" {
				cmdName = "if> "
			}
			fallthrough
		case IF_EQ:
			if cmdName == "" {
				cmdName = "if= "
			}
			fallthrough
		case MOV:
			if cmdName == "" {
				cmdName = "mov "
			}
			fallthrough
		case ADD:
			if cmdName == "" {
				cmdName = "add "
			}
			fallthrough
		case SUB:
			if cmdName == "" {
				cmdName = "sub "
			}
			fallthrough
		default:
			disasm.asm += cmdName
			i++
			disasm.asm += disasm.operandToString(int(i), varMap, false) + " "
			i += 2
			disasm.asm += disasm.operandToString(int(i), varMap, false)
		}
		disasm.asm += "\n"
		i++
	}

	return nil
}

func (disasm *Disassembler) WriteAsm(fileName string) error {
	if filepath.Ext(fileName) != AsmExt {
		return fmt.Errorf("Wrong extension, expeted: %d, actual: ", AsmExt, filepath.Ext(fileName))
	}
	return ioutil.WriteFile(fileName, []byte(disasm.asm), 0644)
}

func (disasm *Disassembler) load(fileName string) error {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return err
	}
	buf := bytes.NewBuffer(data)
	err = struc.Unpack(buf, &disasm.src)
	return err
}

func (disasm *Disassembler) operandToString(i int, varMap map[int]string, isString bool) string {
	if disasm.src.Data[i] < 10 && !isString {
		register := Register(disasm.src.Data[i])
		var name string
		switch register {
		case P1:
			name = "p1"
		case P2:
			name = "p2"
		case P3:
			name = "p3"
		case R1:
			name = "r1"
		case R2:
			name = "r2"
		case R3:
			name = "r3"
		case TMP:
			name = "tmp"
		case BP:
			name = "bp"
		case IP:
			name = "ip"
		}
		return "%" + name
	} else {
		return "!" + varMap[int(disasm.src.Data[i])]
	}
}
