package probe

import (
	"io"
	"fmt"
	"net"
	"time"
	"sync"
	"errors"
	"context"
	"runtime"
	"strings"
	"net/http"
	"net/http/httputil"
)

import (
	"github.com/dlclark/regexp2/v2"
	"github.com/dlclark/regexp2/v2/compat"
)

import (
	"amt/utils/unix"
	"amt/utils/print"
	"amt/utils/filesystem"
)

type Options struct {
	FileName string
	Seconds int
	BatchSize int
	Output string
}

type Show struct {
	IPAddress bool
	StatusCode bool
	Server bool
	XPoweredBy bool
	Location bool
	ContentLength bool
	ContentType bool
	Title bool
}

func buildResult(url string, show Show, address string, response *http.Response) string {
	builder := strings.Builder {}

	builder.WriteString(url)

	if show.IPAddress {
		ip := strings.Split(address, ":")[0]

		builder.WriteString(fmt.Sprintf(" [%s]", ip))
	}

	if show.StatusCode {
		statusCode := response.StatusCode

		builder.WriteString(fmt.Sprintf(" [%d]", statusCode))
	}

	if show.Server {
		server := response.Header.Get("Server")

		builder.WriteString(fmt.Sprintf(" [%s]", server))
	}

	if show.XPoweredBy {
		xPoweredBy := response.Header.Get("X-Powered-By")

		builder.WriteString(fmt.Sprintf(" [%s]", xPoweredBy))
	}

	if show.Location {
		location := response.Header.Get("Location")

		builder.WriteString(fmt.Sprintf(" [%s]", location))
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
			builder.WriteString(fmt.Sprintf(" [%d]", contentLength))
		}
	}

	if show.ContentType {
		contentType := response.Header.Get("Content-Type")

		builder.WriteString(fmt.Sprintf(" [%s]", contentType))
	}

	if show.Title {
		body, err := io.ReadAll(response.Body)

		if err == nil {
			pattern := compat.MustCompile("(?<=<title>)(.*)(?=<\\/title>)", regexp2.RE2)

			title := pattern.FindString(string(body))

			builder.WriteString(fmt.Sprintf(" [%s]", title))
		}
	}

	return builder.String()
}

func sendProbe(url string, timeOut time.Duration, locker *sync.Mutex, show Show, results *[]string) {
	request, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return
	}

	client := http.Client {
		Timeout: timeOut,
		Transport: &http.Transport {
			DialContext: func(ctx context.Context, network, host string) (net.Conn, error) {
				var dialer net.Dialer

				connection, err := dialer.DialContext(ctx, network, host)

				if err == nil {
					request.RemoteAddr = connection.RemoteAddr().String()
				}

				return connection, err
			},
		},
		CheckRedirect: func(newRequest *http.Request, via []*http.Request) error {
			request = newRequest

			return nil
		},
	}

	alreadyFallback := false

	var response *http.Response

	for count := 0; count < 2; count++ {
		var err error

		response, err = client.Do(request)

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

	if request.RemoteAddr == "" {
		resolver := &net.Resolver {
			PreferGo: true,
		}

		resolverCtx, cancel := context.WithTimeout(context.Background(), timeOut)

		defer cancel()

		hostname := strings.Split(url, "/")[2]

		ips, err := resolver.LookupIPAddr(resolverCtx, hostname)

		if err == nil {
			request.RemoteAddr = ips[0].String()
		}
	}

	result := buildResult(url, show, request.RemoteAddr, response)

	locker.Lock()

	fmt.Println(result)

	*results = append(*results, result)

	locker.Unlock()
}

func Run(options Options, show Show) {
	if options.FileName == "" {
		print.Panic(errors.New("Argument -l is required"))
	}

	var urls []string

	lines := filesystem.ReadFile(options.FileName)

	for line := range lines {
		urls = append(urls, line)
	}

	timeOut := time.Duration(options.Seconds) * time.Second

	if runtime.GOOS != "windows" {
		unix.IncreaseUlimit(uint64(options.BatchSize))
	}

	semaphore := make(chan struct{}, options.BatchSize)

	var waitGroup sync.WaitGroup

	var locker sync.Mutex

	var results []string

	for _, url := range urls {
		if len(url) == 0 {
			continue
		}

		semaphore <- struct{}{}

		waitGroup.Add(1)

		go func() {
			defer waitGroup.Done()

			defer func() { <- semaphore }()

			if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
				url = "https://" + url
			}

			sendProbe(url, timeOut, &locker, show, &results)
		} ()
	}

	waitGroup.Wait()

	if options.Output != "" {
		filesystem.WriteResults(options.Output, results)
	}
}