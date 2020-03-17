package main

import (
	"File_Microservice/POST_service"
	"log"
	"net/http"
	"os"
	"time"
)






func main(){
	l := log.New(os.Stdout, "File_Microservice", log.LstdFlags)
	sm := http.NewServeMux()

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
		err := s.ListenAndServe()
		if err != nil{
			l.Printf("Virhe käynnistäessä palvelinta", err)
			os.Exit(1)
		}
	}()

	POST_service.SetupPostRoutes()
}