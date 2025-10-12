package cmd

import (
	"os"
)

type Entry struct {
	name  string
	isDir bool
}

func ListDirectoryContents() ([]Entry, error) {
	entries, err := os.ReadDir(".")

	if err != nil {
		return nil, err
	}

	var names []Entry
	for _, entry := range entries {
		if entry != nil {
			name := entry.Name()
			names = append(names, Entry{name, entry.IsDir()})
		}
	}
	return names, nil
}

func (e Entry) IsGitDir() bool {
	return e.name == ".git"
}
