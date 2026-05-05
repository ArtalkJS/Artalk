package db

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/artalkjs/artalk/v2/internal/config"
	"github.com/artalkjs/artalk/v2/internal/utils"
	gomysql "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

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
	var serverCertPool *x509.CertPool
	var clientCert *tls.Certificate

	if dbConf.ServerCaPath != "" {
		serverCertPool = x509.NewCertPool()
		pem, err := os.ReadFile(dbConf.ServerCaPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read server CA file: %w", err)
		}
		if ok := serverCertPool.AppendCertsFromPEM(pem); !ok {
			return nil, fmt.Errorf("unable to append root cert to pool")
		}
	}

	if dbConf.ClientCertPath != "" && dbConf.ClientKeyPath != "" {
		cert, err := tls.LoadX509KeyPair(dbConf.ClientCertPath, dbConf.ClientKeyPath)
		if err != nil {
			return nil, fmt.Errorf("failed to load client certificate/key: %w", err)
		}
		clientCert = &cert
	}

	// Register custom TLS configurations

	// Handle Google Cloud SQL:
	// - https://cloud.google.com/sql/docs/mysql/samples/cloud-sql-mysql-databasesql-connect-tcp-sslcerts
	// - https://cloud.google.com/sql/docs/mysql/configure-ssl-instance#enforcing-ssl
	if serverCertPool != nil && clientCert != nil {
		gomysql.RegisterTLSConfig("cloudsql", &tls.Config{
			RootCAs:            serverCertPool,
			Certificates:       []tls.Certificate{*clientCert},
			InsecureSkipVerify: true,
			VerifyPeerCertificate: func(rawCerts [][]byte, _ [][]*x509.Certificate) error {
				if len(rawCerts) == 0 {
					return errors.New("no certificates available to verify")
				}
				cert, err := x509.ParseCertificate(rawCerts[0])
				if err != nil {
					return err
				}
				opts := x509.VerifyOptions{Roots: serverCertPool}
				if _, err = cert.Verify(opts); err != nil {
					return err
				}
				return nil
			},
		})
	}
	return gorm.Open(mysql.Open(dsn), gormConfig)
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
