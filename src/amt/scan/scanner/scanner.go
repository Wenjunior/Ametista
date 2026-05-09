package scanner

import (
	"fmt"
	"net"
	"sync"
	"time"
)

type Scanner struct {}

func (self Scanner) scan(ipAddress string, port int, timeOut time.Duration, locker *sync.Mutex, results *[]string) {
	dialer := net.Dialer {
		Timeout: timeOut,
	}

	host := fmt.Sprintf("%s:%d", ipAddress, port)

	connection, err := dialer.Dial("tcp", host)

	if err != nil {
		return
	}

	defer connection.Close()

	locker.Lock()

	fmt.Printf("%d is open\n", port)

	*results = append(*results, host)

	locker.Unlock()
}

func (self Scanner) Run(batchSize int, ipAddress string, ports []int, timeOut time.Duration) []string {
	fmt.Printf("Scanning: %s\n", ipAddress)

	semaphore := make(chan struct{}, batchSize)

	var waitGroup sync.WaitGroup

	var locker sync.Mutex

	var results []string

	for _, port := range ports {
		semaphore <- struct{}{}

		waitGroup.Add(1)

		go func() {
			defer waitGroup.Done()

			defer func() { <- semaphore }()

			self.scan(ipAddress, port, timeOut, &locker, &results)
		} ()
	}

	waitGroup.Wait()

	return results
}