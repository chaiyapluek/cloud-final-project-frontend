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
}

func NewLocationRepository(backendURL string) LocationRepository {
	return &locationRepository{
		backendURL: backendURL,
	}
}

func (r *locationRepository) GetLocations() ([]*entity.Location, error) {
	resp, err := http.Get(r.backendURL + "/locations")
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
	resp, err := http.Get(r.backendURL + "/locations/" + id + "/menus")
	if err != nil {
		log.Println("location repository get /locations/:id error", err)
		return nil, err
	}

	if resp.StatusCode > 299 {
		log.Println("location repository get /locations/:id status code error", resp.StatusCode)
		var errorResp entity.ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errorResp); err != nil {
			log.Println("location repository get /locations/:id decode ErrorResponse error", err)
			return nil, err
		}
		return nil, errors.New(errorResp.Message)
	}

	var location entity.LocationResponse
	if err := json.NewDecoder(resp.Body).Decode(&location); err != nil {
		log.Println("location repository get /locations/:id decode Location error", err)
		return nil, err
	}

	return location.Data, nil
}

func (r *locationRepository) GetLocationItems(locationId, menuId string) (*entity.Menu, error) {
	resp, err := http.Get(r.backendURL + "/locations/" + locationId + "/menus/" + menuId)
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
