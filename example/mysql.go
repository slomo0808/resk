package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"time"
)

func main() {
	db, err := sql.Open("mysql",
		"root:123456@tcp(127.0.0.1:3306)/po?charset=utf8&parseTime=true&loc=Local")
	defer db.Close()
	if err != nil {
		log.Fatal(err)
	}
	db.SetConnMaxLifetime(7 * time.Hour)
	db.SetMaxIdleConns(2)
	db.SetMaxOpenConns(5)

	fmt.Println(db.Query("select now()"))
}
