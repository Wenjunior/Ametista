package filesystem

import (
	"os"
	"fmt"
	"bufio"
)

func ReadFile(fileName string) (<-chan string) {
	line := make(chan string)

	go func() {
		defer close(line)

		file, err := os.Open(fileName)

		if err != nil {
			panic(fmt.Sprintf("Could not open %s: %s", fileName, err.Error()))
		}

		defer file.Close()

		reader := bufio.NewScanner(file)

		for reader.Scan() {
			line <- reader.Text()
		}

		err = reader.Err()

		if err != nil {
			panic(fmt.Sprintf("Could not read line: %s", err.Error()))
		}
	} ()

	return line
}

func WriteResults(fileName string, results []string) {
	file, err := os.Create(fileName)

	if err != nil {
		panic(fmt.Sprintf("Could not save results: %s", err.Error()))
	}

	defer file.Close()

	for _, result := range results {
		_, err = file.WriteString(result + "\n")

		if err != nil {
			panic(fmt.Sprintf("Could not write result in output file: %s", err.Error()))
		}
	}

	file.Sync()
}