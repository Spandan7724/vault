package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"vault/internal/ui"
)

const (
	appName    = "Vault"
	appVersion = "1.0.0"
	appDesc    = "A secure CLI password manager"
)

func main() {
	// Command line flags
	var (
		vaultPath = flag.String("vault", "", "Path to vault file (default: ~/.vault/vault.enc)")
		version   = flag.Bool("version", false, "Show version information")
		help      = flag.Bool("help", false, "Show help information")
	)
	flag.Parse()

	// Handle version flag
	if *version {
		fmt.Printf("%s v%s\n%s\n", appName, appVersion, appDesc)
		os.Exit(0)
	}

	// Handle help flag
	if *help {
		showHelp()
		os.Exit(0)
	}

	// Create and run the application
	app := ui.NewAppModel(*vaultPath)
	
	// Create Bubble Tea program
	program := tea.NewProgram(app)

	// Run the program
	if _, err := program.Run(); err != nil {
		log.Fatal("Error running program:", err)
	}
}

// showHelp displays help information
func showHelp() {
	fmt.Printf(`%s v%s - %s

USAGE:
    vault [OPTIONS]

OPTIONS:
    --vault PATH    Path to vault file (default: ~/.vault/vault.enc)
    --version       Show version information
    --help          Show this help message

FEATURES:
    • Secure AES-256-GCM encryption with PBKDF2 key derivation
    • Master password protection
    • Add, edit, delete, and search passwords
    • Secure password generation
    • Cross-platform clipboard integration
    • Beautiful terminal user interface

KEYBOARD SHORTCUTS:
    Navigation:
        ↑/↓ or j/k    Navigate entries
        Enter         Toggle entry details / Submit forms
        Tab           Navigate form fields
        Esc           Go back / Cancel
        Ctrl+C        Quit application

    Password Management:
        n             Add new password
        e             Edit selected password
        d             Delete selected password
        c             Copy password to clipboard
        /             Search passwords

    Form Actions:
        Ctrl+S        Save password entry
        Ctrl+G        Generate random password
        Ctrl+H        Toggle password visibility

FIRST RUN:
    On first run, you'll be prompted to create a master password.
    This password encrypts your entire vault - keep it safe!

SECURITY:
    • Passwords are encrypted with AES-256-GCM
    • Master password is processed with PBKDF2 (100,000 iterations)
    • Vault file is only readable by the owner (permissions 0600)
    • Sensitive data is cleared from memory when possible

EXAMPLES:
    vault                           # Use default vault location
    vault --vault /path/to/my.enc   # Use custom vault file
    vault --version                 # Show version
    vault --help                    # Show this help

For more information, visit: https://github.com/your-username/vault
`, appName, appVersion, appDesc)
}