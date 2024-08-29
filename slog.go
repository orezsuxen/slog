package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"
)

type progMsg string

type progErrMsg struct{ err error }

type progDoneMsg struct{}

func checkRunning(cmd *exec.Cmd) bool {
	if cmd == nil || cmd.ProcessState != nil && cmd.ProcessState.Exited() || cmd.Process == nil {
		return false
	}
	return true
}

func runner(progName string, progArgs []string, sub chan string) tea.Cmd {
	return func() tea.Msg {
		//setup
		cmd := exec.Command(progName, progArgs...)
		out, err := cmd.StdoutPipe()
		if err != nil {
			return progErrMsg{err}
		}
		cmd.Stdin = os.Stdin
		//execution
		if err := cmd.Start(); err != nil {
			return progErrMsg{err}
		}
		//read prog output
		buf := bufio.NewReader(out)
		for {
			line, _, err := buf.ReadLine() //TODO: should be read bytes
			if err == io.EOF || !checkRunning(cmd) {
				return progDoneMsg{}
			}
			if err != nil {
				return progErrMsg{err}
			}

			sub <- string(line)
		}
	}
}

func waitforProgResponse(sub chan string) tea.Cmd {
	return func() tea.Msg {
		return progMsg(<-sub)
	}
}

//REM: ===== MODEL =====

type model struct {
	done   bool
	sub    chan string
	result string
}

func newModel() model {
	m := model{
		sub: make(chan string),
		// result: new(string),
	}
	return m
}

func (m model) Init() tea.Cmd {
	return tea.Batch(runner("./mirrordir/mirror", []string{}, m.sub), waitforProgResponse(m.sub))
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.done { // ignore input while prog is running
			switch msg.String() {
			case "q", "ctrl+c", "esc":
				return m, tea.Quit
			}
		}
	case progMsg:
		m.result = "Prog is running:\n----------\n"
		m.result += string(msg)
		cmds = append(cmds, waitforProgResponse(m.sub))

	case progErrMsg:
		m.done = true
		m.result = fmt.Sprintf("ERROR:", msg.err.Error()) //TODO: error handling

	case progDoneMsg:
		m.done = true
		m.result = fmt.Sprint("Prog execution Done.\nPress q or ^C or esc to exit.")
	}

	return m, tea.Batch(cmds...)

}

func (m model) View() string {
	return m.result
}

func main() {
	p := tea.NewProgram(newModel())
	if _, err := p.Run(); err != nil {
		fmt.Println("dangit", err)
		os.Exit(1)
	}
}
