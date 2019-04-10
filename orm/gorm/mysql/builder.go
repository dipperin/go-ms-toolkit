package mysql

import "github.com/dipperin/go-ms-toolkit/db-config"

func MakeDBUtil(dbConfig *db_config.DbConfig) DBUtil {
	return newGormMysql(dbConfig, true)
}

func MakeDB(dbConfig *db_config.DbConfig) DB {
	return newGormMysql(dbConfig, false)
}
