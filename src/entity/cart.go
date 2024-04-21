package entity

type CartItemStep struct {
	Step    string   `json:"step"`
	Options []string `json:"options"`
}

type CartItem struct {
	MenuId     string          `json:"menuId"`
	MenuName   string          `json:"menuName"`
	ItemId     int             `json:"itemId"`
	Quantity   int             `json:"quantity"`
	TotalPrice int             `json:"totalPrice"`
	Steps      []*CartItemStep `json:"steps"`
}

type Cart struct {
	CartId       string      `json:"cartId"`
	UserId       string      `json:"userId"`
	LocationId   string      `json:"locationId"`
	LocationName string      `json:"locationName"`
	CartItems    []*CartItem `json:"items"`
}

type CartResponse struct {
	SuccessResponse
	Data *Cart `json:"data"`
}

type AddCartItemRequest struct {
	MenuId     string          `json:"menuId"`
	Quantity   int             `json:"quantity"`
	TotalPrice int             `json:"totalPrice"`
	Steps      []*CartItemStep `json:"steps"`
}

type AddCartItem struct {
	CartId    string `json:"cartId"`
	ItemId    int    `json:"itemId"`
	TotalItem int    `json:"totalItem"`
}

type AddCartItemResponse struct {
	SuccessResponse
	Data *AddCartItem `json:"data"`
}

type DeleteCartItem struct {
	CartId     string `json:"cartId"`
	TotalItem  int    `json:"totalItem"`
	TotalPrice int    `json:"totalPrice"`
}

type DeleteCartItemResponse struct {
	SuccessResponse
	Data *DeleteCartItem `json:"data"`
}

type CheckoutRequest struct {
	CartId  string `json:"cartId"`
	UserId  string `json:"userId"`
	Address string `json:"address"`
}
