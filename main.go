package main

var consoleBufferCellDefault = ConsoleBufferCell{
	Character:  0,
	Foreground: AnsiColorDefault,
	Background: AnsiColorDefault,
}

func GetConsoleBufferCellZero() ConsoleBufferCell {
	return consoleBufferCellDefault
}
