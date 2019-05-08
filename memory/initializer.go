package memory

import (
	"bufio"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var re = regexp.MustCompile(`.*`)

const MAX_IMM_BITS = 14

func check(e error) {
	if e != nil {
		panic(e)
	}
}
func addInstruction(mem []byte, ins string, i *int) {
	mem[(*i)*4], mem[(*i)*4+1], mem[(*i)*4+2], mem[(*i)*4+3] = instructionToBytes(ins)
	*i++
	for j := 0; j < 3; j++ {
		mem[(*i)*4], mem[(*i)*4+1], mem[(*i)*4+2], mem[(*i)*4+3] = instructionToBytes("NOP \n")
		*i++
	}
}

// Instruction memory
func InitializeInstructionMemory(mem []byte) int {
	fileHandle, err := os.Open("./program.asm")
	defer fileHandle.Close()
	fileScanner := bufio.NewScanner(fileHandle)
	check(err)
	i := 0
	for fileScanner.Scan() {
		instruction := fileScanner.Text()
		if !re.MatchString(instruction) {
			panic("[Assembler] Instruction doesn't match pattern!")
		}
		if strings.Contains(instruction, "#repeat") {
			repeat, err := strconv.Atoi((strings.Split(instruction, " ")[1]))
			check(err)
			var instructionsToProcess []string
			//currentPosition := i
			for fileScanner.Scan() {
				instruction := fileScanner.Text()
				if !strings.Contains(instruction, "#endrepeat") {
					instructionsToProcess = append(instructionsToProcess, instruction)
				} else {
					for j := 0; j < repeat; j++ {
						for _, ins := range instructionsToProcess {
							addInstruction(mem, ins, &i)
						}
					}
					break
				}
			}
		} else {
			addInstruction(mem, instruction, &i)
		}

	}
	return i
}

func instructionToBytes(instruction string) (one byte, two byte, three byte, four byte) {
	//fmt.Println(instruction)
	ins := strings.Split(instruction, " ")

	switch ins[0] {
	case "NOP":
		one, two, three, four = 0, 0, 0, 0
	case "ADD":
		if ins[3][:1] != "V" {
			panic("[Assembler] Add is immediate")
		}
		one = 0x01
		rd, err := strconv.ParseUint(ins[1][1:], 10, 5)
		check(err)
		rl, err := strconv.ParseUint(ins[2][1:], 10, 5)
		check(err)
		rr, err := strconv.ParseUint(ins[3][1:], 10, 5)
		check(err)
		// byte length  = 8  registers_bits = 5
		// RD RD RD RD RD RL RL RL
		two = byte(rd<<(8-5) | (rl >> (5 - (8 - 5))))
		// RL RL RR RR RR RR RR --
		three = byte(rl<<(8-2) | rr<<1)
		//
		four = 0
	case "ADDI":
		if ins[3][:1] != "#" {
			panic("[Assembler] Add not immediate")
		}
		one = 0x11
		rd, err := strconv.ParseUint(ins[1][1:], 10, 5)
		check(err)
		rl, err := strconv.ParseUint(ins[2][1:], 10, 5)
		check(err)
		imm, err := strconv.ParseInt(ins[3][1:], 10, MAX_IMM_BITS)
		check(err)
		// byte length  = 8  registers_bits = 5  IMM = 14
		// RD RD RD RD RD RL RL RL
		two = byte(rd<<(8-5) | (rl >> (5 - (8 - 5))))
		// RL RL IMM IMM IMM IMM IMM IMM
		three = byte(rl<<(8-2) | uint64((imm>>8)&0x3F))
		four = byte(imm)
	case "ADD255":
		if ins[3][:1] != "#" {
			panic("[Assembler] Add not immediate")
		}
		one = 0x12
		rd, err := strconv.ParseUint(ins[1][1:], 10, 5)
		check(err)
		rl, err := strconv.ParseUint(ins[2][1:], 10, 5)
		check(err)
		imm, err := strconv.ParseInt(ins[3][1:], 10, MAX_IMM_BITS)
		check(err)
		// byte length  = 8  registers_bits = 5  IMM = 14
		// RD RD RD RD RD RL RL RL
		two = byte(rd<<(8-5) | (rl >> (5 - (8 - 5))))
		// RL RL IMM IMM IMM IMM IMM IMM
		three = byte(rl<<(8-2) | uint64((imm>>8)&0x3F))
		four = byte(imm)
	case "XOR255":
		if ins[3][:1] != "#" {
			panic("[Assembler] Add not immediate")
		}
		one = 0x13
		rd, err := strconv.ParseUint(ins[1][1:], 10, 5)
		check(err)
		rl, err := strconv.ParseUint(ins[2][1:], 10, 5)
		check(err)
		imm, err := strconv.ParseInt(ins[3][1:], 10, MAX_IMM_BITS)
		check(err)
		// byte length  = 8  registers_bits = 5  IMM = 14
		// RD RD RD RD RD RL RL RL
		two = byte(rd<<(8-5) | (rl >> (5 - (8 - 5))))
		// RL RL IMM IMM IMM IMM IMM IMM
		three = byte(rl<<(8-2) | uint64((imm>>8)&0x3F))
		four = byte(imm)
	case "LOAD":
		rd, err := strconv.ParseUint(ins[1][1:], 10, 5)
		check(err)
		if ins[2][:1] == "#" {
			one = 0x31
			imm, err := strconv.ParseInt(ins[3][1:], 10, 15)
			check(err)
			two = byte(rd<<(8-5) | uint64(imm))
		} else if ins[2][:1] == "V" {
			one = 0x41
			rl, err := strconv.ParseUint(ins[2][1:], 10, 5)
			check(err)
			two = byte(rd<<(8-5) | (rl >> (5 - (8 - 5))))
			// RL RL RR RR RR RR RR --
			three = byte(rl << (8 - 2))
		} else {
			panic("[Assembler] Unknown source on load")
		}
	case "STORE":
		rd, err := strconv.ParseUint(ins[1][1:], 10, 5)
		check(err)
		if ins[2][:1] == "#" {
			one = 0x32
			imm, err := strconv.ParseInt(ins[3][1:], 10, 15)
			check(err)
			two = byte(rd<<(8-5) | uint64(imm))
		} else if ins[2][:1] == "V" {
			one = 0x42
			rl, err := strconv.ParseUint(ins[2][1:], 10, 5)
			check(err)
			two = byte(rd<<(8-5) | (rl >> (5 - (8 - 5))))
			// RL RL RR RR RR RR RR --
			three = byte(rl << (8 - 2))
		} else {
			panic("[Assembler] Unknown source on store")
		}
	case "//":
	case "":

	default:
		fmt.Println("Error on: ", instruction)
		panic("[Assembler] Instruction not supported")

	}

	return
}

// Main memory Image

func InitializeMainMemory(memory []byte) {
	// You can register another format here
	image.RegisterFormat("png", "png", png.Decode, png.DecodeConfig)

	file, err := os.Open("./res/lenna.png")

	if err != nil {
		fmt.Println("Error: File could not be opened")
		os.Exit(1)
	}

	defer file.Close()

	err = getPixels(file, memory)

	if err != nil {
		fmt.Println("Error: Image could not be decoded")
		os.Exit(1)
	}
}

// Get the bi-dimensional pixel array
func getPixels(file io.Reader, mem []byte) error {
	img, _, err := image.Decode(file)

	if err != nil {
		return err
	}

	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	mem[0] = byte(width)
	mem[1] = byte(width >> 8)
	mem[2] = byte(height)
	mem[3] = byte(height >> 8)
	fmt.Println("Image Width:", (int(mem[0]) + int(mem[1])<<8), "Height:", (int(mem[2]) + int(mem[3])<<8))

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			index := ((y * width) + x) + 4
			mem[index] = rgbaToByte(img.At(x, y).RGBA())
		}
	}

	return nil
}

func rgbaToByte(r uint32, g uint32, b uint32, a uint32) byte {
	y := 0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b)
	pixel := color.Gray{uint8(y / 256)}
	return byte(pixel.Y)
}

func SaveImage(mem []byte, out string) {
	width := (int(mem[0]) + int(mem[1])<<8)
	height := (int(mem[2]) + int(mem[3])<<8)
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			index := ((y * width) + x) + 4
			rgb := uint8(mem[index])
			color_in := color.RGBA{rgb, rgb, rgb, 255}
			img.Set(x, y, color_in)
		}
	}

	// Save image
	toImg, err := os.OpenFile(out, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		panic("Cannot open file")
	}

	png.Encode(toImg, img)

}
