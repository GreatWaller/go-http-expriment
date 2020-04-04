package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"test/handlers"
	"time"
)

func main() {
	l := log.New(os.Stdout, "product-api", log.LstdFlags)
	hh := handlers.NewHello(l)

	sm := http.NewServeMux()
	sm.Handle("/", hh)

	s := &http.Server{
		Addr:         ":9090",
		Handler:      sm,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	// start the server
	go func() {
		l.Println("Starting server on port 9090")

		err := s.ListenAndServe()
		if err != nil {
			l.Printf("Error starting server: %s\n", err)
			os.Exit(1)
		}
	}()

	// trap sigterm or interupt and gracefully shutdown the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	// Block until a signal is received.
	sig := <-c
	log.Println("Got signal:", sig)

	// gracefully shutdown the server, waiting max 30 seconds for current operations to complete
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(ctx)

	// http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
	// 	log.Println("Hello World")
	// 	d, err := ioutil.ReadAll(r.Body)
	// 	if err != nil {
	// 		http.Error(rw, "Oops", http.StatusBadRequest)
	// 		return
	// 	}
	// 	fmt.Fprintf(rw, "hello %s", d)
	// })

	// http.HandleFunc("/goodbye", func(http.ResponseWriter, *http.Request) {
	// 	log.Println("Goodbye World")
	// })

	// http.ListenAndServe(":9090", nil)
}
