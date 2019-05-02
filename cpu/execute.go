package cpu

type Execute struct {
	InDecodeControlSignals DecodeControlSignals `json:"in_decode_control_signals"`
	Rd1                    uint32               `json:"rd_1"`
	Rd2                    uint32               `json:"rd_2"`
}

func (exec *Execute) Run(done chan string) {

	done <- "execute"
}

func (exec *Execute) UpdateInRegisters(inControlSignals DecodeControlSignals) {
	exec.InDecodeControlSignals = inControlSignals
}
