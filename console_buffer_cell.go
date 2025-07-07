package quickconsole

func (cell ConsoleBufferCell) WithCharacter(c rune) ConsoleBufferCell {
	return ConsoleBufferCell{
		Character:  c,
		Foreground: cell.Foreground,
		Background: cell.Background,
	}
}

func (cell ConsoleBufferCell) WithForeGround(color int) ConsoleBufferCell {
	return ConsoleBufferCell{
		Character:  cell.Character,
		Foreground: color,
		Background: cell.Background,
	}
}

func (cell ConsoleBufferCell) WithBackground(color int) ConsoleBufferCell {
	return ConsoleBufferCell{
		Character:  cell.Character,
		Foreground: cell.Foreground,
		Background: color,
	}
}

func (cell ConsoleBufferCell) OverrideDefaults(foreground, background int) ConsoleBufferCell {
	newForeground := cell.Foreground
	newBackground := cell.Background
	if newForeground == AnsiColorDefault {
		newForeground = foreground
	}
	if newBackground == AnsiColorDefault {
		newBackground = background
	}
	return ConsoleBufferCell{
		Character:  cell.Character,
		Foreground: newForeground,
		Background: newBackground,
	}
}

func (cell ConsoleBufferCell) Equals(otherCell ConsoleBufferCell) bool {
	return cell.Character == otherCell.Character && cell.Foreground == otherCell.Foreground && cell.Background == otherCell.Background
}
