package cpu

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

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
	fmt.Println("Starting Processor...")
	for {
		done := make(chan string)
		proc.Clock = <-proc.InClock
		go proc.Fetch.Run(done, proc.InstructionsMemory)
		go proc.Decode.Run(done, proc.Registers)
		go proc.Execute.Run(done)
		go proc.Memory.Run(done)
		go proc.Writeback.Run(done, proc.Registers)

		//JOIN
		for i := 0; i < 5; i++ {
			<-done
		}
		proc.saveState()

		//Update Input Registers
		proc.PC += 4
		proc.Fetch.UpdateInRegisters(proc.PC)
		proc.Decode.UpdateInRegisters(proc.Fetch.Instruction)
		proc.Execute.UpdateInRegisters(proc.Decode.OutControlSignals, proc.Decode.Rd1, proc.Decode.Rd2, proc.Decode.Immediate)
		proc.Memory.UpdateInRegisters(proc.Execute.OutControlSignals, proc.Execute.ALUResult)
		proc.Writeback.UpdateInRegisters(proc.Memory.OutControlSignals, proc.Memory.ALUResult)
	}
}

func (proc *Processor) saveState() {
	instant, err := json.Marshal(proc)
	catch(err)
	//f, err := os.OpenFile("now.json", os.O_APPEND|os.O_WRONLY, 0600)
	f, err := os.Create("now.json")
	catch(err)
	w := bufio.NewWriter(f)
	_, err = w.WriteString(string(instant) + "\n")
	catch(err)
	err = w.Flush()
	catch(err)
	//fmt.Println("JSON Saved!")
}

func catch(err error) {
	if err != nil {
		fmt.Println(err)
		return
	}
}
