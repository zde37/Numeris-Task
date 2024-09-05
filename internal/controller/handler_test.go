package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	mocked "github.com/zde37/Numeris-Task/internal/mock"
	"github.com/zde37/Numeris-Task/internal/models"
	"github.com/zde37/Numeris-Task/internal/service"
	"go.uber.org/mock/gomock"
)

func TestCreateInvoice(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockInvoiceService := mocked.NewMockInvoiceService(ctrl)
	srv := &service.Service{
		Invoice: mockInvoiceService,
	}
	handler := NewHandlerImpl("dev", srv)

	t.Run("successful invoice creation", func(t *testing.T) {
		req := models.CreateInvoiceRequest{
			Invoice: models.InvoiceInfo{
				Status:             string(models.InvoiceStatusPending),
				SenderID:           uuid.New().String(),
				IssueDate:          time.Now().Format("2006-01-02"),
				DueDate:            time.Now().Format("2006-01-02"),
				TotalAmount:        10,
				DiscountPercentage: 100,
				DiscountedAmount:   1000,
				FinalAmount:        9000,
				Currency:           "NGN",
				Notes:              "Test invoice",
			},
			InvoiceItems: []models.InvoiceItemDetails{
				{
					Name:        "Test Item",
					Description: "Test Description",
					Quantity:    1,
					UnitPrice:   10.0,
				},
			},
			CustomerID:      uuid.New().String(),
			PaymentMethodID: uuid.New().String(),
		}
		expectedInvoiceID := uuid.New()

		mockInvoiceService.EXPECT().
			CreateInvoice(gomock.Any(), req).
			Return(expectedInvoiceID, nil)
 
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		jsonData, _ := json.Marshal(req)
		c.Request, _ = http.NewRequest(http.MethodPost, "/invoices", bytes.NewBuffer(jsonData))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.CreateInvoice(c)

		require.Equal(t, http.StatusCreated, w.Code)
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		require.Equal(t, expectedInvoiceID.String(), response["invoice_id"])
	})

	t.Run("invalid request body", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request, _ = http.NewRequest(http.MethodPost, "/invoices", bytes.NewBufferString("invalid json"))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.CreateInvoice(c)

		require.Equal(t, http.StatusBadRequest, w.Code)
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		require.Contains(t, response["error"], "invalid character")
	})

	t.Run("service error", func(t *testing.T) {
		req := models.CreateInvoiceRequest{
			Invoice: models.InvoiceInfo{
				Status:             string(models.InvoiceStatusPending),
				SenderID:           uuid.New().String(),
				IssueDate:          time.Now().Format("2006-01-02"),
				DueDate:            time.Now().Format("2006-01-02"),
				TotalAmount:        10,
				DiscountPercentage: 100,
				DiscountedAmount:   1000,
				FinalAmount:        9000,
				Currency:           "NGN",
				Notes:              "Test invoice",
			},
			InvoiceItems: []models.InvoiceItemDetails{
				{
					Name:        "Test Item",
					Description: "Test Description",
					Quantity:    1,
					UnitPrice:   10.0,
				},
			},
			CustomerID:      uuid.New().String(),
			PaymentMethodID: uuid.New().String(),
		}
		expectedError := errors.New("service error")

		mockInvoiceService.EXPECT().
			CreateInvoice(gomock.Any(), req).
			Return(uuid.Nil, expectedError)

		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		jsonData, _ := json.Marshal(req)
		c.Request, _ = http.NewRequest(http.MethodPost, "/invoices", bytes.NewBuffer(jsonData))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.CreateInvoice(c)

		require.Equal(t, http.StatusInternalServerError, w.Code)
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		require.Equal(t, expectedError.Error(), response["error"])
	})
}

func TestGetInvoiceDetails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockInvoiceService := mocked.NewMockInvoiceService(ctrl)
	srv := &service.Service{
		Invoice: mockInvoiceService,
	}
	handler := NewHandlerImpl("prod", srv)

	t.Run("successful invoice details retrieval", func(t *testing.T) {
		invoiceID := uuid.New()
		invoice := models.Invoice{
			Status: string(models.InvoiceStatusPending),
		}
		expectedDetails := &models.InvoiceDetails{
			Invoice: invoice,
		}

		mockInvoiceService.EXPECT().
			GetInvoiceDetails(gomock.Any(), invoiceID).
			Return(expectedDetails, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "invoiceID", Value: invoiceID.String()}}

		handler.GetInvoiceDetails(c)

		require.Equal(t, http.StatusOK, w.Code)
		response := &models.InvoiceDetails{
			Invoice: invoice,
		}
		err := json.Unmarshal(w.Body.Bytes(), response)
		require.NoError(t, err)
		require.Equal(t, expectedDetails, response)
	})

	t.Run("invalid invoice ID", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "invoiceID", Value: "invalid-uuid"}}

		handler.GetInvoiceDetails(c)

		require.Equal(t, http.StatusBadRequest, w.Code)
		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		require.Equal(t, "Invalid invoice ID", response["error"])
	})

	t.Run("service error", func(t *testing.T) {
		invoiceID := uuid.New()
		expectedError := errors.New("service error")

		mockInvoiceService.EXPECT().
			GetInvoiceDetails(gomock.Any(), invoiceID).
			Return(&models.InvoiceDetails{}, expectedError)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "invoiceID", Value: invoiceID.String()}}

		handler.GetInvoiceDetails(c)

		require.Equal(t, http.StatusInternalServerError, w.Code)
		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		require.Equal(t, expectedError.Error(), response["error"])
	})
}

func TestAddInvoiceActivity(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockInvoiceService := mocked.NewMockInvoiceService(ctrl)
	srv := &service.Service{
		Invoice: mockInvoiceService,
	}
	handler := NewHandlerImpl("dev", srv)

	t.Run("successful activity addition", func(t *testing.T) {
		req := models.AddInvoiceActivityRequest{
			InvoiceID:   uuid.New().String(),
			UserID:      uuid.New().String(),
			Title:       "Test Title",
			Description: "Test Desc",
		}
		expectedActivityID := uuid.New()

		mockInvoiceService.EXPECT().
			AddInvoiceActivity(gomock.Any(), gomock.Eq(req)).
			Return(expectedActivityID, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		jsonData, _ := json.Marshal(req)
		c.Request, _ = http.NewRequest(http.MethodPost, "/invoices/activity", bytes.NewBuffer(jsonData))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.AddInvoiceActivity(c)

		require.Equal(t, http.StatusCreated, w.Code)
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		require.Equal(t, expectedActivityID.String(), response["activity_id"])
	})

	t.Run("invalid request body", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request, _ = http.NewRequest(http.MethodPost, "/invoices/activity", bytes.NewBufferString("invalid json"))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.AddInvoiceActivity(c)

		require.Equal(t, http.StatusBadRequest, w.Code)
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		require.Contains(t, response["error"], "invalid character")
	})

	t.Run("service error", func(t *testing.T) {
		req := models.AddInvoiceActivityRequest{
			InvoiceID:   uuid.New().String(),
			UserID:      uuid.New().String(),
			Title:       "Test Title",
			Description: "Test Desc",
		}
		expectedError := errors.New("service error")

		mockInvoiceService.EXPECT().
			AddInvoiceActivity(gomock.Any(), req).
			Return(uuid.Nil, expectedError)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		jsonData, _ := json.Marshal(req)
		c.Request, _ = http.NewRequest(http.MethodPost, "/invoices/activity", bytes.NewBuffer(jsonData))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.AddInvoiceActivity(c)

		require.Equal(t, http.StatusInternalServerError, w.Code)
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		require.Equal(t, expectedError.Error(), response["error"])
	})
}

func TestGetTotalByStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockInvoiceService := mocked.NewMockInvoiceService(ctrl)
	srv := &service.Service{
		Invoice: mockInvoiceService,
	}
	handler := NewHandlerImpl("dev", srv)

	t.Run("successful total retrieval", func(t *testing.T) {
		status := models.InvoiceStatusPending
		expectedTotal := float64(1000)
		expectedCount := int(5)

		mockInvoiceService.EXPECT().
			GetTotalByStatus(gomock.Any(), status).
			Return(expectedTotal, expectedCount, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "status", Value: string(status)}}

		handler.GetTotalByStatus(c)

		require.Equal(t, http.StatusOK, w.Code)
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		require.Equal(t, expectedTotal, response["total_amount"])
		require.Equal(t, float64(expectedCount), response["count"])
	})

	t.Run("invalid status", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "status", Value: "invalid_status"}}

		handler.GetTotalByStatus(c)

		require.Equal(t, http.StatusBadRequest, w.Code)
		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		require.Contains(t, response["error"], "invalid_status")
	})

	t.Run("service error", func(t *testing.T) {
		status := models.InvoiceStatusPaid
		expectedError := errors.New("service error")

		mockInvoiceService.EXPECT().
			GetTotalByStatus(gomock.Any(), status).
			Return(float64(0), int(0), expectedError)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "status", Value: string(status)}}

		handler.GetTotalByStatus(c)

		require.Equal(t, http.StatusInternalServerError, w.Code)
		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		require.Equal(t, expectedError.Error(), response["error"])
	})
}

func TestGetRecentInvoices(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockInvoiceService := mocked.NewMockInvoiceService(ctrl)
	srv := &service.Service{
		Invoice: mockInvoiceService,
	}
	handler := NewHandlerImpl("dev", srv)

	t.Run("successful recent invoices retrieval", func(t *testing.T) {
		senderID := uuid.New()
		limit := int32(10)
		page := int32(1)
		expectedInvoices := []models.Invoice{
			{Status: string(models.InvoiceStatusPending)},
			{Status: string(models.InvoiceStatusPaid)},
		}

		mockInvoiceService.EXPECT().
			GetRecentInvoices(gomock.Any(), senderID, page, limit).
			Return(expectedInvoices, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "senderID", Value: senderID.String()}}
		c.Request, _ = http.NewRequest(http.MethodGet, "/invoices/recent?limit=10&page=1", nil)

		handler.GetRecentInvoices(c)

		require.Equal(t, http.StatusOK, w.Code)
		var response []models.Invoice
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		require.Equal(t, expectedInvoices, response)
	})

	t.Run("invalid sender ID", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "senderID", Value: "invalid-uuid"}}

		handler.GetRecentInvoices(c)

		require.Equal(t, http.StatusBadRequest, w.Code)
		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		require.Equal(t, "Invalid sender ID", response["error"])
	})

	t.Run("service error", func(t *testing.T) {
		senderID := uuid.New()
		expectedError := errors.New("service error")

		mockInvoiceService.EXPECT().
			GetRecentInvoices(gomock.Any(), senderID, gomock.Any(), gomock.Any()).
			Return(nil, expectedError)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "senderID", Value: senderID.String()}}
		c.Request, _ = http.NewRequest(http.MethodGet, "/invoices/recent", nil)

		handler.GetRecentInvoices(c)

		require.Equal(t, http.StatusInternalServerError, w.Code)
		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		require.Equal(t, expectedError.Error(), response["error"])
	})

	t.Run("pagination parameters", func(t *testing.T) {
		senderID := uuid.New()
		limit := int32(20)
		page := int32(2)
		expectedInvoices := []models.Invoice{}

		mockInvoiceService.EXPECT().
			GetRecentInvoices(gomock.Any(), senderID, page, limit).
			Return(expectedInvoices, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "senderID", Value: senderID.String()}}
		c.Request, _ = http.NewRequest(http.MethodGet, fmt.Sprintf("/invoices/recent?limit=%d&page=%d", limit, page), nil)

		handler.GetRecentInvoices(c)

		require.Equal(t, http.StatusOK, w.Code)
		var response []models.Invoice
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		require.Equal(t, expectedInvoices, response)
	})
}

func TestGetRecentActivities(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockInvoiceService := mocked.NewMockInvoiceService(ctrl)
	srv := &service.Service{
		Invoice: mockInvoiceService,
	}
	handler := NewHandlerImpl("dev", srv)

	t.Run("successful recent activities retrieval", func(t *testing.T) {
		userID := uuid.New()
		limit := int32(10)
		page := int32(1)
		expectedActivities := []models.RecentActivity{
			{Title: "Activity 1"},
			{Title: "Activity 2"},
		}

		mockInvoiceService.EXPECT().
			GetRecentActivities(gomock.Any(), userID, page, limit).
			Return(expectedActivities, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "userID", Value: userID.String()}}
		c.Request, _ = http.NewRequest(http.MethodGet, "/activities/recent?limit=10&page=1", nil)

		handler.GetRecentActivities(c)

		require.Equal(t, http.StatusOK, w.Code)
		var response []models.RecentActivity
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		require.Equal(t, expectedActivities, response)
	})

	t.Run("invalid user ID", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "userID", Value: "invalid-uuid"}}

		handler.GetRecentActivities(c)

		require.Equal(t, http.StatusBadRequest, w.Code)
		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		require.Equal(t, "Invalid user ID", response["error"])
	})

	t.Run("service error", func(t *testing.T) {
		userID := uuid.New()
		expectedError := errors.New("service error")

		mockInvoiceService.EXPECT().
			GetRecentActivities(gomock.Any(), userID, gomock.Any(), gomock.Any()).
			Return(nil, expectedError)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "userID", Value: userID.String()}}
		c.Request, _ = http.NewRequest(http.MethodGet, "/activities/recent", nil)

		handler.GetRecentActivities(c)

		require.Equal(t, http.StatusInternalServerError, w.Code)
		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		require.Equal(t, expectedError.Error(), response["error"])
	})

	t.Run("pagination parameters", func(t *testing.T) {
		userID := uuid.New()
		limit := int32(20)
		page := int32(2)
		expectedActivities := []models.RecentActivity{}

		mockInvoiceService.EXPECT().
			GetRecentActivities(gomock.Any(), userID, page, limit).
			Return(expectedActivities, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "userID", Value: userID.String()}}
		c.Request, _ = http.NewRequest(http.MethodGet, "/activities/recent?limit=20&page=2", nil)

		handler.GetRecentActivities(c)

		require.Equal(t, http.StatusOK, w.Code)
		var response []models.RecentActivity
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		require.Equal(t, expectedActivities, response)
	})
}

func TestGetInvoiceActivities(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockInvoiceService := mocked.NewMockInvoiceService(ctrl)
	srv := &service.Service{
		Invoice: mockInvoiceService,
	}
	handler := NewHandlerImpl("dev", srv)

	t.Run("successful activities retrieval", func(t *testing.T) {
		userID := uuid.New()
		invoiceID := uuid.New()
		limit := int32(10)
		page := int32(1)
		expectedActivities := []models.InvoiceActivity{
			{Title: "Activity 1"},
			{Title: "Activity 2"},
		}

		mockInvoiceService.EXPECT().
			GetInvoiceActivities(gomock.Any(), userID, invoiceID, page, limit).
			Return(expectedActivities, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{
			{Key: "userID", Value: userID.String()},
			{Key: "invoiceID", Value: invoiceID.String()},
		}
		c.Request, _ = http.NewRequest(http.MethodGet, "/invoices/activities?limit=10&page=1", nil)

		handler.GetInvoiceActivities(c)

		require.Equal(t, http.StatusOK, w.Code)
		var response []models.InvoiceActivity
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		require.Equal(t, expectedActivities, response)
	})

	t.Run("invalid user ID", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{
			{Key: "userID", Value: "invalid-uuid"},
			{Key: "invoiceID", Value: uuid.New().String()},
		}

		handler.GetInvoiceActivities(c)

		require.Equal(t, http.StatusBadRequest, w.Code)
		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		require.Equal(t, "Invalid user ID", response["error"])
	})

	t.Run("invalid invoice ID", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{
			{Key: "userID", Value: uuid.New().String()},
			{Key: "invoiceID", Value: "invalid-uuid"},
		}

		handler.GetInvoiceActivities(c)

		require.Equal(t, http.StatusBadRequest, w.Code)
		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		require.Equal(t, "Invalid invoice ID", response["error"])
	})

	t.Run("service error", func(t *testing.T) {
		userID := uuid.New()
		invoiceID := uuid.New()
		expectedError := errors.New("service error")

		mockInvoiceService.EXPECT().
			GetInvoiceActivities(gomock.Any(), userID, invoiceID, gomock.Any(), gomock.Any()).
			Return(nil, expectedError)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{
			{Key: "userID", Value: userID.String()},
			{Key: "invoiceID", Value: invoiceID.String()},
		}
		c.Request, _ = http.NewRequest(http.MethodGet, "/invoices/activities", nil)

		handler.GetInvoiceActivities(c)

		require.Equal(t, http.StatusInternalServerError, w.Code)
		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		require.Equal(t, expectedError.Error(), response["error"])
	})
}

func TestCreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mocked.NewMockUserService(ctrl)
	srv := &service.Service{
		User: mockUserService,
	}
	handler := NewHandlerImpl("dev", srv)

	t.Run("successful user creation", func(t *testing.T) {
		req := models.CreateUserRequest{
			Username:          "John Doe",
			Email:             "john@example.com",
			Password:          "password123",
			FirstName:         "TEst 1",
			LastName:          "Test 2",
			ProfilePictureURL: "Pic 1",
			PhoneNumber:       "+1111111111",
			Address:           "Test Address",
		}
		expectedUserID := uuid.New()

		mockUserService.EXPECT().
			CreateUser(gomock.Any(), req).
			Return(expectedUserID, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		jsonData, _ := json.Marshal(req)
		c.Request, _ = http.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(jsonData))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.CreateUser(c)

		require.Equal(t, http.StatusCreated, w.Code)
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		require.Equal(t, expectedUserID.String(), response["user_id"])
	})

	t.Run("invalid request body", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request, _ = http.NewRequest(http.MethodPost, "/users", bytes.NewBufferString("invalid json"))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.CreateUser(c)

		require.Equal(t, http.StatusBadRequest, w.Code)
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		require.Contains(t, response["error"], "invalid character")
	})

	t.Run("service error", func(t *testing.T) {
		req := models.CreateUserRequest{
			Username:          "John Doe",
			Email:             "john@example.com",
			Password:          "password123",
			FirstName:         "TEst 1",
			LastName:          "Test 2",
			ProfilePictureURL: "Pic 1",
			PhoneNumber:       "+1111111111",
			Address:           "Test Address",
		}
		expectedError := errors.New("service error")

		mockUserService.EXPECT().
			CreateUser(gomock.Any(), req).
			Return(uuid.Nil, expectedError)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		jsonData, _ := json.Marshal(req)
		c.Request, _ = http.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(jsonData))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.CreateUser(c)

		require.Equal(t, http.StatusInternalServerError, w.Code)
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		require.Equal(t, expectedError.Error(), response["error"])
	})
}

func TestAddPaymentMethod(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mocked.NewMockUserService(ctrl)
	srv := &service.Service{
		User: mockUserService,
	}
	handler := NewHandlerImpl("dev", srv)

	t.Run("successful payment method addition", func(t *testing.T) {
		req := models.AddPaymentMethodRequest{
			UserID:        uuid.New().String(),
			AccountName:   "Account 1",
			BankName:      "Bank 1",
			AccountNumber: "4111111111111111",
			BankAddress:   "Bank Address 1",
			SwiftCode:     "Swift code 1",
		}
		expectedPaymentMethodID := uuid.New()

		mockUserService.EXPECT().
			AddPaymentMethod(gomock.Any(), req).
			Return(expectedPaymentMethodID, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		jsonData, _ := json.Marshal(req)
		c.Request, _ = http.NewRequest(http.MethodPost, "/payment-methods", bytes.NewBuffer(jsonData))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.AddPaymentMethod(c)

		require.Equal(t, http.StatusCreated, w.Code)
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		require.Equal(t, expectedPaymentMethodID.String(), response["payment_method_id"])
	})

	t.Run("invalid request body", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request, _ = http.NewRequest(http.MethodPost, "/payment-methods", bytes.NewBufferString("invalid json"))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.AddPaymentMethod(c)

		require.Equal(t, http.StatusBadRequest, w.Code)
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		require.Contains(t, response["error"], "invalid character")
	})

	t.Run("service error", func(t *testing.T) {
		req := models.AddPaymentMethodRequest{
			UserID:        uuid.New().String(),
			AccountName:   "Account 1",
			BankName:      "Bank 1",
			AccountNumber: "4111111111111111",
			BankAddress:   "Bank Address 1",
			SwiftCode:     "Swift code 1",
		}
		expectedError := errors.New("service error")

		mockUserService.EXPECT().
			AddPaymentMethod(gomock.Any(), req).
			Return(uuid.Nil, expectedError)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		jsonData, _ := json.Marshal(req)
		c.Request, _ = http.NewRequest(http.MethodPost, "/payment-methods", bytes.NewBuffer(jsonData))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.AddPaymentMethod(c)

		require.Equal(t, http.StatusInternalServerError, w.Code)
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		require.Equal(t, expectedError.Error(), response["error"])
	})
}
 
func TestAddCustomer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mocked.NewMockUserService(ctrl)
	srv := &service.Service{
		User: mockUserService,
	}
	handler := NewHandlerImpl("dev", srv)

	t.Run("successful customer addition", func(t *testing.T) {
		req := models.AddCustomerRequest{
			Name:    "John Doe",
			Email:   "john@example.com",
			PhoneNumber:   "+1234567890",
			Address: "123 Main St",
		}
		expectedCustomerID := uuid.New()

		mockUserService.EXPECT().
			AddCustomer(gomock.Any(), req).
			Return(expectedCustomerID, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		jsonData, _ := json.Marshal(req)
		c.Request, _ = http.NewRequest(http.MethodPost, "/customers", bytes.NewBuffer(jsonData))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.AddCustomer(c)

		require.Equal(t, http.StatusCreated, w.Code)
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		require.Equal(t, expectedCustomerID.String(), response["customer_id"])
	})

	t.Run("invalid request body", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request, _ = http.NewRequest(http.MethodPost, "/customers", bytes.NewBufferString("invalid json"))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.AddCustomer(c)

		require.Equal(t, http.StatusBadRequest, w.Code)
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		require.Contains(t, response["error"], "invalid character")
	})

	t.Run("service error", func(t *testing.T) {
		req := models.AddCustomerRequest{
			Name:    "John Doe",
			Email:   "john@example.com",
			PhoneNumber:   "+1234567890",
			Address: "123 Main St",
		}
		expectedError := errors.New("service error")

		mockUserService.EXPECT().
			AddCustomer(gomock.Any(), req).
			Return(uuid.Nil, expectedError)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		jsonData, _ := json.Marshal(req)
		c.Request, _ = http.NewRequest(http.MethodPost, "/customers", bytes.NewBuffer(jsonData))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.AddCustomer(c)

		require.Equal(t, http.StatusInternalServerError, w.Code)
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		require.Equal(t, expectedError.Error(), response["error"])
	})
}

func TestGetPaginationParams(t *testing.T) {
	handler := &handlerImpl{}

	t.Run("default values", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodGet, "/test", nil)

		limit, page := handler.getPaginationParams(c)

		require.Equal(t, int32(10), limit)
		require.Equal(t, int32(1), page)
	})

	t.Run("custom valid values", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodGet, "/test?limit=20&page=2", nil)

		limit, page := handler.getPaginationParams(c)

		require.Equal(t, int32(20), limit)
		require.Equal(t, int32(2), page)
	})

	t.Run("invalid limit", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodGet, "/test?limit=invalid&page=2", nil)

		limit, page := handler.getPaginationParams(c)

		require.Equal(t, int32(10), limit)
		require.Equal(t, int32(2), page)
	})

	t.Run("invalid page", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodGet, "/test?limit=20&page=invalid", nil)

		limit, page := handler.getPaginationParams(c)

		require.Equal(t, int32(20), limit)
		require.Equal(t, int32(1), page)
	})

	t.Run("negative values", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodGet, "/test?limit=-5&page=-1", nil)

		limit, page := handler.getPaginationParams(c)

		require.Equal(t, int32(-5), limit)
		require.Equal(t, int32(-1), page)
	})

	t.Run("zero values", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodGet, "/test?limit=0&page=0", nil)

		limit, page := handler.getPaginationParams(c)

		require.Equal(t, int32(0), limit)
		require.Equal(t, int32(0), page)
	})
}

func TestHelloWorld(t *testing.T) {
	handler := &handlerImpl{}

	t.Run("successful hello world response", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		handler.HelloWorld(c)

		require.Equal(t, http.StatusOK, w.Code)
		require.Equal(t, "Hello from Numeris Book", w.Body.String())
	})

	t.Run("correct content type", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		handler.HelloWorld(c)

		require.Equal(t, "text/plain; charset=utf-8", w.Header().Get("Content-Type"))
	})

	t.Run("no additional headers", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		handler.HelloWorld(c)

		require.Len(t, w.Header(), 1) // Only Content-Type should be present
	})
}

func TestRegisterRoutesUnknownPath(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	handler := &handlerImpl{router: router}

	handler.registerRoutes()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/v1/unknown", nil)
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusNotFound, w.Code)
}
