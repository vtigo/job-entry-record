package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m model) View() string {
	return "'Ctrl + c' to quit"
}

func main() {
	state := NewState("data")
	state.LoadEntries("test")
	// state.ListEntries()

	model := model{}
	p := tea.NewProgram(model)
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
