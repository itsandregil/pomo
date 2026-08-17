package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
)

func main() {
	p := tea.NewProgram(NewModel())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "running program: %s", err)
		os.Exit(1)
	}
}
