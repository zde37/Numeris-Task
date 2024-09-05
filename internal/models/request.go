package models

type AddInvoiceActivityRequest struct {
	InvoiceID   string `json:"invoice_id" binding:"required"`
	UserID      string `json:"user_id" binding:"required"`
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
}

type CreateUserRequest struct {
	Username          string `json:"username" binding:"required"`
	Email             string `json:"email" binding:"required"`
	Password          string `json:"password" binding:"required"`
	FirstName         string `json:"first_name" binding:"required"`
	LastName          string `json:"last_name" binding:"required"`
	ProfilePictureURL string `json:"profile_picture_url"`
	PhoneNumber       string `json:"phone_number"`
	Address           string `json:"address"`
}

type AddPaymentMethodRequest struct {
	UserID        string `json:"user_id"  binding:"required"`
	AccountName   string `json:"account_name"  binding:"required"`
	AccountNumber string `json:"account_number"  binding:"required"`
	BankName      string `json:"bank_name"  binding:"required"`
	BankAddress   string `json:"bank_address" binding:"required"`
	SwiftCode     string `json:"swift_code" binding:"required"`
}

type InvoiceInfo struct {
	SenderID           string  `json:"sender_id" binding:"required"`
	IssueDate          string  `json:"issue_date" binding:"required"`
	DueDate            string  `json:"due_date" binding:"required"`
	TotalAmount        float64 `json:"total_amount" binding:"required"`
	DiscountPercentage float64 `json:"discount_percentage" binding:"required"`
	DiscountedAmount   float64 `json:"discounted_amount" binding:"required"`
	FinalAmount        float64 `json:"final_amount" binding:"required"`
	Status             string  `json:"status" binding:"required"`
	Currency           string  `json:"currency" binding:"required"`
	Notes              string  `json:"notes" binding:"required"`
}

type AddCustomerRequest struct {
	Name        string `json:"name" binding:"required"`
	Email       string `json:"email" binding:"required"`
	PhoneNumber string `json:"phone_number" binding:"required"`
	Address     string `json:"address" binding:"required"`
}

type InvoiceItemDetails struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description" binding:"required"`
	Quantity    int     `json:"quantity" binding:"required"`
	UnitPrice   float64 `json:"unit_price" binding:"required"`
	TotalPrice  float64 `json:"total_price" binding:"required"`
}

type CreateInvoiceRequest struct {
	Invoice         InvoiceInfo          `json:"invoice" binding:"required"`
	CustomerID      string               `json:"customer_id" binding:"required"`
	PaymentMethodID string               `json:"payment_method_id" binding:"required"`
	InvoiceItems    []InvoiceItemDetails `json:"invoice_items" binding:"required"`
}
