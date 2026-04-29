package anubis_db

import (
	"fmt"
	"net/http"
	jsoniter "github.com/json-iterator/go"
)

type AnubisDB struct {}

func (a AnubisDB) Search(domain string) ([]string, error) {
	url := fmt.Sprintf("https://anubisdb.com/anubis/subdomains/%s", domain)

	response, err := http.Get(url)

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