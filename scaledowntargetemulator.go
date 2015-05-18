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
		log.Printf("Health changed to %q", health)
	} else {
		log.Printf("Health requested. Returning '%q'", health)
	}
	w.Write([]byte(health))
}

func handler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Default handler invoked")
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func exitHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var delay time.Duration
	delayStr := r.Form.Get("delay")
	if delayStr == "" {
		delay = time.Second
	} else {
		var err error
		delay, err = time.ParseDuration(delayStr)
		if err != nil {
			delay = time.Second
		}
	}
	log.Printf("Exiting with exit code 0 after %q s", delay.Seconds())
	fmt.Fprintf(w, "Exiting with exit code 0 after %f s", delay.Seconds())

	go func() {
		time.Sleep(delay)
		log.Printf("Exiting now")
		os.Exit(0)
	}()
}

func preStopHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var delay time.Duration
	delayStr := r.Form.Get("delay")
	if delayStr == "" {
		delay = time.Second
	} else {
		var err error
		delay, err = time.ParseDuration(delayStr)
		if err != nil {
			delay = time.Second
		}
	}
	log.Printf("PreStop handler invoked. The HTTP response will be returned in %q s", delay.Seconds())
	fmt.Fprintf(w, "preStop handler invoked; waiting for %f s", delay.Seconds())

	networkCheckHandler(w, r)
	time.Sleep(delay)
	networkCheckHandler(w, r)
}

func networkCheckHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Checking network connectivity: performing GET on http://www.google.com")
	resp, err := http.DefaultClient.Get("http://www.google.com")
	if (err != nil) {
		log.Printf("Error performing GET request on www.google.com: %q", err)
		fmt.Fprintf(w, "Error performing GET request on www.google.com: %q", err)
	}

	log.Printf("HTTP status code returned by www.google.com: %d", resp.StatusCode)
	fmt.Fprintf(w, "HTTP status code returned by www.google.com: %d", resp.StatusCode)
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
	log.Printf("ScaleDownTargetEmulator listening on port 8080")
	addSignalTrap()
	http.HandleFunc("/", handler)
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/exit", exitHandler)
	http.HandleFunc("/preStop", preStopHandler)
	http.HandleFunc("/checkNetwork", networkCheckHandler)
	http.ListenAndServe(":8080", nil)
}