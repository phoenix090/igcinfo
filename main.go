package main

import (
	"igcinfo/api"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {

	api.Start = time.Now()
	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "8080"
	}

	http.HandleFunc("/api/", api.Index)
	http.HandleFunc("/api/igc", api.RegAndShowTrackIds)
	http.HandleFunc("/api/igc/", api.ShowTrackInfo)
	err := http.ListenAndServe(":"+port, nil)
	log.Fatalf("Server error: %s", err)
}
