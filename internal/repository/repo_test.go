package repository

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"github.com/zde37/Numeris-Task/internal/helpers"
	"github.com/zde37/Numeris-Task/internal/models"
)

type testID struct {
	customerID      uuid.UUID
	senderID        uuid.UUID
	invoiceID       uuid.UUID
	paymentMethodID uuid.UUID
}

type InvoiceRepoTestSuite struct {
	suite.Suite
	ctx                context.Context
	dbPool             *pgxpool.Pool
	pgContainer        testcontainers.Container
	pgConnectionString string
	migrationURL       string
	repo               *Repository
	ids                testID
}

func (suite *InvoiceRepoTestSuite) SetupSuite() {
	suite.ctx = context.Background()
	containerPort := "5432"
	req := testcontainers.ContainerRequest{
		Image: "postgres:16-alpine",
		Env: map[string]string{
			"POSTGRES_USER":     "postgres_user",
			"POSTGRES_PASSWORD": "postgres_user",
			"POSTGRES_DB":       "Numeris_User_DB",
		},
		ExposedPorts: []string{containerPort + "/tcp"}, 
		WaitingFor: wait.ForAll(
			wait.ForLog("database system is ready to accept connections").WithOccurrence(2),
			wait.ForListeningPort(nat.Port(containerPort)),
		).WithDeadline(5 * time.Minute),
	}

	container, err := testcontainers.GenericContainer(suite.ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	suite.NoError(err)

	host, err := container.Host(suite.ctx)
	suite.NoError(err)

	mappedPort, err := container.MappedPort(suite.ctx, nat.Port(containerPort))
	suite.NoError(err)

	connStr := fmt.Sprintf("postgresql://postgres_user:postgres_user@%s:%s/Numeris_User_DB?sslmode=disable", host, mappedPort.Port())

	dbPool, err := pgxpool.New(suite.ctx, connStr)
	suite.NoError(err)

	err = dbPool.Ping(suite.ctx)
	suite.NoError(err)

	suite.pgContainer = container
	suite.pgConnectionString = connStr
	suite.dbPool = dbPool
	suite.migrationURL = "file://../../migrations"
	suite.repo = NewRepository(suite.dbPool)

	migration, err := migrate.New(suite.migrationURL, suite.pgConnectionString)
	suite.NoError(err)

	err = migration.Up()
	suite.NoError(err)

	suite.setupTestData()
}

func (suite *InvoiceRepoTestSuite) TearDownSuite() {
	if suite.dbPool != nil {
		suite.dbPool.Close()
	}
	if suite.pgContainer != nil {
		err := suite.pgContainer.Terminate(suite.ctx)
		suite.NoError(err)
	}
}

func (suite *InvoiceRepoTestSuite) setupTestData() {
	// create new customer
	customer := models.Customer{
		CustomerID:  uuid.New(),
		Name:        "Name 1",
		Email:       "Email 1",
		PhoneNumber: "Phone Number 1",
		Address:     "Address 1",
	}
	customerID, err := suite.repo.User.AddCustomer(suite.ctx, customer)
	suite.NoError(err)
	suite.Equal(customerID, customer.CustomerID)
	suite.ids.customerID = customerID

	// create user
	user := models.User{
		UserID:            uuid.New(),
		Username:          "Username 2",
		Email:             "Email 2",
		Password:          "Password 2",
		FirstName:         "First name 2",
		LastName:          "Last name 2",
		ProfilePictureURL: "Profile pic 2",
		PhoneNumber:       "Phone number 2",
		Address:           "Address",
	}
	userID, err := suite.repo.User.CreateUser(suite.ctx, user)
	suite.NoError(err)
	suite.Equal(userID, user.UserID)
	suite.ids.senderID = userID

	// add payment method
	paymentMethod := models.UserPaymentMethod{
		PaymentMethodID: uuid.New(),
		UserID:          suite.ids.senderID,
		AccountName:     "Account Name 2",
		AccountNumber:   "Account Number 2",
		BankName:        "Bank Name 2",
		BankAddress:     "Bank Address 2",
		SwiftCode:       "Swift Code 2",
	}
	paymentMethodID, err := suite.repo.User.AddPaymentMethod(suite.ctx, paymentMethod)
	suite.NoError(err)
	suite.Equal(paymentMethodID, paymentMethod.PaymentMethodID)
	suite.ids.paymentMethodID = paymentMethodID

	// create invoice
	invoiceID := uuid.New()
	invoice := models.Invoice{
		InvoiceID:          invoiceID,
		InvoiceNumber:      helpers.RandomNumber(1000000000, 9999999999),
		SenderID:           suite.ids.senderID,
		CustomerID:         suite.ids.customerID,
		IssueDate:          time.Now(),
		DueDate:            time.Now().Add(30 * time.Minute),
		TotalAmount:        10000,
		DiscountPercentage: 10,
		DiscountedAmount:   1000,
		FinalAmount:        9000,
		Status:             string(models.InvoiceStatusPaid),
		Currency:           "NGN",
		Notes:              "Thanks for your patronage",
	}

	items := []models.InvoiceItem{
		{
			ItemID:      uuid.New(),
			InvoiceID:   invoiceID,
			Name:        "Item 1",
			Description: "Description 1 ",
			Quantity:    1,
			UnitPrice:   100,
			TotalPrice:  100,
		},
	}

	paymentInfo := models.PaymentInformation{
		PaymentInfoID:   uuid.New(),
		InvoiceID:       invoiceID,
		PaymentMethodID: suite.ids.paymentMethodID,
	}
	id, err := suite.repo.Invoice.CreateInvoice(suite.ctx, invoice, items, suite.ids.customerID, paymentInfo)
	suite.NoError(err)
	suite.Equal(id, invoiceID)
	suite.ids.invoiceID = id

	// add invoice activity
	activity := models.InvoiceActivity{
		ActivityID:  uuid.New(),
		InvoiceID:   suite.ids.invoiceID,
		UserID:      suite.ids.senderID,
		Title:       "Payment Confirmed",
		Description: "You confirmed payment",
	}
	activityID, err := suite.repo.Invoice.AddInvoiceActivity(suite.ctx, activity)
	suite.NoError(err)
	suite.Equal(activityID, activity.ActivityID)
}

func (suite *InvoiceRepoTestSuite) TestGetTotalByStatus() {
	totalAmount, count, err := suite.repo.Invoice.GetTotalByStatus(suite.ctx, models.InvoiceStatusPaid)
	suite.Require().NoError(err)
	suite.Equal(float64(9000), totalAmount)
	suite.Equal(1, count)
}

func (suite *InvoiceRepoTestSuite) TestGetRecentInvoices() {
	invoices, err := suite.repo.Invoice.GetRecentInvoices(suite.ctx, suite.ids.senderID, 5, 0)
	suite.Require().NoError(err)
	suite.Len(invoices, 1)
	suite.NotEmpty(invoices[0])
	suite.Equal(suite.ids.customerID, invoices[0].CustomerID)
	suite.Equal(suite.ids.senderID, invoices[0].SenderID)
	suite.Equal(suite.ids.invoiceID, invoices[0].InvoiceID)
	suite.Equal(float64(10000), invoices[0].TotalAmount)
	suite.Equal(float64(10), invoices[0].DiscountPercentage)
	suite.Equal(float64(1000), invoices[0].DiscountedAmount)
	suite.Equal(float64(9000), invoices[0].FinalAmount)
	suite.Equal(models.InvoiceStatusPaid, models.InvoiceStatus(invoices[0].Status))
	suite.Equal("NGN", invoices[0].Currency)
	suite.Equal("Thanks for your patronage", invoices[0].Notes)
}

func (suite *InvoiceRepoTestSuite) TestGetInvoiceDetails() {
	invoice, err := suite.repo.Invoice.GetInvoiceDetails(suite.ctx, suite.ids.invoiceID)
	suite.Require().NoError(err)
	suite.NotEmpty(invoice)
	suite.Len(invoice.Activities, 2)
	suite.Len(invoice.Items, 1)
	suite.Equal(suite.ids.customerID, invoice.Invoice.CustomerID)
	suite.Equal(suite.ids.senderID, invoice.Invoice.SenderID)
	suite.Equal(suite.ids.invoiceID, invoice.Invoice.InvoiceID)
	suite.Equal(float64(10000), invoice.Invoice.TotalAmount)
	suite.Equal(float64(10), invoice.Invoice.DiscountPercentage)
	suite.Equal(float64(1000), invoice.Invoice.DiscountedAmount)
	suite.Equal(float64(9000), invoice.Invoice.FinalAmount)
	suite.Equal(models.InvoiceStatusPaid, models.InvoiceStatus(invoice.Invoice.Status))
	suite.Equal("NGN", invoice.Invoice.Currency)
	suite.Equal("Thanks for your patronage", invoice.Invoice.Notes)
}

func (suite *InvoiceRepoTestSuite) TestGetRecentActivities() {
	activities, err := suite.repo.Invoice.GetRecentActivities(suite.ctx, suite.ids.senderID, 5, 0)
	suite.Require().NoError(err)
	suite.Len(activities, 1)
	suite.NotEmpty(activities[0])
	suite.Equal(suite.ids.senderID, activities[0].UserID)
	suite.Equal("Invoice Creation", activities[0].Title)
}

func (suite *InvoiceRepoTestSuite) TestGetInvoiceActivities() {
	activities, err := suite.repo.Invoice.GetInvoiceActivities(suite.ctx, suite.ids.senderID, suite.ids.invoiceID, 5, 0)
	suite.Require().NoError(err)
	suite.Len(activities, 2)
	suite.NotEmpty(activities[0])
	suite.NotEmpty(activities[1])
	suite.Equal(suite.ids.invoiceID, activities[1].InvoiceID)
	suite.Equal("Payment Confirmed", activities[0].Title)
	suite.Equal("Invoice Creation", activities[1].Title)
}

func TestInvoiceRepoSuite(t *testing.T) {
	suite.Run(t, new(InvoiceRepoTestSuite))
}
