package ui

import (
	"crypto/rand"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"vault/internal/models"
)

// FormModel represents the add/edit password form state
type FormModel struct {
	inputs       []textinput.Model
	focusIndex   int
	isEdit       bool
	entryID      string
	error        string
	showPassword bool
}

// Form input indices
const (
	titleInput = iota
	usernameInput
	passwordInput
	urlInput
	notesInput
)

// FormResult represents the result of form submission
type FormResult struct {
	Title     string
	Username  string
	Password  string
	URL       string
	Notes     string
	IsEdit    bool
	EntryID   string
	Cancelled bool
}

// NewFormModel creates a new form model
func NewFormModel(isEdit bool, entry *models.PasswordEntry) FormModel {
	m := FormModel{
		inputs:     make([]textinput.Model, 5),
		isEdit:     isEdit,
		focusIndex: 0,
	}

	if entry != nil {
		m.entryID = entry.ID
	}

	// Initialize inputs
	m.inputs[titleInput] = textinput.New()
	m.inputs[titleInput].Placeholder = "title"
	m.inputs[titleInput].CharLimit = 100
	m.inputs[titleInput].Width = 40

	m.inputs[usernameInput] = textinput.New()
	m.inputs[usernameInput].Placeholder = "username"
	m.inputs[usernameInput].CharLimit = 100
	m.inputs[usernameInput].Width = 40

	m.inputs[passwordInput] = textinput.New()
	m.inputs[passwordInput].Placeholder = "password"
	m.inputs[passwordInput].EchoMode = textinput.EchoPassword
	m.inputs[passwordInput].EchoCharacter = '*'
	m.inputs[passwordInput].CharLimit = 200
	m.inputs[passwordInput].Width = 40

	m.inputs[urlInput] = textinput.New()
	m.inputs[urlInput].Placeholder = "url"
	m.inputs[urlInput].CharLimit = 200
	m.inputs[urlInput].Width = 40

	m.inputs[notesInput] = textinput.New()
	m.inputs[notesInput].Placeholder = "notes"
	m.inputs[notesInput].CharLimit = 500
	m.inputs[notesInput].Width = 40

	// Pre-fill form if editing
	if isEdit && entry != nil {
		m.inputs[titleInput].SetValue(entry.Title)
		m.inputs[usernameInput].SetValue(entry.Username)
		m.inputs[passwordInput].SetValue(entry.Password)
		m.inputs[urlInput].SetValue(entry.URL)
		m.inputs[notesInput].SetValue(entry.Notes)
	}

	// Focus first input
	m.inputs[titleInput].Focus()

	return m
}

func (m FormModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m FormModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit

		case "esc":
			return m, func() tea.Msg {
				return FormResult{Cancelled: true}
			}

		case "ctrl+s":
			return m.handleSubmit()

		case "ctrl+g":
			return m.generatePassword()

		case "ctrl+h":
			// Toggle password visibility
			m.showPassword = !m.showPassword
			if m.showPassword {
				m.inputs[passwordInput].EchoMode = textinput.EchoNormal
			} else {
				m.inputs[passwordInput].EchoMode = textinput.EchoPassword
			}

		case "tab", "shift+tab", "up", "down":
			s := msg.String()

			if s == "up" || s == "shift+tab" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			if m.focusIndex > len(m.inputs)-1 {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = len(m.inputs) - 1
			}

			for i := 0; i < len(m.inputs); i++ {
				if i == m.focusIndex {
					m.inputs[i].Focus()
				} else {
					m.inputs[i].Blur()
				}
			}

		case "enter":
			if m.focusIndex == len(m.inputs)-1 {
				return m.handleSubmit()
			} else {
				m.focusIndex++
				if m.focusIndex > len(m.inputs)-1 {
					m.focusIndex = 0
				}
				for i := 0; i < len(m.inputs); i++ {
					if i == m.focusIndex {
						m.inputs[i].Focus()
					} else {
						m.inputs[i].Blur()
					}
				}
			}
		}
	}

	m.inputs[m.focusIndex], cmd = m.inputs[m.focusIndex].Update(msg)
	return m, cmd
}

func (m FormModel) handleSubmit() (tea.Model, tea.Cmd) {
	m.error = ""

	title := strings.TrimSpace(m.inputs[titleInput].Value())
	username := strings.TrimSpace(m.inputs[usernameInput].Value())
	password := strings.TrimSpace(m.inputs[passwordInput].Value())
	url := strings.TrimSpace(m.inputs[urlInput].Value())
	notes := strings.TrimSpace(m.inputs[notesInput].Value())

	if title == "" {
		m.error = "title required"
		m.focusIndex = titleInput
		return m, nil
	}

	if password == "" {
		m.error = "password required"
		m.focusIndex = passwordInput
		return m, nil
	}

	if url != "" && !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		url = "https://" + url
	}

	return m, func() tea.Msg {
		return FormResult{
			Title:     title,
			Username:  username,
			Password:  password,
			URL:       url,
			Notes:     notes,
			IsEdit:    m.isEdit,
			EntryID:   m.entryID,
			Cancelled: false,
		}
	}
}

func (m FormModel) generatePassword() (tea.Model, tea.Cmd) {
	password := GenerateSecurePassword(16)
	m.inputs[passwordInput].SetValue(password)
	return m, nil
}

func (m FormModel) View() string {
	var s strings.Builder

	if m.isEdit {
		s.WriteString(TitleStyle.Render("edit password") + "\n\n")
	} else {
		s.WriteString(TitleStyle.Render("new password") + "\n\n")
	}

	fields := []string{"title", "username", "password", "url", "notes"}
	for i, field := range fields {
		s.WriteString(AccentStyle.Render(field + ":") + "\n")
		
		// Special handling for password field to show visibility indicator
		if i == passwordInput {
			s.WriteString(m.inputs[i].View())
			if m.showPassword {
				s.WriteString(" " + HelpStyle.Render("(visible)"))
			} else {
				s.WriteString(" " + HelpStyle.Render("(hidden)"))
			}
		} else {
			s.WriteString(m.inputs[i].View())
		}
		s.WriteString("\n\n")
	}

	if m.error != "" {
		s.WriteString(ErrorStyle.Render(m.error))
		s.WriteString("\n\n")
	}

	help := AccentStyle.Render("ctrl+s") + ": save • " + AccentStyle.Render("ctrl+g") + ": generate • " + AccentStyle.Render("ctrl+h") + ": toggle password • " + AccentStyle.Render("tab") + ": next • " + AccentStyle.Render("esc") + ": cancel"
	s.WriteString(HelpStyle.Render(help))

	return s.String()
}

func (m FormModel) SetError(err string) FormModel {
	m.error = err
	return m
}

// GenerateSecurePassword generates a cryptographically secure password
func GenerateSecurePassword(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()-_=+"
	
	if length < 8 {
		length = 8
	}
	if length > 128 {
		length = 128
	}

	password := make([]byte, length)
	charsetLen := len(charset)

	for i := range password {
		randomBytes := make([]byte, 1)
		rand.Read(randomBytes)
		password[i] = charset[int(randomBytes[0])%charsetLen]
	}

	return string(password)
}