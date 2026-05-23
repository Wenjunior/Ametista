package spider

import (
	"runtime"
)

import (
	"amt/utils/ulimit"
	"amt/utils/filesystem"
)

type Options struct {
	URL string
	BatchSize int
	FileName string
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
		ulimit.IncreaseUlimit(uint64(options.BatchSize))
	}
}