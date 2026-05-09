package certificate_search

import (
	"time"
	"strings"
	"net/http"
)

import (
	jsoniter "github.com/json-iterator/go"
)

type CertificateSearch struct {}

type certificateObj struct {
	NameValue string `json:"name_value"`
	CommonName string `json:"common_name"`
}

func (self CertificateSearch) Search(domain string, timeOut time.Duration) ([]string, error) {
	client := http.Client {
		Timeout: timeOut,
	}

	url := "https://crt.sh/?q=" + domain + "&output=json"

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

func (self CertificateSearch) GetName() string {
	return "Certificate Search"
}