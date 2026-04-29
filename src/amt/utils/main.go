package utils

import (
	"os"
	"fmt"
	"bufio"
	"strings"
)

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

func Eprintln(err string) {
	msg := fmt.Sprintf("%s\n", err)

	fmt.Fprintf(os.Stderr, msg)
}

// I known that panic already exists, but it don't like debug info being showed to the final user.
func Panic(err error) {
	Eprintln(err.Error())

	os.Exit(1)
}

func BufferedPrint(items []string) {
	writer := bufio.NewWriter(os.Stdout)

	for _, item := range items {
		writer.WriteString(fmt.Sprintf("%s\n", item))
	}

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