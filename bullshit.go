package main

import (
	"bufio"
	_ "embed"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var (
	words  = make(map[string][]string)
	noends = make(map[string]bool)
)

func fileExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// locateFile determines the path to the data file based on environment variables or defaults.
func locateFile() string {
	if path, exists := os.LookupEnv("BULLSHIT_FILE"); exists {
		if fileExist(path) {
			return path
		}
	}
	if confdir, err := os.UserConfigDir(); err == nil {
		path := filepath.Join(confdir, "/bullshit.txt")
		if fileExist(path) {
			return path
		}
	}
	if home, err := os.UserHomeDir(); err == nil {
		path := filepath.Join(home, ".config/bullshit.txt")
		if fileExist(path) {
			return path
		}
	}
	return "/usr/share/bullshit.txt"
}

// loadData loads words and their categories from the specified file.
func loadData(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	var currentCategory string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "%") {
			currentCategory = line[1:]
			words[currentCategory] = []string{}
		} else if strings.HasPrefix(line, "!") {
			word := line[1:]
			words[currentCategory] = append(words[currentCategory], word)
			noends[word] = true
		} else {
			words[currentCategory] = append(words[currentCategory], line)
		}
	}

	return scanner.Err()
}

// randomChoice selects a random element from a slice.
func randomChoice(options []string) string {
	if len(options) == 0 {
		return ""
	}
	return options[rand.Intn(len(options))]
}

// generateBullshit generates and returns a single nonsense sentence.
func generateBullshit() string {
	var result []string
	totalWords := rand.Intn(8) + 3
	outputCount := 0

	// Add starting words
	numStarts := rand.Intn(4)
	for i := 0; i < numStarts; i++ {
		result = append(result, randomChoice(words["start"]))
		outputCount++
	}

	// Add main words with optional suffixes
	remaining := min(totalWords-outputCount, 3)
	numWords := 0
	if outputCount < totalWords {
		numWords = rand.Intn(remaining + 1)
	}
	for i := 0; i < numWords; i++ {
		word := randomChoice(words["word"])
		suffix := ""
		if rand.Float64() < 0.2 {
			suffix = randomChoice(words["suffix"])
		}
		result = append(result, word+suffix)
		outputCount++
	}

	// Add protocol section
	if rand.Float64() > 0.5 {
		numProtocols := rand.Intn(4)
		for i := 0; i < numProtocols; i++ {
			result = append(result, randomChoice(words["protocol"]))
			if i != numProtocols-1 {
				result = append(result, "over")
			}
		}
		outputCount++
	}

	// Add more words
	remaining = min(totalWords-outputCount, 3)
	numMoreWords := 0
	if outputCount < totalWords {
		numMoreWords = rand.Intn(remaining + 1)
	}
	if outputCount+numMoreWords <= 1 {
		numMoreWords += 2
	}
	for i := 0; i < numMoreWords; i++ {
		word := randomChoice(words["word"])
		suffix := ""
		if rand.Float64() < 0.2 {
			suffix = randomChoice(words["suffix"])
		}
		result = append(result, word+suffix)
		outputCount++
	}

	// Optionally add an ending
	if rand.Float64() < 0.1 || noends[result[len(result)-1]] {
		result = append(result, randomChoice(words["end"]))
	}

	return strings.Join(result, " ")
}

// min returns the smaller of two integers.
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

//go:embed help.txt
var helpmsg string

func usage(exit int) {
	fmt.Print(helpmsg)
	os.Exit(exit)
}

func main() {
	file := ""
	times := 1
	var err error
	for i := 1; i < len(os.Args); i++ {
		switch os.Args[i] {
		case "-f", "--file":
			if i == len(os.Args)-1 {
				fmt.Fprintf(os.Stderr, "error: `%s` requires an argument\n", os.Args[i])
				os.Exit(1)
			}
			i++
			file = os.Args[i]
		default:
			if os.Args[i][0] == '-' {
				fmt.Fprintf(os.Stderr, "error: unknown option `%s`\n", os.Args[i])
				os.Exit(1)
			}
			times, err = strconv.Atoi(os.Args[i])
			if err != nil {
				fmt.Fprintf(os.Stderr, "error: invalid integer `%s`\n", os.Args[i])
				os.Exit(1)
			}
		}
	}

	// Locate the file if not provided
	if file == "" {
		file = locateFile()
	}

	// Load data
	err = loadData(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: unable to load file at %s: %v\n", file, err)
		os.Exit(1)
	}

	// Generate sentences
	for i := 0; i < times; i++ {
		fmt.Println(generateBullshit())
	}
}
