package controller

import "github.com/gin-gonic/gin"

type Handler interface {
	HelloWorld(ctx *gin.Context)
	CreateInvoice(ctx *gin.Context)
	GetInvoiceDetails(ctx *gin.Context)
	AddInvoiceActivity(ctx *gin.Context)
	GetTotalByStatus(ctx *gin.Context)
	GetRecentInvoices(ctx *gin.Context)
	GetRecentActivities(ctx *gin.Context)
	GetInvoiceActivities(ctx *gin.Context)
	CreateUser(ctx *gin.Context)
	AddPaymentMethod(ctx *gin.Context)
	AddCustomer(ctx *gin.Context)
	GetRouter() *gin.Engine 
}
