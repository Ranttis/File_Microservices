package GET_service

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
)

func getFile(w http.ResponseWriter, r *http.Request) {

}


func SetupGetRoutes() {
	http.HandleFunc("/GET", getFile)
	http.ListenAndServe(":8080", nil)
}