package controller

import (
	"bytes"
	"log"
	"net/http"
	"sort"
	"strconv"

	"dev.chaiyapluek.cloud.final.frontend/src/service"
	"dev.chaiyapluek.cloud.final.frontend/template/pages/cart"
	"github.com/labstack/echo/v4"
)

type cartController struct {
	cartService service.CartService
}

func NewCartController(cartService service.CartService) *cartController {
	return &cartController{
		cartService: cartService,
	}
}

func (c *cartController) GetCartPage(e echo.Context) error {

	isLogin, ok := e.Get("isLogin").(bool)
	if !ok || !isLogin {
		return e.Redirect(http.StatusFound, "/login")
	}

	locationId := e.QueryParam("locationId")
	if locationId == "" {
		return e.HTML(http.StatusBadRequest, "Bad request")
	}

	userId := e.Get("userId").(string)
	cartProps, err := c.cartService.GetUserCart(userId, locationId)
	if err != nil {
		return e.String(400, "Error")
	}

	for _, item := range cartProps.CartItems {
		sort.Slice(item.Steps, func(i, j int) bool {
			return item.Steps[i].Step < item.Steps[j].Step
		})
	}

	return cart.Cart(*cartProps, len(cartProps.CartItems) > 0).Render(e.Request().Context(), e.Response().Writer)
}

func (c *cartController) DeleteCartItem(e echo.Context) error {
	isLogin, ok := e.Get("isLogin").(bool)
	if !ok || !isLogin {
		return e.Redirect(http.StatusFound, "/login")
	}

	cartId := e.Param("cartId")
	itemId := e.Param("itemId")
	log.Println("delete cart item", cartId, itemId)
	itemIdInt, err := strconv.Atoi(itemId)
	if err != nil {
		return e.String(http.StatusBadRequest, "Invalid item id")
	}

	totalPrice, totalItem, err := c.cartService.DeleteCartItem(cartId, itemIdInt)
	if err != nil {
		return e.String(http.StatusBadRequest, "Error")
	}

	priceWriter := bytes.NewBufferString("")
	cart.Price(totalPrice, totalPrice, totalItem).Render(e.Request().Context(), priceWriter)
	buttonWriter := bytes.NewBufferString("")
	cart.SubmitButton(totalItem > 0).Render(e.Request().Context(), buttonWriter)

	return e.HTML(http.StatusOK, priceWriter.String()+buttonWriter.String())
}

type CheckoutRequest struct {
	CartId  string `form:"cartId"`
	Address string `form:"address"`
}

func (c *cartController) Checkout(e echo.Context) error {
	isLogin, ok := e.Get("isLogin").(bool)
	if !ok || !isLogin {
		return e.Redirect(http.StatusFound, "/login")
	}

	var req CheckoutRequest
	if err := e.Bind(&req); err != nil {
		return e.String(http.StatusBadRequest, "Invalid request")
	}

	userId := e.Get("userId").(string)

	err := c.cartService.Checkout(req.CartId, userId, req.Address)
	if err != nil {
		return e.String(http.StatusBadRequest, "Error")
	}

	return cart.CheckoutResponse().Render(e.Request().Context(), e.Response().Writer)
}
