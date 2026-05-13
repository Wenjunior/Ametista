package filesystem

import (
	"os"
	"fmt"
	"bufio"
)

import (
	"amt/utils/print"
)

func ReadFile(fileName string) (<-chan string) {
	line := make(chan string)

	go func() {
		defer close(line)

		file, err := os.Open(fileName)

		if err != nil {
			print.Panic(fmt.Errorf("Could not open %s: %s", fileName, err.Error()))
		}

		defer file.Close()

		reader := bufio.NewScanner(file)

		for reader.Scan() {
			line <- reader.Text()
		}

		err = reader.Err()

		if err != nil {
			print.Panic(fmt.Errorf("Could not read line: %s", err.Error()))
		}
	} ()

	return line
}

func WriteResults(fileName string, results []string) {
	file, err := os.Create(fileName)

	if err != nil {
		print.Panic(fmt.Errorf("Could not save results: %s", err.Error()))
	}

	defer file.Close()

	for _, result := range results {
		_, err = file.WriteString(result + "\n")

		if err != nil {
			print.Panic(fmt.Errorf("Could not write result in output file: %s", err.Error()))
		}
	}

	file.Sync()
}