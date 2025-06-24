# Vault - Secure CLI Password Manager

A fast, secure, and beautiful command-line password manager built in Go with Bubble Tea.


![Vault CLI Password Manager Interface](images/Screenshot%202025-06-25%20012753.png)



## Features

- **Military-grade Security**: AES-256-GCM encryption with PBKDF2 key derivation (100,000+ iterations)
- **Master Password Protection**: Single password protects your entire vault
- **Smart Search**: Fuzzy search through all your password entries
- **Lightning Fast**: Built in Go for maximum performance
- **Beautiful Interface**: Modern terminal UI with intuitive navigation
- **Clipboard Integration**: Secure password copying across platforms
- **Password Generation**: Cryptographically secure random passwords
- **Cross-Platform**: Works on Linux, macOS, and Windows
- **Memory Safety**: Sensitive data cleared from memory automatically

## Quick Start

### Installation

#### Pre-built Binaries (Recommended)

Download the latest release for your platform from [GitHub Releases](https://github.com/Spandan7724/vault/releases):

- **Linux (x64)**: `vault-vX.X.X-linux-amd64.tar.gz`
- **Linux (ARM64)**: `vault-vX.X.X-linux-arm64.tar.gz`  
- **macOS (Intel)**: `vault-vX.X.X-darwin-amd64.tar.gz`
- **macOS (Apple Silicon)**: `vault-vX.X.X-darwin-arm64.tar.gz`
- **Windows (x64)**: `vault-vX.X.X-windows-amd64.zip`
- **Windows (ARM64)**: `vault-vX.X.X-windows-arm64.zip`

**Quick Install:**
```bash
# Linux/macOS example (replace with your platform and version)
curl -L https://github.com/Spandan7724/vault/releases/download/v1.0.0/vault-v1.0.0-linux-amd64.tar.gz | tar xz
sudo mv vault /usr/local/bin/
```

#### From Source
```bash
git clone https://github.com/Spandan7724/vault.git
cd vault
go build -o vault .
```

#### System-Wide Installation

After building from source, you can install vault system-wide for convenient access from any directory:

**Linux & macOS:**
```bash
# Make executable and install to system PATH
sudo cp vault /usr/local/bin/vault
sudo chmod +x /usr/local/bin/vault

# Alternative: Install to user bin (add ~/.local/bin to PATH)
mkdir -p ~/.local/bin
cp vault ~/.local/bin/vault
chmod +x ~/.local/bin/vault
```

**Windows:**
```cmd
# Option 1: Copy to a directory in your PATH
copy vault.exe C:\Windows\System32\

# Option 2: Add vault directory to PATH environment variable
# 1. Right-click "This PC" → Properties → Advanced System Settings
# 2. Click "Environment Variables"
# 3. Edit "Path" variable and add the vault directory
# 4. Restart your terminal
```

**Verify Installation:**
```bash
# Should work from any directory
vault --version
```

#### Using Go Install
```bash
go install github.com/Spandan7724/vault@latest
```

### First Run

1. Start the application:
   ```bash
   vault
   ```
   (or `./vault` if not installed system-wide)

2. On first run, you'll be prompted to create a master password
3. Enter a strong master password (minimum 8 characters)
4. Confirm your password
5. Your encrypted vault is now ready!

## Usage

### Starting the Application

```bash
# Use default vault location (~/.vault/vault.enc)
vault

# Specify custom vault file
vault --vault /path/to/my-vault.enc

# Show version
vault --version

# Show help
vault --help
```

**Note:** If you haven't installed vault system-wide, use `./vault` instead of `vault`.

### Navigation & Keyboard Shortcuts

#### General Navigation
- `↑/↓` or `j/k` - Navigate entries
- `Enter` - Toggle entry details / Submit forms
- `Tab` - Navigate form fields
- `Esc` - Go back / Cancel
- `Ctrl+C` - Quit application

#### Password Management
- `n` - Add new password
- `e` - Edit selected password
- `d` - Delete selected password
- `c` - Copy password to clipboard
- `/` - Search passwords

#### Form Actions
- `Ctrl+S` - Save password entry
- `Ctrl+G` - Generate random password
- `Ctrl+H` - Toggle password visibility

### Managing Passwords

#### Adding a New Password
1. Press `n` from the main list
2. Fill in the form:
   - **Title** (required): e.g., "Gmail", "Bank Account"
   - **Username**: your username or email
   - **Password** (required): enter manually or generate with `Ctrl+G`
   - **URL**: website URL (optional)
   - **Notes**: additional information (optional)
3. Press `Ctrl+S` to save

#### Editing a Password
1. Select the password entry with `↑/↓`
2. Press `e` to edit
3. Modify the fields as needed
4. Press `Ctrl+S` to save changes

#### Deleting a Password
1. Select the password entry with `↑/↓`
2. Press `d` to delete
3. Confirm with `y` or cancel with `n`

#### Searching Passwords
1. Press `/` to open search
2. Type your search query
3. Search matches title, username, URL, and notes
4. Press `Enter` to apply search or `Esc` to cancel

#### Copying Passwords
1. Select the password entry with `↑/↓`
2. Press `c` to copy the password to clipboard
3. The password is now ready to paste elsewhere


##  Advanced Configuration

### Custom Vault Location
Store your vault in a custom location:
```bash
vault --vault /secure/path/my-vault.enc
```

### Environment Variables
- `DEBUG=1` - Enable debug logging to `debug.log`

##  Development

### Prerequisites
- Go 1.19+
- Git

### Building from Source
```bash
git clone https://github.com/Spandan7724/vault.git
cd vault
go mod download
go build -o vault .
```

### Running Tests
```bash
go test ./...
```

##  Troubleshooting

### Common Issues

**"Failed to copy password" error**
- On Linux: Install `xclip` or `xsel`
  ```bash
  # Ubuntu/Debian
  sudo apt install xclip
  
  # Arch Linux
  sudo pacman -S xclip
  ```

**"Permission denied" when reading vault**
- Check file permissions: `ls -la ~/.vault/vault.enc`
- Ensure you own the file: `chown $USER ~/.vault/vault.enc`

**"Invalid master password" error**
- Ensure you're entering the correct master password
- Check for caps lock or keyboard layout issues
- If forgotten, the vault cannot be recovered (this is by design for security)

**Vault file corruption**
- Restore from backup if available
- Check disk space and file system integrity



##  Acknowledgments

- [Charm](https://charm.sh/) for the amazing Bubble Tea framework
