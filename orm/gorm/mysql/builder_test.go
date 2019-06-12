package mysql

import (
	"github.com/stretchr/testify/assert"
	"testing"

	db_config "github.com/dipperin/go-ms-toolkit/db-config"
)

func TestMakeDBUtil(t *testing.T) {
	conf := db_config.NewDbConfig()
	conf.DbName = "hahaha_test"

	assert.NotNil(t, MakeDBUtil(conf))
}

func TestMakeDB(t *testing.T) {
	conf := db_config.NewDbConfig()
	conf.DbName = "hahaha_test"

	utilDB := MakeDBUtil(conf)
	assert.NotNil(t, utilDB)

	CreateDB()

	db := MakeDB(conf)
	assert.NotNil(t, db)

	DropDB()
}
