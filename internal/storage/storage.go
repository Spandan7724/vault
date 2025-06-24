package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"vault/internal/crypto"
	"vault/internal/models"
)

const (
	DefaultVaultFile = "vault.enc"
	VaultPermissions = 0600 // Owner read/write only
)

// Storage handles encrypted vault persistence
type Storage struct {
	filePath string
}

// NewStorage creates a new storage instance
func NewStorage(filePath string) *Storage {
	if filePath == "" {
		// Use default location in user's home directory
		homeDir, err := os.UserHomeDir()
		if err != nil {
			// Fallback to current directory
			filePath = DefaultVaultFile
		} else {
			filePath = filepath.Join(homeDir, ".vault", DefaultVaultFile)
		}
	} else if !filepath.IsAbs(filePath) {
		// Make relative paths absolute based on home directory
		homeDir, err := os.UserHomeDir()
		if err == nil {
			filePath = filepath.Join(homeDir, filePath)
		}
	}
	return &Storage{filePath: filePath}
}

// EnsureVaultDir creates the vault directory if it doesn't exist
func (s *Storage) EnsureVaultDir() error {
	dir := filepath.Dir(s.filePath)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("failed to create vault directory: %w", err)
	}
	return nil
}

// VaultExists checks if the vault file exists
func (s *Storage) VaultExists() bool {
	_, err := os.Stat(s.filePath)
	return !os.IsNotExist(err)
}

// SaveVault encrypts and saves the vault to disk
func (s *Storage) SaveVault(vault *models.Vault, masterPassword string) error {
	// Ensure vault directory exists
	if err := s.EnsureVaultDir(); err != nil {
		return err
	}

	// Marshal vault to JSON
	jsonData, err := json.Marshal(vault)
	if err != nil {
		return fmt.Errorf("failed to marshal vault: %w", err)
	}

	// Derive encryption key from master password
	key := crypto.DeriveKey(masterPassword, vault.Salt)
	defer crypto.SecureWipe(key)

	// Encrypt the vault data
	encryptedData, err := crypto.Encrypt(jsonData, key)
	if err != nil {
		return fmt.Errorf("failed to encrypt vault: %w", err)
	}

	// Create the file format: [32-byte salt][encrypted data]
	fileData := make([]byte, len(vault.Salt)+len(encryptedData))
	copy(fileData[:len(vault.Salt)], vault.Salt)
	copy(fileData[len(vault.Salt):], encryptedData)

	// Write data to file with secure permissions
	if err := os.WriteFile(s.filePath, fileData, VaultPermissions); err != nil {
		return fmt.Errorf("failed to write vault file: %w", err)
	}

	// Clear sensitive data from memory
	crypto.SecureWipe(jsonData)

	return nil
}

// LoadVault loads and decrypts the vault from disk
func (s *Storage) LoadVault(masterPassword string) (*models.Vault, error) {
	// Check if vault file exists
	if !s.VaultExists() {
		return nil, fmt.Errorf("vault file does not exist: %s", s.filePath)
	}

	// Read file data
	fileData, err := os.ReadFile(s.filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read vault file: %w", err)
	}

	// File format: [32-byte salt][encrypted data]
	if len(fileData) < crypto.SaltLength {
		return nil, fmt.Errorf("invalid vault file format")
	}

	// Extract salt and encrypted data
	salt := fileData[:crypto.SaltLength]
	encryptedData := fileData[crypto.SaltLength:]

	// Derive encryption key from master password and salt
	key := crypto.DeriveKey(masterPassword, salt)
	defer crypto.SecureWipe(key)

	// Decrypt the vault data
	jsonData, err := crypto.Decrypt(encryptedData, key)
	if err != nil {
		return nil, fmt.Errorf("invalid master password or corrupted vault")
	}

	// Parse the decrypted JSON
	var vault models.Vault
	if err := json.Unmarshal(jsonData, &vault); err != nil {
		crypto.SecureWipe(jsonData)
		return nil, fmt.Errorf("corrupted vault data")
	}

	// Clear sensitive data from memory
	crypto.SecureWipe(jsonData)

	return &vault, nil
}


// CreateNewVault creates a new encrypted vault with a random salt
func (s *Storage) CreateNewVault(masterPassword string) (*models.Vault, error) {
	// Generate a random salt
	salt, err := crypto.GenerateSalt()
	if err != nil {
		return nil, fmt.Errorf("failed to generate salt: %w", err)
	}

	// Create new vault with the salt
	vault := models.NewVault(salt)

	// Save the empty vault
	if err := s.SaveVault(vault, masterPassword); err != nil {
		return nil, fmt.Errorf("failed to save new vault: %w", err)
	}

	return vault, nil
}

// GetVaultPath returns the path to the vault file
func (s *Storage) GetVaultPath() string {
	return s.filePath
}

// DeleteVault removes the vault file from disk
func (s *Storage) DeleteVault() error {
	if !s.VaultExists() {
		return nil // Already deleted
	}
	
	if err := os.Remove(s.filePath); err != nil {
		return fmt.Errorf("failed to delete vault file: %w", err)
	}
	
	return nil
}