package helpers

import (
	"fmt"
	"strconv"
	"time"

	"github.com/zde37/Numeris-Task/internal/models"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/exp/rand"
)

// HashPassword generates a bcrypt hash of the provided password string. It returns the hashed password and an error if the hashing fails.
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hashedPassword), nil
}

// CheckPassword compares a plaintext password with a hashed password and returns an error if they do not match. 
func CheckPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// RandomNumber generates a random number between min and max (inclusive)
func RandomNumber(min, max int64) string {
	source := rand.NewSource(uint64(time.Now().UnixNano()))
	r := rand.New(source)
	result := r.Int63n(max-min+1) + min
	return strconv.FormatInt(result, 10)
}

// ValidateInvoiceStatus checks if the provided invoice status is one of the valid statuses (paid, draft, overdue, or pending)
func ValidateInvoiceStatus(status string) error {
	if status != string(models.InvoiceStatusPaid) && status != string(models.InvoiceStatusDraft) &&
		status != string(models.InvoiceStatusOverDue) && status != string(models.InvoiceStatusPending) {
		return fmt.Errorf("invalid invoice status: %s", status)
	}
	return nil
}
