package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoad(t *testing.T) {
	t.Run("load config with all fields", func(t *testing.T) {
		config := Load("production", "localhost:8080", "postgres://user:pass@host:5432/db")

		require.Equal(t, "production", config.Environment)
		require.Equal(t, "localhost:8080", config.HTTPServerAddr)
		require.Equal(t, "postgres://user:pass@host:5432/db", config.DSN) 
	})

	t.Run("load config with empty fields", func(t *testing.T) {
		config := Load("", "", "")

		require.Empty(t, config.Environment)
		require.Empty(t, config.HTTPServerAddr)
		require.Empty(t, config.DSN)

	})

	t.Run("load config with mixed empty and non-empty fields", func(t *testing.T) {
		config := Load("development", "", "sqlite://test.db")

		require.Equal(t, "development", config.Environment)
		require.Empty(t, config.HTTPServerAddr)
		require.Equal(t, "sqlite://test.db", config.DSN) 
	})

	t.Run("load config with special characters", func(t *testing.T) {
		config := Load("test!@#$%^&*()", "127.0.0.1:3000", "mysql://root:p@ssw0rd@localhost/testdb")

		require.Equal(t, "test!@#$%^&*()", config.Environment)
		require.Equal(t, "127.0.0.1:3000", config.HTTPServerAddr)
		require.Equal(t, "mysql://root:p@ssw0rd@localhost/testdb", config.DSN) 
	})
}
