package utils

import (
	"os"
	"fmt"
	"bufio"
	"regexp"
	"strings"
)

// https://www.dolthub.com/blog/2024-02-23-colors-in-golang/
const (
	RED = "\033[31m"
	GREEN = "\033[32m"
	YELLOW = "\033[33m"

	RESET = "\033[0m"
)

// It needs to be like that because later,
// i will use this function to read files with literally millions of lines.
func ReadFile(fileName string) (<-chan string, <-chan error) {
	line := make(chan string)

	errChan := make(chan error, 1)

	go func() {
		defer close(line)

		defer close(errChan)

		file, err := os.Open(fileName)

		if err != nil {
			errChan <- err

			return
		}

		defer file.Close()

		reader := bufio.NewScanner(file)

		for reader.Scan() {
			line <- reader.Text()
		}

		err = reader.Err()

		if err != nil {
			errChan <- err

			return
		}
	} ()

	return line, errChan
}

func Eprintln(err string, color string) {
	msg := fmt.Sprintf("%s%s%s\n", color, err, RESET)

	fmt.Fprintf(os.Stderr, msg)
}

// I known that panic already exists, but it don't like debug info being showed to the final user.
func Panic(err error) {
	Eprintln(err.Error(), RED)

	os.Exit(1)
}

func RemoveDuplicatedStrings(items []string) []string {
	keys := make(map[string] bool)

	result := []string{}

	for _, item := range items {
		value := keys[item]

		if !value {
			keys[item] = true

			result = append(result, item)
		}
	}

	return result
}

func RetainSpecificStrings(items []string, expression string) []string {
	regex, err := regexp.Compile(expression)

	if err != nil {
		Panic(err)
	}

	var result []string

	for _, item := range items {
		if regex.MatchString(item) {
			result = append(result, item)
		}
	}

	return result
}

func BufferedPrint(items []string) {
	writer := bufio.NewWriter(os.Stdout)

	writer.WriteString(GREEN)

	for _, item := range items {
		writer.WriteString(fmt.Sprintf("%s\n", item))
	}

	writer.WriteString(RESET)

	writer.Flush()
}

func WriteResults(fileName string, results []string) {
	file, err := os.Create(fileName)

	if err != nil {
		Panic(err)
	}

	defer file.Close()

	_, err = file.WriteString(strings.Join(results[:], "\n"))

	if err != nil {
		Panic(err)
	}
}