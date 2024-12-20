package models

import (
	"time"

	"gorm.io/gorm"
)

type Order struct {
	ID             string              `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	UserID         string              `gorm:"type:uuid;not null" json:"user_id"`
	User           *User               `json:"user" gorm:"foreignKey:UserID"`
	AddressID      string              `gorm:"type:uuid;not null" json:"address_id"`
	Address        *Address            `json:"address" gorm:"foreignKey:AddressID"`
	Status         OrderStatus         `gorm:"type:varchar(255);not null" json:"status"`
	ShippingID     *string             `gorm:"type:varchar(255)" json:"shipping_id"`
	TrackingID     *string             `gorm:"type:varchar(255)" json:"tracking_id"`
	WaybillID      *string             `gorm:"type:varchar(255)" json:"waybill_id"`
	ShippingStatus OrderShippingStatus `gorm:"type:varchar(255)" json:"shipping_status"`
	ShippingPrice  *int                `gorm:"type:integer;" json:"shipping_price"`
	TotalPrice     int                 `gorm:"type:integer;not null" json:"total_price"`
	CourierCompany *string             `gorm:"type:varchar(255)" json:"courier_company"`
	CourierType    *string             `gorm:"type:varchar(255)" json:"courier_type"`
	CreatedAt      time.Time           `json:"created_at"`
	UpdatedAt      time.Time           `json:"updated_at"`
	DeletedAt      gorm.DeletedAt      `gorm:"index" json:"deleted_at"`
	OrderItems     []OrderItem         `json:"order_items" gorm:"foreignKey:OrderID"`
	PaymentID      *string             `gorm:"type:uuid" json:"payment_id"`
	Payment        *Payment            `json:"payment" gorm:"foreignKey:PaymentID"`
	MerchantID     string              `gorm:"type:uuid;not null" json:"merchant_id"`
	Merchant       *User               `json:"merchant" gorm:"foreignKey:MerchantID"`
}

type OrderStatus string
type OrderShippingStatus string

const (
	OrderStatusWaiting    OrderStatus = "waiting_for_payment"
	OrderStatusProcessing OrderStatus = "processing"
	OrderStatusShipped    OrderStatus = "shipped"
	OrderStatusDelivered  OrderStatus = "delivered"
	OrderStatusCancelled  OrderStatus = "cancelled"
	OrderStatusCompleted  OrderStatus = "completed"
)

const (
	ShippingStatusConfirmed       OrderShippingStatus = "confirmed"
	ShippingStatusScheduled       OrderShippingStatus = "scheduled"
	ShippingStatusAllocated       OrderShippingStatus = "allocated"
	ShippingStatusPickedUp        OrderShippingStatus = "picking_up"
	ShippingStatusPicked          OrderShippingStatus = "picked"
	ShippingStatusCancelled       OrderShippingStatus = "cancelled"
	ShippingStatusOnHold          OrderShippingStatus = "on_hold"
	ShippingStatusDroppingOff     OrderShippingStatus = "dropping_off"
	ShippingStatusReturnInTransit OrderShippingStatus = "return_in_transit"
	ShippingStatusReturned        OrderShippingStatus = "returned"
	ShippingStatusRejected        OrderShippingStatus = "rejected"
	ShippingStatusDisposed        OrderShippingStatus = "disposed"
	ShippingStatusCourierNotFound OrderShippingStatus = "courier_not_found"
	ShippingStatusDelivered       OrderShippingStatus = "delivered"
)
