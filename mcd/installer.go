package mcd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func Install(filePath string, alias string) {
	fileDir := filepath.Dir(filePath)
	fileName := filepath.Base(filePath)
	backupFileName := fmt.Sprintf("%s.%d.mcd-install-backup", fileName, time.Now().UnixMilli())
	backupFilePath := filepath.Join(fileDir, backupFileName)

	copyError := copyFile(filePath, backupFilePath)

	if copyError != nil {
		log.Fatal(copyError)
	}

	b, readError := os.ReadFile(filePath)

	if readError != nil {
		log.Fatal(readError)
	}

	fileContents := string(b)
	lines := strings.Split(strings.TrimSpace(fileContents), "\n")
	modifiedLines := make([]string, 0, len(lines)+128)

	started := false
	ended := false

	for _, line := range lines {
		if !started && strings.TrimSpace(line) == "# start: mono-cd" {
			started = true
			continue
		} else if started && strings.TrimSpace(line) == "# end: mono-cd" {
			ended = true
			continue
		}

		if !started || ended {
			modifiedLines = append(modifiedLines, line)
		}
	}

	shellFunctionLines := strings.Split(
		strings.Replace(ShellFunction, "mcd", alias, 1),
		"\n",
	)

	modifiedLines = append(
		modifiedLines,
		shellFunctionLines...,
	)

	writeError := os.WriteFile(
		filePath,
		[]byte(
			strings.TrimSpace(strings.Join(modifiedLines, "\n"))+"\n",
		),
		0644,
	)

	if writeError != nil {
		log.Fatal(writeError)
	}

	backupRemovalError := os.Remove(backupFilePath)

	if backupRemovalError != nil {
		log.Fatal(backupRemovalError)
	}
}
