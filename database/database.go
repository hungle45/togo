package database

import "gorm.io/gorm"

type Database interface {
	GetConn() *gorm.DB
}
