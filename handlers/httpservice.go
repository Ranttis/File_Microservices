package handlers


import (
	"log"
	"net/http"
	"regexp"
	"strconv"

	"File_Microservices/Main/data"
)

// Products is a http.Handler
type Files struct {
	l *log.Logger
}

// NewProducts creates a products handler with the given logger
func NewFiles(l *log.Logger) *Files {
	return &Files{l}
}

// ServeHTTP is the main entry point for the handler and staisfies the http.Handler
// interface
func (f *Files) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	// handle the request for a list of products
	if r.Method == http.MethodGet {
		f.getFiles(rw, r)
		return
	}

	if r.Method == http.MethodPost {
		f.addProduct(rw, r)
		return
	}

	if r.Method == http.MethodPut {
		f.l.Println("PUT", r.URL.Path)
		// expect the id in the URI
		reg := regexp.MustCompile(`/([0-9]+)`)
		g := reg.FindAllStringSubmatch(r.URL.Path, -1)

		if len(g) != 1 {
			f.l.Println("Invalid URI more than one id")
			http.Error(rw, "Invalid URI", http.StatusBadRequest)
			return
		}

		if len(g[0]) != 2 {
			f.l.Println("Invalid URI more than one capture group")
			http.Error(rw, "Invalid URI", http.StatusBadRequest)
			return
		}

		idString := g[0][1]
		id, err := strconv.Atoi(idString)
		if err != nil {
			f.l.Println("Invalid URI unable to convert to numer", idString)
			http.Error(rw, "Invalid URI", http.StatusBadRequest)
			return
		}

		f.deleteFile(id, rw, r)
		return
	}

	// catch all
	// if no method is satisfied return an error
	rw.WriteHeader(http.StatusMethodNotAllowed)
}

// getProducts returns the products from the data store
func (p *Files) getFiles(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle GET Products")

	// fetch the products from the datastore
	lp := data.GetFiles()

	// serialize the list to JSON
	err := lp.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}
}

func (f *Files) addProduct(rw http.ResponseWriter, r *http.Request) {
	f.l.Println("Handle POST Product")

	file := &data.File{}

	err := file.FromJSON(r.Body)
	if err != nil {
		http.Error(rw, "Unable to unmarshal json", http.StatusBadRequest)
	}

	data.AddFile(file)
}

func (f Files) deleteFile(id int, rw http.ResponseWriter, r*http.Request) {
	f.l.Println("Handle PUT Product")

	file := &data.File{}

	err := file.FromJSON(r.Body)
	if err != nil {
		http.Error(rw, "Unable to unmarshal json", http.StatusBadRequest)
	}

	err = data.DeleteFile(id, file)
	if err == data.ErrProductNotFound {
		http.Error(rw, "Product not found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(rw, "Product not found", http.StatusInternalServerError)
		return
	}
}