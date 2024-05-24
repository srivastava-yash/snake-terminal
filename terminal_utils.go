package main

import (
	"bufio"
	"fmt"
	"os"

	"golang.org/x/term"
)

var screen = bufio.NewWriter(os.Stdout)

func hideCursor() {
	fmt.Fprint(screen, "\033[?25l")
}

func showCursor() {
	fmt.Fprint(screen, "\033[?25h")
}

func moveCursor(pos [2]int) {
	fmt.Fprintf(screen, "\033[%d;%dH", pos[1], pos[0])
}

func clear() {
	fmt.Fprint(screen, "\033[2J")
}

func print_string(str string) {
	fmt.Fprint(screen, str)
}

func display() {
	screen.Flush()
}

func get_terminal_window_size() (int, int) {
	width, height, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		panic(err)
	}

	return width, height
}
