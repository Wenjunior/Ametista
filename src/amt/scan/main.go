package scan

import (
	"fmt"
	"errors"
	"strconv"
	"strings"
)

import (
	"amt/utils"
	"amt/scan/scanner"
)

type ScanOptions struct {
	TimeOut int
	BatchSize int
	FileName string
	Targets []string
	Patterns []string
}

func parsePatterns(patterns []string) []int {
	var ports []int

	for _, pattern := range patterns {
		if strings.Contains(pattern, "-") {
			parts := strings.Split(pattern, "-")

			firstNumber, err := strconv.Atoi(parts[0])

			if err != nil {
				utils.Panic(err)
			}

			lastNumber, err := strconv.Atoi(parts[1])

			if err != nil {
				utils.Panic(err)
			}

			if firstNumber > lastNumber {
				utils.Panic(errors.New(fmt.Sprintf("%s is greater than %s", firstNumber, lastNumber)))
			}

			for port := firstNumber; port <= lastNumber; port++ {
				ports = append(ports, port)
			}
		} else {
			port, err := strconv.Atoi(pattern)

			if err != nil {
				utils.Panic(err)
			}

			ports = append(ports, port)
		}
	}

	return ports
}

func Run(options ScanOptions) {
	targets := options.Targets

	if options.FileName != "" {
		lines, errChan := utils.ReadFile(options.FileName)

		for line := range lines {
			_ = append(targets, line)
		}

		err := <- errChan

		if err != nil {
			utils.Panic(err)
		}
	}

	ports := parsePatterns(options.Patterns)

	scanner := scanner.Scanner {}

	for _, target := range targets {
		scanner.Run(options.BatchSize, target, ports, options.TimeOut)
	}
}