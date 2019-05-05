package util

import (
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

func GopathDir() string {
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		return filepath.Join(HomeDir(), "go")
	}
	return gopath
}

func HomeDir() string {
	if home := os.Getenv("HOME"); home != "" {
		return home
	}
	if usr, err := user.Current(); err == nil {
		return usr.HomeDir
	}
	return ""
}

func FileExist(filePath string) bool {
	_, err := os.Stat(filePath)
	if err != nil && os.IsNotExist(err) {
		return false
	}

	return true
}

func StringInSlice(str string, sli []string) bool {
	for _, s := range sli {
		if s == str {
			return true
		}
	}
	return false
}

func ToCamelCase(str string) (result string) {
	ss := strings.Split(str, "_")
	for _, s := range ss {
		if len(s) > 0 {
			result += strings.ToUpper(string(s[0])) + s[1:]
		}
	}
	return
}

func StringContainsSliceByLower(str string, sli []string) bool {
	str = strings.ToLower(str)
	for _, s := range sli {
		if strings.Contains(str, strings.ToLower(s)) {
			return true
		}
	}
	return false
}