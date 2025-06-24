package models

import (
	"crypto/rand"
	"encoding/hex"
	"strings"
	"time"
)

// PasswordEntry represents a single password entry in the vault
type PasswordEntry struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Username  string    `json:"username,omitempty"`
	Password  string    `json:"password"`
	URL       string    `json:"url,omitempty"`
	Notes     string    `json:"notes,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Vault represents the entire password vault
type Vault struct {
	Entries []PasswordEntry `json:"entries"`
	Salt    []byte          `json:"salt"`
}

// NewPasswordEntry creates a new password entry with generated ID and timestamps
func NewPasswordEntry(title, username, password, url, notes string) *PasswordEntry {
	now := time.Now()
	return &PasswordEntry{
		ID:        generateID(),
		Title:     title,
		Username:  username,
		Password:  password,
		URL:       url,
		Notes:     notes,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// Update updates the password entry fields and timestamp
func (p *PasswordEntry) Update(title, username, password, url, notes string) {
	p.Title = title
	p.Username = username
	p.Password = password
	p.URL = url
	p.Notes = notes
	p.UpdatedAt = time.Now()
}

// MatchesSearch checks if the entry matches a search query
func (p *PasswordEntry) MatchesSearch(query string) bool {
	query = strings.ToLower(query)
	return strings.Contains(strings.ToLower(p.Title), query) ||
		strings.Contains(strings.ToLower(p.Username), query) ||
		strings.Contains(strings.ToLower(p.URL), query) ||
		strings.Contains(strings.ToLower(p.Notes), query)
}

// NewVault creates a new empty vault with the given salt
func NewVault(salt []byte) *Vault {
	return &Vault{
		Entries: make([]PasswordEntry, 0),
		Salt:    salt,
	}
}

// AddEntry adds a new password entry to the vault
func (v *Vault) AddEntry(entry *PasswordEntry) {
	v.Entries = append(v.Entries, *entry)
}

// UpdateEntry updates an existing entry in the vault
func (v *Vault) UpdateEntry(id string, title, username, password, url, notes string) bool {
	for i := range v.Entries {
		if v.Entries[i].ID == id {
			v.Entries[i].Update(title, username, password, url, notes)
			return true
		}
	}
	return false
}

// DeleteEntry removes an entry from the vault by ID
func (v *Vault) DeleteEntry(id string) bool {
	for i, entry := range v.Entries {
		if entry.ID == id {
			v.Entries = append(v.Entries[:i], v.Entries[i+1:]...)
			return true
		}
	}
	return false
}

// GetEntry retrieves an entry by ID
func (v *Vault) GetEntry(id string) (*PasswordEntry, bool) {
	for _, entry := range v.Entries {
		if entry.ID == id {
			return &entry, true
		}
	}
	return nil, false
}

// SearchEntries returns entries that match the search query
func (v *Vault) SearchEntries(query string) []PasswordEntry {
	if query == "" {
		return v.Entries
	}

	var matches []PasswordEntry
	for _, entry := range v.Entries {
		if entry.MatchesSearch(query) {
			matches = append(matches, entry)
		}
	}
	return matches
}

// generateID creates a random hex ID for password entries
func generateID() string {
	bytes := make([]byte, 8)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}