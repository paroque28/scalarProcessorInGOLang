package cpu

import (
	"fmt"
)

type Decode struct {
	InInstruction     [4]byte              `json:"in_instruction"`
	OpCode            byte                 `json:"op_code"`
	Funct             byte                 `json:"funct"`
	Ra1               byte                 `json:"ra_1"`
	Ra2               byte                 `json:"ra_2"`
	Rd1               uint64               `json:"rd_1"`
	Rd2               uint64               `json:"rd_2"`
	Immediate         int64                `json:"immediate"`
	OutControlSignals DecodeControlSignals `json:"control_signals"`
}

type DecodeControlSignals struct {
	MemWriteAddress     uint64 `json:"mem_write_address"`
	WriteAddress        byte   `json:"write_address"`
	ALUControl          byte   `json:"alu_control"`
	ALUSrcReg           bool   `json:"alu_src_reg"`
	MemWriteEnable      bool   `json:"mem_write_enable"`
	MemToReg            bool   `json:"mem_to_reg"`
	RegisterWriteEnable bool   `json:"register_write_enable"`
}

//opcodes
const AL_REG byte = 0x0
const AL_IMM byte = 0x1

//const SHIFTS byte = 0x2
//const MEM_IMM byte = 0x3
const MEM_REG byte = 0x4

//funct AL
const NOP byte = 0x0
const ADD byte = 0x1

//funct mem
const LOAD_64 byte = 0x1
const STORE_64 byte = 0x2

func (deco *Decode) Run(done chan string, registers []uint64) {
	deco.OpCode = byte((deco.InInstruction[0] >> 4) & 0xF)
	deco.Funct = byte(deco.InInstruction[0] & 0xF)
	deco.OutControlSignals.WriteAddress = byte((deco.InInstruction[1] >> 3) & 0x1F)
	deco.Ra1 = byte((deco.InInstruction[1]&0x7)<<2 | ((deco.InInstruction[2] >> 6) & 0x3))
	deco.Ra2 = byte(deco.InInstruction[2]>>1) & 0x1F
	//fmt.Println(deco.Ra1, deco.Ra2)
	// Get data from registers
	deco.Rd1 = registers[deco.Ra1]
	deco.Rd2 = registers[deco.Ra2]
	deco.Immediate = int64(int64(deco.InInstruction[3]) | (int64(deco.InInstruction[2]) << 8) | (int64(deco.InInstruction[1]>>6) << 16))
	deco.Immediate = (deco.Immediate << 50) >> 50

	switch deco.OpCode {
	case AL_REG:
		deco.registerOperation(deco.Ra1, deco.Ra2)
	case AL_IMM:
		deco.immediateOperation(deco.Ra1)
	case MEM_REG:
		deco.memRegOperation(deco.Ra1, registers)
	default:
		panic("[Deco] Not supported instruction")
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
		deco.setControlSignals(ALU_NOP, true, false, false, false)
		//fmt.Println("[Deco] NOP")
	case ADD:
		deco.setControlSignals(ALU_ADD, true, false, false, true)
		fmt.Println("[Deco] ADD", "V", deco.OutControlSignals.WriteAddress, "V", ra1, "V", ra2)
	default:
		panic("[Deco] Not supported Reg instruction")
	}
}

func (deco *Decode) immediateOperation(ra1 byte) {
	switch deco.Funct {
	case NOP:
		deco.setControlSignals(ALU_NOP, false, false, false, false)
		fmt.Println("[Deco] NOPI")
	case ADD:
		deco.setControlSignals(ALU_ADD, false, false, false, true)
		fmt.Println("[Deco] ADDI", "V", deco.OutControlSignals.WriteAddress, "V", ra1, "#", deco.Immediate)
	default:
		fmt.Println("Funct: ", deco.Funct)
		panic("[Deco] Not supported Imm instruction")
	}
}

func (deco *Decode) memRegOperation(ra1 byte, registers []uint64) {
	switch deco.Funct {
	case LOAD_64:
		deco.setControlSignals(ALU_BUFFER, true, false, true, true)
		fmt.Printf("[Deco] LOAD V%d, V%d : V[%x]=M[%x]\n", deco.OutControlSignals.WriteAddress, ra1, deco.OutControlSignals.WriteAddress, deco.Rd1)
	case STORE_64:
		deco.setControlSignals(ALU_BUFFER, true, true, false, false)
		deco.OutControlSignals.MemWriteAddress = uint64(registers[deco.OutControlSignals.WriteAddress])
		fmt.Printf("[Deco] STORE V%d  V%d : M[%x] = %x \n", deco.OutControlSignals.WriteAddress, deco.Ra1, deco.OutControlSignals.MemWriteAddress, deco.Rd1)

	default:
		fmt.Println("Funct: ", deco.Funct)
		panic("[Deco] Not supported Mem instruction")
	}
}

func (deco *Decode) setControlSignals(aluControl byte,
	aluSrcReg bool,
	memWriteEnable bool,
	memToReg bool,
	registerWriteEnable bool) {

	deco.OutControlSignals.ALUControl = aluControl
	deco.OutControlSignals.ALUSrcReg = aluSrcReg
	deco.OutControlSignals.MemToReg = memToReg
	deco.OutControlSignals.MemWriteEnable = memWriteEnable
	deco.OutControlSignals.RegisterWriteEnable = registerWriteEnable
}
