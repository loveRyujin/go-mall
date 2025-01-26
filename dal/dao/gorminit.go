package dao

import (
	"github.com/loveRyujin/go-mall/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	// _DbMaster 主库实例
	_DbMaster *gorm.DB
	// _DbSlave 从库实例
	_DbSlave *gorm.DB
)

func DB() *gorm.DB {
	return _DbSlave
}

func DBMaster() *gorm.DB {
	return _DbMaster
}

func init() {
	_DbMaster = initDB(&config.Database.Master)
	_DbSlave = initDB(&config.Database.Slave)
}

func initDB(option *config.DbConnectOption) *gorm.DB {
	db, err := gorm.Open(mysql.Open(option.Dsn), &gorm.Config{
		Logger: NewGormLogger(),
	})
	if err != nil {
		panic(err)
	}
	sqlDb, _ := db.DB()
	// 设置连接池参数
	sqlDb.SetMaxIdleConns(option.MaxIdleConn)
	sqlDb.SetMaxOpenConns(option.MaxOpenConn)
	sqlDb.SetConnMaxLifetime(option.MaxLifeTime)

	if err = sqlDb.Ping(); err != nil {
		panic(err)
	}
	return db
}
