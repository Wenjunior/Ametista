package scan

import (
	"fmt"
	"errors"
	"runtime"
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
	Output string
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

	if len(options.Patterns) == 0 {
		// Source: https://exposure.shodan.io/#/

		options.Patterns = []string{"21-22", "25", "53", "80", "110", "123", "143", "161", "179", "443", "465", "500", "541", "554", "587", "646", "888", "993", "1024", "1701", "1723", "1801", "1900", "2000", "2082-2083", "2087", "3306", "4567", "5001", "5060", "5353", "5683", "5985", "7170", "7547", "7676", "8008-8010", "8080-8081", "8085", "8089", "8159", "8291", "8443", "9000", "9080", "9100", "10443", "30005-30006", "37777", "49152", "49501-49502", "50001", "50080", "50805", "51005", "58000", "58603"}
	}

	ports := parsePatterns(options.Patterns)

	if runtime.GOOS != "windows" {
		utils.IncreaseUlimit(uint64(options.BatchSize))
	}

	scanner := scanner.Scanner {}

	var results []string

	for _, target := range targets {
		result := scanner.Run(options.BatchSize, target, ports, options.TimeOut)

		results = append(results[:], result[:]...)
	}

	if options.Output != "" {
		utils.WriteResults(options.Output, results)
	}
}