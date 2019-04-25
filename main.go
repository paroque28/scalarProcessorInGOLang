package main

import (
	"fmt"
	"scalarProcessor/memory"
)


func main() {
	//Create memories
	mainMemory := make([]byte,131032);
	instructionsMemory := make([]byte,512);

	//Initialize memories
	memory.InitializeInstructionMemory(instructionsMemory);
	memory.InitializeMainMemory(mainMemory)
	fmt.Print(instructionsMemory[0]);


	//Save image
	memory.SaveImage(mainMemory,"result.png")


}
