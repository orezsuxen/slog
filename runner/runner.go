package runner

import (
	"bufio"
	"io"
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

func readInput(cmd *exec.Cmd, out io.ReadCloser, outChan chan string, msgChan chan tea.Msg) {
	buf := bufio.NewReader(out)
	for {
		line, _, err := buf.ReadLine()
		if err == io.EOF || !CheckRunning(cmd) {
			msgChan <- ProgDoneMsg{}
		}
		if err != nil {
			msgChan <- ProgErrMsg{err}
		}
		outChan <- string(line)
	}
}

func Run(
	progName string,
	progArgs []string,
	fromProg chan string,
	toProg chan bool,
	progIn io.Reader,
) tea.Cmd {
	return func() tea.Msg {
		//setup
		cmd := exec.Command(progName, progArgs...)
		out, err := cmd.StdoutPipe()
		if err != nil {
			return ProgErrMsg{err}
		}
		cmd.Stdin = progIn
		//execution
		if err := cmd.Start(); err != nil {
			return ProgErrMsg{err}
		}
		//read prog output
		msgChan := make(chan tea.Msg)
		outputChan := make(chan string)
		go readInput(cmd, out, outputChan, msgChan)
		//send msgs back
		for {
			select {
			case <-toProg:
				cmd.Process.Kill()
			case outString := <-outputChan:
				fromProg <- outString
			case retMsg := <-msgChan:
				return retMsg
			}
		}
	}
}
