package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/zde37/Numeris-Task/internal/helpers"
	"github.com/zde37/Numeris-Task/internal/models"
	"github.com/zde37/Numeris-Task/internal/service"
)

type handlerImpl struct {
	service *service.Service
	router  *gin.Engine
}

// NewHandlerImpl creates a new instance of the handlerImpl struct, which implements the Handler interface. 
func NewHandlerImpl(environment string, service *service.Service) Handler {
	h := &handlerImpl{
		service: service,
		router:  gin.Default(),
	}

	if environment == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}

	h.registerRoutes()
	return h
}

// GetRouter returns the underlying Gin router instance.
func (h *handlerImpl) GetRouter() *gin.Engine {
	return h.router
}

// registerRoutes sets up the routes for the Gin router. It creates a v1 group and registers the following routes:
//
// GET /v1/hello-world - Handles the "Hello World" request.
// POST /v1/invoices - Handles the creation of a new invoice.
// POST /v1/user - Handles the creation of a new user.
// POST /v1/payment - Handles the addition of a new payment method.
// POST /v1/customer - Handles the addition of a new customer.
// GET /v1/invoices/:invoiceID - Handles the retrieval of invoice details.
// POST /v1/invoices/activity - Handles the addition of a new invoice activity.
// GET /v1/invoices/total/:status - Handles the retrieval of the total invoices by status.
// GET /v1/invoices/recent/:senderID - Handles the retrieval of the most recent invoices for a given sender.
// GET /v1/activities/recent/:userID - Handles the retrieval of the most recent activities for a given user.
// GET /v1/invoices/:invoiceID/activities/:userID - Handles the retrieval of the activities for a given invoice and user.
func (h *handlerImpl) registerRoutes() {
	v1 := h.router.Group("v1")
	{
		v1.GET("/hello-world", h.HelloWorld)
		v1.POST("/invoices", h.CreateInvoice)
		v1.POST("/user", h.CreateUser)
		v1.POST("/payment", h.AddPaymentMethod)
		v1.POST("/customer", h.AddCustomer)
		v1.GET("/invoices/:invoiceID", h.GetInvoiceDetails)
		v1.POST("/invoices/activity", h.AddInvoiceActivity)
		v1.GET("/invoices/total/:status", h.GetTotalByStatus)
		v1.GET("/invoices/recent/:senderID", h.GetRecentInvoices)
		v1.GET("/activities/recent/:userID", h.GetRecentActivities)
		v1.GET("/invoices/:invoiceID/activities/:userID", h.GetInvoiceActivities)
	}
}

// HelloWorld is a handler function that responds with a "Hello from Numeris Book" message. 
func (h *handlerImpl) HelloWorld(c *gin.Context) {
	c.String(http.StatusOK, "Hello from Numeris Book")
}

// CreateInvoice is a handler function that creates a new invoice. 
func (h *handlerImpl) CreateInvoice(ctx *gin.Context) {
	var req models.CreateInvoiceRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	invoiceID, err := h.service.Invoice.CreateInvoice(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"invoice_id": invoiceID})
}

// GetInvoiceDetails is a handler function that retrieves the details of an invoice. 
func (h *handlerImpl) GetInvoiceDetails(ctx *gin.Context) {
	invoiceID, err := uuid.Parse(ctx.Param("invoiceID"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid invoice ID"})
		return
	}

	details, err := h.service.Invoice.GetInvoiceDetails(ctx, invoiceID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, details)
}

// AddInvoiceActivity is a handler function that adds a new activity to an invoice. 
func (h *handlerImpl) AddInvoiceActivity(ctx *gin.Context) {
	var activity models.AddInvoiceActivityRequest
	if err := ctx.ShouldBind(&activity); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	activityID, err := h.service.Invoice.AddInvoiceActivity(ctx, activity)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"activity_id": activityID})
}

// GetTotalByStatus is a handler function that retrieves the total amount and count of invoices by a given status. 
func (h *handlerImpl) GetTotalByStatus(ctx *gin.Context) {
	if err := helpers.ValidateInvoiceStatus(ctx.Param("status")); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	status := models.InvoiceStatus(ctx.Param("status"))

	totalAmount, count, err := h.service.Invoice.GetTotalByStatus(ctx, status)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"total_amount": totalAmount, "count": count})
}

// GetRecentInvoices is a handler function that retrieves the most recent invoices for a given sender. 
func (h *handlerImpl) GetRecentInvoices(ctx *gin.Context) {
	senderID, err := uuid.Parse(ctx.Param("senderID"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid sender ID"})
		return
	}

	limit, page := h.getPaginationParams(ctx)

	invoices, err := h.service.Invoice.GetRecentInvoices(ctx, senderID, page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, invoices)
}

// GetRecentActivities is a handler function that retrieves the recent activities for a given user. 
func (h *handlerImpl) GetRecentActivities(ctx *gin.Context) {
	userID, err := uuid.Parse(ctx.Param("userID"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	limit, page := h.getPaginationParams(ctx)

	activities, err := h.service.Invoice.GetRecentActivities(ctx, userID, page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, activities)
}

// GetInvoiceActivities is a handler function that retrieves the recent activities for a given invoice and user. 
func (h *handlerImpl) GetInvoiceActivities(ctx *gin.Context) {
	userID, err := uuid.Parse(ctx.Param("userID"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	invoiceID, err := uuid.Parse(ctx.Param("invoiceID"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid invoice ID"})
		return
	}

	limit, page := h.getPaginationParams(ctx)

	activities, err := h.service.Invoice.GetInvoiceActivities(ctx, userID, invoiceID, page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, activities)
}

// CreateUser is a handler function that creates a new user. 
func (h *handlerImpl) CreateUser(ctx *gin.Context) {
	var req models.CreateUserRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := h.service.User.CreateUser(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"user_id": userID})
}

// AddPaymentMethod is a handler function that adds a new payment method for a user 
func (h *handlerImpl) AddPaymentMethod(ctx *gin.Context) {
	var req models.AddPaymentMethodRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	paymentMethodID, err := h.service.User.AddPaymentMethod(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"payment_method_id": paymentMethodID})
}

// AddCustomer is a handler function that creates a new customer. 
func (h *handlerImpl) AddCustomer(ctx *gin.Context) {
	var req models.AddCustomerRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	customerID, err := h.service.User.AddCustomer(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"customer_id": customerID})
}

// getPaginationParams is a helper function that extracts the limit and page
// parameters from the request context. If the parameters are not provided,
// it uses default values of 10 for limit and 1 for page. 
func (h *handlerImpl) getPaginationParams(ctx *gin.Context) (limit, page int32) {
	limitStr := ctx.DefaultQuery("limit", "10")
	pageStr := ctx.DefaultQuery("page", "1")

	limit64, err := strconv.ParseInt(limitStr, 10, 32)
	if err != nil {
		limit = 10
	} else {
		limit = int32(limit64)
	}

	page64, err := strconv.ParseInt(pageStr, 10, 32)
	if err != nil {
		page = 1
	} else {
		page = int32(page64)
	}

	return limit, page
}
