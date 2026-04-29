package sub

import (
	"fmt"
	"sort"
	"sync"
)

import (
	"amt/utils"
	"amt/sub/sources"
	"amt/sub/sources/anubis_db"
	"amt/sub/sources/hackertarget"
)

type SubOptions struct {
	Domains []string
	FileName string
	TimeOut int
	Output string
}

func runSource(waitGroup *sync.WaitGroup, source sources.Source, domain string, timeOut int, locker *sync.Mutex, foundSubdomains *[]string) {
	defer waitGroup.Done()

	subdomains, err := source.Search(domain, timeOut)

	locker.Lock()

	if err != nil {
		utils.Eprintln(fmt.Sprintf("Could not search on %s", source.GetName()), utils.YELLOW)

		locker.Unlock()

		return
	}

	*foundSubdomains = append(*foundSubdomains, subdomains...)

	locker.Unlock()
}

func enumerateSubdomains(domain string, timeOut int) []string {
	fmt.Println(fmt.Sprintf("Enumerating subdomains for %s", domain))

	sources_ := []sources.Source {
		anubis_db.AnubisDB {},
		hackertarget.HackerTarget {},
	}

	var waitGroup sync.WaitGroup

	var locker sync.Mutex

	var subdomains []string

	for _, source := range sources_ {
		waitGroup.Add(1)

		go runSource(&waitGroup, source, domain, timeOut, &locker, &subdomains)
	}

	waitGroup.Wait()

	subdomains = utils.RetainSpecificStrings(subdomains, fmt.Sprintf("[0-9a-z-.]+%s", domain))

	subdomains = utils.RemoveDuplicatedStrings(subdomains)

	sort.Strings(subdomains)

	return subdomains
}

func Run(options SubOptions) {
	domains := options.Domains

	if options.FileName != "" {
		lines, errChan := utils.ReadFile(options.FileName)

		for line := range lines {
			_ = append(domains, line)
		}

		err := <- errChan

		if err != nil {
			utils.Panic(err)
		}
	}

	var results []string

	for _, domain := range domains {
		result := enumerateSubdomains(domain, options.TimeOut)

		utils.BufferedPrint(result)

		results = append(results[:], result[:]...)
	}

	if options.Output != "" {
		utils.WriteResults(options.Output, results)
	}
}