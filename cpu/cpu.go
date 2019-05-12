package cpu

import (
	"fmt"
)

const DEBUG = 1

type Processor struct {
	MainMemory         []byte      `json:"-"`
	InstructionsMemory []byte      `json:"-"`
	InClock            chan uint64 `json:"-"`
	Clock              uint64      `json:"clk"`
	PC                 uint        `json:"pc"`
	Fetch              *Fetch      `json:"fetch"`
	Decode             *Decode     `json:"decode"`
	Execute            *Execute    `json:"execute"`
	Memory             *Memory     `json:"memory"`
	Writeback          *WriteBack  `json:"writeBack"`
	Registers          []uint64    `json:"registers"`
}

func (proc *Processor) Init(clock chan uint64, mainMemory []byte, instructionsMemory []byte) {
	proc.MainMemory = mainMemory
	proc.InstructionsMemory = instructionsMemory
	proc.InClock = clock
	proc.Registers = make([]uint64, 32)
	proc.Fetch = new(Fetch)
	proc.Decode = new(Decode)
	proc.Execute = new(Execute)
	proc.Memory = new(Memory)
	proc.Writeback = new(WriteBack)
}

func (proc Processor) Start() {
	if DEBUG > 0 {
		fmt.Println("Starting Processor...")
	}
	for {
		done := make(chan string)
		proc.Clock = <-proc.InClock
		//println(proc.Clock)
		go proc.Fetch.Run(done, proc.InstructionsMemory)
		go proc.Decode.Run(done, proc.Registers)
		go proc.Execute.Run(done)
		go proc.Memory.Run(done, proc.MainMemory)
		go proc.Writeback.Run(done, proc.Registers)

		//JOIN
		for i := 0; i < 5; i++ {
			<-done
		}

		//Update Input Registers
		proc.PC += 4
		proc.Fetch.UpdateInRegisters(proc.PC)
		proc.Decode.UpdateInRegisters(proc.Fetch.Instruction)
		proc.Execute.UpdateInRegisters(proc.Decode.OutControlSignals, proc.Decode.Rd1, proc.Decode.Rd2, proc.Decode.Immediate)
		proc.Memory.UpdateInRegisters(proc.Execute.OutControlSignals, proc.Execute.ALUResult)
		proc.Writeback.UpdateInRegisters(proc.Memory.OutControlSignals, proc.Memory.ALUResult)
	}
}
