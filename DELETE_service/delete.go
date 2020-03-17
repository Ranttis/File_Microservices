package DELETE_service

import (
	"log"
	"net/http"
	"regexp"
	"strconv"
)



type Files struct {
	l *log.Logger
}


func NewFiles(l *log.Logger) *Files {
	return &Files{l}
}

func (f *Files) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodDelete{
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
}



func (f *Files) deleteFile(id int, rw http.ResponseWriter, r *http.Request) {
	lp := 

}

func SetupDelRoutes() {
	http.HandleFunc("/DELETE", DeleteFile)
	http.ListenAndServe(":8080", nil)
}
