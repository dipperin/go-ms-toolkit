package mysql

import (
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"

	dbConfig "github.com/dipperin/go-ms-toolkit/db-config"
)

func Test_newGormMysql(t *testing.T) {
	conf := dbConfig.NewDbConfig()
	conf.DbName = "hahaha_test"

	gm := newGormMysql(conf, true)
	assert.NotNil(t, gm)
	gm.CreateDB()
	gm.GetUtilDB()
	defer gm.DropDB()

	gm2 := newGormMysql(conf, false)
	assert.NotNil(t, gm2)
	assert.NotNil(t, gm2.GetDB())
	gm2.ClearAllData()
}

func Test_gormMysql_Create(t *testing.T) {
	conf := dbConfig.NewDbConfig()
	conf.DbName = "hahaha_test"

	gm := newGormMysql(conf, true)
	assert.NotNil(t, gm)
	gm.CreateDB()
	gm.GetUtilDB()
	defer gm.DropDB()

	type Test struct {
		Data string `gorm:"data"`
	}
	gm2 := newGormMysql(conf, false)
	gm2.GetDB().AutoMigrate(&Test{})
	assert.NotNil(t, gm2)
	assert.NoError(t, gm2.Create(&Test{}))
	assert.NoError(t, gm2.Create([]*Test{}))
	assert.NoError(t, gm2.Create([]*Test{{Data: "1"}, {Data: "2"}, {Data: "3"}}))
	gm2.ClearAllData()
}
