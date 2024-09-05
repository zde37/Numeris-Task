package service

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	mocked "github.com/zde37/Numeris-Task/internal/mock"
	"github.com/zde37/Numeris-Task/internal/models"
	"go.uber.org/mock/gomock"
)

func TestGetInvoiceDetails(t *testing.T) {
	ctx := context.Background()
	invoiceID := uuid.New()
	mockInvoiceDetails := &models.InvoiceDetails{
		Invoice: models.Invoice{
			InvoiceID: invoiceID,
			Status:    string(models.InvoiceStatusPaid),
		},
		Items: []models.InvoiceItem{
			{ItemID: uuid.New(), Name: "Test Item"},
		},
		Activities: []models.InvoiceActivity{
			{ActivityID: uuid.New(), Title: "Test Activity"},
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocked.NewMockInvoiceRepository(ctrl)

	t.Run("successful retrieval", func(t *testing.T) {
		repo.EXPECT().
			GetInvoiceDetails(gomock.Any(), invoiceID).
			Times(1).
			Return(mockInvoiceDetails, nil)

		service := newInvoiceServiceImpl(repo)
		details, err := service.GetInvoiceDetails(ctx, invoiceID)
		require.NoError(t, err)
		require.NotNil(t, details)
		require.Equal(t, mockInvoiceDetails, details)
	})

	t.Run("invoice not found", func(t *testing.T) {
		repo.EXPECT().
			GetInvoiceDetails(gomock.Any(), invoiceID).
			Times(1).
			Return(nil, sql.ErrNoRows)

		service := newInvoiceServiceImpl(repo)
		details, err := service.GetInvoiceDetails(ctx, invoiceID)
		require.Error(t, err)
		require.Nil(t, details)
		require.Equal(t, sql.ErrNoRows, err)
	})

	t.Run("database error", func(t *testing.T) {
		expectedErr := errors.New("database connection error")
		repo.EXPECT().
			GetInvoiceDetails(gomock.Any(), invoiceID).
			Times(1).
			Return(nil, expectedErr)

		service := newInvoiceServiceImpl(repo)
		details, err := service.GetInvoiceDetails(ctx, invoiceID)
		require.Error(t, err)
		require.Nil(t, details)
		require.Equal(t, expectedErr, err)
	})

	t.Run("invalid invoice ID", func(t *testing.T) {
		invalidID := uuid.Nil
		repo.EXPECT().
			GetInvoiceDetails(gomock.Any(), invalidID).
			Times(1).
			Return(nil, errors.New("invalid invoice ID"))

		service := newInvoiceServiceImpl(repo)
		details, err := service.GetInvoiceDetails(ctx, invalidID)
		require.Error(t, err)
		require.Nil(t, details)
		require.Contains(t, err.Error(), "invalid invoice ID")
	})
}

func TestGetTotalByStatus(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocked.NewMockInvoiceRepository(ctrl)

	t.Run("successful retrieval", func(t *testing.T) {
		expectedTotal := 1000.0
		expectedCount := 5
		repo.EXPECT().
			GetTotalByStatus(gomock.Any(), models.InvoiceStatusPaid).
			Times(1).
			Return(expectedTotal, expectedCount, nil)

		service := newInvoiceServiceImpl(repo)
		total, count, err := service.GetTotalByStatus(ctx, models.InvoiceStatusPaid)
		require.NoError(t, err)
		require.Equal(t, expectedTotal, total)
		require.Equal(t, expectedCount, count)
	})

	t.Run("zero invoices", func(t *testing.T) {
		repo.EXPECT().
			GetTotalByStatus(gomock.Any(), models.InvoiceStatusPending).
			Times(1).
			Return(0.0, 0, nil)

		service := newInvoiceServiceImpl(repo)
		total, count, err := service.GetTotalByStatus(ctx, models.InvoiceStatusPending)
		require.NoError(t, err)
		require.Equal(t, 0.0, total)
		require.Equal(t, 0, count)
	})

	t.Run("database error", func(t *testing.T) {
		expectedErr := errors.New("database connection error")
		repo.EXPECT().
			GetTotalByStatus(gomock.Any(), models.InvoiceStatusOverDue).
			Times(1).
			Return(0.0, 0, expectedErr)

		service := newInvoiceServiceImpl(repo)
		total, count, err := service.GetTotalByStatus(ctx, models.InvoiceStatusOverDue)
		require.Error(t, err)
		require.Equal(t, 0.0, total)
		require.Equal(t, 0, count)
		require.Equal(t, expectedErr, err)
	})

	t.Run("invalid status", func(t *testing.T) {
		invalidStatus := models.InvoiceStatus("INVALID")
		repo.EXPECT().
			GetTotalByStatus(gomock.Any(), invalidStatus).
			Times(1).
			Return(0.0, 0, errors.New("invalid status"))

		service := newInvoiceServiceImpl(repo)
		total, count, err := service.GetTotalByStatus(ctx, invalidStatus)
		require.Error(t, err)
		require.Equal(t, 0.0, total)
		require.Equal(t, 0, count)
		require.Contains(t, err.Error(), "invalid status")
	})
}

func TestGetRecentInvoices(t *testing.T) {
	ctx := context.Background()
	senderID := uuid.New()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocked.NewMockInvoiceRepository(ctrl)

	t.Run("successful retrieval", func(t *testing.T) {
		expectedInvoices := []models.Invoice{
			{InvoiceID: uuid.New(), SenderID: senderID},
			{InvoiceID: uuid.New(), SenderID: senderID},
		}
		repo.EXPECT().
			GetRecentInvoices(gomock.Any(), senderID, int32(10), int32(0)).
			Times(1).
			Return(expectedInvoices, nil)

		service := newInvoiceServiceImpl(repo)
		invoices, err := service.GetRecentInvoices(ctx, senderID, 1, 10)
		require.NoError(t, err)
		require.Equal(t, expectedInvoices, invoices)
	})

	t.Run("empty result", func(t *testing.T) {
		repo.EXPECT().
			GetRecentInvoices(gomock.Any(), senderID, int32(10), int32(90)).
			Times(1).
			Return([]models.Invoice{}, nil)

		service := newInvoiceServiceImpl(repo)
		invoices, err := service.GetRecentInvoices(ctx, senderID, 10, 10)
		require.NoError(t, err)
		require.Empty(t, invoices)
	})

	t.Run("database error", func(t *testing.T) {
		expectedErr := errors.New("database connection error")
		repo.EXPECT().
			GetRecentInvoices(gomock.Any(), senderID, int32(10), int32(0)).
			Times(1).
			Return(nil, expectedErr)

		service := newInvoiceServiceImpl(repo)
		invoices, err := service.GetRecentInvoices(ctx, senderID, 1, 10)
		require.Error(t, err)
		require.Nil(t, invoices)
		require.Equal(t, expectedErr, err)
	})
}

func TestGetRecentActivities(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocked.NewMockInvoiceRepository(ctrl)

	t.Run("successful retrieval", func(t *testing.T) {
		expectedActivities := []models.RecentActivity{
			{ActivityID: uuid.New(), UserID: userID},
			{ActivityID: uuid.New(), UserID: userID},
		}
		repo.EXPECT().
			GetRecentActivities(gomock.Any(), userID, int32(10), int32(0)).
			Times(1).
			Return(expectedActivities, nil)

		service := newInvoiceServiceImpl(repo)
		activities, err := service.GetRecentActivities(ctx, userID, 1, 10)
		require.NoError(t, err)
		require.Equal(t, expectedActivities, activities)
	})

	t.Run("empty result", func(t *testing.T) {
		repo.EXPECT().
			GetRecentActivities(gomock.Any(), userID, int32(10), int32(90)).
			Times(1).
			Return([]models.RecentActivity{}, nil)

		service := newInvoiceServiceImpl(repo)
		activities, err := service.GetRecentActivities(ctx, userID, 10, 10)
		require.NoError(t, err)
		require.Empty(t, activities)
	})

	t.Run("database error", func(t *testing.T) {
		expectedErr := errors.New("database connection error")
		repo.EXPECT().
			GetRecentActivities(gomock.Any(), userID, int32(10), int32(0)).
			Times(1).
			Return(nil, expectedErr)

		service := newInvoiceServiceImpl(repo)
		activities, err := service.GetRecentActivities(ctx, userID, 1, 10)
		require.Error(t, err)
		require.Nil(t, activities)
		require.Equal(t, expectedErr, err)
	})
}

func TestGetInvoiceActivities(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	invoiceID := uuid.New()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocked.NewMockInvoiceRepository(ctrl)

	t.Run("successful retrieval", func(t *testing.T) {
		expectedActivities := []models.InvoiceActivity{
			{ActivityID: uuid.New(), InvoiceID: invoiceID, UserID: userID},
			{ActivityID: uuid.New(), InvoiceID: invoiceID, UserID: userID},
		}
		repo.EXPECT().
			GetInvoiceActivities(gomock.Any(), userID, invoiceID, int32(10), int32(0)).
			Times(1).
			Return(expectedActivities, nil)

		service := newInvoiceServiceImpl(repo)
		activities, err := service.GetInvoiceActivities(ctx, userID, invoiceID, 1, 10)
		require.NoError(t, err)
		require.Equal(t, expectedActivities, activities)
	})

	t.Run("empty result", func(t *testing.T) {
		repo.EXPECT().
			GetInvoiceActivities(gomock.Any(), userID, invoiceID, int32(10), int32(90)).
			Times(1).
			Return([]models.InvoiceActivity{}, nil)

		service := newInvoiceServiceImpl(repo)
		activities, err := service.GetInvoiceActivities(ctx, userID, invoiceID, 10, 10)
		require.NoError(t, err)
		require.Empty(t, activities)
	})

	t.Run("database error", func(t *testing.T) {
		expectedErr := errors.New("database connection error")
		repo.EXPECT().
			GetInvoiceActivities(gomock.Any(), userID, invoiceID, int32(10), int32(0)).
			Times(1).
			Return(nil, expectedErr)

		service := newInvoiceServiceImpl(repo)
		activities, err := service.GetInvoiceActivities(ctx, userID, invoiceID, 1, 10)
		require.Error(t, err)
		require.Nil(t, activities)
		require.Equal(t, expectedErr, err)
	})

	t.Run("invalid user ID", func(t *testing.T) {
		invalidUserID := uuid.Nil
		repo.EXPECT().
			GetInvoiceActivities(gomock.Any(), invalidUserID, invoiceID, int32(10), int32(0)).
			Times(1).
			Return(nil, errors.New("invalid user ID"))

		service := newInvoiceServiceImpl(repo)
		activities, err := service.GetInvoiceActivities(ctx, invalidUserID, invoiceID, 1, 10)
		require.Error(t, err)
		require.Nil(t, activities)
		require.Contains(t, err.Error(), "invalid user ID")
	})

	t.Run("invalid invoice ID", func(t *testing.T) {
		invalidInvoiceID := uuid.Nil
		repo.EXPECT().
			GetInvoiceActivities(gomock.Any(), userID, invalidInvoiceID, int32(10), int32(0)).
			Times(1).
			Return(nil, errors.New("invalid invoice ID"))

		service := newInvoiceServiceImpl(repo)
		activities, err := service.GetInvoiceActivities(ctx, userID, invalidInvoiceID, 1, 10)
		require.Error(t, err)
		require.Nil(t, activities)
		require.Contains(t, err.Error(), "invalid invoice ID")
	})
}

func TestAddInvoiceActivity(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocked.NewMockInvoiceRepository(ctrl)

	t.Run("successful addition", func(t *testing.T) {
		validInvoiceID := uuid.New()
		validUserID := uuid.New()
		expectedActivityID := uuid.New()

		request := models.AddInvoiceActivityRequest{
			InvoiceID:   validInvoiceID.String(),
			UserID:      validUserID.String(),
			Title:       "Test Activity",
			Description: "Test Description",
		}

		repo.EXPECT().
			AddInvoiceActivity(gomock.Any(), gomock.Any()).
			DoAndReturn(func(_ context.Context, activity models.InvoiceActivity) (uuid.UUID, error) {
				require.Equal(t, validInvoiceID, activity.InvoiceID)
				require.Equal(t, validUserID, activity.UserID)
				require.Equal(t, request.Title, activity.Title)
				require.Equal(t, request.Description, activity.Description)
				return expectedActivityID, nil
			})

		service := newInvoiceServiceImpl(repo)
		activityID, err := service.AddInvoiceActivity(ctx, request)
		require.NoError(t, err)
		require.Equal(t, expectedActivityID, activityID)
	})

	t.Run("invalid invoice ID", func(t *testing.T) {
		request := models.AddInvoiceActivityRequest{
			InvoiceID:   "invalid-uuid",
			UserID:      uuid.New().String(),
			Title:       "Test Activity",
			Description: "Test Description",
		}

		service := newInvoiceServiceImpl(repo)
		activityID, err := service.AddInvoiceActivity(ctx, request)
		require.Error(t, err)
		require.Equal(t, uuid.Nil, activityID)
		require.Contains(t, err.Error(), "invalid invoice id")
	})

	t.Run("invalid user ID", func(t *testing.T) {
		request := models.AddInvoiceActivityRequest{
			InvoiceID:   uuid.New().String(),
			UserID:      "invalid-uuid",
			Title:       "Test Activity",
			Description: "Test Description",
		}

		service := newInvoiceServiceImpl(repo)
		activityID, err := service.AddInvoiceActivity(ctx, request)
		require.Error(t, err)
		require.Equal(t, uuid.Nil, activityID)
		require.Contains(t, err.Error(), "invalid user id")
	})

	t.Run("repository error", func(t *testing.T) {
		validInvoiceID := uuid.New()
		validUserID := uuid.New()
		expectedError := errors.New("repository error")

		request := models.AddInvoiceActivityRequest{
			InvoiceID:   validInvoiceID.String(),
			UserID:      validUserID.String(),
			Title:       "Test Activity",
			Description: "Test Description",
		}

		repo.EXPECT().
			AddInvoiceActivity(gomock.Any(), gomock.Any()).
			Return(uuid.Nil, expectedError)

		service := newInvoiceServiceImpl(repo)
		activityID, err := service.AddInvoiceActivity(ctx, request)
		require.Error(t, err)
		require.Equal(t, uuid.Nil, activityID)
		require.Equal(t, expectedError, err)
	})
}

func TestCreateInvoice(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocked.NewMockInvoiceRepository(ctrl)

	t.Run("successful creation", func(t *testing.T) {
		expectedInvoiceID := uuid.New()
		validRequest := models.CreateInvoiceRequest{
			Invoice: models.InvoiceInfo{
				SenderID:           uuid.New().String(),
				TotalAmount:        1000,
				DiscountPercentage: 10,
				DiscountedAmount:   100,
				FinalAmount:        900,
				Status:             string(models.InvoiceStatusPending),
				Currency:           "USD",
				Notes:              "Test invoice",
				IssueDate:          "2023-05-01",
				DueDate:            "2023-05-31",
			},
			CustomerID:      uuid.New().String(),
			PaymentMethodID: uuid.New().String(),
			InvoiceItems: []models.InvoiceItemDetails{
				{
					Name:        "Item 1",
					Description: "Description 1",
					Quantity:    2,
					UnitPrice:   500,
					TotalPrice:  1000,
				},
			},
		}

		repo.EXPECT().
			CreateInvoice(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(expectedInvoiceID, nil)

		service := newInvoiceServiceImpl(repo)
		invoiceID, err := service.CreateInvoice(ctx, validRequest)
		require.NoError(t, err)
		require.Equal(t, expectedInvoiceID, invoiceID)
	})

	t.Run("invalid sender ID", func(t *testing.T) {
		invalidRequest := models.CreateInvoiceRequest{
			Invoice: models.InvoiceInfo{
				SenderID: "invalid-uuid",
			},
			CustomerID:      uuid.New().String(),
			PaymentMethodID: uuid.New().String(),
		}

		service := newInvoiceServiceImpl(repo)
		invoiceID, err := service.CreateInvoice(ctx, invalidRequest)
		require.Error(t, err)
		require.Equal(t, uuid.Nil, invoiceID)
		require.Contains(t, err.Error(), "invalid sender id")
	})

	t.Run("invalid customer ID", func(t *testing.T) {
		invalidRequest := models.CreateInvoiceRequest{
			Invoice: models.InvoiceInfo{
				SenderID: uuid.New().String(),
			},
			CustomerID:      "invalid-uuid",
			PaymentMethodID: uuid.New().String(),
		}

		service := newInvoiceServiceImpl(repo)
		invoiceID, err := service.CreateInvoice(ctx, invalidRequest)
		require.Error(t, err)
		require.Equal(t, uuid.Nil, invoiceID)
		require.Contains(t, err.Error(), "invalid customer id")
	})

	t.Run("invalid invoice status", func(t *testing.T) {
		invalidRequest := models.CreateInvoiceRequest{
			Invoice: models.InvoiceInfo{
				SenderID: uuid.New().String(),
				Status:   "INVALID_STATUS",
			},
			CustomerID:      uuid.New().String(),
			PaymentMethodID: uuid.New().String(),
		}

		service := newInvoiceServiceImpl(repo)
		invoiceID, err := service.CreateInvoice(ctx, invalidRequest)
		require.Error(t, err)
		require.Equal(t, uuid.Nil, invoiceID)
		require.Contains(t, err.Error(), "invalid invoice status")
	})

	t.Run("invalid issue date format", func(t *testing.T) {
		invalidRequest := models.CreateInvoiceRequest{
			Invoice: models.InvoiceInfo{
				SenderID:  uuid.New().String(),
				Status:    string(models.InvoiceStatusPending),
				IssueDate: "01-05-2023",
				DueDate:   "2023-05-31",
			},
			CustomerID:      uuid.New().String(),
			PaymentMethodID: uuid.New().String(),
		}

		service := newInvoiceServiceImpl(repo)
		invoiceID, err := service.CreateInvoice(ctx, invalidRequest)
		require.Error(t, err)
		require.Equal(t, uuid.Nil, invoiceID)
		require.Contains(t, err.Error(), "issue date has invalid date format")
	})

	t.Run("invalid due date format", func(t *testing.T) {
		invalidRequest := models.CreateInvoiceRequest{
			Invoice: models.InvoiceInfo{
				SenderID:  uuid.New().String(),
				Status:    string(models.InvoiceStatusPending),
				IssueDate: "2023-05-01",
				DueDate:   "31-05-2023",
			},
			CustomerID:      uuid.New().String(),
			PaymentMethodID: uuid.New().String(),
		}

		service := newInvoiceServiceImpl(repo)
		invoiceID, err := service.CreateInvoice(ctx, invalidRequest)
		require.Error(t, err)
		require.Equal(t, uuid.Nil, invoiceID)
		require.Contains(t, err.Error(), "due date has invalid date format")
	})

	t.Run("invalid payment method ID", func(t *testing.T) {
		invalidRequest := models.CreateInvoiceRequest{
			Invoice: models.InvoiceInfo{
				SenderID:  uuid.New().String(),
				Status:    string(models.InvoiceStatusPending),
				IssueDate: "2023-05-01",
				DueDate:   "2023-05-31",
			},
			CustomerID:      uuid.New().String(),
			PaymentMethodID: "invalid-uuid",
		}

		service := newInvoiceServiceImpl(repo)
		invoiceID, err := service.CreateInvoice(ctx, invalidRequest)
		require.Error(t, err)
		require.Equal(t, uuid.Nil, invoiceID)
		require.Contains(t, err.Error(), "invalid payment method id")
	})

	t.Run("repository error", func(t *testing.T) {
		validRequest := models.CreateInvoiceRequest{
			Invoice: models.InvoiceInfo{
				SenderID:  uuid.New().String(),
				Status:    string(models.InvoiceStatusPending),
				IssueDate: "2023-05-01",
				DueDate:   "2023-05-31",
			},
			CustomerID:      uuid.New().String(),
			PaymentMethodID: uuid.New().String(),
		}

		expectedError := errors.New("repository error")
		repo.EXPECT().
			CreateInvoice(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(uuid.Nil, expectedError)

		service := newInvoiceServiceImpl(repo)
		invoiceID, err := service.CreateInvoice(ctx, validRequest)
		require.Error(t, err)
		require.Equal(t, uuid.Nil, invoiceID)
		require.Equal(t, expectedError, err)
	})
}
