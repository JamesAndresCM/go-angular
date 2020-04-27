package database

import "database/sql"

func InitDB() *sql.DB {
	connectionString := "root:@tcp(localhost:3306)/northwind"

	databaseConnection, err := sql.Open("mysql", connectionString)
	if err != nil {
		panic(err.Error())
	}
	return databaseConnection

}
