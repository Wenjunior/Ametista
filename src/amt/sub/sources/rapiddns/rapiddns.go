package rapiddns

import (
	"io"
	"fmt"
	"time"
	"net/http"
)

import (
	"amt/utils"
)

type RapidDNS struct {}

func (self RapidDNS) Search(domain string, timeOut time.Duration) ([]string, error) {
	client := http.Client {
		Timeout: timeOut,
	}

	url := fmt.Sprintf("https://rapiddns.io/subdomain/%s", domain)

	response, err := client.Get(url)

	if err != nil {
		return []string{}, err
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)

	if err != nil {
		return []string{}, err
	}

	expression := fmt.Sprintf("[0-9a-z-.]+%s", domain)

	subdomains := utils.FindSpecificStrings(string(body), expression)

	return subdomains, nil
}

func (self RapidDNS) GetName() string {
	return "RapidDNS"
}