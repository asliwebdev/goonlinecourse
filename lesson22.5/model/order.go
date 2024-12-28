package model

import "time"

type PaymentMethod string

const (
	PaymentMethodCard   PaymentMethod = "CARD"
	PaymentMethodPayPal PaymentMethod = "PAYPAL"
	PaymentMethodCash   PaymentMethod = "CASH"
)

type Order struct {
	ID              string        `json:"id"`
	CustomerID      string        `json:"customer_id"`
	Products        []string      `json:"products"`
	TotalAmount     float64       `json:"total_amount"`
	Status          string        `json:"status"`
	ShippingAddress string        `json:"shipping_address"`
	PaymentMethod   PaymentMethod `json:"payment_method"`
	CreatedAt       time.Time     `json:"created_at"`
	UpdatedAt       time.Time     `json:"updated_at"`
}
