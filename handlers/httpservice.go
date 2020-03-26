package handlers

import (
	"log"							// Yksinkertainen lokipaketti
	"net/http"						// Paketti HTTP client-server implementaatioon
	"regexp"						// Regular expression haku
	"strconv"						// String conversion
	"File_Microservices/Main/data"	// Testidatan implementaatio
)

// Files toimii http.Handler tyyppinä, jota tarvitaan palvelimen käynnistyksessä
// http.Handler vastaa HTTP kutsuihin ja se wrappaa ServerHTTP metodin.
type Files struct {
	lgr *log.Logger
}

// NewFiles luo handlerin tiedostoille halutulla loggerilla
func NewFiles(l *log.Logger) *Files {
	return &Files{l}
}

// ServeHTTP täyttää http.Handler interfacen ehdot ja toimii aloituspisteenä handlerille
// ServeHTTP käsittelee saapuvat HTTP kutsut (POST,DELETE,GET)
// f on ServeHTTP vastaanottaja (receiver), eli ServeHTTP:n sisällä päästään käsiksi f:ään.
// Argumenttina ottaa http.ResponseWriterin, joka rakentaa HTTP vastauksen (response)
// ja http.Request pointerin, joka on serverin vastaanottama tai clientin lähettämä pyyntö.
// ServeHTTP ei palauta mitään
func (f *Files) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet { 													// Tarkastetaan onko saapuva pyyntö GET
		f.getFiles(rw, r)																// Jos pyyntö on GET, kutsutaan getFiles funktiota
		return																			// ServeHTTP ei palauta mitään
	}

	if r.Method == http.MethodPost {													// Tarkastetaan onko saapuva pyyntö POST
		f.addFile(rw, r)																// Jos pyyntö on POST, kutsutaan addFile funktiota
		return																			// ServeHTTP ei palauta mitään
	}

	if r.Method == http.MethodDelete {													// Tarkastetaan onko saapuva pyyntö DELETE
		f.lgr.Println("DELETE", r.URL.Path)
		reg := regexp.MustCompile(`/([0-9]+)`)										// Luodaan Regural expression, jonka avulla  tarkistetaan, että URLissa on vain yksi ID
		g := reg.FindAllStringSubmatch(r.URL.Path, -1)								// Etsitään URLista regexin avulla ID

		if len(g) != 1 {																// Jos löytyy useampi, kuin yksi ID, kirjataan virhe
			f.lgr.Println("URI ei kelpaa, liian monta ID:tä")						// Tulostetaan käyttäjälle virheteksti
			http.Error(rw, "URI ei kelpaa", http.StatusBadRequest)				// Lähetetään BadRequest (status 400) virhe
				return																	// ServeHTTP ei palauta mitään
			}

		if len(g[0]) != 2 {																// Tarkastetaan, ettei regexin capture grouppeja ole liikaa
			f.lgr.Println("URI ei kelpaa, liian monta capture grouppia")			// Tulostetaan käyttäjälle virheteksti
			http.Error(rw, "URI ei kelpaa", http.StatusBadRequest)				// Lähetetään BadRequest (status 400) virhe
			return																		// ServeHTTP ei palauta mitään
		}

		idString := g[0][1]																// Luodaan IDstä string muuttuja
		id, err := strconv.Atoi(idString)												// Muutetaan string intiksi
		if err != nil {																	// Jos muutos ei onnistu
			f.lgr.Println("URI ei kelpaa, ei voitu muuttaa numeroksi", idString)	// Näytetään virheteksti käyttäjälle
			http.Error(rw, "URI ei kelpaa", http.StatusBadRequest)				// Lähetetään BadRequest (status 400) virhe
			return																		// ServeHTTP ei palauta mitään
		}
		f.deleteFile(id, rw, r)															// Jos ID löytyy oikein, kutsutaan deleteFile funktiota.
		return																			// ServeHTTP ei palauta mitään
	}

	// Jos metodi on muu kuin POST, GET, DELETE kirjoitetaan virhe
	rw.WriteHeader(http.StatusMethodNotAllowed)
}



// getFiles funktio hakee kaikki tiedostot testidata.go tiedostosta
// f on getFilesin vastaanottaja (receiver), eli getFilesin sisällä päästään käsiksi f:ään.
// Argumentteina samat ResponseWriter ja Request, kuin ServeHTTPllä
func (f *Files) getFiles(rw http.ResponseWriter, r *http.Request) {
	f.lgr.Println("Handle GET")

	lp := data.GetFiles()																	// Noudetaan tiedostot testidata.gosta

	err := lp.ToJSON(rw)																	// Serialisoidaan saatu lista JSONiksi
	if err != nil {																			// Tarkastetaan virheiden varalta
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)		// Jos JSONin muodostamisessa on ongelma, näytetään virhe käyttäjälle
	}
}

// addFile funktio luo tiedoston saadusta datasta
// f on addFilen vastaanottaja (receiver), eli addFilen:n sisällä päästään käsiksi f:ään.
// Argumentteina samat ResponseWriter ja Request, kuin getFilesillä
func (f *Files) addFile(rw http.ResponseWriter, r *http.Request) {
	f.lgr.Println("Handle POST Product")

	file := &data.File{}																	// Luodaan file muuttuja
	err := file.FromJSON(r.Body)															// Muodostetaan sen data JSONista
	if err != nil {																			// Tarkastetaan virheiden varalta
		http.Error(rw, "Unable to unmarshal json", http.StatusBadRequest)				// Jos JSONin muodostamisessa on ongelma, näytetään virhe käyttäjälle
	}

	data.AddFile(file)																		// Kutsutaan testidata.gon AddFile funktiota, ja annetaan sille luotu file parametriksi
}

// deletefile funktio poistaa annetulla IDllä löytyvän tiedoston lokaalista polusta ja data storesta
// f on deleteFilen vastaanottaja (receiver), eli deleteFilen sisällä päästään käsiksi f:ään.
// Argumentteina samat ResponseWriter ja Request, kuin getFilesillä ja lisäksi postettavan tiedoston ID
func (f Files) deleteFile(id int, rw http.ResponseWriter, r*http.Request) {
	f.lgr.Println("Handle DELETE Product")

	file := &data.File{}																	// Luodaan file muuttuja

	err := data.DeleteFile(id, file)														// Kutsutaan testidata.gon DeleteFile funktiota ja annetaan parametreiksi poistettavan tiedoston ID ja file
	if err == data.ErrFileNotFound{															// Tarkastetaan virheiden varalta, jos tiedostoa ei löydy
		http.Error(rw, "File not found", http.StatusNotFound)							// Näytetään virhe; status 404, NotFound
		return																				// deleteFile ei palauta mitään
	}

	if err != nil{																			// Jos virhe on jokin muu
		http.Error(rw, "File not found", http.StatusInternalServerError)				// Näytetään virhe; status 500, InternalServerError
		return																				// deleteFile ei palauta mitään
	}
}
