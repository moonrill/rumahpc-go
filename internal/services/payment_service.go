package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/moonrill/rumahpc-api/config"
	"github.com/moonrill/rumahpc-api/internal/models"
	"github.com/moonrill/rumahpc-api/types"
	"github.com/moonrill/rumahpc-api/utils"
	"github.com/xendit/xendit-go/v6/invoice"
	"github.com/xendit/xendit-go/v6/payout"
)

func CreateBuyNowXenditInvoice(order *models.Order, orderItem *models.OrderItem) (*invoice.Invoice, error) {
	description := "Payment for order " + order.ID
	totalPrice := order.TotalPrice + *order.ShippingPrice
	externalId := utils.GenerateShortExternalID()

	var product models.Product
	err := config.DB.Preload("Category").First(&product, "id = ?", orderItem.ProductID).Error
	if err != nil {
		return nil, err
	}

	invoiceRequest := invoice.CreateInvoiceRequest{
		ExternalId:  externalId,
		Amount:      float64(totalPrice),
		PayerEmail:  &order.User.Email,
		Description: &description,
		Customer: &invoice.CustomerObject{
			Id:           *invoice.NewNullableString(&order.User.ID),
			GivenNames:   *invoice.NewNullableString(&order.User.Name),
			Email:        *invoice.NewNullableString(&order.User.Email),
			MobileNumber: *invoice.NewNullableString(&order.User.PhoneNumber),
		},
		Items: []invoice.InvoiceItem{
			{
				ReferenceId: &product.ID,
				Name:        product.Name,
				Price:       float32(product.Price),
				Quantity:    float32(orderItem.Quantity),
				Category:    &product.Category.Name,
			},
		},
		Fees: []invoice.InvoiceFee{
			{
				Type:  "Shipping Fee",
				Value: float32(*order.ShippingPrice),
			},
		},
	}

	inv, r, xerr := config.XenditClient.InvoiceApi.CreateInvoice(context.Background()).
		CreateInvoiceRequest(invoiceRequest).
		Execute()

	if xerr != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InvoiceApi.CreateInvoice``: %v\n", xerr.Error())

		b, _ := json.Marshal(xerr.FullError())
		fmt.Fprintf(os.Stderr, "Full Error Struct: %v\n", string(b))

		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)

		return nil, xerr
	}

	if r.StatusCode != 200 {
		return nil, fmt.Errorf("failed to create invoice, status code: %d", r.StatusCode)
	}

	// Create Payment
	payment := models.Payment{
		ExternalID: externalId,
		UserID:     order.User.ID,
		Amount:     totalPrice,
		Status:     models.PaymentPending,
	}

	err = config.DB.Create(&payment).Error
	if err != nil {
		return nil, err
	}

	// Update Order
	order.PaymentID = &payment.ID
	err = config.DB.Save(order).Error
	if err != nil {
		return nil, err
	}

	return inv, nil
}

func CreateCartCheckoutXenditInvoice(orders []*models.Order, orderItems []*models.OrderItem) (*invoice.Invoice, error) {
	// Calculate total price accross all orders
	var totalPrice int
	for _, order := range orders {
		totalPrice += order.TotalPrice + *order.ShippingPrice
	}

	externalId := utils.GenerateShortExternalID()

	var invoiceItems []invoice.InvoiceItem
	var shippingTotal int

	var user *models.User

	for _, order := range orders {
		user = order.User
		shippingTotal += *order.ShippingPrice
	}

	for _, item := range orderItems {
		var product models.Product
		err := config.DB.Preload("Category").First(&product, "id = ?", item.ProductID).Error
		if err != nil {
			return nil, err
		}

		invoiceItems = append(invoiceItems, invoice.InvoiceItem{
			ReferenceId: &item.ProductID,
			Name:        product.Name,
			Price:       float32(product.Price),
			Quantity:    float32(item.Quantity),
			Category:    &product.Category.Name,
		})
	}

	// Ensure we have a user
	if user == nil {
		return nil, errors.New("no user associated with orders")
	}

	// Prepare invoice description
	description := fmt.Sprintf("Cart Checkout - %d Items", len(orderItems))

	invoiceRequest := invoice.CreateInvoiceRequest{
		ExternalId:  externalId,
		Amount:      float64(totalPrice),
		PayerEmail:  &user.Email,
		Description: &description,
		Customer: &invoice.CustomerObject{
			Id:           *invoice.NewNullableString(&user.ID),
			GivenNames:   *invoice.NewNullableString(&user.Name),
			Email:        *invoice.NewNullableString(&user.Email),
			MobileNumber: *invoice.NewNullableString(&user.PhoneNumber),
		},
		Items: invoiceItems,
		Fees: []invoice.InvoiceFee{
			{
				Type:  "Shipping Fee",
				Value: float32(shippingTotal),
			},
		},
	}

	inv, r, xerr := config.XenditClient.InvoiceApi.CreateInvoice(context.Background()).
		CreateInvoiceRequest(invoiceRequest).
		Execute()

	if xerr != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `InvoiceApi.CreateInvoice``: %v\n", xerr.Error())

		b, _ := json.Marshal(xerr.FullError())
		fmt.Fprintf(os.Stderr, "Full Error Struct: %v\n", string(b))

		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)

		return nil, xerr
	}

	if r.StatusCode != 200 {
		return nil, fmt.Errorf("failed to create invoice, status code: %d", r.StatusCode)
	}

	// Create payment record
	payment := models.Payment{
		ExternalID: externalId,
		UserID:     user.ID,
		Amount:     totalPrice,
		Status:     models.PaymentPending,
	}

	err := config.DB.Create(&payment).Error
	if err != nil {
		return nil, fmt.Errorf("failed to create payment record: %v", err)
	}

	// Update orders
	for _, order := range orders {
		order.PaymentID = &payment.ID
		if err := config.DB.Save(order).Error; err != nil {
			// Log error but don't stop the process
			fmt.Fprintf(os.Stderr, "Failed to update order with invoice ID: %v\n", err)
		}
	}

	return inv, nil
}

func HandleXenditCallback(callback *invoice.InvoiceCallback) error {
	// Get External ID
	externalId := callback.GetExternalId()

	// Search payment by external ID
	var payment models.Payment
	err := config.DB.First(&payment, "external_id = ?", externalId).Error

	if err != nil {
		return err
	}

	// Update payment
	if callback.GetStatus() == "PAID" {
		payment.Status = models.PaymentPaid
		payment.PaymentDate = callback.GetPaidAt()
		payment.PaymentMethod = callback.GetPaymentChannel()
	} else {
		payment.Status = models.PaymentExpired
	}

	err = config.DB.Save(&payment).Error
	if err != nil {
		return err
	}

	// Create Shipping Orders
	if payment.Status == models.PaymentPaid {
		var orders []models.Order
		err := config.DB.Preload("Merchant").Preload("User").Preload("Address").Where("payment_id = ?", payment.ID).Find(&orders).Error

		if err != nil {
			return err
		}

		for _, order := range orders {
			var biteShipItems []types.BiteshipItem
			var orderItems []models.OrderItem
			err := config.DB.Preload("Product").Where("order_id = ?", order.ID).Find(&orderItems).Error
			if err != nil {
				return err
			}

			// Get order items
			for _, item := range orderItems {
				var product models.Product
				err := config.DB.First(&product, "id = ?", item.ProductID).Error
				if err != nil {
					return err
				}

				// Append to biteShipItems
				biteShipItems = append(biteShipItems, types.BiteshipItem{
					Name:        product.Name,
					Description: product.Description,
					Value:       product.Price,
					Quantity:    item.Quantity,
					Weight:      product.Weight * 1000,
				})
			}

			// Get Merchant Address
			var merchantAddress models.Address
			err = config.DB.First(&merchantAddress, "user_id = ?", order.MerchantID).Error
			if err != nil {
				return err
			}

			// Create shipping order
			response, err := CreateShippingOrder(&types.ShippingOrderRequest{
				OriginContactName:       order.Merchant.Name,
				OriginContactPhone:      order.Merchant.PhoneNumber,
				OriginAddress:           ConvertAddressToString(&merchantAddress),
				OriginNote:              *order.Address.Note,
				OriginPostalCode:        merchantAddress.ZipCode,
				DestinationContactName:  order.User.Name,
				DestinationContactEmail: order.User.Email,
				DestinationContactPhone: order.User.PhoneNumber,
				DestinationAddress:      ConvertAddressToString(order.Address),
				DestinationNote:         *order.Address.Note,
				DestinationPostalCode:   order.Address.ZipCode,
				Items:                   biteShipItems,
				CourierCompany:          *order.CourierCompany,
				CourierType:             *order.CourierType,
				DeliveryType:            "now",
			})

			if err != nil {
				return err
			}

			// Update order
			order.Status = models.OrderStatusProcessing
			order.ShippingStatus = models.ShippingStatusConfirmed
			order.TrackingID = &response.TrackingID
			order.WaybillID = &response.WaybillID
			order.ShippingID = &response.ShippingOrderID
			err = config.DB.Save(&order).Error
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func CreatePayoutOrder(order *models.Order) error {
	var merchant models.User
	err := config.DB.First(&merchant, "id = ?", order.MerchantID).Error
	if err != nil {
		return utils.ErrNotFound
	}

	accoutnHolderName := *payout.NewNullableString(merchant.AccountName)

	payoutRequest := *payout.NewCreatePayoutRequest(
		order.ID,
		*merchant.PaymentChannel,
		payout.DigitalPayoutChannelProperties{
			AccountHolderName: accoutnHolderName,
			AccountNumber:     *merchant.AccountNumber,
		},
		float32(order.TotalPrice),
		"IDR",
	)

	res, r, respErr := config.XenditClient.PayoutApi.CreatePayout(context.Background()).
		IdempotencyKey(utils.GenerateIdempotencyKey()).
		CreatePayoutRequest(payoutRequest).
		Execute()

	// Handle API call errors
	if respErr != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `PayoutApi.CreatePayout``: %v\n", err)

		b, _ := json.Marshal(err)
		fmt.Fprintf(os.Stderr, "Full Error Struct: %v\n", string(b))

		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}

	// Response logging
	fmt.Fprintf(os.Stdout, "Response from `PayoutApi.CreatePayout`: %v\n", res)
	return nil
}
