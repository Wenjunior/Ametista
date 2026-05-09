package anubis_db

import (
	"time"
	"net/http"
)

import (
	jsoniter "github.com/json-iterator/go"
)

type AnubisDB struct {}

func (self AnubisDB) Search(domain string, timeOut time.Duration) ([]string, error) {
	client := http.Client {
		Timeout: timeOut,
	}

	url := "https://anubisdb.com/anubis/subdomains/" + domain

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

func (self AnubisDB) GetName() string {
	return "AnubisDB"
}