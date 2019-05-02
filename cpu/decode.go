package cpu

import (
	"fmt"
)

type Decode struct {
	InInstruction     [4]byte              `json:"in_instruction"`
	OpCode            byte                 `json:"op_code"`
	Funct             byte                 `json:"funct"`
	Rd1               uint64               `json:"rd_1"`
	Rd2               uint64               `json:"rd_2"`
	Immediate         int64                `json:"immediate"`
	OutControlSignals DecodeControlSignals `json:"control_signals"`
}

type DecodeControlSignals struct {
	WriteAddress        byte `json:"write_address"`
	ALUControl          byte `json:"alu_control"`
	ALUSrcReg           bool `json:"alu_src_reg"`
	MemWriteEnable      bool `json:"mem_write_enable"`
	MemToReg            bool `json:"mem_to_reg"`
	RegisterWriteEnable bool `json:"register_write_enable"`
}

//opcodes
const AL_REG byte = 0x0
const AL_IMM byte = 0x1

//funct
const NOP byte = 0x0
const ADD byte = 0x1

func (deco *Decode) Run(done chan string, registers []uint64) {
	deco.OpCode = byte((deco.InInstruction[0] >> 4) & 0xF)
	deco.Funct = byte(deco.InInstruction[0] & 0xF)
	deco.OutControlSignals.WriteAddress = byte((deco.InInstruction[1] >> 3) & 0x1F)
	ra1 := byte((deco.InInstruction[1]&0x7)<<2 | ((deco.InInstruction[2] >> 6) & 0x3))
	ra2 := byte(deco.InInstruction[2]>>1) & 0x1F
	// Get data from registers
	deco.Rd1 = registers[ra1]
	deco.Rd2 = registers[ra2]
	deco.Immediate = int64(int64(deco.InInstruction[3]) | (int64(deco.InInstruction[2]) << 8) | (int64(deco.InInstruction[1]>>6) << 16))
	deco.Immediate = (deco.Immediate << 50) >> 50

	switch deco.OpCode {
	case AL_REG:
		deco.registerOperation(ra1, ra2)
	case AL_IMM:
		deco.immediateOperation(ra1)
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

func (deco *Decode) registerOperation(ra1 byte, ra2 byte) {
	switch deco.Funct {
	case NOP:
		fmt.Println("NOP")
	case ADD:
		fmt.Println("ADD", "V", deco.OutControlSignals.WriteAddress, "V", ra1, "V", ra2)
	default:
		panic("Not supported Reg instruction")
	}
}

func (deco *Decode) immediateOperation(ra1 byte) {
	switch deco.Funct {
	case NOP:
		fmt.Println("NOPI")
	case ADD:
		fmt.Println("ADDI", "V", deco.OutControlSignals.WriteAddress, "V", ra1, "#", deco.Immediate)
	default:
		fmt.Println("Funct: ", deco.Funct)
		panic("Not supported Imm instruction")
	}
}

func (deco *Decode) setControlSignals(writeAddress byte,
	aluControl byte,
	aluSrcReg bool,
	memWriteEnable bool,
	memToReg bool,
	registerWriteEnable bool) {
	deco.OutControlSignals.WriteAddress = writeAddress
	deco.OutControlSignals.ALUControl = aluControl
	deco.OutControlSignals.ALUSrcReg = aluSrcReg
	deco.OutControlSignals.MemToReg = memToReg
	deco.OutControlSignals.MemWriteEnable = memWriteEnable
	deco.OutControlSignals.RegisterWriteEnable = registerWriteEnable
}
