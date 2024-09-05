package helpers

import (
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zde37/Numeris-Task/internal/models"
)

func TestHashPassword(t *testing.T) {
	t.Run("hash password successfully", func(t *testing.T) {
		password := "securePassword123"
		hashedPassword, err := HashPassword(password)

		require.NoError(t, err)
		require.NotEmpty(t, hashedPassword)
		require.NotEqual(t, password, hashedPassword)
	})

	t.Run("hash empty password", func(t *testing.T) {
		hashedPassword, err := HashPassword("")

		require.NoError(t, err)
		require.NotEmpty(t, hashedPassword)
	})

	t.Run("hash long password", func(t *testing.T) {
		longPassword := strings.Repeat("a", 1000)
		hashedPassword, err := HashPassword(longPassword)

		require.Error(t, err)
		require.Empty(t, hashedPassword)
		require.NotEqual(t, longPassword, hashedPassword)
	})

	t.Run("hash password with special characters", func(t *testing.T) {
		specialPassword := "!@#$%^&*()_+{}[]|\\:;\"'<>,.?/~`"
		hashedPassword, err := HashPassword(specialPassword)

		require.NoError(t, err)
		require.NotEmpty(t, hashedPassword)
		require.NotEqual(t, specialPassword, hashedPassword)
	})

	t.Run("hash unicode password", func(t *testing.T) {
		unicodePassword := "パスワード123"
		hashedPassword, err := HashPassword(unicodePassword)

		require.NoError(t, err)
		require.NotEmpty(t, hashedPassword)
		require.NotEqual(t, unicodePassword, hashedPassword)
	})
}

func TestCheckPassword(t *testing.T) {
	t.Run("correct password", func(t *testing.T) {
		password := "correctPassword123"
		hashedPassword, _ := HashPassword(password)

		err := CheckPassword(password, hashedPassword)

		require.NoError(t, err)
	})

	t.Run("incorrect password", func(t *testing.T) {
		password := "correctPassword123"
		hashedPassword, _ := HashPassword(password)

		err := CheckPassword("wrongPassword", hashedPassword)

		require.Error(t, err)
	})

	t.Run("empty password", func(t *testing.T) {
		hashedPassword, _ := HashPassword("somePassword")

		err := CheckPassword("", hashedPassword)

		require.Error(t, err)
	})

	t.Run("empty hashed password", func(t *testing.T) {
		err := CheckPassword("somePassword", "")

		require.Error(t, err)
	})

	t.Run("invalid hashed password format", func(t *testing.T) {
		err := CheckPassword("somePassword", "invalidHashedPassword")

		require.Error(t, err)
	})

	t.Run("case sensitivity", func(t *testing.T) {
		password := "CaseSensitivePassword"
		hashedPassword, _ := HashPassword(password)

		err := CheckPassword("casesensitivepassword", hashedPassword)

		require.Error(t, err)
	})
}

func TestRandomNumber(t *testing.T) {
	t.Run("generates number within range", func(t *testing.T) {
		min := int64(1)
		max := int64(100)
		result, err := strconv.ParseInt(RandomNumber(min, max), 10, 64)
		require.NoError(t, err)
		require.GreaterOrEqual(t, result, min)
		require.LessOrEqual(t, result, max)
	})

	t.Run("generates different numbers", func(t *testing.T) {
		min := int64(1)
		max := int64(1000000)
		results := make(map[string]bool)
		for i := 0; i < 100; i++ {
			num := RandomNumber(min, max)
			results[num] = true
		}
		require.Greater(t, len(results), 1)
	})

	t.Run("handles min equal to max", func(t *testing.T) {
		min := int64(42)
		max := int64(42)
		result, err := strconv.ParseInt(RandomNumber(min, max), 10, 64)
		require.NoError(t, err)
		require.Equal(t, min, result)
	})

	t.Run("handles negative numbers", func(t *testing.T) {
		min := int64(-100)
		max := int64(-1)
		result, err := strconv.ParseInt(RandomNumber(min, max), 10, 64)
		require.NoError(t, err)
		require.GreaterOrEqual(t, result, min)
		require.LessOrEqual(t, result, max)
	})

	t.Run("handles large numbers", func(t *testing.T) {
		min := int64(1000000000)
		max := int64(9999999999)
		result, err := strconv.ParseInt(RandomNumber(min, max), 10, 64)
		require.NoError(t, err)
		require.GreaterOrEqual(t, result, min)
		require.LessOrEqual(t, result, max)
	})
}

func TestValidateInvoiceStatus(t *testing.T) {
	t.Run("valid status: paid", func(t *testing.T) {
		err := ValidateInvoiceStatus(string(models.InvoiceStatusPaid))
		require.NoError(t, err)
	})

	t.Run("valid status: draft", func(t *testing.T) {
		err := ValidateInvoiceStatus(string(models.InvoiceStatusDraft))
		require.NoError(t, err)
	})

	t.Run("valid status: overdue", func(t *testing.T) {
		err := ValidateInvoiceStatus(string(models.InvoiceStatusOverDue))
		require.NoError(t, err)
	})

	t.Run("valid status: pending", func(t *testing.T) {
		err := ValidateInvoiceStatus(string(models.InvoiceStatusPending))
		require.NoError(t, err)
	})

	t.Run("invalid status", func(t *testing.T) {
		err := ValidateInvoiceStatus("invalid_status")
		require.Error(t, err)
		require.Contains(t, err.Error(), "invalid invoice status: invalid_status")
	})

	t.Run("empty status", func(t *testing.T) {
		err := ValidateInvoiceStatus("")
		require.Error(t, err)
		require.Contains(t, err.Error(), "invalid invoice status: ")
	})

	t.Run("case sensitivity", func(t *testing.T) {
		err := ValidateInvoiceStatus("PAID")
		require.Error(t, err)
		require.Contains(t, err.Error(), "invalid invoice status: PAID")
	})

	t.Run("whitespace in status", func(t *testing.T) {
		err := ValidateInvoiceStatus(" paid ")
		require.Error(t, err)
		require.Contains(t, err.Error(), "invalid invoice status:  paid ")
	})
}
