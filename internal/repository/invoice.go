package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/zde37/Numeris-Task/internal/models"
)

type invoiceRepoImpl struct {
	DBPool *pgxpool.Pool
}

// newInvoiceRepoImpl creates a new instance of the invoiceRepoImpl struct, which is used to interact with the
// invoice-related data in the database.  
func newInvoiceRepoImpl(dbPool *pgxpool.Pool) *invoiceRepoImpl {
	return &invoiceRepoImpl{
		DBPool: dbPool,
	}
}

// CreateInvoice creates a new invoice in the database, including the invoice details, invoice items, payment information, and related activities. 
func (i *invoiceRepoImpl) CreateInvoice(ctx context.Context, invoice models.Invoice, items []models.InvoiceItem, customerID uuid.UUID, paymentInfo models.PaymentInformation) (uuid.UUID, error) {
	tx, err := i.DBPool.Begin(ctx)
	if err != nil {
		return uuid.Nil, err
	}
	defer tx.Rollback(ctx)

	// insert invoice
	query1 := `
        INSERT INTO invoices (invoice_id, invoice_number, sender_id, customer_id, issue_date, due_date, 
                              total_amount, discount_percentage, discounted_amount, final_amount, status, 
                              currency, notes)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
        RETURNING invoice_id`

	err = tx.QueryRow(ctx, query1,
		invoice.InvoiceID, invoice.InvoiceNumber, invoice.SenderID, customerID,
		invoice.IssueDate, invoice.DueDate, invoice.TotalAmount, invoice.DiscountPercentage,
		invoice.DiscountedAmount, invoice.FinalAmount, invoice.Status, invoice.Currency, invoice.Notes,
	).Scan(&invoice.InvoiceID)
	if err != nil {
		return uuid.Nil, err
	}

	// insert invoice items
	for _, item := range items {
		_, err = tx.Exec(ctx, `
            INSERT INTO invoice_items (item_id, invoice_id, name, description, quantity, unit_price, total_price)
            VALUES ($1, $2, $3, $4, $5, $6, $7)`,
			item.ItemID, invoice.InvoiceID, item.Name, item.Description, item.Quantity, item.UnitPrice, item.TotalPrice,
		)
		if err != nil {
			return uuid.Nil, err
		}
	}

	// insert payment information
	_, err = tx.Exec(ctx, `
        INSERT INTO payment_information (payment_info_id, invoice_id, payment_method_id)
        VALUES ($1, $2, $3)`,
		paymentInfo.PaymentInfoID, invoice.InvoiceID, paymentInfo.PaymentMethodID,
	)
	if err != nil {
		return uuid.Nil, err
	}

	// create invoice activity
	activityID := uuid.New()
	_, err = tx.Exec(ctx, `
        INSERT INTO invoice_activities (activity_id, invoice_id, user_id, title, description)
        VALUES ($1, $2, $3, $4, $5)`,
		activityID, invoice.InvoiceID, invoice.SenderID, "Invoice Creation", fmt.Sprintf("Created invoice %s", invoice.InvoiceNumber),
	)
	if err != nil {
		return uuid.Nil, err
	}

	// add to recent activities
	_, err = tx.Exec(ctx, `
        INSERT INTO recent_activities (activity_id, user_id, title, description)
        VALUES ($1, $2, $3, $4)`,
		activityID, invoice.SenderID, "Invoice Creation", fmt.Sprintf("Created invoice %s", invoice.InvoiceNumber),
	)
	if err != nil {
		return uuid.Nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return uuid.Nil, err
	}

	return invoice.InvoiceID, nil
}

// GetInvoiceDetails retrieves the details of an invoice, including the invoice information, invoice items, and invoice activities. 
func (i *invoiceRepoImpl) GetInvoiceDetails(ctx context.Context, invoiceID uuid.UUID) (*models.InvoiceDetails, error) {
	var details models.InvoiceDetails

	// get invoice information
	err := i.DBPool.QueryRow(ctx, `
        SELECT i.invoice_id, i.invoice_number, i.sender_id, i.customer_id, i.issue_date, i.due_date, 
               i.total_amount, i.discount_percentage, i.discounted_amount, i.final_amount, i.status, 
               i.currency, i.notes, i.created_at, i.updated_at,
               s.first_name || ' ' || s.last_name AS sender_name, s.email AS sender_email, s.phone_number AS sender_phone_number, s.address AS sender_address,
               c.name AS customer_name, c.email AS customer_email, c.phone_number AS customer_phone_number,
               pm.payment_method_id, pm.user_id, pm.account_name, pm.account_number, pm.bank_name, pm.bank_address, pm.swift_code
        FROM invoices i
        JOIN users s ON i.sender_id = s.user_id
        JOIN customers c ON i.customer_id = c.customer_id
        LEFT JOIN payment_information pi ON i.invoice_id = pi.invoice_id
        LEFT JOIN user_payment_methods pm ON pi.payment_method_id = pm.payment_method_id
        WHERE i.invoice_id = $1`,
		invoiceID,
	).Scan(
		&details.Invoice.InvoiceID, &details.Invoice.InvoiceNumber, &details.Invoice.SenderID, &details.Invoice.CustomerID,
		&details.Invoice.IssueDate, &details.Invoice.DueDate, &details.Invoice.TotalAmount, &details.Invoice.DiscountPercentage,
		&details.Invoice.DiscountedAmount, &details.Invoice.FinalAmount, &details.Invoice.Status, &details.Invoice.Currency,
		&details.Invoice.Notes, &details.Invoice.CreatedAt, &details.Invoice.UpdatedAt, &details.SenderName, &details.SenderEmail,
		&details.SenderPhoneNumber, &details.SenderAddress, &details.CustomerName, &details.CustomerEmail, &details.CustomerPhoneNumber,
		&details.PaymentInformation.PaymentMethodID, &details.PaymentInformation.UserID, &details.PaymentInformation.AccountName,
		&details.PaymentInformation.AccountNumber, &details.PaymentInformation.BankName, &details.PaymentInformation.BankAddress,
		&details.PaymentInformation.SwiftCode,
	)
	if err != nil {
		return nil, err
	}

	// get invoice items
	rows, err := i.DBPool.Query(ctx, `
        SELECT item_id, invoice_id, name, description, quantity, unit_price, total_price
        FROM invoice_items
        WHERE invoice_id = $1`,
		invoiceID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item models.InvoiceItem
		err := rows.Scan(&item.ItemID, &item.InvoiceID, &item.Name, &item.Description, &item.Quantity, &item.UnitPrice, &item.TotalPrice)
		if err != nil {
			return nil, err
		}
		details.Items = append(details.Items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	// get invoice activities
	rows, err = i.DBPool.Query(ctx, `
        SELECT activity_id, invoice_id, user_id, title, description, created_at
        FROM invoice_activities
        WHERE invoice_id = $1
        ORDER BY created_at`,
		invoiceID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var activity models.InvoiceActivity
		err := rows.Scan(&activity.ActivityID, &activity.InvoiceID, &activity.UserID, &activity.Title, &activity.Description, &activity.CreatedAt)
		if err != nil {
			return nil, err
		}
		details.Activities = append(details.Activities, activity)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &details, nil
}

// AddInvoiceActivity adds a new activity to an invoice. 
func (i *invoiceRepoImpl) AddInvoiceActivity(ctx context.Context, activity models.InvoiceActivity) (uuid.UUID, error) {
	query := `
        INSERT INTO invoice_activities (activity_id, invoice_id, user_id, title, description)
        VALUES ($1, $2, $3, $4, $5)
		RETURNING activity_id
	`
	err := i.DBPool.QueryRow(ctx, query, activity.ActivityID, activity.InvoiceID, activity.UserID,
		activity.Title, activity.Description).Scan(&activity.ActivityID)
	if err != nil {
		return uuid.Nil, err
	}
	return activity.ActivityID, nil
}

// GetTotalByStatus retrieves the total amount and count of invoices with the specified status. 
func (i *invoiceRepoImpl) GetTotalByStatus(ctx context.Context, status models.InvoiceStatus) (totalAmount float64, count int, err error) {
	query := `SELECT COUNT(*) as count, COALESCE(SUM(final_amount), 0) as total_amount FROM invoices WHERE status = $1`

	err = i.DBPool.QueryRow(ctx, query, status).Scan(&count, &totalAmount)
	if err != nil {
		return 0, 0, err
	}

	return totalAmount, count, nil
}

// GetRecentInvoices retrieves a list of the most recent invoices for the specified sender, with optional pagination. 
func (i *invoiceRepoImpl) GetRecentInvoices(ctx context.Context, senderID uuid.UUID, limit, offset int32) ([]models.Invoice, error) {
	query := `
        SELECT invoice_id, invoice_number, sender_id, customer_id, issue_date, due_date, 
               total_amount, discount_percentage, discounted_amount, final_amount, status, 
               currency, notes, created_at, updated_at 
        FROM invoices 
        WHERE sender_id = $1 
        ORDER BY created_at DESC 
        LIMIT $2 OFFSET $3`

	rows, err := i.DBPool.Query(ctx, query, senderID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	invoices := []models.Invoice{}
	for rows.Next() {
		var invoice models.Invoice
		err := rows.Scan(
			&invoice.InvoiceID, &invoice.InvoiceNumber, &invoice.SenderID, &invoice.CustomerID,
			&invoice.IssueDate, &invoice.DueDate, &invoice.TotalAmount, &invoice.DiscountPercentage,
			&invoice.DiscountedAmount, &invoice.FinalAmount, &invoice.Status, &invoice.Currency,
			&invoice.Notes, &invoice.CreatedAt, &invoice.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		invoices = append(invoices, invoice)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return invoices, nil
}

// GetRecentActivities retrieves a list of recent activities for the specified user, with pagination. 
func (i *invoiceRepoImpl) GetRecentActivities(ctx context.Context, userID uuid.UUID, limit, offset int32) ([]models.RecentActivity, error) {
	query := `
        SELECT activity_id, user_id, title, description, created_at 
        FROM recent_activities 
        WHERE user_id = $1 
        ORDER BY created_at DESC 
        LIMIT $2 OFFSET $3`

	rows, err := i.DBPool.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	activities := []models.RecentActivity{}
	for rows.Next() {
		var activity models.RecentActivity
		err := rows.Scan(
			&activity.ActivityID, &activity.UserID, &activity.Title, &activity.Description, &activity.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		activities = append(activities, activity)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return activities, nil
}

// GetInvoiceActivities retrieves the recent activities associated with a specific invoice for a given user. 
func (i *invoiceRepoImpl) GetInvoiceActivities(ctx context.Context, userID, invoiceID uuid.UUID, limit, offset int32) ([]models.InvoiceActivity, error) {
	query := `
		SELECT activity_id, invoice_id, user_id, title, description, created_at 
		FROM invoice_activities 
		WHERE user_id = $1 AND invoice_id = $2
		ORDER BY created_at DESC 
		LIMIT $3 OFFSET $4
	`

	rows, err := i.DBPool.Query(ctx, query, userID, invoiceID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	activities := []models.InvoiceActivity{}
	for rows.Next() {
		var activity models.InvoiceActivity
		if err := rows.Scan(&activity.ActivityID, &activity.InvoiceID, &activity.UserID, &activity.Title, &activity.Description, &activity.CreatedAt); err != nil {
			return nil, err
		}
		activities = append(activities, activity)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return activities, nil
}
