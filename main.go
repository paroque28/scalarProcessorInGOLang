package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"scalarProcessor/cpu"
	"scalarProcessor/memory"
)

const saveFile = "cpu.json"
const allSnapshots = false

func initJSON() {
	f, err := os.Create(saveFile)
	catch(err)
	if allSnapshots {
		w := bufio.NewWriter(f)
		_, err = w.WriteString("[\n")
		err = w.Flush()
		catch(err)
	}
}
func saveState(proc *cpu.Processor) {
	var f *os.File
	instant, err := json.Marshal(proc)
	catch(err)
	if allSnapshots {
		f, err = os.OpenFile(saveFile, os.O_APPEND|os.O_WRONLY, 0600)
		instant = append(instant, byte(','))
	} else {
		f, err = os.Create(saveFile)
	}
	catch(err)
	w := bufio.NewWriter(f)
	_, err = w.WriteString(string(instant) + "\n")
	catch(err)
	err = w.Flush()
	catch(err)
}
func endJSON() {
	if allSnapshots {
		f, err := os.OpenFile(saveFile, os.O_APPEND|os.O_WRONLY, 0600)
		catch(err)
		w := bufio.NewWriter(f)
		_, err = w.WriteString("{}]")
		catch(err)
		err = w.Flush()
		catch(err)
	}
}
func catch(err error) {
	if err != nil {
		fmt.Println(err)
		return
	}
}
func main() {
	//init JSON
	initJSON()
	//Create memories
	mainMemory := make([]byte, 131032)
	instructionsMemory := make([]byte, 100000000)

	//Initialize memories
	numberOfInstructions := memory.InitializeInstructionMemory(instructionsMemory)
	memory.InitializeMainMemory(mainMemory)

	//Create clock
	clock := make(chan uint64)

	//Create CPU
	processor := new(cpu.Processor)
	processor.Init(clock, mainMemory, instructionsMemory)

	go processor.Start()

	fmt.Printf("Running for %d cycles\n", numberOfInstructions)
	for i := uint64(0); i < uint64(numberOfInstructions)+5; i++ {
		//fmt.Scanln()
		clock <- i
		//time.Sleep(1 * time.Millisecond)
		if i%2000 == 0 {
			saveState(processor)
			memory.SaveImage(mainMemory, "result.png")
		}

	}
	endJSON()
	//Save image
	memory.SaveImage(mainMemory, "result.png")

}
