package main

// A simple program demonstrating the textarea component from the Bubbles
// component library.

import (
	tea "github.com/charmbracelet/bubbletea"
	"log"
	"main/interactive"
)

func main() {
	p := tea.NewProgram(interactive.InitialModel())

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
