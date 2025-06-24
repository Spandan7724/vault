package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"

	"golang.org/x/crypto/pbkdf2"
)

const (
	KeyLength    = 32          // AES-256 key length
	SaltLength   = 32          // Salt length for PBKDF2
	Iterations   = 100000      // PBKDF2 iterations
	NonceLength  = 12          // GCM nonce length
)

var (
	ErrInvalidKeyLength = errors.New("invalid key length")
	ErrInvalidNonce     = errors.New("invalid nonce length")
	ErrDecryption       = errors.New("decryption failed")
)

// GenerateSalt creates a cryptographically secure random salt
func GenerateSalt() ([]byte, error) {
	salt := make([]byte, SaltLength)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return nil, fmt.Errorf("failed to generate salt: %w", err)
	}
	return salt, nil
}

// DeriveKey derives an encryption key from a password using PBKDF2
func DeriveKey(password string, salt []byte) []byte {
	return pbkdf2.Key([]byte(password), salt, Iterations, KeyLength, sha256.New)
}

// Encrypt encrypts plaintext using AES-256-GCM
func Encrypt(plaintext []byte, key []byte) ([]byte, error) {
	if len(key) != KeyLength {
		return nil, ErrInvalidKeyLength
	}

	// Create AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	// Create GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	// Generate random nonce
	nonce := make([]byte, NonceLength)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("failed to generate nonce: %w", err)
	}

	// Encrypt and authenticate the data
	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return ciphertext, nil
}

// Decrypt decrypts ciphertext using AES-256-GCM
func Decrypt(ciphertext []byte, key []byte) ([]byte, error) {
	if len(key) != KeyLength {
		return nil, ErrInvalidKeyLength
	}

	if len(ciphertext) < NonceLength {
		return nil, ErrInvalidNonce
	}

	// Create AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	// Create GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	// Extract nonce and encrypted data
	nonce := ciphertext[:NonceLength]
	encrypted := ciphertext[NonceLength:]

	// Decrypt and verify the data
	plaintext, err := gcm.Open(nil, nonce, encrypted, nil)
	if err != nil {
		return nil, ErrDecryption
	}

	return plaintext, nil
}

// SecureWipe overwrites sensitive data in memory
func SecureWipe(data []byte) {
	for i := range data {
		data[i] = 0
	}
}

// SecureWipeString overwrites sensitive string data in memory
func SecureWipeString(s *string) {
	if s == nil {
		return
	}
	// Convert to byte slice and wipe
	b := []byte(*s)
	SecureWipe(b)
	*s = ""
}