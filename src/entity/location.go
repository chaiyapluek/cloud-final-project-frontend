package entity

import "time"

type SuccessResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	Code      int        `json:"code"`
	Message   string     `json:"message"`
	Path      string     `json:"path"`
	Timestamp *time.Time `json:"timestamp"`
}

type AllLocationResponse struct {
	SuccessResponse
	Data []*Location `json:"data"`
}

type LocationResponse struct {
	SuccessResponse
	Data *Location `json:"data"`
}

type MenuResponse struct {
	SuccessResponse
	Data *Menu `json:"data"`
}

type Location struct {
	ID    string  `json:"id"`
	Name  string  `json:"name"`
	Menus []*Menu `json:"menus"`
}

type Menu struct {
	ID             string  `json:"id"`
	Name           string  `json:"name"`
	Description    string  `json:"description"`
	Price          int     `json:"price"`
	IconImage      string  `json:"iconImage"`
	ThumbnailImage string  `json:"thumbnailImage"`
	Steps          []*Step `json:"steps"`
}

type Step struct {
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Type        string     `json:"type"`
	Required    bool       `json:"required"`
	Min         int        `json:"min"`
	Max         int        `json:"max"`
	Options     []*Options `json:"options"`
}

type Options struct {
	Name  string
	Value string
	Price int
}
