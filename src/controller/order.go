package controller

import (
	"strconv"

	"dev.chaiyapluek.cloud.final.frontend/src/entity"
	"dev.chaiyapluek.cloud.final.frontend/src/service"
	"dev.chaiyapluek.cloud.final.frontend/template/pages/order"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

type orderController struct {
	locationService service.LocationService
	sessionService  service.SessionService
	cartService     service.CartService
}

func NewOrderController(locationService service.LocationService, sessionService service.SessionService, cartService service.CartService) *orderController {
	return &orderController{
		locationService: locationService,
		sessionService:  sessionService,
		cartService:     cartService,
	}
}

type OrderSession struct {
	Order       order.OrderProps
	Preferences map[string]interface{}
}

func (c *orderController) HandleOrderSubmit(e echo.Context) error {

	isLogin, ok := e.Get("isLogin").(bool)
	if !ok || !isLogin {
		return e.Redirect(302, "/login")
	}

	e.Request().ParseForm()
	sess, _ := session.Get("sessionid", e)
	id := sess.Values["id"].(string)
	o := c.sessionService.GetSessionDetail(id)
	if o == nil {
		return e.String(400, "Invalid session")
	}

	mapp := map[string]map[string]int{}
	for _, v := range o.CurrentMenu.Steps {
		mapp[v.FormName] = map[string]int{}
		for _, i := range v.Items {
			mapp[v.FormName][i.Value] = i.Price
		}
	}

	// calculate price
	price := o.CurrentMenu.MenuPrice
	for key, val := range mapp {
		formVals, ok := e.Request().Form[key]
		if ok {
			for _, v := range formVals {
				p, ok := val[v]
				if ok {
					price += p
				}
			}
		}
	}

	// check form completion
	isComplete := true
	for _, g := range o.CurrentMenu.Steps {
		formVals, ok := e.Request().Form[g.FormName]
		if !ok && g.Required && g.Min > 0 {
			isComplete = false
			break
		}

		if g.Required && (len(formVals) < g.Min || len(formVals) > g.Max) {
			isComplete = false
			break
		}
	}

	isValid := true
	quantity := 1
	if _, ok := e.Request().Form["quantity"]; ok {
		q := e.Request().FormValue("quantity")
		var qty int
		var err error
		qty, err = strconv.Atoi(q)
		if err != nil {
			isValid = false
		}
		if qty < 0 {
			isValid = false
		}
		quantity = qty
	}

	if !isValid || quantity == 0 {
		return e.HTML(400, "Invalid order")
	}

	if !isComplete {
		return e.HTML(400, "Incomplete order")
	}

	steps := []*entity.CartItemStep{}
	for key, val := range mapp {
		formVals, ok := e.Request().Form[key]
		if ok {
			options := []string{}
			for _, v := range formVals {
				_, ok := val[v]
				if ok {
					options = append(options, v)
				}
			}
			steps = append(steps, &entity.CartItemStep{
				Step:    key,
				Options: options,
			})
		}
	}

	userId := e.Get("userId").(string)
	_, err := c.cartService.AddCartItem(userId, o.CurrentMenu.LocationId, &entity.AddCartItemRequest{
		MenuId:     o.CurrentMenu.MenuId,
		Quantity:   quantity,
		TotalPrice: price,
		Steps:      steps,
	})
	if err != nil {
		return e.String(400, err.Error())
	}

	e.Response().Header().Set("hx-redirect", "/location/"+o.CurrentMenu.LocationId)
	return nil
}

func (c *orderController) HandleOrderUpdate(e echo.Context) error {

	isLogin, ok := e.Get("isLogin").(bool)
	if !ok || !isLogin {
		return order.PleaseLoginButton().Render(e.Request().Context(), e.Response().Writer)
	}

	e.Request().ParseForm()

	sess, _ := session.Get("sessionid", e)
	id := sess.Values["id"].(string)
	o := c.sessionService.GetSessionDetail(id)
	if o == nil {
		return e.String(400, "Invalid session")
	}

	// calculate price
	price := o.CurrentMenu.MenuPrice
	for _, g := range o.CurrentMenu.Steps {
		for _, i := range g.Items {
			if _, ok := e.Request().Form[g.FormName]; ok {
				for _, v := range e.Request().Form[g.FormName] {
					if v == i.Value {
						price += i.Price
					}
				}
			}
		}
	}

	// update preferences
	for _, g := range o.CurrentMenu.Steps {
		if _, ok := e.Request().Form[g.FormName]; ok {
			o.Preferences[g.FormName] = e.Request().Form[g.FormName]
		}
	}

	// check form completion
	isComplete := true
	for _, g := range o.CurrentMenu.Steps {
		formVals, ok := e.Request().Form[g.FormName]
		if !ok && g.Required && g.Min > 0 {
			isComplete = false
			break
		}

		if g.Required && (len(formVals) < g.Min || len(formVals) > g.Max) {
			isComplete = false
			break
		}
	}

	// quantity
	if _, ok := e.Request().Form["quantity"]; ok {
		q := e.Request().FormValue("quantity")
		var qty int
		var err error
		qty, err = strconv.Atoi(q)
		if err != nil {
			qty = 1
		}
		if qty < 0 {
			qty = 0
		}
		o.CurrentMenu.Quantity = qty
		price *= qty
	}

	sess.Values["preferences"] = o
	sess.Save(e.Request(), e.Response())

	if o.CurrentMenu.Quantity == 0 {
		return order.BackButton().Render(e.Request().Context(), e.Response().Writer)
	}
	if isComplete {
		return order.CompleteButton(price).Render(e.Request().Context(), e.Response().Writer)
	}
	return order.IncompleteButton(price).Render(e.Request().Context(), e.Response().Writer)
}
