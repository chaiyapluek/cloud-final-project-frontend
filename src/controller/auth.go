package controller

import (
	"bytes"
	"net/http"
	"time"

	"dev.chaiyapluek.cloud.final.frontend/src/service"
	"dev.chaiyapluek.cloud.final.frontend/template/component/button"
	"dev.chaiyapluek.cloud.final.frontend/template/component/errorResponse"
	"dev.chaiyapluek.cloud.final.frontend/template/pages/login"
	"dev.chaiyapluek.cloud.final.frontend/template/pages/register"
	"github.com/labstack/echo/v4"
)

type authController struct {
	authService service.AuthService
}

func NewAuthController(authService service.AuthService) *authController {
	return &authController{
		authService: authService,
	}
}

func (c *authController) GetLoginPage(e echo.Context) error {
	e.Response().Header().Set("cache-control", "no-cache")

	isLogin, ok := e.Get("isLogin").(bool)
	if ok && isLogin {
		e.Redirect(http.StatusFound, "/location")
		return nil
	}

	return login.LoginAttempt().Render(e.Request().Context(), e.Response().Writer)
}

func (c *authController) GetRegisterPage(e echo.Context) error {
	e.Response().Header().Set("cache-control", "no-cache")

	isLogin, ok := e.Get("isLogin").(bool)
	if ok && isLogin {
		e.Redirect(http.StatusFound, "/location")
		return nil
	}

	return register.Register().Render(e.Request().Context(), e.Response().Writer)
}

type loginAttemptRequest struct {
	Email    string `form:"email"`
	Password string `form:"password"`
}

type loginRequest struct {
	AttemptId string `form:"attemptId"`
	Email     string `form:"email"`
	Code      string `form:"code"`
}

type registerAttemptRequest struct {
	Name            string `form:"display-name"`
	Email           string `form:"email"`
	Password        string `form:"password"`
	ConfirmPassword string `form:"confirm-password"`
}

type registerRequest struct {
	AttemptId string `form:"attemptId"`
	Email     string `form:"email"`
	Code      string `form:"code"`
}

func (c *authController) HandleLoginAttempt(e echo.Context) error {
	var req loginAttemptRequest
	if err := e.Bind(&req); err != nil {
		return err
	}

	attemptId, err := c.authService.LoginAttempt(req.Email, req.Password)
	if err != nil {
		textWriter := bytes.NewBufferString("")
		errorResponse.ErrorText(err.Error()).Render(e.Request().Context(), textWriter)
		buttonWriter := bytes.NewBufferString("")
		button.NextButton("Continue", true).Render(e.Request().Context(), buttonWriter)
		e.Response().Header().Set("hx-reswap", "innerHTML")
		return e.HTML(400, textWriter.String()+buttonWriter.String())
	}

	return login.LoginCode(attemptId, req.Email).Render(e.Request().Context(), e.Response().Writer)
}

func (c *authController) HandleLogin(e echo.Context) error {
	var req loginRequest
	if err := e.Bind(&req); err != nil {
		return err
	}
	accessToken, _, err := c.authService.Login(req.AttemptId, req.Email, req.Code)

	if err != nil {
		textWriter := bytes.NewBufferString("")
		errorResponse.ErrorText(err.Error()).Render(e.Request().Context(), textWriter)
		buttonWriter := bytes.NewBufferString("")
		button.NextButton("Login", true).Render(e.Request().Context(), buttonWriter)
		e.Response().Header().Set("hx-reswap", "innerHTML")
		return e.HTML(400, textWriter.String()+buttonWriter.String())
	}

	e.SetCookie(&http.Cookie{
		Name:     "accessToken",
		Value:    accessToken,
		Path:     "/",
		Expires:  time.Now().Add(time.Hour * 24),
		HttpOnly: true,
	})

	e.Response().Header().Set("hx-redirect", "/location")
	return nil
}

func (c *authController) HandleRegisterAttempt(e echo.Context) error {
	var req registerAttemptRequest
	if err := e.Bind(&req); err != nil {
		return err
	}

	if req.Password != req.ConfirmPassword {
		textWriter := bytes.NewBufferString("")
		errorResponse.ErrorText("Password does not match").Render(e.Request().Context(), textWriter)
		buttonWriter := bytes.NewBufferString("")
		button.NextButton("Register", true).Render(e.Request().Context(), buttonWriter)
		e.Response().Header().Set("hx-reswap", "innerHTML")
		return e.HTML(400, textWriter.String()+buttonWriter.String())
	}

	attemptId, err := c.authService.RegisterAttempt(req.Name, req.Email, req.Password)
	if err != nil {
		textWriter := bytes.NewBufferString("")
		errorResponse.ErrorText(err.Error()).Render(e.Request().Context(), textWriter)
		buttonWriter := bytes.NewBufferString("")
		button.NextButton("Continue", true).Render(e.Request().Context(), buttonWriter)
		e.Response().Header().Set("hx-reswap", "innerHTML")
		return e.HTML(400, textWriter.String()+buttonWriter.String())
	}

	return register.RegisterCode(attemptId, req.Email).Render(e.Request().Context(), e.Response().Writer)
}

func (c *authController) HandleRegister(e echo.Context) error {
	var req registerRequest
	if err := e.Bind(&req); err != nil {
		return err
	}

	accessToken, _, err := c.authService.Register(req.AttemptId, req.Email, req.Code)
	if err != nil {
		textWriter := bytes.NewBufferString("")
		errorResponse.ErrorText(err.Error()).Render(e.Request().Context(), textWriter)
		buttonWriter := bytes.NewBufferString("")
		button.NextButton("Register", true).Render(e.Request().Context(), buttonWriter)
		e.Response().Header().Set("hx-reswap", "innerHTML")
		return e.HTML(400, textWriter.String()+buttonWriter.String())
	}

	e.SetCookie(&http.Cookie{
		Name:     "accessToken",
		Value:    accessToken,
		Path:     "/",
		Expires:  time.Now().Add(time.Hour * 24),
		HttpOnly: true,
	})

	e.Response().Header().Set("hx-redirect", "/location")
	return nil
}

func (c *authController) Logout(e echo.Context) error {
	isLogin, ok := e.Get("isLogin").(bool)
	if !ok || !isLogin {
		return e.Redirect(http.StatusFound, "/login")
	}

	e.SetCookie(&http.Cookie{
		Name:     "accessToken",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
	})

	e.Response().Header().Set("cache-control", "no-cache")
	return e.Redirect(301, "/login")
}
