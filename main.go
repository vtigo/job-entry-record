package main

import (
	"fmt"
	"log"

	tea "github.com/charmbracelet/bubbletea"
)

type appState int
const (
	appStateMenu appState = iota
	appStateList
	appStateCreate
)

type model struct {
	entries  []Entry
	cursor 	 int
	appState appState
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
		case "up":
			if m.cursor > 0 {
				m.cursor --
			}
		case "down":
			if m.appState == appStateMenu {
				if m.cursor < 1 {
					m.cursor ++
				}
			}
		}
	}

	return m, nil
}

func (m model) View() string {
	// The header
	s := "J A R\n\n"


	if m.appState == appStateMenu {
		choices := []string{"List entries.", "Create entry."}
		for i, choice := range choices {
			cursor := " "
			if m.cursor == i {
				cursor = ">"
			}

			s += fmt.Sprintf("%s %s\n", cursor, choice)
		}
	}

	// The footer
	s += "\nPress 'Ctrl + c' to quit.\n"

	// Send the UI for rendering
	return s
}

func main() {
	storage := NewCsvStorage("data", "test")
	entries, err := storage.GetAll()
	if err != nil {
		log.Fatalf("Failed to get records from storage: %v", err)
	}

	model := model{entries: entries, appState: appStateMenu }
	p := tea.NewProgram(model)
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
