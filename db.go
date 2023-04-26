package main

import (
	"database/sql"
	"fmt"
	"github.com/go-faker/faker/v4"
	_ "github.com/mattn/go-sqlite3"
	"strings"
)

var db *sql.DB

type Customer struct {
	Id                  int
	Host, App, Instance string
}

func Init() error {
	var err error
	db, err = sql.Open("sqlite3", "customers.db")

	if err != nil {
		return fmt.Errorf("could not init database: %w", err)
	}

	return nil
}

func CreateDB() error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS "customers" (
    "id" integer primary key autoincrement not null,
    "host" varchar not null,
    "app" varchar not null,
    "instance" varchar not null
);`)

	if err != nil {
		return fmt.Errorf("could not create database: %w", err)
	}

	_, err = db.Exec(`CREATE UNIQUE INDEX IF NOT EXISTS "customers_host_unique" on "customers" ("host");`)

	if err != nil {
		return fmt.Errorf("could not create index on table: %w", err)
	}

	return nil
}

func PopulateDB() error {
	insert := "INSERT INTO customers(host, app, instance) VALUES "
	const row = "(?, ?, ?)"
	var inserts []string
	var vals []interface{}

	for i := 0; i < 5000; i++ {
		inserts = append(inserts, row)
		vals = append(vals, faker.DomainName(), fmt.Sprintf("app-%s", faker.Word()), faker.UUIDDigit())
	}

	insert += strings.Join(inserts, ",")

	stmt, err := db.Prepare(insert)

	if err != nil {
		return fmt.Errorf("could not prepare query: %w", err)
	}

	_, err = stmt.Exec(vals...)

	if err != nil {
		return fmt.Errorf("could not insert rows: %w", err)
	}

	return nil
}

func Find(host string) (*Customer, error) {
	row := db.QueryRow("select id, host, app, instance from customers where host = ?", host)

	customer := Customer{}
	err := row.Scan(&customer.Id, &customer.Host, &customer.App, &customer.Instance)

	if err != nil {
		return nil, fmt.Errorf("could not parse query result: %w", err)
	}

	return &customer, nil
}
