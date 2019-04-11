package mysql

import (
	"fmt"
	"github.com/dipperin/go-ms-toolkit/db-config"
	"github.com/dipperin/go-ms-toolkit/log"
	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
	"strings"
	"time"
)

type gormMysql struct {
	dbConfig *db_config.DbConfig
	db       *gorm.DB
	utilDB   *gorm.DB
}

func (gm *gormMysql) CreateDB() {
	createDbSQL := "CREATE DATABASE IF NOT EXISTS " + gm.dbConfig.DbName + " DEFAULT CHARSET utf8 COLLATE utf8_general_ci;"

	err := gm.utilDB.Exec(createDbSQL).Error
	if err != nil {
		fmt.Println("创建失败：" + err.Error() + " sql:" + createDbSQL)
		return
	}
	fmt.Println(gm.dbConfig.DbName + "数据库创建成功")
}

func (gm *gormMysql) DropDB() {
	dropDbSQL := "DROP DATABASE IF EXISTS " + gm.dbConfig.DbName + ";"

	err := gm.utilDB.Exec(dropDbSQL).Error
	if err != nil {
		fmt.Println("删除失败：" + err.Error() + " sql:" + dropDbSQL)
		return
	}
	fmt.Println(gm.dbConfig.DbName + "数据库删除成功")
}

func (gm *gormMysql) GetDB() *gorm.DB {
	return gm.db
}

func (gm *gormMysql) GetUtilDB() *gorm.DB {
	log.QyLogger.Info("init db connection: ", zap.String("db_host", gm.dbConfig.Host),
		zap.String("db_name", gm.dbConfig.DbName), zap.String("user", gm.dbConfig.Username))

	openedDb, err := gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", gm.dbConfig.Username, gm.dbConfig.Password, gm.dbConfig.Host, gm.dbConfig.Port, gm.dbConfig.DbName))
	if err != nil {
		panic("数据库连接出错：" + err.Error())
	}
	openedDb.DB().SetMaxIdleConns(gm.dbConfig.MaxIdleConns)
	openedDb.DB().SetMaxOpenConns(gm.dbConfig.MaxOpenConns)
	// 避免久了不使用，导致连接被mysql断掉的问题
	openedDb.DB().SetConnMaxLifetime(time.Hour * 2)
	// 如果不是生产数据库则打开详细日志
	//if !strings.Contains(dbConfig.DbName, "prod") {
	if substr(gm.dbConfig.DbName, len(gm.dbConfig.DbName)-4, 4) != "prod" {
		openedDb.LogMode(true)
	}

	return openedDb
}

func (gm *gormMysql) ClearAllData() {
	if strings.Contains(gm.dbConfig.DbName, "test") {
		tmpDb := gm.db
		if tmpDb == nil {
			panic("尚未初始化数据库, 清空数据库失败")
		}
		if rs, err := tmpDb.Raw("show tables;").Rows(); err == nil {
			var tName string
			for rs.Next() {
				if err := rs.Scan(&tName); err != nil || tName == "" {
					fmt.Println("表名获取失败", err, tName)
					panic("表名获取失败")
				}
				if err := tmpDb.Exec(fmt.Sprintf("delete from %s", tName)).Error; err != nil {
					panic("清空表数据失败:" + err.Error())
				}
			}
		} else {
			panic("表名列表获取失败：" + err.Error())
		}
	} else {
		panic("非法操作！在非测试环境下调用了清空所有数据的方法")
	}
}

func newGormMysql(dbConfig *db_config.DbConfig, forUtil bool) *gormMysql {
	gm := &gormMysql{dbConfig: dbConfig}

	if forUtil {
		gm.initCdDb()
		return gm
	}

	// init db
	gm.initGormDB()

	return gm
}

func (gm *gormMysql) initGormDB() {
	if gm.db != nil {
		panic("gorm db should nil")
	}

	log.QyLogger.Info("init db connection: ", zap.String("db_host", gm.dbConfig.Host),
		zap.String("db_name", gm.dbConfig.DbName), zap.String("user", gm.dbConfig.Username))

	openedDb, err := gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", gm.dbConfig.Username, gm.dbConfig.Password, gm.dbConfig.Host, gm.dbConfig.Port, gm.dbConfig.DbName))
	if err != nil {
		panic("数据库连接出错：" + err.Error())
	}
	openedDb.DB().SetMaxIdleConns(gm.dbConfig.MaxIdleConns)
	openedDb.DB().SetMaxOpenConns(gm.dbConfig.MaxOpenConns)
	// 避免久了不使用，导致连接被mysql断掉的问题
	openedDb.DB().SetConnMaxLifetime(time.Hour * 2)
	// 如果不是生产数据库则打开详细日志
	//if !strings.Contains(dbConfig.DbName, "prod") {
	if substr(gm.dbConfig.DbName, len(gm.dbConfig.DbName)-4, 4) != "prod" {
		openedDb.LogMode(true)
	}

	gm.db = openedDb
}

func (gm *gormMysql) initCdDb() {
	if gm.db != nil {
		panic("gorm db should nil")
	}

	cStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", gm.dbConfig.Username, gm.dbConfig.Password, gm.dbConfig.Host, gm.dbConfig.Port, "information_schema")
	openedDb, err := gorm.Open("mysql", cStr)
	if err != nil {
		fmt.Println(cStr)
		panic("连接数据库出错:" + err.Error())
	}

	gm.utilDB = openedDb
}

func substr(str string, start, length int) string {
	rs := []rune(str)
	rl := len(rs)
	end := 0

	if start < 0 {
		start = rl - 1 + start
	}
	end = start + length

	if start > end {
		start, end = end, start
	}

	if start < 0 {
		start = 0
	}
	if start > rl {
		start = rl
	}
	if end < 0 {
		end = 0
	}
	if end > rl {
		end = rl
	}
	return string(rs[start:end])
}
