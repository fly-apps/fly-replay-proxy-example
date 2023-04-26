package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
)

// Details on fly-replay:
// https://fly.io/docs/reference/dynamic-request-routing/
func RunServer(addr string) {
	log.Printf("running web server, listening on %s", addr)
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		// We'll find a customer based on `request.Host`
		customer, err := Find(request.Host)

		if err != nil {
			log.Printf("error: %v", err)
			writer.Header().Set("fly-replay", "app=our-default-app")
		} else {
			// Tell Fly-Proxy to replay the request on customer's app
			//writer.Header().Set("fly-replay", fmt.Sprintf("app=%s", customer.App))

			// Or instead, perhaps send to a specific VM
			writer.Header().Set("fly-replay", fmt.Sprintf("instance=%s", customer.Instance))
		}
	})

	err := http.ListenAndServe(addr, nil)

	if errors.Is(err, http.ErrServerClosed) {
		log.Println("serer closing")
	} else if err != nil {
		log.Printf("error starting server: %s", err)
		os.Exit(1)
	}
}
