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

const MAX_UNSIGNED_IMM = 16384

func check(e error) {
	if e != nil {
		panic(e)
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
			panic("Instruction doesn't match pattern!")
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
							mem[i*4], mem[i*4+1], mem[i*4+2], mem[i*4+3] = instructionToBytes(ins)
							i++
						}
					}
					break
				}
			}
		} else {
			mem[i*4], mem[i*4+1], mem[i*4+2], mem[i*4+3] = instructionToBytes(instruction)
			i++
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
		one = 0x01
		rd, err := strconv.Atoi(ins[1][1:])
		if rd > 31 || rd < 0 {
			panic("RD higher than 32")
		}
		check(err)
		rl, err := strconv.Atoi(ins[2][1:])
		if rd > 31 || rd < 0 {
			panic("RL higher than 32")
		}
		check(err)
		rr, err := strconv.Atoi(ins[3][1:])
		if rd > 31 || rd < 0 {
			panic("RR higher than 32")
		}
		check(err)
		// byte length  = 8  registers_bits = 5
		// RD RD RD RD RD RL RL RL
		two = byte(rd<<(8-5) | (rl >> (5 - (8 - 5))))
		// RL RL RR RR RR RR RR --
		three = byte(rl<<(8-2) | rr<<1)
		//
		four = 0
	case "ADDI":
		one = 0x11
		rd, err := strconv.Atoi(ins[1][1:])
		if rd > 31 || rd < 0 {
			panic("RD higher than 32")
		}
		check(err)
		rl, err := strconv.Atoi(ins[2][1:])
		if rd > 31 || rd < 0 {
			panic("RL higher than 32")
		}
		check(err)
		imm, err := strconv.Atoi(ins[3][1:])
		if rd > MAX_UNSIGNED_IMM/2 || rd < -MAX_UNSIGNED_IMM/2 {
			panic("Imm higher than MAX")
		}
		check(err)
		// byte length  = 8  registers_bits = 5  IMM = 14
		// RD RD RD RD RD RL RL RL
		two = byte(rd<<(8-5) | (rl >> (5 - (8 - 5))))
		// RL RL IMM IMM IMM IMM IMM IMM
		three = byte(rl<<(8-2) | ((imm >> 8) & 0x3F))
		four = byte(imm)
	default:
		fmt.Println("Error on: ", instruction)
		panic("Instruction not supported")

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
