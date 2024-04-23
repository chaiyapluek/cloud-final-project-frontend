package repository

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"dev.chaiyapluek.cloud.final.frontend/src/entity"
)

type CartPepository interface {
	GetUserCart(userId string, locationId string) (*entity.Cart, error)
	AddCartItem(cartId string, item *entity.AddCartItemRequest) (*entity.AddCartItem, error)
	DeleteCartItem(cartId string, itemId int) (*entity.DeleteCartItem, error)
	Checkout(req *entity.CheckoutRequest) error
}

type cartRepository struct {
	backendURL string
	apiKey     string
}

func NewCartRepository(backendURL string, apiKey string) CartPepository {
	return &cartRepository{
		backendURL: backendURL,
		apiKey:     apiKey,
	}
}

func (r *cartRepository) GetUserCart(userId string, locationId string) (*entity.Cart, error) {
	endpoint := "/users/" + userId + "/carts?locationId=" + locationId
	httpReq, err := http.NewRequest("GET", r.backendURL+endpoint, nil)
	if err != nil {
		log.Println("cart repository get "+endpoint+" error", err)
		return nil, err
	}
	httpReq.Header.Set("X-API-KEY", r.apiKey)
	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		log.Println("cart repository get "+endpoint+" error", err)
		return nil, err
	}

	if resp.StatusCode > 299 {
		log.Println("cart repository get "+endpoint+" status code error", resp.StatusCode)
		var errorResp entity.ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errorResp); err != nil {
			log.Println("cart repository get "+endpoint+" decode ErrorResponse error", err)
			return nil, err
		}
		return nil, errors.New(errorResp.Message)
	}

	var cartResp entity.CartResponse
	if err := json.NewDecoder(resp.Body).Decode(&cartResp); err != nil {
		log.Println("cart repository get "+endpoint+" decode CartResponse error", err)
		return nil, err
	}

	return cartResp.Data, nil
}

func (r *cartRepository) AddCartItem(cartId string, item *entity.AddCartItemRequest) (*entity.AddCartItem, error) {
	endpoint := "/carts/" + cartId + "/items"
	reqBody, err := json.Marshal(item)
	if err != nil {
		log.Println("cart repository add item to cart marshal error", err)
		return nil, err
	}
	httpReq, err := http.NewRequest("POST", r.backendURL+endpoint, bytes.NewBuffer(reqBody))
	if err != nil {
		log.Println("cart repository add item to cart create request error", err)
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("X-API-KEY", r.apiKey)

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		log.Println("cart repository add item to cart error", err)
		return nil, err
	}

	if resp.StatusCode > 299 {
		log.Println("cart repository add item to cart status code error", resp.StatusCode)
		var errorResp entity.ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errorResp); err != nil {
			log.Println("cart repository add item to cart decode ErrorResponse error", err)
			return nil, err
		}
		return nil, errors.New(errorResp.Message)
	}

	var addCartItemResp entity.AddCartItemResponse
	if err := json.NewDecoder(resp.Body).Decode(&addCartItemResp); err != nil {
		log.Println("cart repository add item to cart decode AddCartItemResponse error", err)
		return nil, err
	}

	return addCartItemResp.Data, nil
}

func (r *cartRepository) DeleteCartItem(cartId string, itemId int) (*entity.DeleteCartItem, error) {
	endpoint := fmt.Sprintf("/carts/%s/items/%d", cartId, itemId)
	req, err := http.NewRequest(http.MethodDelete, r.backendURL+endpoint, nil)
	if err != nil {
		log.Println("cart repository delete item from cart create request error", err)
		return nil, err
	}
	req.Header.Set("X-API-KEY", r.apiKey)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("cart repository delete item from cart error", err)
		return nil, err
	}

	if resp.StatusCode > 299 {
		log.Println("cart repository delete item from cart status code error", resp.StatusCode)
		var errorResp entity.ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errorResp); err != nil {
			log.Println("cart repository delete item from cart decode ErrorResponse error", err)
			return nil, err
		}
		return nil, errors.New(errorResp.Message)
	}

	var deleteCartItemResp entity.DeleteCartItemResponse
	if err := json.NewDecoder(resp.Body).Decode(&deleteCartItemResp); err != nil {
		log.Println("cart repository delete item from cart decode DeleteCartItemResponse error", err)
		return nil, err
	}

	return deleteCartItemResp.Data, nil
}

func (r *cartRepository) Checkout(req *entity.CheckoutRequest) error {
	endpoint := "/checkout"

	reqBody, err := json.Marshal(req)
	if err != nil {
		log.Println("cart repository checkout marshal error", err)
		return err
	}

	httpReq, err := http.NewRequest(http.MethodPost, r.backendURL+endpoint, bytes.NewBuffer(reqBody))
	if err != nil {
		log.Println("cart repository checkout create request error", err)
		return err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("X-API-KEY", r.apiKey)

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		log.Println("cart repository checkout error", err)
		return err
	}

	if resp.StatusCode > 299 {
		log.Println("cart repository checkout status code error", resp.StatusCode)
		var errorResp entity.ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errorResp); err != nil {
			log.Println("cart repository checkout decode ErrorResponse error", err)
			return err
		}
		return errors.New(errorResp.Message)
	}

	return nil
}
