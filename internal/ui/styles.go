package ui

import (
	"github.com/charmbracelet/lipgloss"
)

// Bubble Tea inspired color palette
var (
	Primary   = lipgloss.Color("#ff006e")  // Bright pink
	Secondary = lipgloss.Color("99")  // Light purple  
	Accent    = lipgloss.Color("86")  // Cyan
	Success   = lipgloss.Color("42")  // Green
	Warning   = lipgloss.Color("214") // Orange
	Error     = lipgloss.Color("196") // Red
	Muted     = lipgloss.Color("241") // Gray
	Subtle    = lipgloss.Color("245") // Light gray
)

// Simple styles with better colors
var (
	HelpStyle = lipgloss.NewStyle().Foreground(Muted)
	ErrorStyle = lipgloss.NewStyle().Foreground(Error)
	SuccessStyle = lipgloss.NewStyle().Foreground(Success)
	HighlightStyle = lipgloss.NewStyle().Foreground(Primary)
	AccentStyle = lipgloss.NewStyle().Foreground(Accent)
	TitleStyle = lipgloss.NewStyle().Foreground(Primary).Bold(true)
)