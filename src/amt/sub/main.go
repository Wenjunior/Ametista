package sub

import (
	"fmt"
	"sort"
	"sync"
	"time"
)

import (
	"amt/utils/print"
	"amt/utils/strutils"
	"amt/utils/filesystem"
)

import (
	"amt/sub/sources"
	"amt/sub/sources/myssl"
	"amt/sub/sources/rapiddns"
	"amt/sub/sources/anubis_db"
	"amt/sub/sources/hudson_rock"
	"amt/sub/sources/cert_spotter"
	"amt/sub/sources/hackertarget"
	"amt/sub/sources/certificate_search"
)

type Options struct {
	Domain string
	FileName string
	Seconds int
	Output string
}

func runSource(waitGroup *sync.WaitGroup, source sources.Source, domain string, timeOut time.Duration, locker *sync.Mutex, foundSubdomains *[]string) {
	defer waitGroup.Done()

	subdomains, err := source.Search(domain, timeOut)

	locker.Lock()

	if err != nil {
		print.Eprintln("Could not search on " + source.GetName() + ": " + err.Error())

		locker.Unlock()

		return
	}

	*foundSubdomains = append(*foundSubdomains, subdomains...)

	locker.Unlock()
}

func enumerateSubdomains(domain string, timeOut time.Duration) []string {
	fmt.Printf("Enumerating subdomains for %s\n", domain)

	sources := []sources.Source {
		myssl.MySSL {},
		rapiddns.RapidDNS {},
		anubis_db.AnubisDB {},
		hudson_rock.HudsonRock {},
		cert_spotter.CertSpotter {},
		hackertarget.HackerTarget {},
		certificate_search.CertificateSearch {},
	}

	var waitGroup sync.WaitGroup

	var locker sync.Mutex

	var subdomains []string

	for _, source := range sources {
		waitGroup.Add(1)

		go runSource(&waitGroup, source, domain, timeOut, &locker, &subdomains)
	}

	waitGroup.Wait()

	subdomains = strutils.Retain(subdomains, "^[0-9a-z-.]+" + domain + "$")

	subdomains = strutils.RemoveDuplicated(subdomains)

	sort.Strings(subdomains)

	return subdomains
}

func Run(options Options) {
	domains := []string{options.Domain}

	if options.FileName != "" {
		lines := filesystem.ReadFile(options.FileName)

		for line := range lines {
			domains = append(domains, line)
		}
	}

	timeOut := time.Duration(options.Seconds) * time.Second

	var results []string

	for _, domain := range domains {
		if len(domain) == 0 {
			continue
		}

		result := enumerateSubdomains(domain, timeOut)

		print.BufferedPrint(result)

		results = append(results[:], result[:]...)
	}

	fmt.Printf("%d subdomains was discovered\n", len(results))

	if options.Output != "" {
		filesystem.WriteResults(options.Output, results)
	}
}