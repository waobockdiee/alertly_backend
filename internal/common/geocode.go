package common

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
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
	// ✅ OPTIMIZACIÓN: Timeout de 5 segundos para evitar bloqueos
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	baseURL := "https://nominatim.openstreetmap.org/reverse"
	params := url.Values{}
	params.Add("format", "json")
	params.Add("lat", fmt.Sprintf("%f", lat))
	params.Add("lon", fmt.Sprintf("%f", lon))
	params.Add("addressdetails", "1")

	requestURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	// ✅ OPTIMIZACIÓN: Request con timeout y User-Agent
	req, err := http.NewRequestWithContext(ctx, "GET", requestURL, nil)
	if err != nil {
		return "", "", "", "", fmt.Errorf("error creating request: %w", err)
	}

	// ✅ OPTIMIZACIÓN: User-Agent requerido por Nominatim
	req.Header.Set("User-Agent", "Alertly/1.0")

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", "", "", "", fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	// ✅ OPTIMIZACIÓN: Verificar status code
	if resp.StatusCode != http.StatusOK {
		return "", "", "", "", fmt.Errorf("nominatim returned status: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", "", "", "", fmt.Errorf("error reading response: %w", err)
	}

	var result NominatimResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return "", "", "", "", fmt.Errorf("error parsing response: %w", err)
	}

	// Si city está vacío, se puede utilizar town
	city := result.Address.City
	if city == "" {
		city = result.Address.Town
	}

	return result.Address.Road, city, result.Address.State, result.Address.Postcode, nil
}
