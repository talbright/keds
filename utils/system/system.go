package system

import "os"

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
