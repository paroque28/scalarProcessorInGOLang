package cpu

import "fmt"

type Processor struct {
	MainMemory         []byte
	InstructionsMemory []byte
	Clock              chan uint64
}

func (proc Processor) Start() {
	for {
		tick := <-proc.Clock
		fmt.Println(tick)
		go fetch()
		go decode()
		go execute()
		go writeback()
		go mem()
	}
}
