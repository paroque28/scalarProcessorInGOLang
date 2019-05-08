package cpu

import "fmt"

type Memory struct {
	InExecuteControlSignals ExecuteControlSignals `json:"in_execute_control_signals"`
	OutControlSignals       MemControlSignals     `json:"control_signals"`
	InALUResult             uint64                `json:"alu_result"`
	ALUResult               uint64                `json:"alu_result"`
	MemoryData              uint64                `json:"memory_data"`
}

type MemControlSignals struct {
	WriteAddress        byte `json:"write_address"`
	MemToReg            bool `json:"mem_to_reg"`
	RegisterWriteEnable bool `json:"register_write_enable"`
}

const BITS64_BYTES = 64 / 8

func (mem *Memory) Run(done chan string, mainMemory []byte) {
	if mem.InExecuteControlSignals.MemToReg {
		//fmt.Println("[Mem] Mem Position to Read:", mem.InALUResult)
		mem.ALUResult = 0
		for i := uint64(0); i < BITS64_BYTES; i++ {
			fmt.Printf("[Mem] LOAD:  M[%x]=%x\n", mem.InALUResult+i, mainMemory[mem.InALUResult+i])
			mem.ALUResult |= uint64(mainMemory[mem.InALUResult+i]) << (8 * i)
		}
		fmt.Printf("[Mem] LOAD:  V[%x]=%x\n", mem.InExecuteControlSignals.WriteAddress, mem.ALUResult)
	} else {
		mem.ALUResult = mem.InALUResult
	}

	if mem.InExecuteControlSignals.MemWriteEnable {
		for i := uint64(0); i < BITS64_BYTES; i++ {
			dataIn := byte(mem.InALUResult >> (8 * i))
			baseAddress := mem.InExecuteControlSignals.MemWriteAddress
			fmt.Printf("[Mem] STORE: original byte[%x]: % x\n", baseAddress+i, mainMemory[baseAddress+i])
			fmt.Printf("[Mem] STORE: new      byte[%x]: % x\n", baseAddress+i, dataIn)
			mainMemory[baseAddress+i] = dataIn
		}
	}

	mem.setControlSignals(mem.InExecuteControlSignals.MemWriteEnable,
		mem.InExecuteControlSignals.MemToReg,
		mem.InExecuteControlSignals.RegisterWriteEnable,
		mem.InExecuteControlSignals.WriteAddress)
	done <- "mem"
}

func (mem *Memory) UpdateInRegisters(inControlSignals ExecuteControlSignals, ALUREsult uint64) {
	mem.InALUResult = ALUREsult
	mem.InExecuteControlSignals = inControlSignals
}

func (mem *Memory) setControlSignals(memWriteEnable bool,
	memToReg bool,
	registerWriteEnable bool,
	writeAddress byte) {

	mem.OutControlSignals.MemToReg = memToReg
	mem.OutControlSignals.RegisterWriteEnable = registerWriteEnable
	mem.OutControlSignals.WriteAddress = writeAddress
}
