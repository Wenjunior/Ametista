package rapiddns

import (
	"io"
	"time"
	"net/http"
)

import (
	"amt/utils/strutils"
)

type RapidDNS struct {}

func (self RapidDNS) Search(domain string, timeOut time.Duration) ([]string, error) {
	client := http.Client {
		Timeout: timeOut,
	}

	url := "https://rapiddns.io/subdomain/" + domain

	response, err := client.Get(url)

	if err != nil {
		return []string{}, err
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)

	if err != nil {
		return []string{}, err
	}

	expression := "[0-9a-z-.]+" + domain

	subdomains := strutils.FindAll(string(body), expression)

	return subdomains, nil
}

func (self RapidDNS) GetName() string {
	return "RapidDNS"
}