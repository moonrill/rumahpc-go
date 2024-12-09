package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/moonrill/rumahpc-api/config"
	"github.com/moonrill/rumahpc-api/internal/models"
	"github.com/moonrill/rumahpc-api/utils"
	"github.com/xendit/xendit-go/v6/invoice"
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
		totalPrice += order.TotalPrice
	}

	externalId := utils.GenerateShortExternalID()

	var invoiceItems []invoice.InvoiceItem
	var invoiceFees []invoice.InvoiceFee

	var user *models.User

	for index, order := range orders {
		invoiceFees = append(invoiceFees, invoice.InvoiceFee{
			Type:  fmt.Sprintf("Shipping Fee - Order #%d", index+1),
			Value: float32(*order.ShippingPrice),
		})

		if user == nil {
			user = order.User
		}
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
		Fees:  invoiceFees,
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
