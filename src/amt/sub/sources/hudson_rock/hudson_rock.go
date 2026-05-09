package hudson_rock

import (
	"time"
	"strings"
	"net/http"
)

import (
	jsoniter "github.com/json-iterator/go"
)

type HudsonRock struct {}

type responseData struct {
	Data struct {
		ClientsURLs[] struct {
			URL string `json:"url"`
		} `json:"clients_urls"`
		EmployeesURLs[] struct {
			URL string `json:"url"`
		} `json:"employees_urls"`
	} `json:"data"`
}

func (self HudsonRock) Search(domain string, timeOut time.Duration) ([]string, error) {
	client := http.Client {
		Timeout: timeOut,
	}

	url := "https://cavalier.hudsonrock.com/api/json/v2/osint-tools/urls-by-domain?domain=" + domain

	response, err := client.Get(url)

	if err != nil {
		return []string{}, err
	}

	defer response.Body.Close()

	var json = jsoniter.ConfigCompatibleWithStandardLibrary

	decoder := json.NewDecoder(response.Body)

	var responseData responseData

	err = decoder.Decode(&responseData)

	if err != nil {
		return []string{}, err
	}

	var subdomains []string

	for _, record := range append(responseData.Data.ClientsURLs, responseData.Data.EmployeesURLs...) {
		subdomain := strings.Split(record.URL, "/")[2]

		subdomains = append(subdomains, subdomain)
	}

	return subdomains, nil
}

func (self HudsonRock) GetName() string {
	return "Hudson Rock"
}