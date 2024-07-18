package main

import (
	"context"
	"embed"
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/christitustech/linutil/cmd/ui"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

//go:embed scripts/*
var scripts embed.FS

type script struct {
	name, description, command string
}

func (s script) Title() string       { return s.name }
func (s script) Description() string { return s.description }
func (s script) Command() string     { return s.command }

func (s script) FilterValue() string {
	return s.name + " " + s.description
}

type model struct {
	list      list.Model
	spinner   spinner.Model
	executing bool
	err       error
	success   bool
	resultMsg string
	cancel    context.CancelFunc
	ctx       context.Context
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			if m.executing {
				m.cancel()
				m.executing = false
				return m, m.delayReturnToMenu()
			}
			return m, tea.Quit
		case "enter":
			cmd := strings.TrimSpace(m.list.SelectedItem().(script).command)
			if cmd != "" {
				m.executing = true
				m.success = false
				m.resultMsg = ""
				m.ctx, m.cancel = context.WithCancel(context.Background())
				return m, tea.Batch(
					m.executeScript(cmd),
					m.spinner.Tick,
				)
			}
		}

	case executionState:
		m.executing = false
		if err := msg.err; err != nil {
			m.err = err
			m.success = false
			m.resultMsg = ui.ErrStyle.Render(fmt.Sprintf("Error: %v", err)) + ui.HelpStyle.Render("Enter: retry | C-c/q: quit") + ui.StdErrStyleTitle.Render("\n\nstderr:\n") + ui.StdErrStyle.Render(msg.stderr)
			return m, nil
		}
		m.success = true
		m.resultMsg = ui.SuccessStyle.Render("Execution of script was successful!")
		return m, m.delayReturnToMenu()

	case returnToMenuMsg:
		m.executing = false
		m.success = false
		m.resultMsg = ""

	case tea.WindowSizeMsg:
		h, v := ui.DocStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	if m.executing {
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {

	if m.executing {
		return fmt.Sprintf("%s please wait", m.spinner.View())
	}

	if m.success {
		return m.resultMsg
	}

	if m.err != nil {
		return m.resultMsg
	}

	return ui.DocStyle.Render(m.list.View())
}

func loadScript(scriptPath string) string {
	cnt, err := scripts.ReadFile(scriptPath)
	if err != nil {
		log.Fatal(err)
	}
	return string(cnt)
}

type executionState struct {
	err    error
	stderr string
}

func (m model) executeScript(command string) tea.Cmd {
	return func() tea.Msg {
		cmd := exec.CommandContext(m.ctx, "bash", "-c", command)
		// cmd.Stdout = os.Stdout

		var stderr strings.Builder
		cmd.Stderr = &stderr

		err := cmd.Run()
		if err != nil {
			if err.Error() == "signal: killed" {
				return executionState{err: fmt.Errorf("Execution was cancelled by user"), stderr: stderr.String()}
			}
			return executionState{err: err, stderr: stderr.String()}
		}

		return executionState{nil, ""}
	}
}

type returnToMenuMsg struct{}

func (m model) delayReturnToMenu() tea.Cmd {
	// use waitfor to show delay  below message
	return func() tea.Msg {
		// time.Sleep(1 * time.Second)
		return returnToMenuMsg{}
	}
}

func main() {
	items := []list.Item{
		// script{name: "󰚰 Full System Update", description: "Updates your OS", command: loadScript("scripts/system-update.sh")},
		script{name: "󰚰 Full System Update", description: "Updates your OS", command: "notify-send 'hello'"},
		script{name: " Setup Bash Prompt", description: "Installs Bash Configuration from titus's repo", command: "curl -fsSL https://raw.githubusercontent.com/ChrisTitusTech/mybash/main/setup.sh | sh"},
		script{name: " Setup Neovim", description: "Setup titus's neovim config", command: "curl -fsSL https://raw.githubusercontent.com/ChrisTitusTech/neovim/main/setup.sh | sh"},
		script{name: " Build Prerequisites", description: "install all required build Prerequisites", command: loadScript("scripts/system-setup/1-compile-setup.sh")},
		script{name: " Gaming Dependencies", description: "Install gaming Dependencies", command: loadScript("scripts/system-setup/2-gaming-setup.sh")},
		script{name: " Alacritty Setup", description: "Titus's Alacritty configuration", command: loadScript("scripts/dotfiles/alacritty-setup.sh")},
		script{name: " Kitty Setup", description: "Titus's kitty configuration", command: loadScript("scripts/dotfiles/kitty-setup.sh")},
		script{name: " Rofi Setup", description: "Titus's Rofi configuration", command: loadScript("scripts/dotfiles/rofi-setup.sh")},
		script{name: "test cancel", description: "check cancel", command: "sleep 3 && pip install shit"},
	}

	s := spinner.New()
	s.Spinner = spinner.Points

	m := model{
		list:    list.New(items, list.NewDefaultDelegate(), 0, 0),
		spinner: s,
	}
	m.list.Title = "Scripts"

	m.executing = false
	m.success = false
	m.resultMsg = ""

	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		panic(err)
	}
}
