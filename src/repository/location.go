package repository

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"dev.chaiyapluek.cloud.final.frontend/src/entity"
)

type LocationRepository interface {
	GetLocations() ([]*entity.Location, error)
	GetLocationByID(id string) (*entity.Location, error)
	GetLocationItems(locationId, menuId string) (*entity.Menu, error)
}

type locationRepository struct {
	backendURL string
	apiKey     string
}

func NewLocationRepository(backendURL string, apiKey string) LocationRepository {
	return &locationRepository{
		backendURL: backendURL,
		apiKey:     apiKey,
	}
}

func (r *locationRepository) GetLocations() ([]*entity.Location, error) {
	httpReq, err := http.NewRequest("GET", r.backendURL+"/locations", nil)
	if err != nil {
		log.Println("location repository get /locations error", err)
		return nil, err
	}
	httpReq.Header.Set("X-API-KEY", r.apiKey)
	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		log.Println("location repository get /locations error", err)
		return nil, err
	}

	if resp.StatusCode > 299 {
		log.Println("location repository get /locations status code error", resp.StatusCode)
		var errorResp entity.ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errorResp); err != nil {
			log.Println("location repository get /locations decode ErrorResponse error", err)
			return nil, err
		}
		return nil, errors.New(errorResp.Message)
	}

	var locationResp entity.AllLocationResponse
	if err := json.NewDecoder(resp.Body).Decode(&locationResp); err != nil {
		log.Println("location repository get /locations decode LocationResponse error", err)
		return nil, err
	}

	return locationResp.Data, nil
}

func (r *locationRepository) GetLocationByID(id string) (*entity.Location, error) {
	httpReq, err := http.NewRequest("GET", r.backendURL+"/locations/"+id+"/menus", nil)
	if err != nil {
		log.Println("location repository get /locations/:id/menus error", err)
		return nil, err
	}
	httpReq.Header.Set("X-API-KEY", r.apiKey)
	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		log.Println("location repository get /locations/:id/menus error", err)
		return nil, err
	}

	if resp.StatusCode > 299 {
		log.Println("location repository get /locations/:id/menus status code error", resp.StatusCode)
		var errorResp entity.ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errorResp); err != nil {
			log.Println("location repository get /locations/:id/menus decode ErrorResponse error", err)
			return nil, err
		}
		return nil, errors.New(errorResp.Message)
	}

	var location entity.LocationResponse
	if err := json.NewDecoder(resp.Body).Decode(&location); err != nil {
		log.Println("location repository get /locations/:id/menus decode Location error", err)
		return nil, err
	}

	return location.Data, nil
}

func (r *locationRepository) GetLocationItems(locationId, menuId string) (*entity.Menu, error) {
	httpReq, err := http.NewRequest("GET", r.backendURL+"/locations/"+locationId+"/menus/"+menuId, nil)
	if err != nil {
		log.Println("location repository get /locations/:locationId/menus/:menuId error", err)
		return nil, err
	}
	httpReq.Header.Set("X-API-KEY", r.apiKey)
	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		log.Println("location repository get /locations/:locationId/menus/:menuId error", err)
		return nil, err
	}

	if resp.StatusCode > 299 {
		log.Println("location repository get /locations/:locationId/menus/:menuId status code error", resp.StatusCode)
		var errorResp entity.ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errorResp); err != nil {
			log.Println("location repository get /locations/:locationId/menus/:menuId decode ErrorResponse error", err)
			return nil, err
		}
		return nil, errors.New(errorResp.Message)
	}

	var menu entity.MenuResponse
	if err := json.NewDecoder(resp.Body).Decode(&menu); err != nil {
		log.Println("location repository get /locations/:locationId/menus/:menuId decode Menu error", err)
		return nil, err
	}

	return menu.Data, nil
}
