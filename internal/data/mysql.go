package data

import (
	"github.com/Krados/shortenurl/internal/conf"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func NewMySQL(cfg *conf.Config) *gorm.DB {
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN: cfg.Data.Mysql.DSN, // DSN data source name
	}), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // use singular table name, table for `User` would be `user` with this option enabled
		},
	})
	if err != nil {
		panic("failed to connect database")
	}
	return db
}
