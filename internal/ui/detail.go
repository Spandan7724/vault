package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"vault/internal/models"
)

// DetailModel represents the password detail view
type DetailModel struct {
	entry          models.PasswordEntry
	showPassword   bool
	width          int
	height         int
}

// DetailResult represents actions from the detail view
type DetailResult struct {
	Action    string // "back", "copy", "edit", "delete"
	EntryID   string
	Entry     *models.PasswordEntry
}

// NewDetailModel creates a new detail model
func NewDetailModel(entry models.PasswordEntry) DetailModel {
	return DetailModel{
		entry:        entry,
		showPassword: false,
	}
}

func (m DetailModel) Init() tea.Cmd {
	return nil
}

func (m DetailModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "esc", "backspace":
			return m, func() tea.Msg {
				return DetailResult{Action: "back"}
			}

		case "enter", " ":
			m.showPassword = !m.showPassword

		case "c":
			return m, func() tea.Msg {
				return DetailResult{
					Action:  "copy",
					EntryID: m.entry.ID,
					Entry:   &m.entry,
				}
			}

		case "e":
			return m, func() tea.Msg {
				return DetailResult{
					Action:  "edit",
					EntryID: m.entry.ID,
					Entry:   &m.entry,
				}
			}

		case "d":
			return m, func() tea.Msg {
				return DetailResult{
					Action:  "delete",
					EntryID: m.entry.ID,
					Entry:   &m.entry,
				}
			}
		}
	}

	return m, nil
}

func (m DetailModel) View() string {
	var s strings.Builder

	s.WriteString(TitleStyle.Render(m.entry.Title) + "\n\n")

	// Entry details
	if m.entry.Username != "" {
		s.WriteString(AccentStyle.Render("username: ") + m.entry.Username + "\n")
	}

	if m.entry.URL != "" {
		s.WriteString(AccentStyle.Render("url: ") + m.entry.URL + "\n")
	}

	// Password field with toggle
	s.WriteString(AccentStyle.Render("password: "))
	if m.showPassword {
		s.WriteString(m.entry.Password)
		s.WriteString(" " + HelpStyle.Render("(visible)"))
	} else {
		s.WriteString(strings.Repeat("•", len(m.entry.Password)))
		s.WriteString(" " + HelpStyle.Render("(hidden)"))
	}
	s.WriteString("\n")

	if m.entry.Notes != "" {
		s.WriteString(AccentStyle.Render("notes: ") + m.entry.Notes + "\n")
	}

	s.WriteString("\n")

	// Timestamps
	s.WriteString(HelpStyle.Render(fmt.Sprintf("created: %s", m.entry.CreatedAt.Format("2006-01-02 15:04"))) + "\n")
	s.WriteString(HelpStyle.Render(fmt.Sprintf("updated: %s", m.entry.UpdatedAt.Format("2006-01-02 15:04"))) + "\n")

	s.WriteString("\n")

	// Actions
	help := []string{
		AccentStyle.Render("enter") + ": toggle password",
		AccentStyle.Render("c") + ": copy",
		AccentStyle.Render("e") + ": edit",
		AccentStyle.Render("d") + ": delete",
		AccentStyle.Render("esc") + ": back",
	}
	s.WriteString(HelpStyle.Render(strings.Join(help, " • ")))

	return s.String()
}

// TogglePassword toggles password visibility
func (m DetailModel) TogglePassword() DetailModel {
	m.showPassword = !m.showPassword
	return m
}