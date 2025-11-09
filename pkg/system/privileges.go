package system

import (
	"os"
	"os/user"
	"strconv"
)

// IsSudo detects whether the current process is running with sudo privileges
func IsSudo() bool {
	// Method 1: Check SUDO_USER environment variable
	if sudoUser := os.Getenv("SUDO_USER"); sudoUser != "" {
		return true
	}

	// Method 2: Check if EUID is 0 (root)
	// This works for both sudo and direct root login
	currentUser, err := user.Current()
	if err != nil {
		return false
	}

	uid, err := strconv.Atoi(currentUser.Uid)
	if err != nil {
		return false
	}

	return uid == 0
}

// GetCurrentUser returns information about the current user
func GetCurrentUser() (*user.User, error) {
	return user.Current()
}

// GetSudoUser returns the original user when running under sudo
func GetSudoUser() string {
	return os.Getenv("SUDO_USER")
}
