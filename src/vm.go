package src

import (
	"bytes"
	"fmt"
	"github.com/lunixbochs/struc"
	"io/ioutil"
)

type VirtualMachine struct {
	mem Memory
	ip *int64
	sp *int64
}

func (vm *VirtualMachine) Exec(fileName string) error {
	err := vm.load(fileName)
	if err != nil {
		return err
	}

	vm.mem.Data[IP] = vm.mem.CodeOffset
	vm.ip = &vm.mem.Data[IP]
	vm.mem.Data[SP] = vm.mem.StackOffset - 1
	vm.sp = &vm.mem.Data[SP]
	for *vm.ip < vm.mem.StackOffset {
		switch Command(vm.next()) {
		case IN:
			vm.next()
			fmt.Println("~Enter number:")
			fmt.Scanln(&vm.mem.Data[vm.next()])
		case OUT:
			fmt.Println("~Output:")
			if vm.next() == 0 {
				data := vm.mem.Text[vm.next():]
				data = data[:bytes.IndexByte(data, 0)]
				fmt.Println(string(data))
			} else {
				fmt.Println(vm.mem.Data[vm.next()])
			}
		case GO:
			vm.next()
			vm.push(*vm.ip + 1)
			*vm.ip = vm.next()
		case RET:
			if *vm.sp < vm.mem.StackOffset {
				return nil
			}
			vm.pop(vm.ip)
		case IF_LESS:
			vm.next()
			l := vm.next()
			vm.next()
			r := vm.next()
			if vm.mem.Data[l] >= vm.mem.Data[r] {
				for Command(vm.next()) != FI {}
			}
		case IF_MORE:
			vm.next()
			l := vm.next()
			vm.next()
			r := vm.next()
			if vm.mem.Data[l] <= vm.mem.Data[r] {
				for Command(vm.next()) != FI {}
			}
		case IF_EQ:
			vm.next()
			l := vm.next()
			vm.next()
			r := vm.next()
			if vm.mem.Data[l] != vm.mem.Data[r] {
				for Command(vm.next()) != FI {}
			}
		case MOV:
			vm.next()
			receiverIndex := vm.next()
			vm.next()
			valIndex := vm.next()
			vm.mem.Data[receiverIndex] = vm.mem.Data[valIndex]
		case ADD:
			vm.next()
			receiverIndex := vm.next()
			vm.next()
			valIndex := vm.next()
			vm.mem.Data[receiverIndex] += vm.mem.Data[valIndex]
		case SUB:
			vm.next()
			receiverIndex := vm.next()
			vm.next()
			valIndex := vm.next()
			vm.mem.Data[receiverIndex] -= vm.mem.Data[valIndex]
		case PUSH:
			vm.next()
			vm.push(vm.mem.Data[vm.next()])
		case POP:
			vm.next()
			vm.pop(&vm.mem.Data[vm.next()])
		}
	}
	return nil
}

func (vm *VirtualMachine) load(fileName string) error {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return err
	}
	buf := bytes.NewBuffer(data)
	err = struc.Unpack(buf, &vm.mem)
	return err
}

func (vm *VirtualMachine) next() int64 {
	res := vm.mem.Data[*vm.ip]
	*vm.ip++
	return res
}

func (vm *VirtualMachine) push(val int64) {
	*vm.sp++
	vm.mem.Data[*vm.sp] = val
}

func (vm *VirtualMachine) pop(addr *int64) {
	*addr = vm.mem.Data[*vm.sp]
	vm.mem.Data[*vm.sp] = -1
	*vm.sp--
}