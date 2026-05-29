package spider

import (
	"fmt"
	"time"
	"runtime"
)

import (
	"amt/utils/ulimit"
	"amt/spider/crawler"
	"amt/utils/filesystem"
)

type Options struct {
	URL string
	Seconds int
	BatchSize int
	FileName string
	Robots bool
	Output string
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

	crawler := crawler.Crawler {}

	for _, url := range urls {
		if len(url) == 0 {
			continue
		}

		result := crawler.Run(url, options.Robots, timeOut, options.BatchSize)

		results = append(results[:], result[:]...)
	}

	fmt.Printf("%d URL(s) was discovered\n", len(results))

	if options.Output != "" {
		filesystem.WriteResults(options.Output, results)
	}
}