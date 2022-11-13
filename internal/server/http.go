package server

import (
	"fmt"
	"net/http"
	"spotfinder/cmd"

	"github.com/gorilla/mux"
)

type httpServer struct { 
	req []byte
}

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
	fmt.Println(lat)
	long := r.URL.Query().Get("long")

	spot.FindSpots(lat, long, &w)

}

