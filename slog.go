package main

import (
	"fmt"
	"os"
	"strings"

	"slog/runner"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

func progName() string {
	if len(os.Args) <= 1 {
		return ""
	} else {
		return os.Args[1]
	}
}

func progArgs() []string {
	ret := make([]string, 0)
	if len(os.Args) <= 1 {
		return nil
	} else {
		for i, a := range os.Args {
			if i <= 1 {
				continue
			} else {
				ret = append(ret, a)

			}
		}
		return ret
	}
}

// REM: model
type model struct {
	done bool

	fromProg chan string
	toProg   chan bool

	result     string
	progResult string

	spin spinner.Model
}

func newModel() model {
	s := spinner.New()
	s.Spinner = spinner.Line
	m := model{
		fromProg: make(chan string),
		toProg:   make(chan bool),
		spin:     s,
		// result: new(string),
	}
	return m
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		runner.Run(progName(), progArgs(), m.fromProg, m.toProg),
		runner.WaitforProgResponse(m.fromProg),
		m.spin.Tick,
	)
}

// REM: update
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	var cmds []tea.Cmd
	var cmd tea.Cmd

	if !m.done { // while prog is running
		m.result = "Prog is running: press ^C to kill\n"
		m.result += "args are\n:"
		m.result += strings.Join(progArgs(), " ")
		m.result += "\n--------------------------\n"
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.done { // ignore input while prog is running
			switch msg.String() {
			case "q", "ctrl+c", "esc":
				return m, tea.Quit
			}
		} else {
			switch msg.String() {
			case "ctrl+c":
				cmds = append(cmds, runner.TerminateProg(m.toProg))

			}
		}
	case runner.ProgMsg:
		m.progResult = string(msg)
		cmds = append(cmds, runner.WaitforProgResponse(m.fromProg))

	case runner.ProgErrMsg:
		m.done = true
		m.result = fmt.Sprintf("ERROR slog:", msg.Err.Error()) //TODO: error handling

	case runner.ProgDoneMsg:
		m.done = true
		m.result = fmt.Sprint("Prog execution Done.\nPress q or ^C or esc to exit.")
	default:
		m.spin, cmd = m.spin.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)

}

// REM: view
func (m model) View() string {
	if !m.done {
		m.result += m.spin.View()
	}
	m.result += "\n"
	m.result += m.progResult
	return m.result
}

// REM: main
func main() {
	p := tea.NewProgram(newModel())
	if _, err := p.Run(); err != nil {
		fmt.Println("dangit", err)
		os.Exit(1)
	}
}
