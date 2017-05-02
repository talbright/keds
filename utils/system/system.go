package system

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

//Ripped off from https://golang.org/pkg/os/exec
func IsExecutable(file string) bool {
	d, err := os.Stat(file)
	if err != nil {
		return false
	}
	if m := d.Mode(); !m.IsDir() && m&0111 != 0 {
		return true
	}
	return false
}

func Exists(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	}
	return false
}

//Ripped off from github.com/spf13/cobra
func AbsPathify(inPath string) string {

	if strings.HasPrefix(inPath, "$HOME") {
		inPath = UserHomeDir() + inPath[5:]
	}

	if strings.HasPrefix(inPath, "$") {
		end := strings.Index(inPath, string(os.PathSeparator))
		inPath = os.Getenv(inPath[1:end]) + inPath[end:]
	}

	if filepath.IsAbs(inPath) {
		return filepath.Clean(inPath)
	}

	p, err := filepath.Abs(inPath)
	if err == nil {
		return filepath.Clean(p)
	}

	return ""
}

//Ripped off from github.com/spf13/cobra
func UserHomeDir() (home string) {
	if runtime.GOOS == "windows" {
		home = os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
	} else {
		home = os.Getenv("HOME")
	}
	return
}
