package services

import (
	"errors"

	"github.com/google/uuid"
	"github.com/moonrill/rumahpc-api/config"
	"github.com/moonrill/rumahpc-api/internal/models"
	"github.com/moonrill/rumahpc-api/types"
	"gorm.io/gorm"
)

var ErrOrderProductNotFound = errors.New("order product not found")
var ErrOrderAddressNotFound = errors.New("order address not found")
var ErrOrderShipping = errors.New("error shipping order")
var ErrOrderCreate = errors.New("error creating order")

func CreateBuyNowOrder(request *types.BuyNowRequest, userID string) (*models.Order, error) {
	// Get User Address
	var userAddress models.Address

	err := config.DB.First(&userAddress, "id = ?", request.AddressID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrOrderAddressNotFound
		}
		return nil, err
	}

	var product models.Product
	err = config.DB.Preload("Merchant").Preload("Merchant.Addresses").First(&product, "id = ?", request.ProductID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrOrderProductNotFound
		}
		return nil, err
	}

	// Convert Products to Biteship Items
	var items []types.BiteshipItem

	items = append(items, types.BiteshipItem{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Value:       product.Price,
		Quantity:    request.Quantity,
		Weight:      product.Weight * 1000, // Convert to grams
	})

	// Get Merchant Address
	if product.Merchant == nil || product.Merchant.Addresses == nil {
		return nil, ErrOrderAddressNotFound
	}

	merchantAddresses := *product.Merchant.Addresses
	if len(merchantAddresses) == 0 {
		return nil, ErrOrderAddressNotFound
	}

	var merchantAddress models.Address
	err = config.DB.First(&merchantAddress, "id = ?", merchantAddresses[0].ID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrOrderAddressNotFound
		}
		return nil, err
	}

	// Create Shipping Order
	shippingOrder, err := CreateShippingOrder(&types.ShippingOrderRequest{
		OriginContactName:       merchantAddress.ContactName,
		OriginContactPhone:      merchantAddress.ContactNumber,
		OriginAddress:           ConvertAddressToString(&merchantAddress),
		OriginNote:              *merchantAddress.Note,
		OriginPostalCode:        merchantAddress.ZipCode,
		DestinationContactName:  userAddress.ContactName,
		DestinationContactPhone: userAddress.ContactNumber,
		DestinationContactEmail: product.Merchant.Email,
		DestinationAddress:      ConvertAddressToString(&userAddress),
		DestinationPostalCode:   userAddress.ZipCode,
		DestinationNote:         *userAddress.Note,
		CourierCompany:          request.CourierCompany,
		CourierType:             request.CourierType,
		DeliveryType:            "now",
		Items:                   items,
	})

	if err != nil {
		return nil, ErrOrderShipping
	}

	// Calculate Total Price
	totalPrice := 0
	for _, item := range items {
		totalPrice += item.Value * item.Quantity
	}
	totalPrice += shippingOrder.ShippingPrice

	// Create Order
	orderID := uuid.New().String()
	order := &models.Order{
		ID:             orderID,
		UserID:         userID,
		AddressID:      request.AddressID,
		ShippingID:     &shippingOrder.ShippingOrderID,
		Status:         models.OrderStatusWaiting,
		TrackingID:     &shippingOrder.TrackingID,
		WaybillID:      &shippingOrder.WaybillID,
		ShippingStatus: &shippingOrder.ShippingStatus,
		ShippingPrice:  &shippingOrder.ShippingPrice,
		TotalPrice:     totalPrice,
		CourierCompany: &request.CourierCompany,
		CourierType:    &request.CourierType,
	}

	err = config.DB.Create(order).Error
	if err != nil {
		return nil, ErrOrderCreate
	}

	// Create Order Items
	err = CreateOrderItems(orderID, items)
	if err != nil {
		return nil, ErrOrderCreate
	}

	return order, nil

}

func ConvertAddressToString(address *models.Address) string {
	combined := address.Province + ", " + address.City + ", " + address.District + ", " + address.Village + ", " + address.Address
	return combined
}

func CreateOrderItems(orderID string, items []types.BiteshipItem) error {
	for _, item := range items {
		orderItem := &models.OrderItem{
			OrderID:   orderID,
			ProductID: item.ID,
			Quantity:  item.Quantity,
			SubTotal:  item.Value * item.Quantity,
		}

		err := config.DB.Create(orderItem).Error
		if err != nil {
			return err
		}
	}

	return nil
}
