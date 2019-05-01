package cpu

import "fmt"

type Decode struct {
	InInstruction [4]byte `json:"in_instruction"`
	OpCode        uint8   `json:"op_code"`
	Funct         uint8   `json:"funct"`
	Rd1           uint8   `json:"rd_1"`
	Rd2           uint8   `json:"rd_2"`
	//Control Signals
	WriteAddress        uint8 `json:"write_address"`
	ALUControl          uint8 `json:"alu_control"`
	ALUSrcReg           bool  `json:"alu_src_reg"`
	MemWriteEnable      bool  `json:"mem_write_enable"`
	MemToReg            bool  `json:"mem_to_reg"`
	RegisterWriteEnable bool  `json:"register_write_enable"`
}

//opcodes
const AL_REG uint8 = 0x0
const AL_IMM uint8 = 0x1

//funct
const NOP uint8 = 0x0
const ADD uint8 = 0x1

func (deco *Decode) Run(done chan string, registers []uint64) {
	deco.OpCode = byte(deco.InInstruction[0] >> 4)
	deco.Funct = byte(deco.InInstruction[0] & 0xF)
	switch deco.OpCode {
	case AL_REG:
		deco.registerOperation()
	case AL_IMM:
		deco.immediateOperation()
	default:
		panic("Not supported instruction")
	}
	done <- "decode"
}

func (deco *Decode) UpdateInRegisters(instruction [4]byte) {
	for i := uint(0); i < 4; i++ {
		deco.InInstruction[i] = instruction[i]
	}
}

func (deco *Decode) registerOperation() {
	switch deco.Funct {
	case NOP:
		fmt.Println("NOP")
	case ADD:
		fmt.Println("ADD")
	default:
		panic("Not supported Reg instruction")
	}
}

func (deco *Decode) immediateOperation() {
	switch deco.Funct {
	case NOP:
		fmt.Println("NOPI")
	case ADD:
		fmt.Println("ADDI")
	default:
		fmt.Println("Funct: ", deco.Funct)
		panic("Not supported Imm instruction")
	}
}
