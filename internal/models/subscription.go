// internal/models/subscription.go
package models

import "time"

type Subscription struct {
	ID              int64     `json:"id"`
	UserID          int64     `json:"user_id"`
	Name            string    `json:"name"`
	Cost            float64   `json:"cost"`
	Currency        string    `json:"currency"`
	BillingCycle    string    `json:"billing_cycle"` // monthly, yearly, one-time
	NextPaymentDate time.Time `json:"next_payment_date"`
	Category        string    `json:"category,omitempty"`
	IsActive        bool      `json:"is_active"`
	CreatedAt       time.Time `json:"created_at"`
}

// CreateSubscriptionRequest - валидация при создании
type CreateSubscriptionRequest struct {
	Name            string  `json:"name" binding:"required"`
	Cost            float64 `json:"cost" binding:"required,min=0"`
	Currency        string  `json:"currency" binding:"required,len=3"`
	BillingCycle    string  `json:"billing_cycle" binding:"required,oneof=monthly yearly one-time"`
	NextPaymentDate string  `json:"next_payment_date" binding:"required"` // RFC3339: "2024-06-01"
	Category        string  `json:"category,omitempty"`
}
