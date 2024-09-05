package models

import (
	"time"

	"github.com/google/uuid"
)

type InvoiceStatus string

const (
	InvoiceStatusPaid    InvoiceStatus = "paid"
	InvoiceStatusOverDue InvoiceStatus = "overdue"
	InvoiceStatusDraft   InvoiceStatus = "draft"
	InvoiceStatusPending InvoiceStatus = "pending"
)

type User struct {
	UserID            uuid.UUID `json:"user_id"`
	Username          string    `json:"username"`
	Email             string    `json:"email"`
	Password          string    `json:"password"`
	FirstName         string    `json:"first_name"`
	LastName          string    `json:"last_name"`
	ProfilePictureURL string    `json:"profile_picture_url"`
	PhoneNumber       string    `json:"phone_number"`
	Address           string    `json:"address"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type Customer struct {
	CustomerID  uuid.UUID `json:"user_id"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	PhoneNumber string    `json:"phone_number"`
	Address     string    `json:"address"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Invoice struct {
	InvoiceID          uuid.UUID `json:"invoice_id"`
	InvoiceNumber      string    `json:"invoice_number"`
	SenderID           uuid.UUID `json:"sender_id"`
	CustomerID         uuid.UUID `json:"customer_id"`
	IssueDate          time.Time `json:"issue_date"`
	DueDate            time.Time `json:"due_date"`
	TotalAmount        float64   `json:"total_amount"`
	DiscountPercentage float64   `json:"discount_percentage"`
	DiscountedAmount   float64   `json:"discounted_amount"`
	FinalAmount        float64   `json:"final_amount"`
	Status             string    `json:"status"`
	Currency           string    `json:"currency"`
	Notes              string    `json:"notes"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

type InvoiceItem struct {
	ItemID      uuid.UUID `json:"item_id"`
	InvoiceID   uuid.UUID `json:"invoice_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Quantity    int       `json:"quantity"`
	UnitPrice   float64   `json:"unit_price"`
	TotalPrice  float64   `json:"total_price"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type UserPaymentMethod struct {
	PaymentMethodID uuid.UUID `json:"payment_method_id"`
	UserID          uuid.UUID `json:"user_id"`
	AccountName     string    `json:"account_name"`
	AccountNumber   string    `json:"account_number"`
	BankName        string    `json:"bank_name"`
	BankAddress     string    `json:"bank_address"`
	SwiftCode       string    `json:"swift_code"`
	IsDefault       bool      `json:"is_default"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type PaymentInformation struct {
	PaymentInfoID   uuid.UUID `json:"payment_info_id"`
	InvoiceID       uuid.UUID `json:"invoice_id"`
	PaymentMethodID uuid.UUID `json:"payment_method_id"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type InvoiceDetails struct {
	Invoice             Invoice
	SenderName          string
	SenderEmail         string
	SenderPhoneNumber   string
	SenderAddress       string
	CustomerName        string
	CustomerEmail       string
	CustomerPhoneNumber string
	PaymentInformation  UserPaymentMethod
	Items               []InvoiceItem
	Activities          []InvoiceActivity
}

type InvoiceActivity struct {
	ActivityID  uuid.UUID `json:"activity_id"`
	InvoiceID   uuid.UUID `json:"invoice_id"`
	UserID      uuid.UUID `json:"user_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

type RecentActivity struct {
	ActivityID  uuid.UUID `json:"activity_id"`
	UserID      uuid.UUID `json:"user_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}
