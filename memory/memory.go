package memory

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

// Instruction memory
func InitializeInstructionMemory( mem []byte ){
	dat, err := ioutil.ReadFile("./program.asm")
	check(err)
	instructions := strings.Split(string(dat), ";")
	for i := 0; i< len(instructions); i++  {
		fmt.Println(instructions[i])
		mem[i*4], mem[i*4+1], mem[i*4+2], mem[i*4+3] = instructioToBytes(instructions[i])
	}
}


func instructioToBytes( instruction string ) ( one byte, two byte, three byte, four byte ){
	one = 1
	two = 0
	three = 1
	four = 1
	return
}

// Main memory Image

func InitializeMainMemory( memory []byte){
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
func getPixels(file io.Reader, mem []byte) (error) {
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
	fmt.Println("Image Width:",(int(mem[0]) + int(mem[1]) << 8),"Height:", (int(mem[2]) + int(mem[3]) << 8))

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

func SaveImage(mem []byte, out string)  {
	width := (int(mem[0]) + int(mem[1]) << 8)
	height := (int(mem[2]) + int(mem[3]) << 8)
	img := image.NewRGBA(image.Rect(0, 0, width, height))


	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			index := ((y * width) + x) + 4
			rgb:= uint8(mem[index])
			color_in := color.RGBA{rgb,rgb,rgb,255}
			img.Set(x,y,color_in)
		}
	}

	// Save image
	toImg, err := os.OpenFile(out, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		panic("Cannot open file")
	}

	png.Encode(toImg, img)

}