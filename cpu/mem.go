package cpu

type Memory struct {
	rd1 uint32 `json:"rd_1"`
	rd2 uint32 `json:"rd_2"`
}

func (mem *Memory) Run() {
}
