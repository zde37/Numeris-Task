package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/zde37/Numeris-Task/internal/models"
	"github.com/zde37/Numeris-Task/internal/repository"
)

type UserService interface {
	CreateUser(ctx context.Context, data models.CreateUserRequest) (uuid.UUID, error)
	AddCustomer(ctx context.Context, data models.AddCustomerRequest) (uuid.UUID, error)
	AddPaymentMethod(ctx context.Context, data models.AddPaymentMethodRequest) (uuid.UUID, error)
}

type InvoiceService interface {
	CreateInvoice(ctx context.Context, data models.CreateInvoiceRequest) (uuid.UUID, error)
	GetInvoiceDetails(ctx context.Context, invoiceID uuid.UUID) (*models.InvoiceDetails, error)
	AddInvoiceActivity(ctx context.Context, activity models.AddInvoiceActivityRequest) (uuid.UUID, error)
	GetTotalByStatus(ctx context.Context, status models.InvoiceStatus) (totalAmount float64, count int, err error)
	GetRecentInvoices(ctx context.Context, senderID uuid.UUID, page, limit int32) ([]models.Invoice, error)
	GetRecentActivities(ctx context.Context, userID uuid.UUID, page, limit int32) ([]models.RecentActivity, error)
	GetInvoiceActivities(ctx context.Context, userID, invoiceID uuid.UUID, page, limit int32) ([]models.InvoiceActivity, error)
}

type Service struct {
	User    UserService
	Invoice InvoiceService
}

// NewService creates a new instance of the Service struct, which provides access to the
// UserService and InvoiceService implementations. The Service struct is the main entry
// point for interacting with the application's business logic.
func NewService(repo *repository.Repository) *Service {
	return &Service{
		User:    newUserServiceImpl(repo.User),
		Invoice: newInvoiceServiceImpl(repo.Invoice),
	}
}
