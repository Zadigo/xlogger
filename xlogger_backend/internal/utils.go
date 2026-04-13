package internal

import (
	"bufio"
	"log"
	"os"
)

// Reads the log file and returns an array of strings
// representing each line in the file
func ReadFile(path string) ([]string, error) {
	file, err := os.Open(path)

	var logs []string = make([]string, 0)
	if err != nil {
		log.Fatal("❌ Could not open file")
		return logs, err
	}

	defer file.Close()
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		logs = append(logs, line)
	}

	return logs, nil
}
