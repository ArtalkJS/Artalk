package db

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
	"path/filepath"

	"github.com/artalkjs/artalk/v2/internal/config"
	"github.com/artalkjs/artalk/v2/internal/utils"
	gomysql "github.com/go-sql-driver/mysql"
	"github.com/libtnb/sqlite"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

const mysqlCloudSQLTLSConfigName = "cloudsql"

func OpenSQLite(filename string, gormConfig *gorm.Config) (*gorm.DB, error) {
	if filename == "" {
		return nil, fmt.Errorf("please set `db.file` option in config file to specify a sqlite database path")
	}
	if err := utils.EnsureDir(filepath.Dir(filename)); err != nil {
		return nil, err
	}
	return gorm.Open(sqlite.Open(filename), gormConfig)
}

func OpenMySql(dsn string, gormConfig *gorm.Config, dbConf *config.DBConf) (*gorm.DB, error) {
	tlsConfig, configured, err := loadMySQLTLSConfig(dbConf)
	if err != nil {
		return nil, err
	}
	if configured {
		if err := gomysql.RegisterTLSConfig(mysqlCloudSQLTLSConfigName, tlsConfig); err != nil {
			return nil, fmt.Errorf("failed to register MySQL TLS config: %w", err)
		}
	}
	return gorm.Open(mysql.Open(dsn), gormConfig)
}

func loadMySQLTLSConfig(dbConf *config.DBConf) (*tls.Config, bool, error) {
	paths := []string{dbConf.ServerCaPath, dbConf.ClientCertPath, dbConf.ClientKeyPath}
	configuredPaths := 0
	for _, path := range paths {
		if path != "" {
			configuredPaths++
		}
	}
	if configuredPaths == 0 {
		return nil, false, nil
	}
	if configuredPaths != len(paths) {
		return nil, false, fmt.Errorf("server CA, client certificate, and client key must be configured together")
	}

	serverCertPool := x509.NewCertPool()
	pem, err := os.ReadFile(dbConf.ServerCaPath)
	if err != nil {
		return nil, false, fmt.Errorf("failed to read server CA file: %w", err)
	}
	if ok := serverCertPool.AppendCertsFromPEM(pem); !ok {
		return nil, false, fmt.Errorf("unable to append root cert to pool")
	}

	clientCert, err := tls.LoadX509KeyPair(dbConf.ClientCertPath, dbConf.ClientKeyPath)
	if err != nil {
		return nil, false, fmt.Errorf("failed to load client certificate/key: %w", err)
	}

	return &tls.Config{
		RootCAs:      serverCertPool,
		Certificates: []tls.Certificate{clientCert},
		ServerName:   dbConf.Host,
		MinVersion:   tls.VersionTLS12,
	}, true, nil
}

func OpenPostgreSQL(dsn string, gormConfig *gorm.Config) (*gorm.DB, error) {
	return gorm.Open(postgres.New(postgres.Config{
		DSN: dsn,

		// gorm v2 use `pgx` as postgres’s database/sql driver,
		// it enables prepared statement cache by default,
		// disable it when `PrepareStmt` is false by following code:
		PreferSimpleProtocol: !gormConfig.PrepareStmt,
	}), gormConfig)
}

func OpenSqlServer(dsn string, gormConfig *gorm.Config) (*gorm.DB, error) {
	return gorm.Open(sqlserver.Open(dsn), gormConfig)
}
