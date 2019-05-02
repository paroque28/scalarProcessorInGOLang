package cpu

type Execute struct {
	InDecodeControlSignals DecodeControlSignals `json:"in_decode_control_signals"`
	InRd1                  uint64               `json:"rd_1"`
	InRd2                  uint64               `json:"rd_2"`
	Immediate              int64                `json:"immediate"`
}

const (
	ALU_NOP       = iota
	ALU_ADD       = iota
	ALU_XOR       = iota
	ALU_TOTAL_OPS = iota
)

func (exec *Execute) Run(done chan string) {

	done <- "execute"
}

func (exec *Execute) UpdateInRegisters(inControlSignals DecodeControlSignals, rd1 uint64, rd2 uint64, immediate int64) {
	exec.InDecodeControlSignals = inControlSignals
	exec.InRd1 = rd1
	exec.InRd2 = rd2
	exec.Immediate = immediate
}
