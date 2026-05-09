package filesystem

import (
	"os"
	"bufio"
)

import (
	"amt/utils/print"
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

func WriteResults(fileName string, results []string) {
	file, err := os.Create(fileName)

	if err != nil {
		print.Panic(err)
	}

	defer file.Close()

	for _, result := range results {
		_, err = file.WriteString(result + "\n")

		if err != nil {
			print.Panic(err)
		}
	}

	file.Sync()
}