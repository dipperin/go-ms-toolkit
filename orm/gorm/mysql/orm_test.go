package mysql

import (
	"github.com/stretchr/testify/assert"
	"testing"

	db_config "github.com/dipperin/go-ms-toolkit/db-config"
)

func Test_newGormMysql(t *testing.T) {
	conf := db_config.NewDbConfig()
	conf.DbName = "hahaha_test"

	gm := newGormMysql(conf, true)

	assert.NotNil(t, gm)

	gm.CreateDB()
	gm.GetUtilDB()

	gm2 := newGormMysql(conf, false)
	assert.NotNil(t, gm2)

	assert.NotNil(t, gm2.GetDB())

	gm2.ClearAllData()

	gm.DropDB()
}
