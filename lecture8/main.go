package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

func insertData(db *sql.DB) error {
	numbers := []string{
		"1234567890", "123 456 7891", "(123) 456 7892", "(123) 456-7893", "123-456-7894", "123-456-7890", "1234567892", "(123)456-7892",
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	for i, _ := range numbers {
		_, err = db.Query("INSERT INTO phone_number (number) VALUES (?)", numbers[i])
		if err != nil {
			return err
		}
	}

	tx.Commit()

	return nil
}

func main() {
	db, err := sql.Open("mysql", "root@/go_prac")
	if err != nil {
		panic(err)
	}

	createTableQuery := `CREATE TABLE IF NOT EXISTS phone_number( 
		number VARCHAR(20) NOT NULL,
		PRIMARY KEY(number)
	)`
	_, err = db.Query(createTableQuery)
	if err != nil {
		panic(err)
	}

	/*
		err = insertData(db)
		if err != nil {
			panic(err)
		}
	*/

	selectQuery := "SELECT number FROM phone_number"
	rows, err := db.Query(selectQuery)
	if err != nil {
		panic(err)
	}

	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		var number string
		err = rows.Scan(&number)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Number: %s\t", number)

		var newNumber []byte
		for i, _ := range number {
			if number[i] >= '0' && number[i] <= '9' {
				newNumber = append(newNumber, number[i])
			}
		}
		fmt.Printf("newNumber: %s\n", newNumber)
		_, err = db.Query("UPDATE phone_number SET number = ? WHERE number = ?", newNumber, number)
		if err != nil {
			_, err = db.Query("DELETE FROM phone_number WHERE number = ?", number)
			if err != nil {
				panic(err)
			}
		}
	}

	tx.Commit()
}
