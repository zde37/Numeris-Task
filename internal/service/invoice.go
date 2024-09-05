package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/zde37/Numeris-Task/internal/helpers"
	"github.com/zde37/Numeris-Task/internal/models"
	"github.com/zde37/Numeris-Task/internal/repository"
)

type invoiceServiceImpl struct {
	invoice repository.InvoiceRepository
}

// newInvoiceServiceImpl creates a new instance of the invoiceServiceImpl struct, which implements the InvoiceService interface.
// The invoiceServiceImpl struct is responsible for handling invoice-related operations, and it takes an InvoiceRepository
// implementation as a dependency.
func newInvoiceServiceImpl(invoice repository.InvoiceRepository) *invoiceServiceImpl {
	return &invoiceServiceImpl{
		invoice: invoice,
	}
}

// CreateInvoice creates a new invoice with the provided data. 
func (s *invoiceServiceImpl) CreateInvoice(ctx context.Context, data models.CreateInvoiceRequest) (uuid.UUID, error) {
	invoiceID := uuid.New()
	senderID, err := uuid.Parse(data.Invoice.SenderID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid sender id")
	}
	customerID, err := uuid.Parse(data.CustomerID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid customer id")
	}

	if err := helpers.ValidateInvoiceStatus(data.Invoice.Status); err != nil {
		return uuid.Nil, err
	}

	layout := "2006-01-02"
	issueDate, err := time.Parse(layout, data.Invoice.IssueDate)
	if err != nil {
		return uuid.Nil, fmt.Errorf("issue date has invalid date format")
	}
	dueDate, err := time.Parse(layout, data.Invoice.DueDate)
	if err != nil {
		return uuid.Nil, fmt.Errorf("due date has invalid date format")
	}
	invoice := models.Invoice{
		InvoiceID:          invoiceID,
		InvoiceNumber:      helpers.RandomNumber(1000000000, 9999999999),
		SenderID:           senderID,
		CustomerID:         customerID,
		IssueDate:          issueDate,
		DueDate:            dueDate,
		TotalAmount:        data.Invoice.TotalAmount,
		DiscountPercentage: data.Invoice.DiscountPercentage,
		DiscountedAmount:   data.Invoice.DiscountedAmount,
		FinalAmount:        data.Invoice.FinalAmount,
		Status:             data.Invoice.Status,
		Currency:           data.Invoice.Currency,
		Notes:              data.Invoice.Notes,
	}

	items := make([]models.InvoiceItem, 0)
	for _, item := range data.InvoiceItems {
		itemID := uuid.New()
		item := models.InvoiceItem{
			ItemID:      itemID,
			InvoiceID:   invoiceID,
			Name:        item.Name,
			Description: item.Description,
			Quantity:    item.Quantity,
			UnitPrice:   item.UnitPrice,
			TotalPrice:  item.TotalPrice,
		}
		items = append(items, item)
	}

	paymentInfoID := uuid.New()
	paymentMethodID, err := uuid.Parse(data.PaymentMethodID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid payment method id")
	}
	paymentInfo := models.PaymentInformation{
		PaymentInfoID:   paymentInfoID,
		InvoiceID:       invoiceID,
		PaymentMethodID: paymentMethodID,
	}

	return s.invoice.CreateInvoice(ctx, invoice, items, customerID, paymentInfo)
}

// GetInvoiceDetails retrieves the details of an invoice by the given invoice ID. 
func (s *invoiceServiceImpl) GetInvoiceDetails(ctx context.Context, invoiceID uuid.UUID) (*models.InvoiceDetails, error) {
	return s.invoice.GetInvoiceDetails(ctx, invoiceID)
}

// AddInvoiceActivity creates a new invoice activity record. 
func (s *invoiceServiceImpl) AddInvoiceActivity(ctx context.Context, activity models.AddInvoiceActivityRequest) (uuid.UUID, error) {
	invoiceID, err := uuid.Parse(activity.InvoiceID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid invoice id")
	}

	userID, err := uuid.Parse(activity.UserID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid user id")
	}

	return s.invoice.AddInvoiceActivity(ctx, models.InvoiceActivity{
		ActivityID:  uuid.New(),
		InvoiceID:   invoiceID,
		UserID:      userID,
		Title:       activity.Title,
		Description: activity.Description,
	})
}

// GetTotalByStatus retrieves the total amount and count of invoices with the given status.
func (s *invoiceServiceImpl) GetTotalByStatus(ctx context.Context, status models.InvoiceStatus) (totalAmount float64, count int, err error) {
	return s.invoice.GetTotalByStatus(ctx, status)
}

// GetRecentInvoices retrieves the most recent invoices for the given sender ID, paginated by the provided page and limit. 
func (s *invoiceServiceImpl) GetRecentInvoices(ctx context.Context, senderID uuid.UUID, page, limit int32) ([]models.Invoice, error) {
	offset := (page - 1) * limit
	return s.invoice.GetRecentInvoices(ctx, senderID, limit, offset)
}

// GetRecentActivities retrieves the most recent activities for the given user ID, paginated by the provided page and limit. 
func (s *invoiceServiceImpl) GetRecentActivities(ctx context.Context, userID uuid.UUID, page, limit int32) ([]models.RecentActivity, error) {
	offset := (page - 1) * limit
	return s.invoice.GetRecentActivities(ctx, userID, limit, offset)
}

// GetInvoiceActivities retrieves the invoice activities for the given user ID and invoice ID, paginated by the provided page and limit. 
func (s *invoiceServiceImpl) GetInvoiceActivities(ctx context.Context, userID, invoiceID uuid.UUID, page, limit int32) ([]models.InvoiceActivity, error) {
	offset := (page - 1) * limit
	return s.invoice.GetInvoiceActivities(ctx, userID, invoiceID, limit, offset)
}
