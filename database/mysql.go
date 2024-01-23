package database

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"time"
	"togo/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type mysqlDatabase struct {
	DB *gorm.DB
}

func NewMySQLDatabase(cfg *config.Config) Database {
	db, err := openConnection(cfg)
	if err != nil {
		log.Fatalf("Error connecting to MySQL: %s", err)
	}

	if err := createDatabase(db, cfg.Database.DBName); err != nil {
		log.Fatalf("Error creating database: %s", err)
	}

	db, err = connectToDatabase(cfg)
	if err != nil {
		log.Fatalf("Error connecting to the database: %s", err)
	}

	return &mysqlDatabase{DB: db}
}

func openConnection(cfg *config.Config) (*gorm.DB, error) {
	dsn := getDSN(cfg, "")
	return gorm.Open(mysql.Open(dsn), &gorm.Config{})
}

func createDatabase(db *gorm.DB, dbName string) error {
	return db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s;", dbName)).Error
}

func connectToDatabase(cfg *config.Config) (*gorm.DB, error) {
	dsn := getDSN(cfg, cfg.Database.DBName)
	logger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Warn,
			IgnoreRecordNotFoundError: true,
			ParameterizedQueries:      true,
			Colorful:                  false,
		},
	)
	return gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: logger})
}

func getDSN(cfg *config.Config, dbName string) string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=%s",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		dbName,
		url.QueryEscape(cfg.Database.Timezone),
	)
}

func (m *mysqlDatabase) GetConn() *gorm.DB {
	return m.DB
}
