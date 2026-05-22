package scan

import (
	"fmt"
	"net"
	"time"
	"runtime"
	"strconv"
	"strings"
)

import (
	"amt/utils/print"
	"amt/scan/scanner"
	"amt/utils/ulimit"
	"amt/utils/filesystem"
)

type Options struct {
	Seconds int
	BatchSize int
	FileName string
	Target string
	Patterns string
	Output string
}

func parsePatterns(patterns []string) []int {
	var ports []int

	for _, pattern := range patterns {
		if strings.Contains(pattern, "-") {
			parts := strings.Split(pattern, "-")

			firstNumber, err := strconv.Atoi(parts[0])

			if err != nil {
				print.Panic(fmt.Errorf("Could not convert %s: %s", parts[0], err.Error()))
			}

			lastNumber, err := strconv.Atoi(parts[1])

			if err != nil {
				print.Panic(fmt.Errorf("Could not convert %s: %s", parts[1], err.Error()))
			}

			if firstNumber > lastNumber {
				print.Panic(fmt.Errorf("%d is greater than %d", firstNumber, lastNumber))
			}

			for port := firstNumber; port <= lastNumber; port++ {
				ports = append(ports, port)
			}
		} else {
			port, err := strconv.Atoi(pattern)

			if err != nil {
				print.Panic(fmt.Errorf("Could not convert %s: %s", pattern, err.Error()))
			}

			ports = append(ports, port)
		}
	}

	return ports
}

func isHostname(target string) bool {
	if !strings.Contains(target, ":") && strings.ContainsAny(target, "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ") {
		return true
	}

	return false
}

func Run(options Options) {
	targets := []string{options.Target}

	if options.FileName != "" {
		lines := filesystem.ReadFile(options.FileName)

		for line := range lines {
			targets = append(targets, line)
		}
	}

	patterns := strings.Split(options.Patterns, ",")

	if len(patterns) == 1 {
		// Source: https://exposure.shodan.io/#/

		patterns = []string{"21-22", "25", "53", "80", "110", "123", "143", "161", "179", "443", "465", "500", "541", "554", "587", "646", "888", "993", "1024", "1701", "1723", "1801", "1900", "2000", "2082-2083", "2087", "3306", "4567", "5001", "5060", "5353", "5683", "5985", "7170", "7547", "7676", "8008-8010", "8080-8081", "8085", "8089", "8159", "8291", "8443", "9000", "9080", "9100", "10443", "30005-30006", "37777", "49152", "49501-49502", "50001", "50080", "50805", "51005", "58000", "58603"}
	}

	ports := parsePatterns(patterns)

	if runtime.GOOS != "windows" {
		ulimit.IncreaseUlimit(uint64(options.BatchSize))
	}

	timeOut := time.Duration(options.Seconds) * time.Second

	scanner := scanner.Scanner {}

	var results []string

	for _, target := range targets {
		if len(target) == 0 {
			continue
		}

		ipAddress := target

		if isHostname(target) {
			fmt.Printf("Resolving %s\n", target)

			resolvedAddresses, err := net.LookupHost(target)

			if err != nil {
				print.Eprintln("Could not resolve hostname")

				continue
			}

			ipAddress = resolvedAddresses[0]
		}

		result := scanner.Run(options.BatchSize, ipAddress, ports, timeOut)

		results = append(results[:], result[:]...)
	}

	if options.Output != "" {
		filesystem.WriteResults(options.Output, results)
	}
}