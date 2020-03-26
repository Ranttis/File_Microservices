package main

import (
	"File_Microservices/handlers"		// Handlerin implementaatio
	"context"							// Context tyyppi, joka sisältää deadlinet, perumissignaalit jne.
	"log"								// Yksinkertainen lokipaketti
	"net/http"							// Paketti HTTP client-server implementaatioon
	"os"								// OS paketista saadaan työkalut paikallisten tiedostojen luomiseen ja poistamiseen
	"os/signal"							// Saapuvien signaalien käsittely
	"time"								// Paketti josta saadaan tarvittavat työkalut ajan käsittelyyn
)

// main funktio on palvelun aloituspiste
func main(){
	// Oma loggeri, joka käyttää standardiflagejä
	l := log.New(os.Stdout, "File_Microservice", log.LstdFlags)

	fh := handlers.NewFiles(l)													// Luodaan uusi FileHandler fh
	sm := http.NewServeMux()													// Luodaan uusi ServeMux sm
	sm.Handle("/", fh)													// Määritetään handlerin patterni "/", esim. pyyntö localhost:8080/X

	s := http.Server{															// Luodaan uusi palvelin
		Addr:              ":8080",												// osoitteena toimii kaikki sisäverkon osoitteen porttina 8080
		ErrorLog:		   l,													// Uusi loggeri
		Handler:           sm,													// Oma Mux
		ReadTimeout:       5*time.Second,										// Kuinka kauan lukeminen saa kestää ennen timeouttia
		ReadHeaderTimeout: 0,													// Kuinka kauan pyyntöjen headereitten lukeminen saa kestää ennen timeouttia, 0 = ei timeouttia
		WriteTimeout:      5*time.Second,										// Kuinka kauan kirjoittaminen saa kestää ennen timeouttia
		IdleTimeout:       60*time.Second,										// Kuinka kauan palvelin saa olla idlenä ennen timeouttia
	}

	// Anonyymi funktio, jolla käynnistetään palvelin
	go func(){
		l.Println("Käynnistetään palvelin porttiin :8080")
		err := http.ListenAndServe(":8080", sm)							// http.ListenAndServe käynnistää palvelimen annettuun osoitteeseen ja käyttää annettua handleria pyyntöjen käsittelyyn
		if err != nil{															// Tarkastetaan palvelimen käynnistys virheen varalta
			l.Printf("Virhe käynnistäessä palvelinta", err)							// Jos virhe löytyy, ilmoitetaan käyttäjälle ja
			os.Exit(1)													// Sammutetaan sovellus
		}
	}()

	// Pienimuotoinen esimerkki kanavien käytöstä palvelimien yhteydess
	c := make(chan os.Signal, 1)												// Luodaan puskuroitu kanava (jotta ei riskeerata signaalien menetystä), joka kuuntelee os.Signaalia
	signal.Notify(c, os.Interrupt)												// Notify välittää Interrupt signaalin c:lle = käyttäjä painaa Ctrl + C

	// Blokataan, kunnes saadaan signaali
	sig := <-c
	log.Println("Saatu signaali:", sig)

	// Sammutetaan palvelin sulavasti 30sec jälkeen, jotta päällä olevat operaatiot saadaan suoritettua
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(ctx)

}