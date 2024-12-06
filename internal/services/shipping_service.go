package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/moonrill/rumahpc-api/config"
	"github.com/moonrill/rumahpc-api/internal/models"
	"github.com/moonrill/rumahpc-api/types"
	"github.com/moonrill/rumahpc-api/utils"
	"gorm.io/gorm"
)

func GetCouriersRates(data *types.CourierRatesRequest) (*map[string]interface{}, error) {
	url := os.Getenv("BITESHIP_API_URL") + "/rates/couriers"

	body := map[string]interface{}{
		"origin_postal_code":      data.OriginPostalCode,
		"destination_postal_code": data.DestinationPostalCode,
		"couriers":                "gojek,jnt,jne,sicepat,anteraja,paxel,tiki",
		"items":                   data.Items,
	}

	jsonData, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	utils.SetBiteshipHeaders(req)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return nil, err
	}

	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return nil, err
	}

	var rates map[string]interface{}
	err = json.Unmarshal(bodyBytes, &rates)

	if err != nil {
		fmt.Println("Error unmarshalling response body:", err)
		return nil, err
	}

	return &rates, nil
}

func GetBuyNowCouriersRates(request *types.BuyNowCouriersRatesRequest) (*map[string]interface{}, error) {
	// Get User Address
	var userAddress models.Address

	err := config.DB.First(&userAddress, "id = ?", request.AddressID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, utils.ErrNotFound
		}
		return nil, err
	}

	var product models.Product
	err = config.DB.Preload("Merchant").Preload("Merchant.Addresses").First(&product, "id = ?", request.ProductID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, utils.ErrNotFound
		}
		return nil, err
	}

	// Convert Products to Biteship Items
	var items []types.BiteshipItem

	items = append(items, types.BiteshipItem{
		Name:        product.Name,
		Description: product.Description,
		Value:       product.Price,
		Quantity:    request.Quantity,
		Weight:      product.Weight * 1000, // Convert to grams
	})

	// Get Merchant Address
	if product.Merchant == nil || product.Merchant.Addresses == nil {
		return nil, utils.ErrNotFound
	}

	merchantAddresses := *product.Merchant.Addresses
	if len(merchantAddresses) == 0 {
		return nil, utils.ErrNotFound
	}

	var merchantAddress models.Address
	err = config.DB.First(&merchantAddress, "id = ?", merchantAddresses[0].ID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, utils.ErrNotFound
		}
		return nil, err
	}

	// Get Couriers Rates
	return GetCouriersRates(&types.CourierRatesRequest{
		OriginPostalCode:      merchantAddress.ZipCode,
		DestinationPostalCode: userAddress.ZipCode,
		Items:                 items,
	})
}

func GetCartCouriersRates(request *types.CartCouriersRatesRequest) (*map[string]interface{}, error) {
	// Get User Address
	var userAddress models.Address

	err := config.DB.First(&userAddress, "id = ?", request.AddressID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, utils.ErrNotFound
		}
		return nil, err
	}

	var merchantAddress models.Address

	err = config.DB.First(&merchantAddress, "user_id = ?", request.MerchantID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, utils.ErrNotFound
		}
		return nil, err
	}

	// Convert Cart Items to Biteship Items
	var items []types.BiteshipItem

	for _, cartItemID := range request.CartItems {
		var cartItem models.CartItem
		err = config.DB.Preload("Product").First(&cartItem, "id = ?", cartItemID).Error

		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return nil, utils.ErrNotFound
			}
			return nil, err
		}

		items = append(items, types.BiteshipItem{
			Name:        cartItem.Product.Name,
			Description: cartItem.Product.Description,
			Value:       cartItem.Product.Price,
			Quantity:    cartItem.Quantity,
			Weight:      cartItem.Product.Weight * 1000, // Convert to grams
		})
	}

	// Get Couriers Rates
	return GetCouriersRates(&types.CourierRatesRequest{
		OriginPostalCode:      userAddress.ZipCode,
		DestinationPostalCode: merchantAddress.ZipCode,
		Items:                 items,
	})
}

func CreateShippingOrder(request *types.ShippingOrderRequest) (*types.ShippingOrderSuccessResponse, error) {
	url := os.Getenv("BITESHIP_API_URL") + "/orders"

	body := map[string]interface{}{
		"origin_contact_name":       request.OriginContactName,
		"origin_contact_phone":      request.OriginContactPhone,
		"origin_address":            request.OriginAddress,
		"origin_note":               request.OriginNote,
		"origin_postal_code":        request.OriginPostalCode,
		"destination_contact_name":  request.DestinationContactName,
		"destination_contact_phone": request.DestinationContactPhone,
		"destination_contact_email": request.DestinationContactEmail,
		"destination_address":       request.DestinationAddress,
		"destination_postal_code":   request.DestinationPostalCode,
		"destination_note":          request.DestinationNote,
		"courier_company":           request.CourierCompany,
		"courier_type":              request.CourierType,
		"delivery_type":             request.DeliveryType,
		"items":                     request.Items,
	}

	jsonData, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	utils.SetBiteshipHeaders(req)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("error creating shipping order: %s", resp.Status)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return nil, err
	}

	var result map[string]interface{}
	err = json.Unmarshal(bodyBytes, &result)

	if err != nil {
		return nil, err
	}

	courier, ok := result["courier"].(map[string]interface{})
	if !ok {
		// handle the case when "courier" is not a map
		return nil, errors.New("invalid courier data")
	}

	response := &types.ShippingOrderSuccessResponse{
		ShippingOrderID: result["id"].(string),
		TrackingID:      courier["tracking_id"].(string),
		WaybillID:       courier["waybill_id"].(string),
		ShippingStatus:  result["status"].(string),
		ShippingPrice:   int(result["price"].(float64)),
	}

	return response, nil
}
