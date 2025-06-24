package ui

import (
	"fmt"
	"os/exec"
	"runtime"

	"github.com/atotto/clipboard"
	tea "github.com/charmbracelet/bubbletea"
	"vault/internal/models"
	"vault/internal/storage"
)

// AppState represents the current state of the application
type AppState int

const (
	StateLogin AppState = iota
	StateList
	StateDetail
	StateForm
	StateConfirmDelete
)

// AppModel is the main application model
type AppModel struct {
	state         AppState
	storage       *storage.Storage
	vault         *models.Vault
	masterPassword string
	
	// Screen models
	loginModel    LoginModel
	listModel     ListModel
	detailModel   DetailModel
	formModel     FormModel
	
	// Temporary state
	pendingDeleteID string
}

// NewAppModel creates a new application model
func NewAppModel(vaultPath string) AppModel {
	storage := storage.NewStorage(vaultPath)
	isNewVault := !storage.VaultExists()
	
	return AppModel{
		state:       StateLogin,
		storage:     storage,
		loginModel:  NewLoginModel(isNewVault),
		listModel:   NewListModel([]models.PasswordEntry{}),
	}
}

func (m AppModel) Init() tea.Cmd {
	return m.loginModel.Init()
}

func (m AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	}

	switch m.state {
	case StateLogin:
		return m.handleLoginState(msg)
	case StateList:
		return m.handleListState(msg)
	case StateDetail:
		return m.handleDetailState(msg)
	case StateForm:
		return m.handleFormState(msg)
	case StateConfirmDelete:
		return m.handleConfirmDeleteState(msg)
	}

	return m, nil
}

func (m AppModel) handleLoginState(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	model, cmd := m.loginModel.Update(msg)
	m.loginModel = model.(LoginModel)

	if result, ok := msg.(LoginResult); ok && result.Success {
		m.masterPassword = result.Password
		
		if result.IsNewVault {
			vault, err := m.storage.CreateNewVault(m.masterPassword)
			if err != nil {
				m.loginModel = m.loginModel.SetError("Failed to create vault: " + err.Error())
				return m, nil
			}
			m.vault = vault
			m.listModel = m.listModel.SetStatus("vault created")
		} else {
			vault, err := m.storage.LoadVault(m.masterPassword)
			if err != nil {
				m.loginModel = m.loginModel.SetError("Failed to unlock vault: " + err.Error())
				return m, nil
			}
			m.vault = vault
			m.listModel = m.listModel.SetStatus(fmt.Sprintf("%d entries loaded", len(vault.Entries)))
		}

		m.listModel = m.listModel.UpdateEntries(m.vault.Entries)
		m.state = StateList
	}

	return m, cmd
}

func (m AppModel) handleListState(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	model, cmd := m.listModel.Update(msg)
	m.listModel = model.(ListModel)

	if result, ok := msg.(ListResult); ok {
		switch result.Action {
		case ListActionView:
			if result.Entry != nil {
				m.detailModel = NewDetailModel(*result.Entry)
				m.state = StateDetail
				return m, m.detailModel.Init()
			}

		case ListActionAdd:
			m.formModel = NewFormModel(false, nil)
			m.state = StateForm
			return m, m.formModel.Init()

		case ListActionEdit:
			if result.Entry != nil {
				m.formModel = NewFormModel(true, result.Entry)
				m.state = StateForm
				return m, m.formModel.Init()
			}

		case ListActionDelete:
			if result.Entry != nil {
				m.pendingDeleteID = result.EntryID
				m.state = StateConfirmDelete
			}

		case ListActionCopy:
			if result.Entry != nil {
				if err := m.copyToClipboard(result.Entry.Password); err != nil {
					m.listModel = m.listModel.SetStatus("failed to copy password")
				} else {
					m.listModel = m.listModel.SetStatus("password copied")
				}
			}
		}
	}

	return m, cmd
}

func (m AppModel) handleDetailState(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	model, cmd := m.detailModel.Update(msg)
	m.detailModel = model.(DetailModel)

	if result, ok := msg.(DetailResult); ok {
		switch result.Action {
		case "back":
			m.state = StateList

		case "copy":
			if result.Entry != nil {
				if err := m.copyToClipboard(result.Entry.Password); err != nil {
					m.listModel = m.listModel.SetStatus("failed to copy password")
				} else {
					m.listModel = m.listModel.SetStatus("password copied")
				}
			}
			m.state = StateList

		case "edit":
			if result.Entry != nil {
				m.formModel = NewFormModel(true, result.Entry)
				m.state = StateForm
				return m, m.formModel.Init()
			}

		case "delete":
			if result.Entry != nil {
				m.pendingDeleteID = result.EntryID
				m.state = StateConfirmDelete
			}
		}
	}

	return m, cmd
}

func (m AppModel) handleFormState(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	model, cmd := m.formModel.Update(msg)
	m.formModel = model.(FormModel)

	if result, ok := msg.(FormResult); ok {
		if result.Cancelled {
			m.state = StateList
		} else {
			if result.IsEdit {
				success := m.vault.UpdateEntry(result.EntryID, result.Title, result.Username, result.Password, result.URL, result.Notes)
				if success {
					m.listModel = m.listModel.SetStatus("password updated")
				} else {
					m.listModel = m.listModel.SetStatus("failed to update password")
				}
			} else {
				entry := models.NewPasswordEntry(result.Title, result.Username, result.Password, result.URL, result.Notes)
				m.vault.AddEntry(entry)
				m.listModel = m.listModel.SetStatus("password added")
			}

			if err := m.storage.SaveVault(m.vault, m.masterPassword); err != nil {
				m.listModel = m.listModel.SetStatus("failed to save vault")
			}

			m.listModel = m.listModel.UpdateEntries(m.vault.Entries)
			m.state = StateList
		}
	}

	return m, cmd
}

func (m AppModel) handleConfirmDeleteState(msg tea.Msg) (tea.Model, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch keyMsg.String() {
		case "y", "Y":
			if m.vault.DeleteEntry(m.pendingDeleteID) {
				if err := m.storage.SaveVault(m.vault, m.masterPassword); err != nil {
					m.listModel = m.listModel.SetStatus("failed to save vault")
				} else {
					m.listModel = m.listModel.SetStatus("password deleted")
				}
				m.listModel = m.listModel.UpdateEntries(m.vault.Entries)
			} else {
				m.listModel = m.listModel.SetStatus("failed to delete password")
			}
			m.state = StateList
			m.pendingDeleteID = ""

		case "n", "N", "esc":
			m.state = StateList
			m.pendingDeleteID = ""
		}
	}

	return m, nil
}

func (m AppModel) View() string {
	switch m.state {
	case StateLogin:
		return m.loginModel.View()
	case StateList:
		return m.listModel.View()
	case StateDetail:
		return m.detailModel.View()
	case StateForm:
		return m.formModel.View()
	case StateConfirmDelete:
		return m.renderConfirmDelete()
	}
	return ""
}

func (m AppModel) renderConfirmDelete() string {
	entry, found := m.vault.GetEntry(m.pendingDeleteID)
	if !found {
		return "entry not found"
	}

	return fmt.Sprintf(`%s

%s %s
%s %s

%s`, 
		TitleStyle.Render("delete password?"),
		AccentStyle.Render("title:"), entry.Title,
		AccentStyle.Render("username:"), entry.Username,
		HelpStyle.Render("y: delete • n: cancel"))
}

func (m AppModel) copyToClipboard(text string) error {
	// Try the go-clipboard library first (most reliable)
	if err := clipboard.WriteAll(text); err == nil {
		return nil
	}

	// Fallback to system commands
	return m.copyToClipboardSystem(text)
}

func (m AppModel) copyToClipboardSystem(text string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("pbcopy")
	case "linux":
		// Try different clipboard utilities in order of preference
		if _, err := exec.LookPath("wl-copy"); err == nil {
			// Wayland
			cmd = exec.Command("wl-copy")
		} else if _, err := exec.LookPath("xclip"); err == nil {
			// X11 with xclip
			cmd = exec.Command("xclip", "-selection", "clipboard")
		} else if _, err := exec.LookPath("xsel"); err == nil {
			// X11 with xsel
			cmd = exec.Command("xsel", "--clipboard", "--input")
		} else {
			return fmt.Errorf("no clipboard utility found (install xclip, xsel, or wl-clipboard)")
		}
	case "windows":
		cmd = exec.Command("clip")
	default:
		return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}

	// Set up the command to read from stdin
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdin pipe: %w", err)
	}

	// Start the command
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start clipboard command: %w", err)
	}

	// Write the text to the clipboard command
	if _, err := stdin.Write([]byte(text)); err != nil {
		stdin.Close()
		cmd.Wait()
		return fmt.Errorf("failed to write to clipboard: %w", err)
	}

	// Close stdin and wait for the command to complete
	stdin.Close()
	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("clipboard command failed: %w", err)
	}

	return nil
}