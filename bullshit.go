package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"

	"flag"
)

type wordList struct {
	Word     []string
	Start    []string
	End      []string
	Suffix   []string
	Protocol []string
}

func fileExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// defaultInputFile determines the path to the data file based on environment variables or defaults.
func defaultInputFile() string {
	if path := os.Getenv("BULLSHIT_FILE"); path != "" {
		if fileExist(path) {
			return path
		}
	}
	if confdir, err := os.UserConfigDir(); err == nil {
		path := filepath.Join(confdir, "bullshit.txt")
		if fileExist(path) {
			return path
		}
	}
	if home, err := os.UserHomeDir(); err == nil {
		path := filepath.Join(home, ".config", "bullshit.txt")
		if fileExist(path) {
			return path
		}
	}
	return "/usr/share/bullshit.txt"
}

// loadData loads words and their categories from the specified file.
func loadData(filePath string, categories map[string]*[]string, noends map[string]struct{}) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	var current *[]string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		switch {
		case strings.HasPrefix(line, "%"):
			name := line[1:]
			current = categories[name]
		case strings.HasPrefix(line, "!"):
			line = line[1:]
			noends[line] = struct{}{}
			fallthrough
		default:
			*current = append(*current, line)
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
func generateBullshit(words *wordList, noends map[string]struct{}) string {
	var buf strings.Builder

	totalWords := rand.Intn(8) + 3
	outputCount := 0
	lastword := ""
	hassuffix := false

	// Add starting words
	numStarts := rand.Intn(4)
	for range numStarts {
		lastword = randomChoice(words.Start)
		buf.WriteString(lastword)
		buf.WriteRune(' ')
		outputCount++
	}

	// Add main words with optional suffixes
	remaining := min(totalWords-outputCount, 3)
	numWords := 0
	if outputCount < totalWords {
		numWords = rand.Intn(remaining + 1)
	}
	for range numWords {
		lastword = randomChoice(words.Word)
		hassuffix = rand.Float64() < 0.2
		buf.WriteString(lastword)
		if hassuffix {
			suffix := randomChoice(words.Suffix)
			buf.WriteString(suffix)
		}
		buf.WriteRune(' ')
		outputCount++
	}

	// Add protocol section
	if rand.Float64() > 0.5 {
		numProtocols := rand.Intn(4)
		for i := range numProtocols {
			lastword = randomChoice(words.Protocol)
			buf.WriteString(lastword)
			buf.WriteRune(' ')
			if i != numProtocols-1 {
				buf.WriteString("over ")
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
	for range numMoreWords {
		lastword = randomChoice(words.Word)
		hassuffix = rand.Float64() < 0.2
		buf.WriteString(lastword)
		if hassuffix {
			suffix := randomChoice(words.Suffix)
			buf.WriteString(suffix)
		}
		buf.WriteRune(' ')
		outputCount++
	}

	// Optionally add an ending
	_, dontend := noends[lastword]
	if dontend || hassuffix || rand.Float64() < 0.1 {
		fmt.Print(randomChoice(words.End))
	}
	return buf.String()
}

const description = "Generate one or more nonsense phrases by randomly\ncombining words and phrases from a predefined data file.\n" +
	"The phrases are constructed using categories such as starting words,\nsuffixes, protocols, and endings, producing jargon-filled or humorous output."

func main() {
	file := flag.String("input", defaultInputFile(), "input wordlist")
	times := flag.Int("count", 1, "phrases to generate")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "%s\n\n", description)
		flag.PrintDefaults()
	}
	flag.Parse()

	// Load data
	var words wordList
	noends := make(map[string]struct{})
	err := loadData(*file, map[string]*[]string{
		"word":     &words.Word,
		"start":    &words.Start,
		"end":      &words.End,
		"suffix":   &words.Suffix,
		"protocol": &words.Protocol,
	}, noends)

	if err != nil {
		fmt.Fprintf(os.Stderr, "error: unable to load file at %s: %v\n", *file, err)
		os.Exit(1)
	}

	for i := 0; i < *times; i++ {
		fmt.Println(generateBullshit(&words, noends))
	}
}
