package hackertarget

import (
	"io"
	"fmt"
	"time"
	"strings"
	"net/http"
)

type HackerTarget struct {}

func (h HackerTarget) Search(domain string, timeOut int) ([]string, error) {
	client := &http.Client {
		Timeout: time.Duration(timeOut) * time.Second,
	}

	url := fmt.Sprintf("https://api.hackertarget.com/hostsearch/?q=%s", domain)

	response, err := client.Get(url)

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