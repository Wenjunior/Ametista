package scanner

import (
	"fmt"
	"net"
	"sync"
	"time"
	"strings"
)

import (
	"amt/utils"
)

type Scanner struct {}

func (s Scanner) isHostname(target string) bool {
	if !strings.Contains(target, ":") && strings.ContainsAny(target, "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ") {
		return true
	}

	return false
}

func (s Scanner) connectUsingTCP(ipAddress string, port int, timeOut int, locker *sync.Mutex, openPorts *[]int) {
	dialer := net.Dialer {
		Timeout: time.Duration(timeOut) * time.Second,
	}

	connection, err := dialer.Dial("tcp", fmt.Sprintf("%s:%d", ipAddress, port))

	if err != nil {
		return
	}

	defer connection.Close()

	locker.Lock()

	fmt.Println(fmt.Sprintf("%d is open", port))

	*openPorts = append(*openPorts, port)

	locker.Unlock()
}

func (s Scanner) Run(batchSize int, target string, ports []int, timeOut int) {
	ipAddress := target

	if s.isHostname(target) {
		fmt.Println(fmt.Sprintf("Resolving %s ...", target))

		resolvedAddresses, err := net.LookupHost(target)

		if err != nil {
			utils.Eprintln(fmt.Sprintf("Could not resolve hostname", target), utils.YELLOW)

			return
		}

		ipAddress = resolvedAddresses[0]
	}

	semaphore := make(chan struct{}, batchSize)

	var waitGroup sync.WaitGroup

	var locker sync.Mutex

	var openPorts []int

	fmt.Println(fmt.Sprintf("Scanning: %s", ipAddress))

	for _, port := range ports {
		semaphore <- struct{}{}

		waitGroup.Add(1)

		go func() {
			defer waitGroup.Done()

			defer func() { <- semaphore }()

			s.connectUsingTCP(ipAddress, port, timeOut, &locker, &openPorts)
		} ()
	}

	waitGroup.Wait()
}