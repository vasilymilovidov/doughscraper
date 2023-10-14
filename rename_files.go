package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var sharpToFlat = map[string]string{
	"C#": "db",
	"D#": "eb",
	"F#": "gb",
	"G#": "ab",
	"A#": "bb",
}

func RenameFiles(localDir string) error {
	files, err := os.ReadDir(localDir)
	if err != nil {
		return fmt.Errorf("error reading directory: %w", err)
	}

	r, err := regexp.Compile(`-([A-G]#?)(\d*)-`)
	if err != nil {
		return fmt.Errorf("error compiling regex: %w", err)
	}

	for _, file := range files {
		matches := r.FindStringSubmatch(file.Name())
		if len(matches) > 2 {
			note, octave := matches[1], matches[2]
			if strings.Contains(note, "#") {
				note = sharpToFlat[note]
			}
			newFilename := strings.ToLower(note + octave + filepath.Ext(file.Name()))
			err = os.Rename(filepath.Join(localDir, file.Name()), filepath.Join(localDir, newFilename))
			if err != nil {
				return fmt.Errorf("error renaming file: %w", err)
			}
		}
	}

	return nil
}
