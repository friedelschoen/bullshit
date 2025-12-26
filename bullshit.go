package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"flag"
)

var (
	words  = make(map[string][]string)
	noends = make(map[string]bool)
)

func fileExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// defaultInputFile determines the path to the data file based on environment variables or defaults.
func defaultInputFile() string {
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
		switch {
		case strings.HasPrefix(line, "%"):
			currentCategory = line[1:]
		case strings.HasPrefix(line, "!"):
			line = line[1:]
			noends[line] = true
			fallthrough
		default:
			slc, _ := words[currentCategory]
			words[currentCategory] = append(slc, line)
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
func generateBullshit() {
	totalWords := rand.Intn(8) + 3
	outputCount := 0
	lastword := ""
	hassuffix := false

	// Add starting words
	numStarts := rand.Intn(4)
	for i := 0; i < numStarts; i++ {
		lastword = randomChoice(words["start"])
		fmt.Printf("%s ", lastword)
		outputCount++
	}

	// Add main words with optional suffixes
	remaining := min(totalWords-outputCount, 3)
	numWords := 0
	if outputCount < totalWords {
		numWords = rand.Intn(remaining + 1)
	}
	for i := 0; i < numWords; i++ {
		lastword = randomChoice(words["word"])
		hassuffix = rand.Float64() < 0.2
		if hassuffix {
			suffix := randomChoice(words["suffix"])
			fmt.Printf("%s%s ", lastword, suffix)
		} else {
			fmt.Printf("%s ", lastword)
		}
		outputCount++
	}

	// Add protocol section
	if rand.Float64() > 0.5 {
		numProtocols := rand.Intn(4)
		for i := 0; i < numProtocols; i++ {
			lastword = randomChoice(words["protocol"])
			fmt.Printf("%s ", lastword)
			if i != numProtocols-1 {
				fmt.Printf("over ")
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
		lastword = randomChoice(words["word"])
		hassuffix = rand.Float64() < 0.2
		if hassuffix {
			suffix := randomChoice(words["suffix"])
			fmt.Printf("%s%s ", lastword, suffix)
		} else {
			fmt.Printf("%s ", lastword)
		}
		outputCount++
	}

	// Optionally add an ending
	if rand.Float64() < 0.1 || noends[lastword] || hassuffix {
		fmt.Printf("%s\n", randomChoice(words["end"]))
	} else {
		fmt.Println()
	}
}

func printSorted() {
	keys := make([]string, len(words))
	i := 0
	for key := range words {
		keys[i] = key
		i++
	}
	slices.Sort(keys)
	for i, key := range keys {
		values := words[key]
		if i > 0 {
			fmt.Println()
		}
		fmt.Println("%" + key)
		slices.SortFunc(values, func(left, right string) int {
			if noends[left] != noends[right] {
				if noends[left] {
					return 1
				}
				return -1
			}
			return strings.Compare(left, right)
		})
		for _, value := range values {
			if noends[value] {
				fmt.Print("!")
			}
			fmt.Println(value)
		}
	}
}

const description = "Generate one or more nonsense phrases by randomly\ncombining words and phrases from a predefined data file.\n" +
	"The phrases are constructed using categories such as starting words,\nsuffixes, protocols, and endings, producing jargon-filled or humorous output."

func main() {
	file := flag.String("input", defaultInputFile(), "input wordlist")
	times := flag.Int("count", 1, "sentences to generate")
	sort := flag.Bool("sort", false, "sort the wordlist and print it to stdout")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "%s\n\n", description)
		flag.PrintDefaults()
	}
	flag.Parse()

	// Load data
	err := loadData(*file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: unable to load file at %s: %v\n", *file, err)
		os.Exit(1)
	}

	if *sort {
		printSorted()
		return
	}

	for i := 0; i < *times; i++ {
		generateBullshit()
	}
}
