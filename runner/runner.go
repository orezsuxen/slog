package runner

import (
	"bufio"
	"io"
	"os"
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"
)

type ProgMsg string

type ProgErrMsg struct{ Err error }

type ProgDoneMsg struct{}

func CheckRunning(cmd *exec.Cmd) bool {
	if cmd == nil || cmd.ProcessState != nil && cmd.ProcessState.Exited() || cmd.Process == nil {
		return false
	}
	return true
}

func WaitforProgResponse(fromProg chan string) tea.Cmd {
	return func() tea.Msg {
		return ProgMsg(<-fromProg)
	}
}

func TerminateProg(toProg chan bool) tea.Cmd {
	return func() tea.Msg {
		toProg <- true
		return true
	}
}

func Run(progName string, progArgs []string, fromProg chan string, toProg chan bool) tea.Cmd {
	return func() tea.Msg {
		//setup
		cmd := exec.Command(progName, progArgs...)
		out, err := cmd.StdoutPipe()
		if err != nil {
			return ProgErrMsg{err}
		}
		cmd.Stdin = os.Stdin
		//execution
		if err := cmd.Start(); err != nil {
			return ProgErrMsg{err}
		}
		//read prog output
		buf := bufio.NewReader(out)
		for {
			select {
			case <-toProg:
				cmd.Process.Kill()
			default:
				line, _, err := buf.ReadLine() //TODO: should be read bytes
				if err == io.EOF || !CheckRunning(cmd) {
					return ProgDoneMsg{}
				}
				if err != nil {
					return ProgErrMsg{err}
				}
				fromProg <- string(line)
			}

		}
	}
}
