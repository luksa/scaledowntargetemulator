package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"log"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func main() {
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc)
	go func() {
		for {
			s := <-sigc
			log.Printf("Received signal %q", s)
		}
	}()

	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}