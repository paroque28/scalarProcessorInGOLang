package cpu

type Decode struct {
	InInstruction [4]byte `json:"in_instruction"`
	Rd1           uint32  `json:"rd_1"`
	Rd2           uint32  `json:"rd_2"`
}

func (deco *Decode) Run(done chan string, registers []uint64) {

	done <- "decode"
}

func (deco *Decode) UpdateInRegisters(instruction [4]byte) {
	for i := uint(0); i < 4; i++ {
		deco.InInstruction[i] = instruction[i]
	}
}
