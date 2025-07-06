package main

import (
	"fmt"
	"io"
	"strings"
	"unicode"
)

func (buffer ConsoleBuffer) WriteBuffer(writer io.Writer) {
	currentForeground := buffer.CurrentForegroundColor
	currentBackground := buffer.CurrentBackgroundColor
	if currentForeground == AnsiColorDefault {
		currentForeground = AnsiColorBlack
	}
	if currentBackground == AnsiColorDefault {
		currentBackground = AnsiColorBlack
	}

	var builder strings.Builder
	prevForegroundColor := -1
	prevBackgroundColor := -1
	for i, cell := range buffer.Cells {
		if i > 0 && i%buffer.Width == 0 {
			builder.WriteRune('\n')
		}
		foreground := cell.Foreground
		if foreground == AnsiColorDefault {
			foreground = currentForeground
		}
		background := cell.Background
		if background == AnsiColorDefault {
			background = currentBackground
		}

		if foreground != prevForegroundColor {
			prevForegroundColor = foreground
			builder.WriteString(fmt.Sprintf("\x1b[%dm", 30+prevForegroundColor))
		}
		if background != prevBackgroundColor {
			prevBackgroundColor = background
			builder.WriteString(fmt.Sprintf("\x1b[%dm", 40+prevBackgroundColor))
		}

		if unicode.IsControl(cell.Character) {
			builder.WriteRune(' ')
		} else {
			builder.WriteRune(cell.Character)
		}
	}

	_, err := writer.Write([]byte(builder.String()))
	if err != nil {
		panic(err)
	}
}

func (buffer ConsoleBuffer) IsPointOutOfBounds(x, y int) bool {
	return x < 0 || x >= buffer.Width || y < 0 || y >= buffer.Height
}

func (buffer ConsoleBuffer) IsRectOutOfBounds(x, y, width, height int) bool {
	return x < 0 || x >= buffer.Width || y < 0 || y >= buffer.Height ||
		(x+buffer.Width <= 0) || (y+buffer.Height <= 0)
}

func (buffer ConsoleBuffer) IsRectFullyInBounds(x, y int, width int, height int) bool {
	return x >= 0 && x+width <= buffer.Width && y >= 0 && y+height <= buffer.Height
}

func (buffer ConsoleBuffer) DrawTextAtPoint(x, y int, text string) {
	buffer.DrawTextAtPointWithColors(x, y, text, buffer.CurrentForegroundColor, buffer.CurrentBackgroundColor)
}

func (buffer ConsoleBuffer) DrawTextAtPointWithColor(x, y int, text string, color int) {
	buffer.DrawTextAtPointWithColors(x, y, text, color, buffer.CurrentBackgroundColor)
}

func (buffer ConsoleBuffer) DrawTextAtPointWithColors(x, y int, text string, color int, background int) {
	if buffer.IsRectOutOfBounds(x, y, len(text), 0) {
		return
	}
	textArray := []rune(text)
	idx := x + y*buffer.Width
	maxLength := min(len(textArray), buffer.Width-x)
	for i := 0; i < maxLength; i++ {
		if x+i < 0 {
			continue
		}
		buffer.Cells[idx+i] = ConsoleBufferCell{
			Character:  textArray[i],
			Foreground: color,
			Background: background,
		}
	}
}

func (buffer ConsoleBuffer) DrawCell(x, y int, cell ConsoleBufferCell) {
	if buffer.IsPointOutOfBounds(x, y) {
		return
	}
	buffer.Cells[x+y*buffer.Width] = cell
}

func (buffer ConsoleBuffer) DrawRectangle(x, y, width, height int, cell ConsoleBufferCell) {
	if width <= 0 || height <= 0 {
		return
	}
	if buffer.IsRectOutOfBounds(x, y, width, height) {
		return
	}
	cell.OverrideDefaults(buffer.CurrentBackgroundColor, buffer.CurrentForegroundColor)

	i := 0
	for rowIdx := y * buffer.Width; i < buffer.Height; rowIdx, i = rowIdx+buffer.Width, i+1 {
		if rowIdx < 0 || rowIdx > len(buffer.Cells) {
			continue
		}
		for j := 0; j < width; j++ {
			if x+j < 0 || x+j >= buffer.Width {
				continue
			}
			idx := j + x + rowIdx
			buffer.Cells[idx] = cell
		}
	}
}

func (buffer ConsoleBuffer) DrawBoxComplex(x, y, width, height int, cellSides, cellTopBottom, cellCorner ConsoleBufferCell) {
	if width <= 0 || height <= 0 {
		return
	}
	if buffer.IsRectOutOfBounds(x, y, width, height) {
		return
	}

	cellSides = cellSides.OverrideDefaults(buffer.CurrentBackgroundColor, buffer.CurrentForegroundColor)
	cellTopBottom = cellTopBottom.OverrideDefaults(buffer.CurrentBackgroundColor, buffer.CurrentForegroundColor)
	cellCorner = cellCorner.OverrideDefaults(buffer.CurrentBackgroundColor, buffer.CurrentForegroundColor)

	i := 0
	for rowIdx := y * buffer.Width; i < buffer.Height; rowIdx, i = rowIdx+buffer.Width, i+1 {
		if rowIdx < 0 || rowIdx > len(buffer.Cells) {
			continue
		}
		for j := 0; j < width; j++ {
			if x+j < 0 || x+j >= buffer.Width {
				continue
			}
			cell := GetConsoleBufferCellZero()
			if i == 0 && j == 0 {
				cell = cellCorner
			} else if i == 0 && j == width-1 {
				cell = cellCorner
			} else if i == height-1 && j == width-1 {
				cell = cellCorner
			} else if i == height-1 && j == 0 {
				cell = cellCorner
			} else if i == 0 {
				cell = cellTopBottom
			} else if j == 0 {
				cell = cellSides
			} else if i == height-1 {
				cell = cellTopBottom
			} else if j == width-1 {
				cell = cellSides
			}

			if cell.Equals(GetConsoleBufferCellZero()) {
				continue
			}
			idx := j + x + rowIdx
			buffer.Cells[idx] = cell
		}
	}
}
