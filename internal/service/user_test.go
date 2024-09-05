package service

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	mocked "github.com/zde37/Numeris-Task/internal/mock"
	"github.com/zde37/Numeris-Task/internal/models"
	"go.uber.org/mock/gomock"
)

func TestCreateUser(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocked.NewMockUserRepository(ctrl)

	t.Run("successful user creation", func(t *testing.T) {
		expectedUserID := uuid.New()
		createUserRequest := models.CreateUserRequest{
			Username:          "testuser",
			Email:             "test@example.com",
			Password:          "password123",
			FirstName:         "Test",
			LastName:          "User",
			ProfilePictureURL: "http://example.com/profile.jpg",
			PhoneNumber:       "1234567890",
			Address:           "123 Test St",
		}

		repo.EXPECT().
			CreateUser(gomock.Any(), gomock.Any()).
			DoAndReturn(func(_ context.Context, user models.User) (uuid.UUID, error) {
				require.Equal(t, createUserRequest.Username, user.Username)
				require.Equal(t, createUserRequest.Email, user.Email)
				require.NotEqual(t, createUserRequest.Password, user.Password)
				require.Equal(t, createUserRequest.FirstName, user.FirstName)
				require.Equal(t, createUserRequest.LastName, user.LastName)
				require.Equal(t, createUserRequest.ProfilePictureURL, user.ProfilePictureURL)
				require.Equal(t, createUserRequest.PhoneNumber, user.PhoneNumber)
				require.Equal(t, createUserRequest.Address, user.Address)
				return expectedUserID, nil
			})

		service := newUserServiceImpl(repo)
		userID, err := service.CreateUser(ctx, createUserRequest)
		require.NoError(t, err)
		require.Equal(t, expectedUserID, userID)
	})

	t.Run("repository error", func(t *testing.T) {
		createUserRequest := models.CreateUserRequest{
			Username: "testuser",
			Email:    "test@example.com",
			Password: "password123",
		}

		expectedError := errors.New("database error")
		repo.EXPECT().
			CreateUser(gomock.Any(), gomock.Any()).
			Return(uuid.Nil, expectedError)

		service := newUserServiceImpl(repo)
		userID, err := service.CreateUser(ctx, createUserRequest)
		require.Error(t, err)
		require.Equal(t, uuid.Nil, userID)
		require.Equal(t, expectedError, err)
	})

	t.Run("duplicate username", func(t *testing.T) {
		createUserRequest := models.CreateUserRequest{
			Username: "existinguser",
			Email:    "test@example.com",
			Password: "password123",
		}

		expectedError := errors.New("duplicate username")
		repo.EXPECT().
			CreateUser(gomock.Any(), gomock.Any()).
			Return(uuid.Nil, expectedError)

		service := newUserServiceImpl(repo)
		userID, err := service.CreateUser(ctx, createUserRequest)
		require.Error(t, err)
		require.Equal(t, uuid.Nil, userID)
		require.Equal(t, expectedError, err)
	})

	t.Run("duplicate email", func(t *testing.T) {
		createUserRequest := models.CreateUserRequest{
			Username: "testuser",
			Email:    "existing@example.com",
			Password: "password123",
		}

		expectedError := errors.New("duplicate email")
		repo.EXPECT().
			CreateUser(gomock.Any(), gomock.Any()).
			Return(uuid.Nil, expectedError)

		service := newUserServiceImpl(repo)
		userID, err := service.CreateUser(ctx, createUserRequest)
		require.Error(t, err)
		require.Equal(t, uuid.Nil, userID)
		require.Equal(t, expectedError, err)
	})
}

func TestAddPaymentMethod(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocked.NewMockUserRepository(ctrl)

	t.Run("successful payment method addition", func(t *testing.T) {
		expectedPaymentMethodID := uuid.New()
		validUserID := uuid.New()
		addPaymentMethodRequest := models.AddPaymentMethodRequest{
			UserID:        validUserID.String(),
			AccountName:   "John Doe",
			AccountNumber: "1234567890",
			BankName:      "Test Bank",
			BankAddress:   "123 Bank St",
			SwiftCode:     "TESTSWIFT",
		}

		repo.EXPECT().
			AddPaymentMethod(gomock.Any(), gomock.Any()).
			DoAndReturn(func(_ context.Context, paymentMethod models.UserPaymentMethod) (uuid.UUID, error) {
				require.Equal(t, validUserID, paymentMethod.UserID)
				require.Equal(t, addPaymentMethodRequest.AccountName, paymentMethod.AccountName)
				require.Equal(t, addPaymentMethodRequest.AccountNumber, paymentMethod.AccountNumber)
				require.Equal(t, addPaymentMethodRequest.BankName, paymentMethod.BankName)
				require.Equal(t, addPaymentMethodRequest.BankAddress, paymentMethod.BankAddress)
				require.Equal(t, addPaymentMethodRequest.SwiftCode, paymentMethod.SwiftCode)
				return expectedPaymentMethodID, nil
			})

		service := newUserServiceImpl(repo)
		paymentMethodID, err := service.AddPaymentMethod(ctx, addPaymentMethodRequest)
		require.NoError(t, err)
		require.Equal(t, expectedPaymentMethodID, paymentMethodID)
	})

	t.Run("repository error", func(t *testing.T) {
		validUserID := uuid.New()
		addPaymentMethodRequest := models.AddPaymentMethodRequest{
			UserID:        validUserID.String(),
			AccountName:   "John Doe",
			AccountNumber: "1234567890",
			BankName:      "Test Bank",
			BankAddress:   "123 Bank St",
			SwiftCode:     "TESTSWIFT",
		}

		expectedError := errors.New("database error")
		repo.EXPECT().
			AddPaymentMethod(gomock.Any(), gomock.Any()).
			Return(uuid.Nil, expectedError)

		service := newUserServiceImpl(repo)
		paymentMethodID, err := service.AddPaymentMethod(ctx, addPaymentMethodRequest)
		require.Error(t, err)
		require.Equal(t, uuid.Nil, paymentMethodID)
		require.Equal(t, expectedError, err)
	})

	t.Run("empty account name", func(t *testing.T) {
		validUserID := uuid.New()
		addPaymentMethodRequest := models.AddPaymentMethodRequest{
			UserID:        validUserID.String(),
			AccountName:   "",
			AccountNumber: "1234567890",
			BankName:      "Test Bank",
			BankAddress:   "123 Bank St",
			SwiftCode:     "TESTSWIFT",
		}

		repo.EXPECT().
			AddPaymentMethod(gomock.Any(), gomock.Any()).
			Return(uuid.Nil, errors.New("account name cannot be empty"))

		service := newUserServiceImpl(repo)
		paymentMethodID, err := service.AddPaymentMethod(ctx, addPaymentMethodRequest)
		require.Error(t, err)
		require.Equal(t, uuid.Nil, paymentMethodID)
		require.Contains(t, err.Error(), "account name cannot be empty")
	})

	t.Run("empty account number", func(t *testing.T) {
		validUserID := uuid.New()
		addPaymentMethodRequest := models.AddPaymentMethodRequest{
			UserID:        validUserID.String(),
			AccountName:   "John Doe",
			AccountNumber: "",
			BankName:      "Test Bank",
			BankAddress:   "123 Bank St",
			SwiftCode:     "TESTSWIFT",
		}

		repo.EXPECT().
			AddPaymentMethod(gomock.Any(), gomock.Any()).
			Return(uuid.Nil, errors.New("account number cannot be empty"))

		service := newUserServiceImpl(repo)
		paymentMethodID, err := service.AddPaymentMethod(ctx, addPaymentMethodRequest)
		require.Error(t, err)
		require.Equal(t, uuid.Nil, paymentMethodID)
		require.Contains(t, err.Error(), "account number cannot be empty")
	})
}

func TestAddCustomer(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocked.NewMockUserRepository(ctrl)

	t.Run("successful customer addition", func(t *testing.T) {
		expectedCustomerID := uuid.New()
		addCustomerRequest := models.AddCustomerRequest{
			Name:        "John Doe",
			Email:       "john@example.com",
			PhoneNumber: "1234567890",
			Address:     "123 Main St",
		}

		repo.EXPECT().
			AddCustomer(gomock.Any(), gomock.Any()).
			DoAndReturn(func(_ context.Context, customer models.Customer) (uuid.UUID, error) {
				require.Equal(t, addCustomerRequest.Name, customer.Name)
				require.Equal(t, addCustomerRequest.Email, customer.Email)
				require.Equal(t, addCustomerRequest.PhoneNumber, customer.PhoneNumber)
				require.Equal(t, addCustomerRequest.Address, customer.Address)
				require.NotEqual(t, uuid.Nil, customer.CustomerID)
				return expectedCustomerID, nil
			})

		service := newUserServiceImpl(repo)
		customerID, err := service.AddCustomer(ctx, addCustomerRequest)
		require.NoError(t, err)
		require.Equal(t, expectedCustomerID, customerID)
	})

	t.Run("repository error", func(t *testing.T) {
		addCustomerRequest := models.AddCustomerRequest{
			Name:  "Jane Doe",
			Email: "jane@example.com",
		}

		expectedError := errors.New("database error")
		repo.EXPECT().
			AddCustomer(gomock.Any(), gomock.Any()).
			Return(uuid.Nil, expectedError)

		service := newUserServiceImpl(repo)
		customerID, err := service.AddCustomer(ctx, addCustomerRequest)
		require.Error(t, err)
		require.Equal(t, uuid.Nil, customerID)
		require.Equal(t, expectedError, err)
	})

	t.Run("empty name", func(t *testing.T) {
		addCustomerRequest := models.AddCustomerRequest{
			Name:  "",
			Email: "empty@example.com",
		}

		repo.EXPECT().
			AddCustomer(gomock.Any(), gomock.Any()).
			Return(uuid.Nil, errors.New("name cannot be empty"))

		service := newUserServiceImpl(repo)
		customerID, err := service.AddCustomer(ctx, addCustomerRequest)
		require.Error(t, err)
		require.Equal(t, uuid.Nil, customerID)
		require.Contains(t, err.Error(), "name cannot be empty")
	})

	t.Run("invalid email", func(t *testing.T) {
		addCustomerRequest := models.AddCustomerRequest{
			Name:  "Invalid Email",
			Email: "invalid-email",
		}

		repo.EXPECT().
			AddCustomer(gomock.Any(), gomock.Any()).
			Return(uuid.Nil, errors.New("invalid email format"))

		service := newUserServiceImpl(repo)
		customerID, err := service.AddCustomer(ctx, addCustomerRequest)
		require.Error(t, err)
		require.Equal(t, uuid.Nil, customerID)
		require.Contains(t, err.Error(), "invalid email format")
	})
}
