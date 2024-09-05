package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/zde37/Numeris-Task/internal/models"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user models.User) (uuid.UUID, error)
	AddCustomer(ctx context.Context, customer models.Customer) (uuid.UUID, error)
	AddPaymentMethod(ctx context.Context, paymentMethod models.UserPaymentMethod) (uuid.UUID, error)
}

type InvoiceRepository interface {
	GetTotalByStatus(ctx context.Context, status models.InvoiceStatus) (float64, int, error)
	CreateInvoice(ctx context.Context, invoice models.Invoice, items []models.InvoiceItem, customer uuid.UUID, paymentInfo models.PaymentInformation) (uuid.UUID, error)
	GetInvoiceDetails(ctx context.Context, invoiceID uuid.UUID) (*models.InvoiceDetails, error)
	AddInvoiceActivity(ctx context.Context, activity models.InvoiceActivity) (uuid.UUID, error)
	GetRecentInvoices(ctx context.Context, senderID uuid.UUID, limit, offset int32) ([]models.Invoice, error)
	GetRecentActivities(ctx context.Context, userID uuid.UUID, limit, offset int32) ([]models.RecentActivity, error)
	GetInvoiceActivities(ctx context.Context, userID, invoiceID uuid.UUID, limit, offset int32) ([]models.InvoiceActivity, error)
}

type Repository struct {
	User    UserRepository
	Invoice InvoiceRepository
}

// NewRepository creates a new Repository instance that provides access to the User and Invoice repositories.
// The Repository struct is the main entry point for interacting with the application's data storage.
// It takes a *pgxpool.Pool as a parameter, which is used to create the underlying repository implementations.
func NewRepository(dbPool *pgxpool.Pool) *Repository {
	return &Repository{
		User:    newUserRepoImpl(dbPool),
		Invoice: newInvoiceRepoImpl(dbPool),
	}
}
