package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/charmbracelet/bubbles/spinner"
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

func runner(progName string, progArgs []string, fromProg chan string, toProg chan bool) tea.Cmd {
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
			select {
			case <-toProg:
				cmd.Process.Kill()
				// return progDoneMsg{} //TEST:
			default:
				line, _, err := buf.ReadLine() //TODO: should be read bytes
				if err == io.EOF || !checkRunning(cmd) {
					return progDoneMsg{}
				}
				if err != nil {
					return progErrMsg{err}
				}
				fromProg <- string(line)
			}

		}
	}
}

func waitforProgResponse(fromProg chan string) tea.Cmd {
	return func() tea.Msg {
		return progMsg(<-fromProg)
	}
}

func terminateProg(toProg chan bool) tea.Cmd {
	return func() tea.Msg {
		toProg <- true
		return true
	}
}

//REM: ===== MODEL =====

type model struct {
	done     bool
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
		runner("./counterdir/counter", []string{"1000", "10"}, m.fromProg, m.toProg),
		waitforProgResponse(m.fromProg),
		m.spin.Tick,
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	var cmds []tea.Cmd
	var cmd tea.Cmd

	if !m.done { // while prog is running
		m.result = "Prog is running: press ^C to kill\n"
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
				cmds = append(cmds, terminateProg(m.toProg))

			}
		}
	case progMsg:
		m.progResult = string(msg)
		cmds = append(cmds, waitforProgResponse(m.fromProg))

	case progErrMsg:
		m.done = true
		m.result = fmt.Sprintf("ERROR slog:", msg.err.Error()) //TODO: error handling

	case progDoneMsg:
		m.done = true
		m.result = fmt.Sprint("Prog execution Done.\nPress q or ^C or esc to exit.")
	default:
		m.spin, cmd = m.spin.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)

}

func (m model) View() string {
	if !m.done {
		m.result += m.spin.View()
		m.result += "\n"
		m.result += m.progResult
	}
	return m.result
}

func main() {
	p := tea.NewProgram(newModel())
	if _, err := p.Run(); err != nil {
		fmt.Println("dangit", err)
		os.Exit(1)
	}
}
