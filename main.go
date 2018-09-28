package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/marni/goigc"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)
var id = 1
var AllIds []TrackId         //Track ids
var AllTracks []Track
var start time.Time

type TrackId struct {
    Id int
}

type URL struct {
    Url string `json:"url"`
}

type Information struct {
    Uptime string
    Info string
    Version string
}

type Track struct {
    ID int
    HDate time.Time
    Pilot string
    Glider string
    GliderId string
    TrackLength float64
}


func getUptime() (uptime string){
	now := time.Now()
	newTime := now.Sub(start).String()
	hours := int(now.Hour())
	y, m, d := "0", "0", "0"
	if hours > 23 && hours % 24 != 0 {
		// Checking if the days are beneeth 10
		d = strconv.Itoa(hours * 24 % 31)
		if hours * 240 < 10 {
			d = "0" + strconv.Itoa(hours * 24 % 31)
		}
	} else {

	}
	days, _ := strconv.Atoi(d)

	if days > 31 {
		m = strconv.Itoa(days * 31 % 31)
	}
	months, _ := strconv.Atoi(m)
	if months > 12 {
		y = strconv.Itoa(months * 12)
	}
	uptime = "P" + y + "Y" + m + "M" + d + "DT" + newTime

	return uptime
}

/*
** Basepoint of the API- URI. Gives basic info about the API
*/
func Index(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	reg := regexp.MustCompile("^/igcinfoapi.herokuapp.com(/api/)*$")
	parts := reg.FindStringSubmatch(r.URL.Path)
	uptime := getUptime()

    // bytt ut true med statusCode sjekk!

    if parts != nil {
        if r.Method == "GET" {
            json.NewEncoder(w).Encode(Information { 
                Uptime : uptime, Info : "Service for IGC tracks.", Version: "version 1.0",
            })
        } else {
            http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
        }
    } else {
		http.NotFound(w, r)
	}
    
}

/*
** Accepts POST or GET request
** Restores a track when the right igc- url is sent with POST
** Shows slices of IDs of tracks restored in the memory when GET are used
*/
func RegAndShowTrackIds(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    //w.Header().Set("Access-Control-Allow-Origin", "*")

    // bytt ut true med statusCode sjekk!
	if true {
        if r.Method == "POST" {
        
            var url URL
            if err := json.NewDecoder(r.Body).Decode(&url); err != nil {
				http.Error(w, "The api only accepts raw JSON", 404)
            }

            track, err := igc.ParseLocation(url.Url)
            if err != nil {
            	http.Error(w, "Empty or wrong igc url provided", 404)
            } else {
				var trackLen float64
				for i := 0; i < len(track.Points)-1; i++ {
					trackLen += track.Points[i].Distance(track.Points[i+1])
				}
                AllTracks = append(AllTracks, Track { 
                    ID : id, HDate : track.Date,
                    Pilot : track.Pilot, Glider : track.GliderType,
                    GliderId : track.GliderID, TrackLength : trackLen,
                })
				newTrack := TrackId { Id : id }
                AllIds = append(AllIds, newTrack)
                id++
                json.NewEncoder(w).Encode(newTrack)
            }

        } else if r.Method == "GET" {
            //fmt.Println("Status: ", http.statusCode)
            s := make([]int, len(AllIds))
            for index, ids := range AllIds {
                s[index] = ids.Id
            }
            json.NewEncoder(w).Encode(s)
        } else {
            http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
        }  

    } else {
		http.Error(w, "Error!", http.StatusBadRequest)
    }
      
}

/*
** Retrieves a track by its id, Accepts only GET
*/
func ShowTrackInfo(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	var field string
	path := strings.Split(r.URL.Path, "/")

	id, conErr := strconv.Atoi(path[4])
	fmt.Println(conErr)
	if conErr != nil && (len(path) >= 4 && len(path) <= 5){
		http.Error(w, "Error with the given ID, must be integer", 404)
	}
	if len(path) > 5 {
		field = path[5]
	}
	// bytt ut true med statusCode sjekk!
	if true {
        if r.Method == "GET" {
        	track, err := getTrackById(id)
        	if err == nil && field == "" {
				json.NewEncoder(w).Encode(track)
			} else if err == nil && field != "" {
				ShowTrackField(w, r, track, field)
			} else {
				http.Error(w, "Did't find the track with id (" + path[4] + ")", 404)
			}
        } else {
            http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
        }
    }
    
}

/*
** Retrieves the track field, Accepts only GET
*/
func ShowTrackField(w http.ResponseWriter, r *http.Request, obj Track, field string){
	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")

	switch field {
	case "pilot":
		json.NewEncoder(w).Encode(obj.Pilot)
	case "glider":
		json.NewEncoder(w).Encode(obj.Glider)
	case "glider_id":
		json.NewEncoder(w).Encode(obj.GliderId)
	case "calculated total track length":
		json.NewEncoder(w).Encode(obj.TrackLength)
	case "H_date":
		json.NewEncoder(w).Encode(obj.HDate)
	default:
		http.Error(w, "Could't find the field in the record", http.StatusNotFound)
	}
}

// Returns a struct of Track that matches the param id.
func getTrackById(id int) (T Track, err error) {
	for _, T = range AllTracks {
		if id == T.ID {
			return T, nil
		}
	}
	return T , errors.New("No track found")
}



func main() {
	start = time.Date(2018, time.September, 28, 17, 0, 0, 0, time.UTC)
	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "8080"
	}

	http.HandleFunc("/igcinfoapi.herokuapp.com/api/", Index)
    http.HandleFunc("/igcinfoapi.herokuapp.com/api/igc", RegAndShowTrackIds)
	http.HandleFunc("/igcinfoapi.herokuapp.com/api/igc/", ShowTrackInfo)
    err := http.ListenAndServe( ":" + port, nil)
	log.Fatalf("Server error: %s", err)
}