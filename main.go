package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	p := tea.NewProgram(
		initialModel(),
		tea.WithAltScreen(),       // use full terminal screen
		tea.WithMouseCellMotion(), // optional: enable mouse scroll
	)

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running bd-tui: %v\n", err)
		os.Exit(1)
	}
}
