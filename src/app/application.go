package app

import (
	"net/http"

	"github.com/gorilla/mux"
)

var (
	router = mux.NewRouter()
)

func StartApplication(){
	
	mapUrls()
	srv := &http.Server{
		Handler: router,
		Addr: "localhost:8000",
	}

	if err := srv.ListenAndServe(); err != nil {
		panic(err)
	}

}