package camdb

import (
	"go.uber.org/zap"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type Config struct {
	Driver   string `default:"mysql"`
	Host     string
	Database string
	Username string
	Password string
}

type DB struct {
	db  *sqlx.DB
	log *zap.Logger
}

func New(config *Config, logger *zap.Logger) (*DB, error) {
	switch config.Driver {
	case "mysql":
		return newMySQL(config, logger)
	default:
		return nil, &ErrorUnsupportedDriver{msg: config.Driver}
	}
}

func newMySQL(config *Config, logger *zap.Logger) (*DB, error) {
	mysqlConfig := mysql.Config{
		User:   config.Username,
		Passwd: config.Password,
		Net:    "tcp",
		Addr:   config.Host,
		DBName: config.Database,
	}
	db, err := sqlx.Connect(config.Driver, mysqlConfig.FormatDSN())
	if err != nil {
		return nil, err
	}
	return &DB{
		db:  db,
		log: logger,
	}, err
}
