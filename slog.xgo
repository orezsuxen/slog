package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	choices  []string
	cursor   int
	selected map[int]struct{}
}

func initiaModel() model {
	return model{
		choices:  []string{"choice one", "choice two", "choice three"},
		selected: make(map[int]struct{}),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "up", "k":
			if m.cursor > 0 {
				m.cursor -= 1
			}

		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor += 1
			}

		case "enter", " ":
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)

			} else {
				m.selected[m.cursor] = struct{}{}
			}
		}
	}

	return m, nil

}

func (m model) View() string {
	s := "this is the header text !!\n"

	for i, choice := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		checked := " "
		if _, ok := m.selected[i]; ok {
			checked = "x"
		}
		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
	}
	s += "\nPress q to quir.\n"

	return s

}

func main() {
	p := tea.NewProgram(initiaModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("we done fucked up: %v", err)
		os.Exit(1)
	}
}
