package ui

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// LoginModel represents the login screen state
type LoginModel struct {
	passwordInput textinput.Model
	error         string
	isNewVault    bool
	confirmInput  textinput.Model
	focusIndex    int
}

// LoginResult represents the result of a login attempt
type LoginResult struct {
	Password  string
	IsNewVault bool
	Success   bool
}

// NewLoginModel creates a new login model
func NewLoginModel(isNewVault bool) LoginModel {
	passwordInput := textinput.New()
	passwordInput.Placeholder = "master password"
	passwordInput.EchoMode = textinput.EchoPassword
	passwordInput.EchoCharacter = '*'
	passwordInput.Focus()

	confirmInput := textinput.New()
	confirmInput.Placeholder = "confirm password"
	confirmInput.EchoMode = textinput.EchoPassword
	confirmInput.EchoCharacter = '*'

	return LoginModel{
		passwordInput: passwordInput,
		confirmInput:  confirmInput,
		isNewVault:    isNewVault,
		focusIndex:    0,
	}
}

func (m LoginModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m LoginModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit

		case "enter":
			return m.handleSubmit()

		case "tab", "shift+tab", "up", "down":
			if m.isNewVault {
				return m.handleTabulation(msg)
			}
		}
	}

	if m.isNewVault {
		if m.focusIndex == 0 {
			m.passwordInput, cmd = m.passwordInput.Update(msg)
			cmds = append(cmds, cmd)
		} else {
			m.confirmInput, cmd = m.confirmInput.Update(msg)
			cmds = append(cmds, cmd)
		}
	} else {
		m.passwordInput, cmd = m.passwordInput.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m LoginModel) handleSubmit() (tea.Model, tea.Cmd) {
	password := strings.TrimSpace(m.passwordInput.Value())
	
	if password == "" {
		m.error = "Password cannot be empty"
		return m, nil
	}

	if m.isNewVault {
		confirm := strings.TrimSpace(m.confirmInput.Value())
		if confirm == "" {
			m.error = "Please confirm your password"
			return m, nil
		}
		if password != confirm {
			m.error = "Passwords do not match"
			return m, nil
		}
		if len(password) < 8 {
			m.error = "Password must be at least 8 characters"
			return m, nil
		}
	}

	m.error = ""
	return m, func() tea.Msg {
		return LoginResult{
			Password:  password,
			IsNewVault: m.isNewVault,
			Success:   true,
		}
	}
}

func (m LoginModel) handleTabulation(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	s := msg.String()

	if s == "up" || s == "shift+tab" {
		m.focusIndex--
	} else {
		m.focusIndex++
	}

	if m.focusIndex > 1 {
		m.focusIndex = 0
	} else if m.focusIndex < 0 {
		m.focusIndex = 1
	}

	if m.focusIndex == 0 {
		m.passwordInput.Focus()
		m.confirmInput.Blur()
	} else {
		m.passwordInput.Blur()
		m.confirmInput.Focus()
	}

	return m, nil
}

func (m LoginModel) View() string {
	var s strings.Builder

	if m.isNewVault {
		s.WriteString(TitleStyle.Render("vault") + " - create master password\n\n")
	} else {
		s.WriteString(TitleStyle.Render("vault") + " - enter master password\n\n")
	}

	s.WriteString(m.passwordInput.View())
	s.WriteString("\n")

	if m.isNewVault {
		s.WriteString(m.confirmInput.View())
		s.WriteString("\n")
	}

	if m.error != "" {
		s.WriteString("\n")
		s.WriteString(ErrorStyle.Render(m.error))
	}

	s.WriteString("\n\n")
	if m.isNewVault {
		help := AccentStyle.Render("tab") + ": switch • " + AccentStyle.Render("enter") + ": create • " + AccentStyle.Render("ctrl+c") + ": quit"
		s.WriteString(HelpStyle.Render(help))
	} else {
		help := AccentStyle.Render("enter") + ": unlock • " + AccentStyle.Render("ctrl+c") + ": quit"
		s.WriteString(HelpStyle.Render(help))
	}

	return s.String()
}

func (m LoginModel) SetError(err string) LoginModel {
	m.error = err
	return m
}