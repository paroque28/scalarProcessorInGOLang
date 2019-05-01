package cpu

type Execute struct {
	rd1 uint32 `json:"rd_1"`
	rd2 uint32 `json:"rd_2"`
}

func (exec *Execute) Run(done chan string) {

	done <- "execute"
}
