package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/zde37/Numeris-Task/internal/models"
)

type userRepoImpl struct {
	DBPool *pgxpool.Pool
}

// newUserRepoImpl creates a new instance of the userRepoImpl struct, which is used to interact with the user-related data in the database. 
func newUserRepoImpl(dbPool *pgxpool.Pool) *userRepoImpl {
	return &userRepoImpl{
		DBPool: dbPool,
	}
}
 
// CreateUser creates a new user in the database and returns the generated user ID. 
func (u *userRepoImpl) CreateUser(ctx context.Context, user models.User) (uuid.UUID, error) {
	query := `
		INSERT INTO users (user_id, username, email, password, first_name, last_name, profile_picture_url, phone_number, address)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING user_id
	` 
	err := u.DBPool.QueryRow(ctx, query, user.UserID, user.Username, user.Email, user.Password, user.FirstName, user.LastName, 
		user.ProfilePictureURL, user.PhoneNumber, user.Address).Scan(&user.UserID)
	if err != nil {
		return uuid.Nil, err
	}
	return user.UserID, nil
}

// AddCustomer creates a new customer in the database and returns the generated customer ID.
func (u *userRepoImpl) AddCustomer(ctx context.Context, customer models.Customer) (uuid.UUID, error) {
	query := `
        INSERT INTO customers (customer_id, name, email, phone_number, address)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING customer_id`

	err := u.DBPool.QueryRow(ctx, query,
		customer.CustomerID, customer.Name, customer.Email, customer.PhoneNumber,
		customer.Address).Scan(&customer.CustomerID)
	if err != nil {
		return uuid.Nil, err
	}
	return customer.CustomerID, nil
}
 
// AddPaymentMethod creates a new payment method for a user in the database and returns the generated payment method ID. 
func (u *userRepoImpl) AddPaymentMethod(ctx context.Context, paymentMethod models.UserPaymentMethod) (uuid.UUID, error) {
	query := `
		INSERT INTO user_payment_methods (payment_method_id, user_id, account_name, account_number, bank_name, bank_address, swift_code)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING payment_method_id
	`
	err := u.DBPool.QueryRow(ctx, query, paymentMethod.PaymentMethodID, paymentMethod.UserID, paymentMethod.AccountName, paymentMethod.AccountNumber,
		paymentMethod.BankName, paymentMethod.BankAddress, paymentMethod.SwiftCode).Scan(&paymentMethod.PaymentMethodID)
	if err != nil {
		return uuid.Nil, err
	}
	return paymentMethod.PaymentMethodID, nil
}
