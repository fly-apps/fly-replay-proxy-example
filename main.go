package main

import (
	"flag"
	"log"
	"os"
)

func main() {
	var create, populate, run bool
	var addr string
	var somethingToDo bool = false

	// Create DB
	flag.BoolVar(&create, "create", false, "create sqlite database")

	// Populate DB
	flag.BoolVar(&populate, "populate", false, "populate sqlite database with fake data")

	// Run server
	flag.BoolVar(&run, "run", false, `run the "proxy" web server"`)
	flag.StringVar(&addr, "addr", ":8080", "the address to listen on")

	flag.Parse()

	err := Init()

	if err != nil {
		log.Println("could not init db: %v", err)
		os.Exit(1)
	}

	if create {
		somethingToDo = true
		err := CreateDB()

		if err != nil {
			log.Printf("error: %v", err)
			os.Exit(1)
		}
	}

	if populate {
		// todo: Populate the DB (ensure it exists first)
		somethingToDo = true
		err := PopulateDB()

		if err != nil {
			log.Printf("could not populate database: %v", err)
			os.Exit(1)
		}

	}

	if run {
		somethingToDo = true
		RunServer(addr)
	}

	if !somethingToDo {
		log.Println("we didn't find anything to do, how boring")
	}

}
