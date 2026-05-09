package myssl

import (
	"fmt"
	"time"
	"net/http"
	jsoniter "github.com/json-iterator/go"
)

type MySSL struct {}

type certificatesObj struct {
	Certificate[] struct {
		Domain string `json:"domain"`
	} `json:"data"`
}

func (self MySSL) Search(domain string, timeOut time.Duration) ([]string, error) {
	client := http.Client {
		Timeout: timeOut,
	}

	url := fmt.Sprintf("https://myssl.com/api/v1/discover_sub_domain?domain=%s", domain)

	response, err := client.Get(url)

	if err != nil {
		return []string{}, err
	}

	defer response.Body.Close()

	var json = jsoniter.ConfigCompatibleWithStandardLibrary

	decoder := json.NewDecoder(response.Body)

	var certificates certificatesObj

	err = decoder.Decode(&certificates)

	if err != nil {
		return []string{}, err
	}

	var subdomains []string

	for _, certificate := range certificates.Certificate {
		subdomain := certificate.Domain

		subdomains = append(subdomains, subdomain)
	}

	return subdomains, nil
}

func (self MySSL) GetName() string {
	return "MySSL"
}