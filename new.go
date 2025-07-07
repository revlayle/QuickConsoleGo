package quickconsole

func NewConsoleBuffer(height, width int) *ConsoleBuffer {
	return &ConsoleBuffer{
		Cells:  make([]ConsoleBufferCell, height*width),
		Width:  width,
		Height: height,
	}
}
