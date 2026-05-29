package crawler

import (
	"io"
	"fmt"
	"sync"
	"time"
	"slices"
	"strings"
	"net/http"
	urlparser "net/url"
)

import (
	"amt/utils/print"
	"amt/utils/print/colors"
)

import (
	"github.com/PuerkitoBio/goquery"
)

type Crawler struct {}

func (self Crawler) parseRobotsTXT(timeOut time.Duration, baseURL string) ([]string, error) {
	client := http.Client {
		Timeout: timeOut,
	}

	response, err := client.Get(baseURL + "robots.txt")

	if err != nil {
		return []string{}, err
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)

	if err != nil {
		return []string{}, err
	}

	var foundURLs []string

	for line := range strings.Lines(string(body)) {
		line = strings.TrimSpace(line)

		if len(line) == 0 || !strings.HasPrefix(line, "allow:") && !strings.HasPrefix(line, "Allow:") && !strings.HasPrefix(line, "disallow:") && !strings.HasPrefix(line, "Disallow:") {
			continue
		}

		path := strings.TrimPrefix(strings.TrimSpace(strings.Split(line, ":")[1]), "/")

		if strings.Contains(path, "*") {
			continue
		}

		foundURL := baseURL + path

		foundURLs = append(foundURLs, foundURL)
	}

	return foundURLs, nil
}

func (self Crawler) crawl(url string, timeOut time.Duration, baseURL string, locker *sync.Mutex, result *[]string) {
	client := http.Client {
		Timeout: timeOut,
	}

	response, err := client.Get(url)

	if err != nil {
		locker.Lock()

		print.Eprintln("Could not send request to " + url + ": " + err.Error())

		locker.Unlock()

		return
	}

	defer response.Body.Close()

	doc, err := goquery.NewDocumentFromReader(response.Body)

	if err != nil {
		locker.Lock()

		print.Eprintln("Could not read or parse " + url + " HTML: " + err.Error())

		locker.Unlock()

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

func (self Crawler) Run(url string, robots bool, timeOut time.Duration, batchSize int) []string {
	fmt.Println("Crawling " + url)

	parsedURL, err := urlparser.Parse(url)

	if err != nil {
		print.Eprintln("Could not parse " + url + ": " + err.Error())

		return []string{}
	}

	baseURL := parsedURL.Scheme + "://" + parsedURL.Host + "/"

	foundURLs := []string{url}

	if robots {
		result, err := self.parseRobotsTXT(timeOut, baseURL)

		if err == nil {
			foundURLs = append(foundURLs, result...)
		} else {
			print.Eprintln("Could not parse robots.txt: " + err.Error())
		}
	}

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

			locker.Lock()

			print.Cprintln(foundURL, colors.GREEN)

			locker.Unlock()

			self.crawl(foundURL, timeOut, baseURL, &locker, &foundURLs)
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