package main

import (
	"scalarProcessor/cpu"
	"scalarProcessor/memory"
	"time"
)

func main() {
	//Create memories
	mainMemory := make([]byte, 131032)
	instructionsMemory := make([]byte, 512)

	//Initialize memories
	memory.InitializeInstructionMemory(instructionsMemory)
	memory.InitializeMainMemory(mainMemory)

	//Create clock
	clock := make(chan uint64)

	//Create CPU
	processor := new(cpu.Processor)
	processor.Init(clock, mainMemory, instructionsMemory)

	go processor.Start()
	for i := uint64(0); i < 10; i++ {
		time.Sleep(1 * time.Second)
		clock <- i
	}

	//Save image
	memory.SaveImage(mainMemory, "result.png")

}
