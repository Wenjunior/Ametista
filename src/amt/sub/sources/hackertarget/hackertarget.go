package hackertarget

import (
	"io"
	"fmt"
	"strings"
	"net/http"
)

type HackerTarget struct {}

func (h HackerTarget) Search(domain string) ([]string, error) {
	url := fmt.Sprintf("https://api.hackertarget.com/hostsearch/?q=%s", domain)

	response, err := http.Get(url)

	if err != nil {
		return []string{}, err
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)

	if err != nil {
		return []string{}, err
	}

	lines := strings.Split(string(body), "\n")

	var subdomains []string

	for _, line := range lines {
		subdomain := strings.Split(line, ",")[0]

		subdomains = append(subdomains, subdomain)
	}

	return subdomains, nil
}

func (h HackerTarget) GetName() string {
	return "HackerTarget"
}