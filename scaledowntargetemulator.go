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

func addSignalTrap() {
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc)
	go func() {
		for {
			s := <-sigc
			log.Printf("Received signal %q", s)
		}
	}()
}

func main() {
	addSignalTrap()
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}