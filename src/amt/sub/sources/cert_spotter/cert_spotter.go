package cert_spotter

import (
	"time"
	"net/http"
)

import (
	jsoniter "github.com/json-iterator/go"
)

type CertSpotter struct {}

type certificateObj struct {
	DNSNames []string `json:"dns_names"`
}

func (self CertSpotter) Search(domain string, timeOut time.Duration) ([]string, error) {
	client := http.Client {
		Timeout: timeOut,
	}

	url := "https://api.certspotter.com/v1/issuances?domain=" + domain + "&include_subdomains=true&expand=dns_names"

	response, err := client.Get(url)

	if err != nil {
		return []string{}, err
	}

	defer response.Body.Close()

	var json = jsoniter.ConfigCompatibleWithStandardLibrary

	decoder := json.NewDecoder(response.Body)

	var certificates []certificateObj

	err = decoder.Decode(&certificates)

	if err != nil {
		return []string{}, err
	}

	var subdomains []string

	for _, certificate := range certificates {
		subdomains = append(subdomains, certificate.DNSNames...)
	}

	return subdomains, nil
}

func (self CertSpotter) GetName() string {
	return "Cert Spotter"
}