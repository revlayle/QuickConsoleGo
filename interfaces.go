package quickconsole

import "io"

type SystemConsoler interface {
	GetCursorVisible() bool
	SetCursorVisible(visible bool)
	IsKeyAvailable() bool
	ReadKey() byte
	Out() io.Writer
	SetCursorPosition(left int, top int)
}

type ConsoleBufferer interface {
	WriteBuffer(writer io.Writer)
	IsPointOutOfBounds(x, y int) bool
	IsRectOutOfBounds(x, y, width, height int) bool
	IsRectFullyInBounds(x, y, width, height int) bool
	DrawTextAtPoint(x, y int, text string)
	DrawTextAtPointWithColor(x, y int, text string, color int)
	DrawTextAtPointWithColors(x, y int, text string, color int, background int)
	DrawCell(x, y int, cell ConsoleBufferCell)
	DrawRectangle(x, y, width, height int, cell ConsoleBufferCell)
	DrawBox(x, y, width, height int, cell ConsoleBufferCell)
	DrawBoxComplex(x, y, width, height int, cellSides, cellTopBottom, cellCorner ConsoleBufferCell)
	DrawLine(x, y, length int, direction int, cell ConsoleBufferCell)
	GetCellAt(x, y int) (ConsoleBufferCell, error)
	GetStringAt(x, y, length int) string
	DrawBuffer(x, y int, buffer *ConsoleBuffer)
	Scroll(xd, yd int)
	Flip(horizontal, vertical bool)
	Rotate(x, y, width int, clockwise bool)
	Copy(x, y, width, height int) (*ConsoleBuffer, error)
}

type ConsoleBufferCeller interface {
	WithCharacter(c rune) ConsoleBufferCell
	WithForeground(color int) ConsoleBufferCell
	WithBackground(color int) ConsoleBufferCell
	OverrideDefaults(foreground, background int) ConsoleBufferCell
	Equals(cell ConsoleBufferCell) bool
}
