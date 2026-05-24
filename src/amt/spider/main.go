package spider

import (
	"fmt"
	"sync"
	"time"
	"slices"
	"runtime"
	"strings"
	"net/http"
	urlparser "net/url"
)

import (
	"amt/utils/print"
	"amt/utils/ulimit"
	"amt/utils/filesystem"
	"amt/utils/print/colors"
)

import (
	"github.com/PuerkitoBio/goquery"
)

type Options struct {
	URL string
	Seconds int
	BatchSize int
	FileName string
	Output string
}

func scrap(url string, timeOut time.Duration, baseURL string, locker *sync.Mutex, result *[]string) {
	client := http.Client {
		Timeout: timeOut,
	}

	response, err := client.Get(url)

	if err != nil {
		return
	}

	defer response.Body.Close()

	doc, err := goquery.NewDocumentFromReader(response.Body)

	if err != nil {
		return
	}

	foundURLs := []string{}

	doc.Find("a").Each(func(_ int, tag *goquery.Selection) {
		href, wasFound := tag.Attr("href")

		if !wasFound {
			return
		}

		if strings.HasPrefix(href, "http://") || strings.HasPrefix(href, "https://") {
			if strings.HasPrefix(href, baseURL) && !slices.Contains(foundURLs, href) {
				foundURLs = append(foundURLs, href)
			}

			return
		}

		href = strings.TrimPrefix(href, "/")

		href = baseURL + href

		if !slices.Contains(foundURLs, href) {
			foundURLs = append(foundURLs, href)
		}
	})

	locker.Lock()

	for _, foundURL := range foundURLs {
		if !slices.Contains(*result, foundURL) {
			*result = append(*result, foundURL)
		}
	}

	locker.Unlock()
}

func crawl(url string, timeOut time.Duration, batchSize int) []string {
	fmt.Println("Crawling " + url)

	parsedURL, err := urlparser.Parse(url)

	if err != nil {
		print.Eprintln("Could not parse " + url + ": " + err.Error())

		return []string{}
	}

	baseURL := parsedURL.Scheme + "://" + parsedURL.Host + "/"

	foundURLs := []string{url}

	semaphore := make(chan struct{}, batchSize)

	var waitGroup sync.WaitGroup

	var locker sync.Mutex

	index := 0

	for {
		foundURL := foundURLs[index]

		semaphore <- struct{}{}

		waitGroup.Add(1)

		go func() {
			defer waitGroup.Done()

			defer func() { <- semaphore }()

			print.Cprintln(foundURL, colors.GREEN)

			scrap(foundURL, timeOut, baseURL, &locker, &foundURLs)
		} ()

		if index == len(foundURLs) - 1 {
			waitGroup.Wait()

			if index < len(foundURLs) - 1 {
				continue
			}

			break
		}

		index++
	}

	return foundURLs
}

func Run(options Options) {
	urls := []string{options.URL}

	if options.FileName != "" {
		lines := filesystem.ReadFile(options.FileName)

		for line := range lines {
			urls = append(urls, line)
		}
	}

	if runtime.GOOS != "windows" {
		ulimit.Increase(uint64(options.BatchSize))
	}

	timeOut := time.Duration(options.Seconds) * time.Second

	var results []string

	for _, url := range urls {
		if len(url) == 0 {
			continue
		}

		result := crawl(url, timeOut, options.BatchSize)

		results = append(results[:], result[:]...)
	}

	if options.Output != "" {
		filesystem.WriteResults(options.Output, results)
	}
}