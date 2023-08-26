package repository

import (
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
)

const (
	branchTable         = "branch"
	additionalBankTable = "additional_bank"
	bankTable           = "bank"
	contractTable       = "contract"
	clientTable         = "client"
	contractsTable      = "contracts"
	revenueTable        = "revenue"
	usersTable          = "users"
)

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

func NewPostgresDB(cfg Config) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.DBName, cfg.Password, cfg.SSLMode))
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	initDB(db)
	return db, nil
}
func initDB(db *sqlx.DB) {
	dat, err := os.ReadFile("sql/db.sql")
	if err != nil {
		// log.Error("Database", logger.Any("SQLFile", err.Error()))
	}
	_, err = db.Query(string(dat))
	if err != nil {
		// log.Error("Database", logger.Any("query", err.Error()))
	}
}
