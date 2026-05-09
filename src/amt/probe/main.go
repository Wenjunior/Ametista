package probe

import (
	"fmt"
	"time"
	"sync"
	"runtime"
	"strings"
	"net/http"
	"net/http/httputil"
)

import (
	"amt/utils/unix"
	"amt/utils/print"
	"amt/utils/filesystem"
)

type ProbeOptions struct {
	URLs []string
	FileName string
	Seconds int
	BatchSize int
}

type Show struct {
	StatusCode bool
	ContentLength bool
}

func sendProbe(url string, timeOut time.Duration, locker *sync.Mutex, show Show) {
	client := http.Client {
		Timeout: timeOut,
	}

	alreadyFallback := false

	var response *http.Response

	for count := 0; count < 2; count++ {
		var err error

		response, err = client.Get(url)

		if err == nil {
			break
		}

		if strings.Contains(err.Error(), "connection refused") && !alreadyFallback {
			url, _ = strings.CutPrefix(url, "https://")

			url = "http://" + url

			alreadyFallback = true

			continue
		}

		return
	}

	locker.Lock()

	fmt.Print(url)

	if show.StatusCode {
		fmt.Printf(" [%d]", response.StatusCode)
	}

	if show.ContentLength {
		contentLength := int(response.ContentLength)

		if contentLength == -1 {
			dump, err := httputil.DumpResponse(response, true)

			if err == nil {
				contentLength = len(dump)
			}
		}

		if contentLength != -1 {
			fmt.Printf(" [%d]", contentLength)
		}
	}

	fmt.Println()

	locker.Unlock()
}

func Run(options ProbeOptions, show Show) {
	urls := options.URLs

	if options.FileName != "" {
		lines, errChan := filesystem.ReadFile(options.FileName)

		for line := range lines {
			urls = append(urls, line)
		}

		err := <- errChan

		if err != nil {
			print.Panic(err)
		}
	}

	timeOut := time.Duration(options.Seconds) * time.Second

	if runtime.GOOS != "windows" {
		unix.IncreaseUlimit(uint64(options.BatchSize))
	}

	semaphore := make(chan struct{}, options.BatchSize)

	var waitGroup sync.WaitGroup

	var locker sync.Mutex

	for _, url := range urls {
		semaphore <- struct{}{}

		waitGroup.Add(1)

		go func() {
			defer waitGroup.Done()

			defer func() { <- semaphore }()

			if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
				url = "https://" + url
			}

			sendProbe(url, timeOut, &locker, show)
		} ()
	}

	waitGroup.Wait()
}