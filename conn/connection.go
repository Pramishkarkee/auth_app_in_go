package mydb

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"log"
	"fmt"
)

func GetDatabase() *gorm.DB {
	databasename := "golang"
	database := "postgres"
	databasepassword := "root1234"
	databaseurl := "postgres://root:" + databasepassword + "@localhost/" + databasename + "?sslmode=disable"

	connection, err := gorm.Open(database, databaseurl)
	if err != nil {
		log.Fatalln("Invalid database url")
	}
	sqldb := connection.DB()

	err = sqldb.Ping()
	if err != nil {
		log.Fatal("Database connected")
	}
	fmt.Println("Database connection successful.")
	return connection
}

func CloseDatabase(connection *gorm.DB) {
	sqldb := connection.DB()
	sqldb.Close()
}