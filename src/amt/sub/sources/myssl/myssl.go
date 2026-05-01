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

func (m MySSL) Search(domain string, timeOut int) ([]string, error) {
	client := http.Client {
		Timeout: time.Duration(timeOut) * time.Second,
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

func (m MySSL) GetName() string {
	return "MySSL"
}