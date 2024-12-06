package services

import (
	"bytes"
	"encoding/json"
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
		"couriers":                "gojek,jnt,jne,sicepat,anteraja",
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
		OriginPostalCode:      userAddress.ZipCode,
		DestinationPostalCode: merchantAddress.ZipCode,
		Items:                 items,
	})
}
