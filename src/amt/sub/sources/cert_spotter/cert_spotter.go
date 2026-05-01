package cert_spotter

import (
	"fmt"
	"time"
	"net/http"
	jsoniter "github.com/json-iterator/go"
)

type CertSpotter struct {}

type Certificate struct {
	DNSNames []string `json:"dns_names"`
}

func (c CertSpotter) Search(domain string, timeOut int) ([]string, error) {
	client := http.Client {
		Timeout: time.Duration(timeOut) * time.Second,
	}

	url := fmt.Sprintf("https://api.certspotter.com/v1/issuances?domain=%s&include_subdomains=true&expand=dns_names", domain)

	response, err := client.Get(url)

	if err != nil {
		return []string{}, err
	}

	defer response.Body.Close()

	var json = jsoniter.ConfigCompatibleWithStandardLibrary

	decoder := json.NewDecoder(response.Body)

	var certificates []Certificate

	err = decoder.Decode(&certificates)

	if err != nil {
		return []string{}, err
	}

	var subdomains []string

	for _, certificate := range certificates {
		for _, dnsName := range certificate.DNSNames {
			subdomains = append(subdomains, dnsName)
		}
	}

	return subdomains, nil
}

func (c CertSpotter) GetName() string {
	return "Cert Spotter"
}