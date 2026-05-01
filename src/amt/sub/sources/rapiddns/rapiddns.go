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

func (r RapidDNS) Search(domain string, timeOut int) ([]string, error) {
	client := http.Client {
		Timeout: time.Duration(timeOut) * time.Second,
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

func (r RapidDNS) GetName() string {
	return "RapidDNS"
}