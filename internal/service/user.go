package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/zde37/Numeris-Task/internal/helpers"
	"github.com/zde37/Numeris-Task/internal/models"
	"github.com/zde37/Numeris-Task/internal/repository"
)

type userServiceImpl struct {
	User repository.UserRepository
}

// newUserServiceImpl creates a new instance of the userServiceImpl struct,
// which implements the UserService interface. It takes a UserRepository
// as a dependency and returns a pointer to the userServiceImpl struct.
func newUserServiceImpl(user repository.UserRepository) *userServiceImpl {
	return &userServiceImpl{
		User: user,
	}
}

// CreateUser creates a new user in the user repository with the provided data. 
func (u *userServiceImpl) CreateUser(ctx context.Context, data models.CreateUserRequest) (uuid.UUID, error) {
	hashedPassword, err := helpers.HashPassword(data.Password)
	if err != nil {
		return uuid.Nil, err
	}

	return u.User.CreateUser(ctx, models.User{
		UserID:            uuid.New(),
		Username:          data.Username,
		Email:             data.Email,
		Password:          hashedPassword,
		FirstName:         data.FirstName,
		LastName:          data.LastName,
		ProfilePictureURL: data.ProfilePictureURL,
		PhoneNumber:       data.PhoneNumber,
		Address:           data.Address,
	})
}

// AddPaymentMethod creates a new payment method for the specified user in the user repository. 
func (u *userServiceImpl) AddPaymentMethod(ctx context.Context, data models.AddPaymentMethodRequest) (uuid.UUID, error) {
	userID, err := uuid.Parse(data.UserID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid invoice id")
	}

	return u.User.AddPaymentMethod(ctx, models.UserPaymentMethod{
		PaymentMethodID: uuid.New(),
		UserID:          userID,
		AccountName:     data.AccountName,
		AccountNumber:   data.AccountNumber,
		BankName:        data.BankName,
		BankAddress:     data.BankAddress,
		SwiftCode:       data.SwiftCode,
	})
}

// AddCustomer creates a new customer in the user repository with the provided data. 
func (u *userServiceImpl) AddCustomer(ctx context.Context, data models.AddCustomerRequest) (uuid.UUID, error) {
	return u.User.AddCustomer(ctx, models.Customer{
		CustomerID:  uuid.New(),
		Name:        data.Name,
		Email:       data.Email,
		PhoneNumber: data.PhoneNumber,
		Address:     data.Address,
	})
}
