package main

import (
	"fmt"
	"log"
	"strings"
	"time"
	tea "github.com/charmbracelet/bubbletea"
)

// TODO: Implement state machine to manage menu, creating and other states.

type Model struct {
	state     *State
	inputs    []string
	labels    []string
	cursor    int
	finished  bool
	err       error
}

func initialModel() Model {
	return Model{
		state:  NewState("data", "csv"),
		inputs: make([]string, 5),
		labels: []string{
			"Company:",
			"Role:",
			"Status:",
			"Platform:",
			"Contact Replied (y/n):",
		},
		cursor: 0,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			if m.cursor < len(m.inputs)-1 {
				m.cursor++
			} else {
				m.finished = true
				return m, tea.Quit
			}
		case "up":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down":
			if m.cursor < len(m.inputs)-1 {
				m.cursor++
			}
		case "backspace":
			if len(m.inputs[m.cursor]) > 0 {
				m.inputs[m.cursor] = m.inputs[m.cursor][:len(m.inputs[m.cursor])-1]
			}
		default:
			if len(msg.String()) == 1 {
				m.inputs[m.cursor] += msg.String()
			}
		}
	}
	return m, nil
}

func (m Model) View() string {
	if m.finished {
		return "Entry created successfully!\n"
	}

	var s strings.Builder
	s.WriteString("Create New Job Entry\n")
	s.WriteString("====================\n\n")

	for i, label := range m.labels {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}
		s.WriteString(fmt.Sprintf("%s %s %s\n", cursor, label, m.inputs[i]))
	}

	s.WriteString("\nPress Enter to move to next field, Ctrl+C to quit\n")
	return s.String()
}

func main() {
	p := tea.NewProgram(initialModel())
	finalModel, err := p.Run()
	if err != nil {
		log.Fatalf("Error: Could not run tea model: %v", err)
	}

	m := finalModel.(Model)
	err = m.state.LoadEntries("test")
	if err != nil {
		log.Fatalf("Error: Could not load existing entries: %v", err)
	}

	if m.finished && len(m.inputs[0]) > 0 {
		contactReplied := strings.ToLower(m.inputs[4]) == "y" || strings.ToLower(m.inputs[4]) == "yes"
		
		entry := NewEntry(
			m.inputs[0],
			m.inputs[1],
			m.inputs[2],
			m.inputs[3],
			time.Now(),
			contactReplied,
		)
		
		entry.AddToState(m.state)
		err = m.state.SaveEntries("test")
		if err != nil {
			log.Fatalf("Error: Could not save entries after editing: %v", err)
		}
	}
}

