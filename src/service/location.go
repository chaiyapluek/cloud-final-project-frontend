package service

import (
	"fmt"

	"dev.chaiyapluek.cloud.final.frontend/src/repository"
	"dev.chaiyapluek.cloud.final.frontend/template/pages/location"
	"dev.chaiyapluek.cloud.final.frontend/template/pages/menu"
	"dev.chaiyapluek.cloud.final.frontend/template/pages/order"
)

type LocationService interface {
	GetLocations() ([]*location.LocationProps, error)
	GetLocationByID(id string) (*menu.MenuProps, error)
	GetLocationItems(locationId, menuId string) (*order.OrderProps, error)
}

type locationService struct {
	locationRepo repository.LocationRepository
}

func NewLocationService(locationRepo repository.LocationRepository) LocationService {
	return &locationService{
		locationRepo: locationRepo,
	}
}

func (s *locationService) GetLocations() ([]*location.LocationProps, error) {
	locations, err := s.locationRepo.GetLocations()
	if err != nil {
		return nil, err
	}

	resp := []*location.LocationProps{}
	for _, l := range locations {
		resp = append(resp, &location.LocationProps{
			Id:   l.ID,
			Name: l.Name,
		})
	}

	return resp, nil
}

func (s *locationService) GetLocationByID(id string) (*menu.MenuProps, error) {
	location, err := s.locationRepo.GetLocationByID(id)
	if err != nil {
		return nil, err
	}

	menuCards := []*menu.MenuCardProps{}
	for _, m := range location.Menus {
		menuCards = append(menuCards, &menu.MenuCardProps{
			Id:    m.ID,
			Name:  m.Name,
			Price: m.Price,
			Img:   m.IconImage,
		})
	}
	resp := &menu.MenuProps{
		LocationId:   location.ID,
		LocationName: location.Name,
		MenuCards:    menuCards,
	}

	return resp, nil
}

func (s *locationService) GetLocationItems(locationId, menuId string) (*order.OrderProps, error) {
	menu, err := s.locationRepo.GetLocationItems(locationId, menuId)
	if err != nil {
		return nil, err
	}

	steps := []*order.Step{}
	for StepIdx, ig := range menu.Steps {
		items := []*order.Option{}
		for _, i := range ig.Options {
			items = append(items, &order.Option{
				Name:   i.Name,
				Value:  i.Value,
				Price:  i.Price,
				Select: false,
			})
		}
		steps = append(steps, &order.Step{
			Name:        ig.Name,
			Description: ig.Description,
			FormName:    fmt.Sprintf("step-%d", StepIdx),
			Type:        ig.Type,
			Required:    ig.Required,
			Min:         ig.Min,
			Max:         ig.Max,
			Items:       items,
		})
	}

	resp := &order.OrderProps{
		LocationId:      locationId,
		MenuId:          menu.ID,
		MenuName:        menu.Name,
		MenuDescription: menu.Description,
		MenuPrice:       menu.Price,
		MenuImage:       menu.ThumbnailImage,
		Complete:        false,
		Steps:           steps,
	}

	return resp, nil
}
