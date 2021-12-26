package utils

type DataBus struct {
	GPUAllocate chan string
	GPURelease  chan string
}

func NewDataBus() *DataBus {
	return &DataBus{
		GPUAllocate: make(chan string),
		GPURelease:  make(chan string),
	}
}
