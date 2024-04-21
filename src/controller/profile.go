package controller

import (
	"log"

	"dev.chaiyapluek.cloud.final.frontend/src/service"
	"dev.chaiyapluek.cloud.final.frontend/template/pages/profile"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

type profileController struct {
	authService     service.AuthService
	seessionService service.SessionService
}

func NewProfileController(authService service.AuthService, sessionService service.SessionService) *profileController {
	return &profileController{
		authService:     authService,
		seessionService: sessionService,
	}
}

func (c *profileController) GetProfilePage(e echo.Context) error {
	var isLogin bool = false
	var ok bool
	isLogin, ok = e.Get("isLogin").(bool)
	if !ok || !isLogin {
		return profile.Profile(false, "", "").Render(e.Request().Context(), e.Response().Writer)
	}

	sess, _ := session.Get("sessionid", e)
	id := sess.Values["id"].(string)
	detail := c.seessionService.GetSessionDetail(id)
	locationId := ""
	if detail != nil {
		locationId = detail.CurrentLocation
	}

	userId := e.Get("userId").(string)
	name, err := c.authService.GetUserInfo(userId)
	log.Println(userId)
	if err != nil {
		return profile.Profile(false, "", "").Render(e.Request().Context(), e.Response().Writer)
	}
	return profile.Profile(true, name, locationId).Render(e.Request().Context(), e.Response().Writer)
}
