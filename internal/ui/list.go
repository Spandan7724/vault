package ui

import (
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"vault/internal/models"
)

// ListItem represents a password entry in the list
type ListItem struct {
	entry models.PasswordEntry
}

func (i ListItem) FilterValue() string {
	return i.entry.Title + " " + i.entry.Username + " " + i.entry.URL + " " + i.entry.Notes
}

func (i ListItem) Title() string {
	return i.entry.Title
}

func (i ListItem) Description() string {
	if i.entry.Username != "" {
		return i.entry.Username
	}
	if i.entry.URL != "" {
		return i.entry.URL
	}
	return "no description"
}

// ListModel represents the password list view state
type ListModel struct {
	list   list.Model
	vault  *models.Vault
	status string
}

// ListAction represents actions that can be performed on the list
type ListAction int

const (
	ListActionNone ListAction = iota
	ListActionAdd
	ListActionEdit
	ListActionDelete
	ListActionCopy
	ListActionView
)

// ListResult represents the result of a list action
type ListResult struct {
	Action  ListAction
	EntryID string
	Entry   *models.PasswordEntry
}

// NewListModel creates a new list model
func NewListModel(entries []models.PasswordEntry) ListModel {
	items := make([]list.Item, len(entries))
	for i, entry := range entries {
		items[i] = ListItem{entry: entry}
	}

	l := list.New(items, list.NewDefaultDelegate(), 80, 24)
	l.Title = "vault"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(true)
	l.SetShowHelp(false)

	return ListModel{
		list: l,
	}
}

func (m ListModel) Init() tea.Cmd {
	return nil
}

func (m ListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		m.list.SetHeight(msg.Height - 3) // Leave space for status

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "esc":
			// Clear filter if active
			if m.list.FilterState() == list.Filtering {
				m.list.ResetFilter()
				return m, nil
			}

		case "enter":
			if item, ok := m.list.SelectedItem().(ListItem); ok {
				return m, func() tea.Msg {
					return ListResult{
						Action:  ListActionView,
						EntryID: item.entry.ID,
						Entry:   &item.entry,
					}
				}
			}

		case "n":
			return m, func() tea.Msg {
				return ListResult{Action: ListActionAdd}
			}

		case "e":
			if item, ok := m.list.SelectedItem().(ListItem); ok {
				return m, func() tea.Msg {
					return ListResult{
						Action:  ListActionEdit,
						EntryID: item.entry.ID,
						Entry:   &item.entry,
					}
				}
			}

		case "d":
			if item, ok := m.list.SelectedItem().(ListItem); ok {
				return m, func() tea.Msg {
					return ListResult{
						Action:  ListActionDelete,
						EntryID: item.entry.ID,
						Entry:   &item.entry,
					}
				}
			}

		case "c":
			if item, ok := m.list.SelectedItem().(ListItem); ok {
				return m, func() tea.Msg {
					return ListResult{
						Action:  ListActionCopy,
						EntryID: item.entry.ID,
						Entry:   &item.entry,
					}
				}
			}
		}
	}

	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m ListModel) View() string {
	var s strings.Builder

	// Check if list is empty
	if len(m.list.Items()) == 0 {
		s.WriteString(TitleStyle.Render("vault") + "\n\n")
		s.WriteString(HelpStyle.Render("no passwords yet"))
		s.WriteString("\n\n")
	} else {
		s.WriteString(m.list.View())
		s.WriteString("\n")
	}

	// Status line
	if m.status != "" {
		if strings.Contains(m.status, "error") || strings.Contains(m.status, "failed") {
			s.WriteString(ErrorStyle.Render("• " + m.status))
		} else {
			s.WriteString(SuccessStyle.Render("• " + m.status))
		}
		s.WriteString("\n")
	}

	// Help
	help := []string{
		AccentStyle.Render("enter") + ": view",
		AccentStyle.Render("n") + ": new",
		AccentStyle.Render("e") + ": edit", 
		AccentStyle.Render("d") + ": delete",
		AccentStyle.Render("c") + ": copy",
		AccentStyle.Render("/") + ": filter",
		AccentStyle.Render("esc") + ": clear filter",
		AccentStyle.Render("q") + ": quit",
	}
	s.WriteString(HelpStyle.Render(strings.Join(help, " • ")))

	return s.String()
}

// UpdateEntries updates the list with new entries
func (m ListModel) UpdateEntries(entries []models.PasswordEntry) ListModel {
	items := make([]list.Item, len(entries))
	for i, entry := range entries {
		items[i] = ListItem{entry: entry}
	}
	m.list.SetItems(items)
	return m
}

// SetStatus sets a status message
func (m ListModel) SetStatus(status string) ListModel {
	m.status = status
	return m
}

// ClearStatus clears the status message
func (m ListModel) ClearStatus() ListModel {
	m.status = ""
	return m
}

// GetCurrentEntry returns the currently selected entry
func (m ListModel) GetCurrentEntry() *models.PasswordEntry {
	if item, ok := m.list.SelectedItem().(ListItem); ok {
		return &item.entry
	}
	return nil
}