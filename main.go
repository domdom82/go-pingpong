package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

func main() {

	url := "http://localhost:8080"
	port, portSet := os.LookupEnv("PORT")

	if len(os.Args) == 2 {
		url = os.Args[1]
	}

	if !portSet {
		port = "8080"
	}

	// start server
	go func() {
		mux := http.NewServeMux()
		mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte("pong"))
		})

		err := http.ListenAndServe(fmt.Sprintf(":%s", port), mux)

		if err != nil {
			panic(err)
		}
	}()

	endpoint := fmt.Sprintf("%s/ping", url)

	for {
		fmt.Printf("-> %s\n", endpoint)
		start := time.Now()
		r, err := http.Get(endpoint)
		stop := time.Now()
		duration := stop.Sub(start)
		if err != nil {
			fmt.Printf("%s (%v)\n", err, duration)
		} else {
			data, err := ioutil.ReadAll(r.Body)
			if err != nil {
				fmt.Printf("%s (%v)\n", err, duration)
			} else {
				fmt.Printf("<- %s (%v)\n", string(data), duration)
			}
		}
		time.Sleep(1 * time.Second)
	}
}
