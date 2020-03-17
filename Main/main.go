package main

import (
	"File_Microservices/handlers"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)


func main(){
	l := log.New(os.Stdout, "File_Microservice", log.LstdFlags)

	ph := handlers.NewFiles(l)
	sm := http.NewServeMux()
	sm.Handle("/", ph)

	s := http.Server{
		Addr:              ":8080",
		ErrorLog:		   l,
		Handler:           sm,
		ReadTimeout:       5*time.Second,
		ReadHeaderTimeout: 0,
		WriteTimeout:      5*time.Second,
		IdleTimeout:       60*time.Second,
	}

	go func(){
		l.Println("Käynnistetään palvelin")
		err := http.ListenAndServe(":8080", sm)
		if err != nil{
			l.Printf("Virhe käynnistäessä palvelinta", err)
			os.Exit(1)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	// Block until a signal is received.
	sig := <-c
	log.Println("Saatu signaali:", sig)

	// gracefully shutdown the server, waiting max 30 seconds for current operations to complete
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(ctx)

}