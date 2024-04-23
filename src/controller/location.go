package controller

import (
	"encoding/base64"
	"encoding/json"
	"log"

	"dev.chaiyapluek.cloud.final.frontend/src/service"
	"dev.chaiyapluek.cloud.final.frontend/template/pages/location"
	"dev.chaiyapluek.cloud.final.frontend/template/pages/menu"
	"dev.chaiyapluek.cloud.final.frontend/template/pages/order"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

type LocationController struct {
	locationService service.LocationService
	sessionService  service.SessionService
}

func NewLocationController(locationService service.LocationService, sessionService service.SessionService) *LocationController {
	return &LocationController{
		locationService: locationService,
		sessionService:  sessionService,
	}
}

func (c *LocationController) GetLocations(e echo.Context) error {
	locations, err := c.locationService.GetLocations()
	if err != nil {
		log.Println("get all location controller", err)
		return e.String(500, "Internal Server Error")
	}

	return location.LocationList(locations).Render(e.Request().Context(), e.Response().Writer)
}

func (c *LocationController) GetLocationMenu(e echo.Context) error {
	id := e.Param("id")
	location, err := c.locationService.GetLocationByID(id)
	if err != nil {
		log.Println("get location by id controller", err)
		return e.String(500, "Internal Server Error")
	}

	sess, _ := session.Get("sessionid", e)
	id = sess.Values["id"].(string)
	detail := c.sessionService.GetSessionDetail(id)
	if detail != nil {
		detail.CurrentLocation = location.LocationId
	}

	return menu.Menu(location).Render(e.Request().Context(), e.Response().Writer)
}

func (c *LocationController) GetLocationItems(e echo.Context) error {
	locationId := e.Param("locationId")
	menuId := e.Param("menuId")
	preferenceBase64 := e.QueryParam("preference")
	var base64Decoded []byte
	var err error
	if preferenceBase64 != "" {
		base64Decoded, err = base64.StdEncoding.DecodeString(preferenceBase64)
		if err != nil {
			base64Decoded = []byte{}
		}
	}
	preference := map[string]interface{}{}
	if len(base64Decoded) > 0 {
		err = json.Unmarshal(base64Decoded, &preference)
		if err != nil {
			preference = map[string]interface{}{}
		}
	}

	orderProp, err := c.locationService.GetLocationItems(locationId, menuId)
	if err != nil {
		log.Println("get location items controller", err)
		return e.String(500, "Internal Server Error")
	}

	isLogin, ok := e.Get("isLogin").(bool)
	if !ok || !isLogin {
		return order.Order(*orderProp, false).Render(e.Request().Context(), e.Response().Writer)
	}

	if len(preference) > 0 {
		for _, g := range orderProp.Steps {
			if v, ok := preference[g.FormName]; ok {
				preferenceValues := v.([]interface{})
				for _, i := range g.Items {
					for _, pv := range preferenceValues {
						if i.Value == pv {
							i.Select = true
						}
					}
				}
			}
		}
	}

	price := orderProp.MenuPrice
	for _, g := range orderProp.Steps {
		for _, i := range g.Items {
			if _, ok := preference[g.FormName]; ok {
				for _, v := range preference[g.FormName].([]interface{}) {
					if v == i.Value {
						price += i.Price
					}
				}
			}
		}
	}
	orderProp.TotalPrice = price
	isComplete := true
	for _, g := range orderProp.Steps {
		formVals, ok := preference[g.FormName].([]interface{})
		if !ok && g.Required {
			isComplete = false
			break
		}

		if g.Required && (len(formVals) < g.Min || len(formVals) > g.Max) {
			isComplete = false
			break
		}
	}
	orderProp.Complete = isComplete

	sess, _ := session.Get("sessionid", e)
	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400,
		HttpOnly: true,
	}
	id := sess.Values["id"].(string)

	if id == "" {
		id, detail := c.sessionService.NewSessionDetail()
		sess.Values["id"] = id
		detail.CurrentMenu = orderProp
		detail.Preferences = preference
		c.sessionService.UpdateSessionDetail(id, detail)
	} else {
		sessionDetail := c.sessionService.GetSessionDetail(id)
		if sessionDetail == nil {
			id, detail := c.sessionService.NewSessionDetail()
			sess.Values["id"] = id
			detail.CurrentMenu = orderProp
			detail.Preferences = preference
			c.sessionService.UpdateSessionDetail(id, detail)
		} else {
			sessionDetail.CurrentMenu = orderProp
			sessionDetail.Preferences = preference
			c.sessionService.UpdateSessionDetail(id, sessionDetail)
		}
	}
	sess.Save(e.Request(), e.Response())

	return order.Order(*orderProp, true).Render(e.Request().Context(), e.Response().Writer)
}
