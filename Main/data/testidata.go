package data

import (
	"encoding/json" // Jsonien käsittelyyn tarvittava paketti
	"fmt"           // Format, sisältää esim. printit
	"io"            // Input & Output paketista saadaan tarvittavat Readerit ja Writerit
	"os"            // OS paketista saadaan työkalut paikallisten tiedostojen luomiseen ja poistamiseen
	"strconv"       // String conversion
)

type File struct {
	ID      int    `json:"id"`      // Luodun tiedoston ID
	Name    string `json:"name"`    // Luodun tiedoston nimi
	Type    string `json:"type"`    // Luodun tiedoston tyyppi
	Content string `json:"content"` // Luodun tiedoston sisältö
}

// Globaalit muuttujat
var filePath = "data/" // Polku, johon tiedosto tallennetaan. Root on kansio, joka sisältää main.go:n

// Funktio, joka decodaa saadun Json datan.
// f on fromJSONin vastaanottaja (receiver), eli FromJSONin sisällä päästään käsiksi f:ään.
// Argumenttina annetaan io.Reader.
// io.Reader on interface, joka wrappaa Read metodin
// Käytetään httpserice.go:n puolella, kun lisätään uusia tiedostoja
func (f *File) FromJSON(r io.Reader) error {
	fileDecoder := json.NewDecoder(r) // fileDecoder on decoder, joka lukee Json datan Readeriltä r
	return fileDecoder.Decode(f)      // Funktio palauttaa decoodatun tiedoston f sisällön
}

// Files on slice muodossa oleva kokoelma luotuja File-tyypisiä 'objekteja'
type Files []*File

// ToJson funktio sarjastaa (serialize) Files kokoelman Json muotoon.
// f on toJSONin vastaanottaja (receiver), eli toJSONin sisällä päästään käsiksi f:ään.
// Argumenttina annetaan io.Writer.
// io. Writer on interface, joka wrappaa Write metodin
// Käytetään httpserice.go:n puolella, kun haetaan olemassa olevia tiedostoja
func (f *Files) ToJSON(w io.Writer) error {
	fileEncoder := json.NewEncoder(w) // fileEncoder on encoder, joka kirjoittaa encodaatun Jsonin Writteriin w
	return fileEncoder.Encode(f)      // Funktio palauttaa encoodatun datan
}

// GetFiles hakee kaikki tiedostot 'filesList' listasta ja luo ne paikallisesti 'data' kansioon
// Käytetään httpservice.go:ssa
func GetFiles() Files {
	if len(fileList) != 0 { // Tarkastetaan ettei lista ole tyhjä
		for _, f := range fileList { // Loopataan fileListin sisällön läpi
			contentString := "ID: " + strconv.Itoa(f.ID) + "\n" + "name: " + f.Name + "\n" +
				"type: " + f.Type + "\n" + "content: " + f.Content // Luodaan string muuttuja, joka sisältää tiedoston datan, ID, Nimi, Tyyppi ja Sisältö
			content := []byte(contentString)                        // Muokataan luotu string muuttuja byteiksi
			testFile, cErr := os.Create(filePath + f.Name + f.Type) // Luodaan os moduulin avulla uusi tiedosto, haluttuun polkuun
			defer testFile.Close()                                  // Defer avainsanalla varmistetaan, että luotu tiedosto suljetaan, poistuttaessa GetFile funktiosta
			if cErr != nil {                                        // Idiomaattisen Go:n perusteiden mukaan, tarkastetaan, ettei tiedoston luomisessa ole virheitä.
				println(cErr.Error())
			}
			_, _ = testFile.Write(content) // Kirjoitetaan data luotuun tiedostoon
		}
	} else {
		fileList = nil // Jos lista on tyhjä asetetaan se nulliksi virhetilojen välttämiseksi
	}

	println("Amount of files: ", len(fileList))
	return fileList // Palautetaan lista löydetyistä tiedostoista
}

// AddFile luo uuden tiedoston annetulla datalla
// Argumenttina on luotava tiedosto f
func AddFile(f *File) {
	f.ID = getNextID()   // Haetaan seuraava ID, filesListin mukaan
	if fileList != nil { // Jos fileList ei ole tyhjä, lisätään f appendilla listaan
		fileList = append(fileList, f)
	} else { // Jos fileList on tyhjä (null), luodaan tyhjä lista ja lisätään f listaan appendilla
		println(fileList)
		fileList = []*File{}
		fileList = append(fileList, f)
	}
	contentString := "ID: " + strconv.Itoa(f.ID) + "\n" + "name: " + f.Name + // Luodaan string muuttuja, joka sisältää tiedoston datan, ID, Nimi, Tyyppi ja Sisältö
		"\n" + "type: " + f.Type + "\n" + "content: " + f.Content
	println(f.Name, f.Type, f.ID)
	content := []byte(contentString)                       // Muokataan luotu string muuttuja byteiksi
	testFile, err := os.Create(filePath + f.Name + f.Type) // Luodaan os moduulin avulla uusi tiedosto, haluttuun polkuun
	defer testFile.Close()                                 // Defer avainsanalla varmistetaan, että luotu tiedosto suljetaan, poistuttaessa AddFile funktiosta
	if err != nil {                                        // Idiomaattisen Go:n perusteiden mukaan, tarkastetaan, ettei tiedoston luomisessa ole virheitä.
		println(err.Error())
	}
	_, _ = testFile.Write(content) // Kirjoitetaan data luotuun tiedostoon
}

// DeleteFile poistaa tiedoston paikallisista tiedostoista ja fileLististä
// Käytetään httpservice.go:ssa
// Argumenttina annetaan poistettavan tiedoston ID ja tiedosto f
func DeleteFile(id int, f *File) error {
	ff, pos, err := findFile(id) // Etsitään findFile funktiolla haluttu tiedosto ID:n avulla	ff = löydetty tiedosto, pos = positio fileListissä
	if err != nil {              // Tarkistetaan etsintä virheiden varalta
		println(err.Error())
		return err
	}
	var pathToDelete = filePath + ff.Name + ff.Type // Poistettavan tiedoston polku

	if _, err := os.Stat(pathToDelete); err == nil { // os.Stat palauttaa FileInfo kuvauksen tiedostosta
		println("Poistettu: ", pathToDelete)
		var er = os.Remove(pathToDelete) // Poistetaan paikallinen tiedosto os.Remove funktiolla
		if er != nil {                   // Tarkastetaan poistaminen virheiden varalta
			println(er.Error())
			return er
		}
	} else {
		println("Tiedostoa ei löytynyt polusta: ", pathToDelete) // Jos tiedostoa ei löydy, tulostetaan virhe
	}

	fileList[pos] = f                    // Määritetään poistettavan tiedoston paikka listassa
	fileList = RemoveIndex(fileList, id) // Poistetaan elementti listasta

	return nil
}

// RemoveIndex poistaa slicen elementin annetussa indexissä			TODO: INDEXI KUNTOON
func RemoveIndex(f []*File, index int) []*File {
	if index >= len(f)-1 {
		f = f[:index+copy(f[index-1:], f[index:])]
	} else {
		f = f[:index+copy(f[index:], f[index+1:])]
		f = append(f[:index], f[index+1:]...)
	}
	return f
}

// getNextID hakee seuraavan ID:n fileLististä
func getNextID() int {
	fp := fileList[len(fileList)-1] // Haetaan viimeinen ID fileLististä
	return fp.ID + 1                // Lisätään viimeiseen ID:hen 1 ja palautetaan uusi ID
}

// findFile käy fileListin läpi ja etsii tiedostoa annetulla ID:llä
func findFile(id int) (*File, int, error) { // Parametrina id, jota etsitään
	for i, p := range fileList { //													// Loopataan fileList läpi
		if p.ID == id { // Jos ID löytyy, palautetaan
			return p, i, nil // löydetty tiedosto, positio ja null virhe
		}
	}
	return nil, -1, ErrFileNotFound // Jos tiedostoa ei löydy, palautetaan virhe.
}

var ErrFileNotFound = fmt.Errorf("File not found") // Custom virheen implementaatio

// Testidatana toimii lista kovakoodatuista Fileistä.
var fileList = []*File{
	&File{
		ID:   1,
		Name: "testitiedosto",
		Type: ".txt",
		Content: "Lorem ipsum dolor sit amet.",
	},

	&File{
		ID:   2,
		Name: "testitiedosto2",
		Type: ".txt",
		Content: "Lorem ipsum dolor sit amet.",
	},

	&File{
		ID:   3,
		Name: "jsontiedosto1",
		Type: ".json",
		Content:"Lorem ipsum dolor sit amet.",
	},
}
