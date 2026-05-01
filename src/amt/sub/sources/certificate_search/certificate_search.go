package certificate_search

import (
	"fmt"
	"time"
	"strings"
	"net/http"
	jsoniter "github.com/json-iterator/go"
)

type CertificateSearch struct {}

type certificateObj struct {
	NameValue string `json:"name_value"`
	CommonName string `json:"common_name"`
}

func (c CertificateSearch) Search(domain string, timeOut int) ([]string, error) {
	client := http.Client {
		Timeout: time.Duration(timeOut) * time.Second,
	}

	url := fmt.Sprintf("https://crt.sh/?q=%s&output=json", domain)

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
		commonName := certificate.CommonName

		subdomains = append(subdomains, commonName)

		nameValue := certificate.NameValue

		values := strings.SplitSeq(nameValue, "\n")

		for value := range values {
			subdomains = append(subdomains, value)
		}
	}

	return subdomains, nil
}

func (c CertificateSearch) GetName() string {
	return "Certificate Search"
}