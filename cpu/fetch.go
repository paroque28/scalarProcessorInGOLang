package cpu

type Fetch struct {
	InPC        uint    `json:"in_pc"`
	Instruction [4]byte `json:"instruction"`
}

func (fetch *Fetch) Run(done chan string, instructionMemory []byte) {
	for i := uint(0); i < 4; i++ {
		fetch.Instruction[i] = instructionMemory[fetch.InPC+i]
	}

	done <- "fetch"
}

func (fetch *Fetch) UpdateInRegisters(PC uint) {
	fetch.InPC = PC
}
