package quickconsole

import (
	"errors"
	"fmt"
	"io"
	"slices"
	"strings"
	"unicode"
)

func (buffer *ConsoleBuffer) WriteBuffer(writer io.Writer) {
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

func (buffer *ConsoleBuffer) IsPointOutOfBounds(x, y int) bool {
	return x < 0 || x >= buffer.Width || y < 0 || y >= buffer.Height
}

func (buffer *ConsoleBuffer) IsRectOutOfBounds(x, y, width, height int) bool {
	return x < 0 || x >= buffer.Width || y < 0 || y >= buffer.Height ||
		(x+width <= 0) || (y+height <= 0)
}

func (buffer *ConsoleBuffer) IsRectFullyInBounds(x, y int, width int, height int) bool {
	return x >= 0 && x+width <= buffer.Width && y >= 0 && y+height <= buffer.Height
}

func (buffer *ConsoleBuffer) DrawTextAtPoint(x, y int, text string) {
	buffer.DrawTextAtPointWithColors(x, y, text, buffer.CurrentForegroundColor, buffer.CurrentBackgroundColor)
}

func (buffer *ConsoleBuffer) DrawTextAtPointWithColor(x, y int, text string, color int) {
	buffer.DrawTextAtPointWithColors(x, y, text, color, buffer.CurrentBackgroundColor)
}

func (buffer *ConsoleBuffer) DrawTextAtPointWithColors(x, y int, text string, color int, background int) {
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

func (buffer *ConsoleBuffer) DrawCell(x, y int, cell ConsoleBufferCell) {
	if buffer.IsPointOutOfBounds(x, y) {
		return
	}
	buffer.Cells[x+y*buffer.Width] = cell
}

func (buffer *ConsoleBuffer) DrawRectangle(x, y, width, height int, cell ConsoleBufferCell) {
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

func (buffer *ConsoleBuffer) DrawBoxComplex(x, y, width, height int, cellSides, cellTopBottom, cellCorner ConsoleBufferCell) {
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

func (buffer *ConsoleBuffer) DrawLine(x, y, length, direction int, cell ConsoleBufferCell) {
	if length <= 0 {
		return
	}
	boundsHeight := length
	boundsWidth := length
	if direction == LineDirectionHorizontal {
		boundsHeight = 0
	} else {
		boundsWidth = 0
	}
	if buffer.IsRectOutOfBounds(x, y, boundsWidth, boundsHeight) {
		return
	}
	inc := 1
	if direction == LineDirectionVertical {
		inc = buffer.Width
	}
	idx := x + y*buffer.Width
	maxLength := buffer.Width - x
	if direction == LineDirectionVertical {
		maxLength = buffer.Height - y
	}
	maxLength = min(length, maxLength)
	for i := 0; i < maxLength; i++ {
		if idx >= 0 {
			buffer.Cells[idx] = cell
		}
		idx += inc
	}
}

func (buffer *ConsoleBuffer) GetCellAt(x, y int) (ConsoleBufferCell, error) {
	if buffer.IsPointOutOfBounds(x, y) {
		return GetConsoleBufferCellZero(), errors.New("X and/or Y out of bounds of console buffer")
	}
	return buffer.Cells[x+y*buffer.Width], nil
}

func (buffer *ConsoleBuffer) GetStringAt(x, y, length int) string {
	if length <= 0 {
		return ""
	}
	if buffer.IsRectOutOfBounds(x, y, length, 0) {
		return ""
	}
	idx := max(0, x) + y*buffer.Width
	end := idx + min(buffer.Width-x+length+min(0, x))

	var runes []rune
	for _, cell := range buffer.Cells[idx:end] {
		runes = append(runes, cell.Character)
	}
	return strings.TrimSpace(string(runes))
}

func (buffer *ConsoleBuffer) DrawBuffer(x, y int, targetBuffer *ConsoleBuffer) {
	if buffer.IsRectOutOfBounds(x, y, targetBuffer.Width, targetBuffer.Height) {
		return
	}

	rowIdx := y * buffer.Width
	targetBufferLength := len(targetBuffer.Cells)
	for bufferIdx := 0; bufferIdx < targetBufferLength; bufferIdx++ {
		bufferX := bufferIdx % targetBuffer.Width
		if bufferX == 0 && bufferIdx > 0 {
			rowIdx += buffer.Width
		}
		if rowIdx < 0 || rowIdx >= len(buffer.Cells) {
			continue
		}
		if x+bufferX < 0 || x+bufferX >= buffer.Width {
			continue
		}
		sourceCell := targetBuffer.Cells[bufferIdx]
		if sourceCell.Character > '\000' {
			buffer.Cells[rowIdx+x+buffer.Width] = sourceCell
		}
	}
}

func (buffer *ConsoleBuffer) Scroll(xd, yd int) {
	xd = xd % buffer.Width
	yd = yd % buffer.Height
	cellLength := len(buffer.Cells)
	if xd != 0 {
		axd := absInt(xd)
		count := buffer.Height * axd
		tempBuffer := make([]ConsoleBufferCell, count)
		cellOffset := 0
		cellOffset2 := axd
		cellOffset3 := 0

		if xd < 0 {
			cellOffset = buffer.Width - axd
			cellOffset2 = 0
			cellOffset3 = axd
		}
		for i, x := 0, 0; i < cellLength; i, x = i+buffer.Width, x+axd {
			copy(tempBuffer, buffer.Cells[cellOffset:cellOffset+axd])
		}
		copy(buffer.Cells[cellOffset3:cellOffset3+cellLength-axd], buffer.Cells[cellOffset2:cellOffset2+cellLength-axd])
		for i, x := 0, 0; i < cellLength; i, x = i+buffer.Width, x+axd {
			copy(buffer.Cells[cellOffset:cellOffset+axd], tempBuffer)
		}
	}

	if yd != 0 {
		count := absInt(buffer.Width * yd)
		tempBuffer := make([]ConsoleBufferCell, count)
		cellOffset := 0
		cellOffset2 := count
		cellOffset3 := 0
		if yd < 0 {
			cellOffset = cellLength - count
			cellOffset2 = 0
			cellOffset3 = count - yd
		}
		copy(tempBuffer, buffer.Cells[cellOffset:cellOffset+count])
		copy(buffer.Cells[cellOffset2:cellOffset2+cellLength-count], buffer.Cells[cellOffset3:cellOffset3+cellLength-count])
		copy(buffer.Cells[cellOffset:cellOffset+count], tempBuffer)
	}
}

func (buffer *ConsoleBuffer) Flip(horizontal, vertical bool) {
	tempBuffer := make([]ConsoleBufferCell, buffer.Width)
	cellLength := len(buffer.Cells)
	for i, j := 0, cellLength-buffer.Width; i < cellLength; i, j = i+buffer.Width, j-buffer.Width {
		if vertical && i < j {
			copy(tempBuffer, buffer.Cells[i:i+buffer.Width])
			copy(buffer.Cells[i:i+buffer.Width], buffer.Cells[j:j+buffer.Width])
			copy(buffer.Cells[j:j+buffer.Width], tempBuffer)
		}

		if horizontal {
			slices.Reverse(buffer.Cells[i : i+buffer.Width])
		}
	}
}

func (buffer *ConsoleBuffer) Rotate(x, y, width int, clockwise bool) {
	if width <= 0 {
		return
	}
	if !buffer.IsRectOutOfBounds(x, y, width, width) {
		return
	}
	rotateCopy, _ := buffer.Copy(x, y, width, width)
	sourceOffset := 0
	for i := 0; i < width; i++ {
		for j := 0; j < width; j++ {
			var offset int
			if clockwise {
				offset = width - i - 1 - x + (j+y)*buffer.Width
			} else {
				offset = x + i + (width-j-1+y)*buffer.Width
			}
			buffer.Cells[offset] = rotateCopy.Cells[sourceOffset]
			sourceOffset++
		}
	}
}

func (buffer *ConsoleBuffer) Copy(x, y, width, height int) (*ConsoleBuffer, error) {
	if width <= 0 || height <= 0 {
		return nil, errors.New("width and height must be greater than 0")
	}
	if !buffer.IsRectFullyInBounds(x, y, width, height) {
		return nil, errors.New("area represented by x, y, width and height are not fully in-bounds")
	}

	actualWidth := min(width, buffer.Width-x)
	actualHeight := min(height, buffer.Height-y)

	bufferCopy := NewConsoleBuffer(width, height)
	sourceIdx := x + y*buffer.Width
	copyIdx := 0
	for i := 0; i < actualHeight; i++ {
		copy(bufferCopy.Cells[copyIdx:copyIdx+actualWidth], buffer.Cells[sourceIdx:sourceIdx+actualWidth])
	}
	return bufferCopy, nil
}
