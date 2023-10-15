package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func isPitched(name string) bool {
	notes := map[string]bool{"a": true, "b": true, "c": true, "d": true, "e": true, "f": true, "g": true}
	_, exists := notes[strings.ToLower(string(name[0]))]
	return exists
}

func processDirectory(directoryPath string, remoteDir string) (string, error) {
	oneshotFolders := make(map[string][]string)
	pitchedFolders := make(map[string]map[string]string)

	// Walk the directory and add all files to the appropriate map
	err := filepath.Walk(directoryPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			ext := filepath.Ext(info.Name())
			if ext != ".wav" && ext != ".mp3" {
				return nil
			}

			dir := filepath.Dir(path)
			relativeDir, _ := filepath.Rel(directoryPath, dir)
			file := info.Name()

			// If the file is pitched, add it to the pitched folders map
			if isPitched(file) {
				if pitchedFolders[relativeDir] == nil {
					pitchedFolders[relativeDir] = make(map[string]string)
				}

				pitchedFolders[relativeDir][strings.TrimSuffix(file, filepath.Ext(file))] = file

			} else {
				// Otherwise, add it to the oneshot folders map
				oneshotFolders[relativeDir] = append(oneshotFolders[relativeDir], file)
			}
		}
		return nil
	})

	if err != nil {
		return "", err
	}

	builder := &strings.Builder{}
	writer := bufio.NewWriter(builder)

	fmt.Fprint(writer, "{\n")
	fmt.Fprintf(writer, "  \"_base\": \"%s/\",\n", remoteDir)

	// Sort the keys so that the output is deterministic
	oneshotKeys := make([]string, len(oneshotFolders))
	i := 0
	for k := range oneshotFolders {
		oneshotKeys[i] = k
		i++
	}
	sort.Strings(oneshotKeys)

	// Write the oneshot folders to the file
	for j, folder := range oneshotKeys {
		files := oneshotFolders[folder]
		fmt.Fprintf(writer, "  \"%s\": [\n", folder)
		for i, file := range files {
			separator := ","
			if i == len(files)-1 {
				separator = ""
			}
			fmt.Fprintf(writer, "    \"%s/%s\"%s\n", folder, file, separator)
		}
		separator := ","
		if j == len(oneshotKeys)-1 && len(pitchedFolders) == 0 {
			separator = ""
		}
		fmt.Fprintf(writer, "  ]%s\n", separator)
	}

	// Sort the keys so that the output is deterministic
	pitchedKeys := make([]string, len(pitchedFolders))
	i = 0
	for k := range pitchedFolders {
		pitchedKeys[i] = k
		i++
	}
	sort.Strings(pitchedKeys)

	// Write the pitched folders to the file
	for j, folder := range pitchedKeys {
		files := pitchedFolders[folder]
		fmt.Fprintf(writer, "  \"%s\": {\n", folder)
		fmt.Fprintf(writer, "    \"_base\": \"%s/%s/\",\n", remoteDir, folder)

		fileKeys := make([]string, len(files))
		i := 0
		for k := range files {
			fileKeys[i] = k
			i++
		}
		sort.Strings(fileKeys)

		for i, name := range fileKeys {
			file := files[name]
			separator := ","
			if i == len(fileKeys)-1 {
				separator = ""
			}
			fmt.Fprintf(writer, "    \"%s\": \"%s\"%s\n", name, file, separator)
		}
		separator := ","
		if j == len(pitchedKeys)-1 {
			separator = ""
		}
		fmt.Fprintf(writer, "  }%s\n", separator)
	}

	fmt.Fprint(writer, "}\n")
	writer.Flush()

	return builder.String(), nil
}

func GenerateJson(localDir string, remoteDir string) error {
	result, err := processDirectory(localDir, remoteDir)
	if err != nil {
		return fmt.Errorf("an error occurred while processing the directory: %w", err)
	}

	err = os.WriteFile(localDir+"/strudel.json", []byte(result), 0644)
	if err != nil {
		return fmt.Errorf("an error occurred while writing the file: %w", err)
	}
	return nil
}
