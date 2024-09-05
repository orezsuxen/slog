package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"slog/collect"
	"slog/help"
	"slog/pargs"
	"slog/runner"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

// REM: model
type model struct {
	done bool

	fromProg chan string
	toProg   chan bool

	result      string
	progResult  string
	progInRead  io.Reader
	progInWrite io.Writer
	//TODO: something like current state?

	messages collect.Collector

	spin spinner.Model
}

func newModel() model {
	s := spinner.New()
	s.Spinner = spinner.Line
	r, w := io.Pipe()
	m := model{
		fromProg:    make(chan string),
		toProg:      make(chan bool),
		spin:        s,
		progInRead:  r,
		progInWrite: w,
		messages:    collect.New(10),
		// result: new(string),
	}
	return m
}

func (m model) Init() tea.Cmd { // handle debug in init ???
	if pargs.ValidProg() {
		return tea.Batch(
			runner.Run(
				pargs.ProgName(),
				pargs.ProgArgs(),
				m.fromProg,
				m.toProg,
				m.progInRead,
			),
			runner.WaitforProgResponse(m.fromProg),
			m.spin.Tick,
		)
	} else {
		return m.spin.Tick // something like a nil msg ?

	}
}

// REM: update
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	var cmds []tea.Cmd
	var cmd tea.Cmd

	if !m.done { // while prog is running
		m.result = "Running: [" + pargs.ProgName() + "]"
		m.result += "\nwith args: "
		m.result += strings.Join(pargs.ProgArgs(), " ")
		m.result += "\n Press ^C to kill"
	}
	if !pargs.ValidProg() {
		m.done = true
		m.result = "Error: "
		m.result += help.Message()
		m.result += "\nPress q or ^C or esc to exit."
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
				// m.result += "=== KILL REQUEST RECEIVED ===" //DEBUG:
				cmds = append(cmds, runner.TerminateProg(m.toProg))
				cmds = append(cmds, runner.WaitforProgResponse(m.fromProg))
				// m.done = true
			default:
				bf := make([]byte, 4)
				bf[0] = byte(msg.Type & 0xFF)
				if bf[0] == 255 && msg.Runes != nil {
					bf[0] = byte(msg.Runes[0] & 0xFF)
					// bf[1] = byte(msg.Runes[0] & 0xFF00)
				}
				m.progInWrite.Write(bf)
				cmds = append(cmds, runner.WaitforProgResponse(m.fromProg))
			}
		}
	case runner.ProgMsg:
		// m.progResult = string(msg)
		m.messages.Store(string(msg))
		cmds = append(cmds, runner.WaitforProgResponse(m.fromProg))

	case runner.ProgErrMsg:
		m.done = true
		m.result = fmt.Sprintf("Error executing: [", pargs.ProgName(), "] ErrorMsg: ", msg.Err.Error()) //TODO: error handling

	case runner.ProgDoneMsg:
		m.done = true
		m.result = fmt.Sprint("Execution of [", pargs.ProgName(), "] Done.\nPress q or ^C or esc to exit.")
	default:
		m.spin, cmd = m.spin.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)

}

// REM: view
func (m model) View() string {
	ret := "\n"
	if !m.done {
		ret += "["
		ret += m.spin.View()
		ret += "]"
	}
	ret += m.result
	ret += "\n===================\n"
	ret += m.messages.Display(10)
	return ret
}

// REM: main
func main() {
	p := tea.NewProgram(newModel())
	if _, err := p.Run(); err != nil {
		fmt.Println("dangit", err)
		os.Exit(1)
	}
}
