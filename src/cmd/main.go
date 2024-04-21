package main

import (
	"fmt"

	"dev.chaiyapluek.cloud.final.frontend/src/config"
	"dev.chaiyapluek.cloud.final.frontend/src/controller"
	myMiddleware "dev.chaiyapluek.cloud.final.frontend/src/middleware"
	"dev.chaiyapluek.cloud.final.frontend/src/repository"
	"dev.chaiyapluek.cloud.final.frontend/src/service"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	cfg := config.LoadEnv()

	sessionService := service.NewInMemorySessionService()

	locationRepo := repository.NewLocationRepository(cfg.BackendURL)
	locationService := service.NewLocationService(locationRepo)
	locationController := controller.NewLocationController(locationService, sessionService)

	authRepo := repository.NewAuthRepository(cfg.BackendURL)
	authService := service.NewAuthService(authRepo, cfg.AccessTokenKey)
	authController := controller.NewAuthController(authService)

	cartRepo := repository.NewCartRepository(cfg.BackendURL)
	cartService := service.NewCartService(cartRepo)
	cartController := controller.NewCartController(cartService)

	orderController := controller.NewOrderController(locationService, sessionService, cartService)
	profileController := controller.NewProfileController(authService, sessionService)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Static("/static", "static")
	e.Use(session.Middleware(sessions.NewCookieStore([]byte(cfg.SessionKey))))
	e.Use(myMiddleware.NewJWTMiddleware(cfg.AccessTokenKey).Middleware)
	e.Use(myMiddleware.NewDefaultSessionMiddleware(sessionService).Middleware)
	e.GET("", locationController.GetLocations)
	e.GET("/location", locationController.GetLocations)
	e.GET("/location/:id", locationController.GetLocationMenu)
	e.GET("/location/:locationId/menus/:menuId", locationController.GetLocationItems)

	e.GET("/login", authController.GetLoginPage)
	e.GET("/register", authController.GetRegisterPage)
	e.GET("/logout", authController.Logout)
	e.POST("/login", authController.HandleLogin)
	e.POST("/login-attempt", authController.HandleLoginAttempt)
	e.POST("/register-attempt", authController.HandleRegisterAttempt)
	e.POST("/register", authController.HandleRegister)

	e.POST("/order", orderController.HandleOrderSubmit)
	e.PUT("/order", orderController.HandleOrderUpdate)

	e.GET("/cart", cartController.GetCartPage)
	e.DELETE("/cart/:cartId/items/:itemId", cartController.DeleteCartItem)

	e.GET("/profile", profileController.GetProfilePage)

	e.POST("/checkout", cartController.Checkout)

	e.Start(fmt.Sprintf("%s:%s", cfg.Host, cfg.Port))
}
