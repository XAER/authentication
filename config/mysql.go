package config

import (
	"authentication/helpers"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	gorm_sql "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func GetConnection() (*gorm.DB, error) {
	sql_host := helpers.GetEnv("MYSQL_HOST", "error")
	sql_port := helpers.GetEnv("MYSQL_PORT", "3306")
	sql_user := helpers.GetEnv("MYSQL_USER", "error")
	sql_pass := helpers.GetEnv("MYSQL_PASSWORD", "error")
	sql_db := helpers.GetEnv("MYSQL_DB", "error")

	// Format for the SQL CONNECTION: DNS
	// <user>:<password>@tcp(<host>:<port>)/<dbname>
	sqlConnectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", sql_user, sql_pass, sql_host, sql_port, sql_db)

	db, err := gorm.Open(gorm_sql.Open(sqlConnectionString), &gorm.Config{})

	if err != nil {
		log.Fatal(err.Error())
		return nil, err
	}

	return db, nil
}
