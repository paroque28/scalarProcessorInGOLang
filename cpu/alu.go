package cpu

import "fmt"

const (
	ALU_NOP          = iota
	ALU_BUFFER       = iota
	ALU_ADD          = iota
	ALU_AND          = iota
	ALU_OR           = iota
	ALU_ADD255       = iota
	ALU_XOR255       = iota
	ALU_SHUFFLE      = iota
	ALU_UNSHUFFLE    = iota
	ALU_SHUFFLE255   = iota
	ALU_UNSHUFFLE255 = iota
	ALU_ROTATE_LEFT  = iota
	ALU_ROTATE_RIGHT = iota
	ALU_FLIP         = iota
)

func ALU(aluOp byte, a uint64, b uint64) (result uint64) {
	switch aluOp {
	case ALU_BUFFER:
		result = a
	case ALU_NOP:
		result = 0
	case ALU_ADD:
		result = uint64(int(a) + int(b))
	case ALU_AND:
		result = a & b
	case ALU_OR:
		result = a | b
	case ALU_ADD255:
		result = add8Lanes(a, b)
	case ALU_XOR255:
		result = xor8Lanes(a, b)
	case ALU_SHUFFLE:
		result = shuffle(a, b)
	case ALU_UNSHUFFLE:
		result = unshuffle(a, b)
	case ALU_SHUFFLE255:
		result = shuffle8Lanes(a, b)
	case ALU_UNSHUFFLE255:
		result = unshuffle8Lanes(a, b)
	case ALU_ROTATE_LEFT:
		result = rotateLeft(a, b)
	case ALU_ROTATE_RIGHT:
		result = rotateRight(a, b)
	case ALU_FLIP:
		result = flip(a, b)
	default:
		panic("[Exec] ALU operation not implemented")
	}
	return
}

func add8Lanes(a uint64, b uint64) (result uint64) {
	result = 0
	for i := uint(0); i < 64/8; i++ {
		miniA := int8(a >> (8 * i))
		sum := uint8(miniA + int8(b))
		result |= uint64(sum) << (8 * i)
		//fmt.Printf("[ALU] ADD255:  a:%x miniA:%x b:%x sum:%x result:%x\n",a,miniA,b,sum,result)
	}
	return result
}

func xor8Lanes(a uint64, b uint64) (result uint64) {
	result = 0
	for i := uint(0); i < 64/8; i++ {
		miniA := uint8(a >> (8 * i))
		xor := uint8(miniA ^ uint8(b))
		result |= uint64(xor) << (8 * i)
		//fmt.Printf("[ALU] XOR255:  a:%x miniA:%x b:%x sum:%x result:%x\n",a,miniA,b,xor,result)
	}
	return result
}

/*
	From abcdABCD to => aAbBcCdD
*/
func shuffle(a uint64, b uint64) (result uint64) {
	result = 0
	for i := uint(0); i < 32; i++ {
		leftBit := (a >> (63 - i)) & 0x1
		rightBit := (a >> (31 - i)) & 0x1
		twoBits := rightBit | (leftBit << 1)
		//fmt.Printf("[ALU] SHUFFLE:  a:%b left:%b right:%b two:%b result:%b\n", a,leftBit,rightBit,twoBits, result)
		result |= (twoBits) << (i * 2)
	}
	//fmt.Printf("[ALU] SHUFFLE:  a:%x result:%x\n", a, result)
	//fmt.Printf("[ALU] SHUFFLE:  a:%b result:%b\n", a, result)
	return result
}

/*
	From aAbBcCdD to => abcdABCD
*/
func unshuffle(a uint64, b uint64) (result uint64) {
	result = 0
	for i := uint(0); i < 32; i++ {
		twoBits := (a >> (2 * i)) & 0x3
		leftBit := twoBits >> 1
		rightBit := twoBits & 0x1
		//fmt.Printf("[ALU] UNSHUFFLE:  a:%b left:%b right:%b two:%b result:%b\n", a,leftBit,rightBit,twoBits, result)
		result |= (leftBit << (63 - i)) | (rightBit << (31 - i))
	}
	//fmt.Printf("[ALU] UNSHUFFLE:  a:%x result:%x\n", a, result)
	//fmt.Printf("[ALU] UNSHUFFLE:  a:%b result:%b\n", a, result)
	return result
}

/*
	From abcdABCD to => aAbBcCdD
*/
func shuffle8Lanes(a uint64, b uint64) (result uint64) {
	result = 0
	for j := uint(0); j < 64/8; j++ {
		miniA := uint8(a >> (8 * j))
		step := uint8(0)
		for i := uint(0); i < 4; i++ {
			leftBit := (miniA >> (7 - i)) & 0x1
			rightBit := (miniA >> (3 - i)) & 0x1
			twoBits := rightBit | (leftBit << 1)
			step |= (twoBits) << (i * 2)
		}
		result |= uint64(step) << (8 * j)
	}

	//fmt.Printf("[ALU] SHUFFLE255:  a:%x result:%x\n", a, result)
	return result
}

/*
	From aAbBcCdD to => abcdABCD
*/
func unshuffle8Lanes(a uint64, b uint64) (result uint64) {

	result = 0
	for j := uint(0); j < 64/8; j++ {
		miniA := uint8(a >> (8 * j))
		step := uint8(0)
		for i := uint(0); i < 4; i++ {
			twoBits := (miniA >> (2 * i)) & 0x3
			leftBit := twoBits >> 1
			rightBit := twoBits & 0x1
			step |= (leftBit << (7 - i)) | (rightBit << (3 - i))
		}
		result |= uint64(step) << (8 * j)
	}
	//fmt.Printf("[ALU] UNSHUFFLE255:  a:%x result:%x\n", a, result)
	return result
}

func rotateLeft(a uint64, b uint64) (result uint64) {
	result = (a << b) | (a >> (64 - b))
	//fmt.Printf("[ALU] RL:  a:%x b:%x result:%x\n", a, b, result)
	return result
}

func rotateRight(a uint64, b uint64) (result uint64) {
	result = (a >> b) | (a << (64 - b))
	//fmt.Printf("[ALU] RR:  a:%x b:%x result:%x\n", a, b, result)
	return result
}

func flip(a uint64, b uint64) (result uint64) {
	result = 0
	for i := uint(0); i < 64; i++ {
		result |= ((a >> (63 - i)) & 0x1) << i
	}
	fmt.Printf("[ALU] FLIP:  a:%x result:%x\n", a, result)
	return result
}
