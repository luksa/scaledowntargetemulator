package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"log"
	"time"
)

var health string = "ok"

func healthHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	newHealth := r.Form.Get("value")
	if newHealth != "" {
		health = newHealth
	}
	w.Write([]byte(health))
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func exitHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Exiting with exit code 0")
	go func() {
		time.Sleep(1000)
		os.Exit(0)
	}()
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
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/exit", exitHandler)
	http.ListenAndServe(":8080", nil)
}