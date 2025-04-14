package common

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type Address struct {
	Road        string `json:"road"`
	HouseNumber string `json:"house_number"`
	City        string `json:"city"`
	Town        string `json:"town"`
	State       string `json:"state"`
	Country     string `json:"country"`
	Postcode    string `json:"post_code"`
}

type NominatimResponse struct {
	Address Address `json:"address"`
}

func ReverseGeocode(lat, lon float32) (string, string, string, string, error) {
	baseURL := "https://nominatim.openstreetmap.org/reverse"
	params := url.Values{}
	params.Add("format", "json")
	params.Add("lat", fmt.Sprintf("%f", lat))
	params.Add("lon", fmt.Sprintf("%f", lon))
	params.Add("addressdetails", "1")

	requestURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())
	resp, err := http.Get(requestURL)
	if err != nil {
		return "", "", "", "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", "", "", "", err
	}

	var result NominatimResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return "", "", "", "", err
	}

	// Si city está vacío, se puede utilizar town
	city := result.Address.City
	if city == "" {
		city = result.Address.Town
	}

	return result.Address.Road, city, result.Address.State, result.Address.Postcode, nil
}
