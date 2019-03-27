package src

type Memory struct {
	CodeOffset int64
	StackOffset int64
	StackSize int64
	Size int64 `struc:"sizeof=Data"`
	Data []int64
	TextSize int64 `struc:"sizeof=Text"`
	Text []byte
}

type Command int64

const (
	IN Command = iota
	OUT
	GO
	RET
	IF_EQ
	IF_LESS
	IF_MORE
	FI
	MOV
	PUSH
	POP
	ADD
	SUB
)

type Register int64

const (
	P1 Register = iota
	P2
	P3
	R1
	R2
	R3
	TMP
	SP
	BP
	IP
)

// Why not ¯\_(ツ)_/¯
const AsmExt = ".gasm"
const BinExt = ".fuff"

const FailEmoji = "(ノಠ益ಠ)ノ彡┻━┻"