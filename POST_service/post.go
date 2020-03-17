package POST_service

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

func uploadFile(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Latauksen päätepiste saavutettu")

	file, handler, err := r.FormFile("testiTiedosto")
	if err != nil {
		fmt.Printf("Virhe tiedoston latauksessa: %v ", err)
		return
	}
	filename := handler.Filename
	defer file.Close()
	n := 0
	tempFile, err := ioutil.TempFile("data", "tiedosto-*-" + strconv.Itoa(n) +".txt")
	if err != nil {
		fmt.Println(err)
	}
	defer tempFile.Close()

	// read all of the contents of our uploaded file into a
	// byte array
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}
	// write this byte array to our temporary file
	tempFile.Write(fileBytes)
	n++
	// return that we have successfully uploaded our file!
	fmt.Fprintf(w, "Tiedosto ladattu! \n")
	fmt.Fprintf(w, tempFile.Name())
}


func SetupPostRoutes() {
	http.HandleFunc("/POST", uploadFile)
	http.ListenAndServe(":8080", nil)
}