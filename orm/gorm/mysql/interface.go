package mysql

import "github.com/jinzhu/gorm"

type DBUtil interface {
	CreateDB()
	DropDB()
	GetUtilDB() *gorm.DB
}

type DB interface {
	GetDB() *gorm.DB
	ClearAllData()
}
