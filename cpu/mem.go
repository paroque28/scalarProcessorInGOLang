package cpu

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

func (mem *Memory) Run(done chan string) {
	mem.ALUResult = mem.InALUResult
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
