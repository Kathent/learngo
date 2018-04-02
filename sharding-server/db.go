package sharding_server

import "database/sql"
import (
	_ "github.com/go-sql-driver/mysql"
	"fmt"
)

var reservedDb *sql.DB

func InitDb() {
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s)/%s", "root", "123456", "192.168.96.204", "robot")
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		panic(err)
	}

	reservedDb = db
}

func GetDb() *sql.DB {
	return reservedDb
}
