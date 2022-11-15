package server

import (
	"net/http"
	"spotfinder/cmd"
	"regexp"
	"github.com/gorilla/mux"
)

type httpServer struct {}

func NewHTTPServer(addr string) *http.Server {
	httpsrv := newHTTPServer()
	r := mux.NewRouter()

	r.HandleFunc("/", httpsrv.handleConsume).Methods("GET") 
	return &http.Server{
		Addr:    addr,
		Handler: r,
	}

}

func newHTTPServer() *httpServer { 
	return &httpServer{}
}

//anytime something is fetched **GET** 
func (s *httpServer) handleConsume(w http.ResponseWriter, r *http.Request) {

	lat := r.URL.Query().Get("lat")
	lat_match, err := regexp.MatchString("[0-9.-]{1,}", lat)
	if err != nil {
		return
	}

	long := r.URL.Query().Get("long")
	long_match, err := regexp.MatchString("[0-9.-]{1,}", long)
	if err != nil {
		return
	}

	if lat_match && long_match {
		spot.FindSpots(lat, long, &w)
	} else {
		w.Write([]byte("Bad Latitude and Longitude, please try again..."))
	}
}

