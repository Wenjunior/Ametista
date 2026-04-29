package anubis_db

import (
	"fmt"
	"time"
	"net/http"
	jsoniter "github.com/json-iterator/go"
)

type AnubisDB struct {}

func (a AnubisDB) Search(domain string, timeOut int) ([]string, error) {
	client := http.Client {
		Timeout: time.Duration(timeOut) * time.Second,
	}

	url := fmt.Sprintf("https://anubisdb.com/anubis/subdomains/%s", domain)

	response, err := client.Get(url)

	if err != nil {
		return []string{}, err
	}

	defer response.Body.Close()

	var json = jsoniter.ConfigCompatibleWithStandardLibrary

	decoder := json.NewDecoder(response.Body)

	var subdomains []string

	err = decoder.Decode(&subdomains)

	if err != nil {
		return []string{}, err
	}

	return subdomains, nil
}

func (a AnubisDB) GetName() string {
	return "AnubisDB"
}