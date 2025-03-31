package repository

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"task4.2.3/internal/models"
)

type AddressRepository interface {
	Search(query string) ([]models.Address, error)
	Geocode(lat, lng string) ([]models.Address, error)
}

type dadataRepository struct {
	client *http.Client
	apiKey string
}

type daDataRequest struct {
	Query string `json:"query"`
}

type daDataResponse struct {
	Suggestions []struct {
		Value string `json:"value"`
		Data  struct {
			Country string `json:"country"`
			City    string `json:"city"`
			Street  string `json:"street"`
			GeoLat  string `json:"geo_lat"`
			GeoLon  string `json:"geo_lon"`
		} `json:"data"`
	} `json:"suggestions"`
}

func NewDaDataRepository(apikey string, client *http.Client) AddressRepository {
	return &dadataRepository{
		client: client,
		apiKey: apikey,
	}
}

func (r *dadataRepository) makeDaDataRequest(apiURL string, payload interface{}) (*daDataResponse, error) {
	reqBody, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("JSON encoding error: %w", err)
	}

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("error creating HTTP request: %w", err)
	}

	req.Header.Set("Authorization", "Token "+r.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending HTTP request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %w", err)
	}

	var result daDataResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &result, nil
}

func parseDaDataResponse(result *daDataResponse) ([]models.Address, error) {
	if len(result.Suggestions) == 0 {
		return nil, fmt.Errorf("no addresses found")
	}

	addresses := []models.Address{}
	for _, suggestion := range result.Suggestions {
		address := models.Address{
			Street:  suggestion.Data.Street,
			City:    suggestion.Data.City,
			Country: suggestion.Data.Country,
			Lat:     suggestion.Data.GeoLat,
			Lng:     suggestion.Data.GeoLon,
		}
		addresses = append(addresses, address)
	}

	return addresses, nil
}

func (r *dadataRepository) Search(query string) ([]models.Address, error) {
	result, err := r.makeDaDataRequest("https://suggestions.dadata.ru/suggestions/api/4_1/rs/suggest/address", daDataRequest{Query: query})
	if err != nil {
		return nil, err
	}
	return parseDaDataResponse(result)
}

func (r *dadataRepository) Geocode(lat, lng string) ([]models.Address, error) {
	result, err := r.makeDaDataRequest("https://suggestions.dadata.ru/suggestions/api/4_1/rs/geolocate/address", map[string]string{
		"lat": lat,
		"lon": lng,
	})
	if err != nil {
		return nil, err
	}
	return parseDaDataResponse(result)
}
