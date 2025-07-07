package quickconsole

type ConsoleBufferCell struct {
	Character  rune
	Foreground int
	Background int
}

type ConsoleBuffer struct {
	Cells                  []ConsoleBufferCell
	Width                  int
	Height                 int
	CurrentForegroundColor int
	CurrentBackgroundColor int
}
