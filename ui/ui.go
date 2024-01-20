package ui

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/1704mori/docker-init/docker"
	"github.com/1704mori/docker-init/language"
	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	Project    *language.Project
	Inputs     []textinput.Model
	FocusIndex int
	CursorMode cursor.Mode
}

type CommandMsg struct {
	Command string
}

var (
	focusedStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("13"))
	blurredStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("7"))
	cursorStyle         = focusedStyle.Copy()
	noStyle             = lipgloss.NewStyle()
	helpStyle           = blurredStyle.Copy()
	cursorModeHelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))

	focusedButton = focusedStyle.Copy().Render("[ Generate ]")
	blurredButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Generate"))
)

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func initModel() Model {
	m := Model{}
	project := &language.Project{}
	language.DetectLanguage(project)
	m.Project = project
	m.Inputs = make([]textinput.Model, 4)

	if project.Language == "Rust" {
		m.Inputs = make([]textinput.Model, 3)
	}

	var t textinput.Model
	for i := range m.Inputs {
		t = textinput.New()
		t.Cursor.Style = cursorStyle

		if project.Language == "Go" {
			switch i {
			case 0:
				t.Prompt = "What application platform does your project use?  "
				t.SetValue(project.Language)
				t.Focus()
				t.TextStyle = focusedStyle
			case 1:
				t.Prompt = fmt.Sprintf("What version of %s do you want to use?  ", project.Language)
				t.SetValue(strings.Trim(runtime.Version(), "go"))
			case 2:
				t.Prompt = "What's the releative directory (with a leading .) of your main package?  "
			case 3:
				t.Prompt = "What port does your server listen on?  "
			}
		}

		if project.Language == "Node" {
			switch i {
			case 0:
				t.Prompt = "What application platform does your project use?  "
				t.SetValue(project.Language)
				t.Focus()
				t.TextStyle = focusedStyle
			case 1:
				t.Prompt = fmt.Sprintf("What version of %s do you want to use?  ", project.Language)
			case 2:
				t.Prompt = "What commmand do you want to use to start the app?  "
			case 3:
				t.Prompt = "What port does your server listen on?  "
			}
		}

		if project.Language == "Python" {
			switch i {
			case 0:
				t.Prompt = "What application platform does your project use?  "
				t.SetValue(project.Language)
				t.Focus()
				t.TextStyle = focusedStyle
			case 1:
				t.Prompt = fmt.Sprintf("What version of %s do you want to use?  ", project.Language)
			case 2:
				t.Prompt = "What port does your server listen on?  "
			case 3:
				t.Prompt = "What is the command to run your app (e.g: gunicorn 'myapp.example:app' --bind=0.0.0.0:8080)  "
			}
		}

		if project.Language == "Rust" {
			switch i {
			case 0:
				t.Prompt = "What application platform does your project use?  "
				t.SetValue(project.Language)
				t.Focus()
				t.TextStyle = focusedStyle
			case 1:
				t.Prompt = fmt.Sprintf("What version of %s do you want to use?  ", project.Language)
			case 2:
				t.Prompt = "What port does your server listen on?  "
			}
		}

		m.Inputs[i] = t
	}

	return m
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			// Did the user press enter while the submit button was focused?
			// If so, exit.
			if s == "enter" && m.FocusIndex == len(m.Inputs) {
				docker.GenerateFiles(m.Project)
				fmt.Printf("%v .dockerignore\n", focusedStyle.Render("CREATED:"))
				fmt.Printf("%v Dockerfile\n", focusedStyle.Render("CREATED:"))
				fmt.Printf("%v compose.yaml\n", focusedStyle.Render("CREATED:"))
				return m, tea.Quit
			}

			// Cycle indexes
			if s == "up" || s == "shift+tab" {
				m.FocusIndex--
			} else {
				m.FocusIndex++
			}

			if m.FocusIndex > len(m.Inputs) {
				m.FocusIndex = 0
			} else if m.FocusIndex < 0 {
				m.FocusIndex = len(m.Inputs)
			}

			cmds := make([]tea.Cmd, len(m.Inputs))
			for i := 0; i <= len(m.Inputs)-1; i++ {
				if i == m.FocusIndex {
					// Set focused state
					cmds[i] = m.Inputs[i].Focus()
					// m.Inputs[i].PromptStyle = focusedStyle
					m.Inputs[i].TextStyle = focusedStyle
					continue
				}
				// Remove focused state
				m.Inputs[i].Blur()
				// m.Inputs[i].PromptStyle = blurredStyle
				m.Inputs[i].TextStyle = blurredStyle
			}

			return m, tea.Batch(cmds...)
		}
	case CommandMsg:
		switch msg.Command {
		case "generate":
			docker.GenerateFiles(m.Project)
			return m, tea.Quit
		}
	}

	cmd := m.updateInputs(msg)

	return m, cmd
}

func (m *Model) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.Inputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range m.Inputs {
		m.Inputs[i], cmds[i] = m.Inputs[i].Update(msg)
	}

	if m.Project.Language == "Go" {
		m.Project.Language = strings.TrimSpace(m.Inputs[0].Value())
		m.Project.LanguageVersion = m.Inputs[1].Value()
		m.Project.RelativeDir = m.Inputs[2].Value()
		m.Project.Port = m.Inputs[3].Value()
	} else if m.Project.Language == "Node" || m.Project.Language == "Python" {
		m.Project.Language = strings.TrimSpace(m.Inputs[0].Value())
		m.Project.LanguageVersion = m.Inputs[1].Value()
		m.Project.StartCommand = m.Inputs[2].Value()
		m.Project.Port = m.Inputs[3].Value()
	} else { // rust
		m.Project.Language = strings.TrimSpace(m.Inputs[0].Value())
		m.Project.LanguageVersion = m.Inputs[1].Value()
		m.Project.Port = m.Inputs[2].Value()
	}

	return tea.Batch(cmds...)
}

func (m Model) View() string {
	var b strings.Builder

	for i := range m.Inputs {
		b.WriteString(m.Inputs[i].View())
		if i < len(m.Inputs)-1 {
			b.WriteRune('\n')
		}
	}

	button := &blurredButton
	if m.FocusIndex == len(m.Inputs) {
		button = &focusedButton
	}
	fmt.Fprintf(&b, "\n\n%s\n\n", *button)

	return b.String()
}

func Run() error {
	p := tea.NewProgram(initModel())
	_, err := p.Run()
	return err
}
