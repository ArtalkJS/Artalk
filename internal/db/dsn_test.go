package db

import (
	"testing"

	"github.com/artalkjs/artalk/v2/internal/config"
)

func TestGetDsnByConf(t *testing.T) {
	testCases := []struct {
		name     string
		conf     config.DBConf
		expected string
	}{
		{
			name: "SQLite configuration",
			conf: config.DBConf{
				Type: config.TypeSQLite,
				File: "test.db",
			},
			expected: "test.db",
		},
		{
			name: "PostgreSQL configuration",
			conf: config.DBConf{
				Type:     config.TypePostgreSQL,
				Host:     "localhost",
				User:     "user",
				Password: "password",
				Name:     "dbname",
				Port:     5432,
				SSL:      false,
			},
			expected: "host=localhost user=user password=password dbname=dbname port=5432 sslmode=disable",
		},
		{
			name: "MySQL configuration",
			conf: config.DBConf{
				Type:     config.TypeMySql,
				Host:     "localhost",
				User:     "user",
				Password: "password",
				Name:     "dbname",
				Port:     3306,
				Charset:  "utf8",
				SSL:      false,
			},
			expected: "user:password@tcp(localhost:3306)/dbname?charset=utf8&parseTime=True&loc=Local&tls=false",
		},
		{
			name: "SQL Server configuration",
			conf: config.DBConf{
				Type:     config.TypeMSSQL,
				Host:     "localhost",
				User:     "user",
				Password: "password",
				Name:     "dbname",
				Port:     3306,
				Charset:  "utf8",
			},
			expected: "sqlserver://user:password@localhost:3306?database=dbname",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := getDsnByConf(tc.conf)
			if result != tc.expected {
				t.Errorf("Expected DSN: %s, but got: %s", tc.expected, result)
			}
		})
	}
}
