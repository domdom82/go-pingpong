package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptrace"
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
	client := &http.Client{}

	trace := &httptrace.ClientTrace{
		GotConn: func(connInfo httptrace.GotConnInfo) {
			fmt.Printf("-> %s (reused: %v, wasIdle: %v, idleTime: %v)\n", connInfo.Conn.RemoteAddr(), connInfo.Reused, connInfo.WasIdle, connInfo.IdleTime)
		},
	}
	for {
		fmt.Printf("-> %s\n", endpoint)
		start := time.Now()
		req, err := http.NewRequest("GET", endpoint, nil)
		req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))
		if err != nil {
			panic(err)
		}
		r, err := client.Do(req)
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
