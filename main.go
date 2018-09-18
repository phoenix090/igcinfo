package main

import (
    "fmt"
    "log"
    "net/http"
    "encoding/json"
    //"regexp"
    "github.com/marni/goigc"
    "time"
    //"strings"
    //"io/ioutil"
)
 var start = time.Now()

//var reg = regexp.MustCompile("^/igcinfo/api/([a-z- A-Z/]+)$")
var AllIds []TrackId         //Track ids
var AllTracks []Track

type TrackId struct {
    Id string 
}

type URL struct {
    Url string `json:"url"`
}

type Information struct {
    Uptime float64
    Info string
    Version string
}

type Track struct {
    ID string
    H_date string //map, slice?
    Pilot string // Struct?
    Glider string
    Glider_id string
    Track_length string
}


/*
** Basepoint of the API- URI. Gives basic info about the API
*/
func Index(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    if http.StatusOK == 200{
        if r.Method == "GET" {
            s := time.Now().Sub(start).Seconds()
            //fmt.Println(time.Now().UTC().Format(time.RFC3339))
            //elapsed := t.Sub(start)
            //fmt.Printf("type: %T", elapsed)       //elapsed er float64
            //fmt.Println("\nTime elapsed: ", elapsed)
            json.NewEncoder(w).Encode(Information { 
                Uptime : s, Info : "Service for IGC tracks.", Version: "version 1.0", 
            })
        } else {
            http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
        }
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
    if http.StatusOK == 200 {
        if r.Method == "POST" {
        
            var url URL
            if err := json.NewDecoder(r.Body).Decode(&url); err != nil {
                fmt.Fprintf(w,"The api only accepts raw json!") 
                panic(err)
            }

            track, err := igc.ParseLocation(url.Url)
            if err != nil {
                fmt.Errorf("Problem reading the track", err)
                fmt.Fprintf(w, "Wrong url!")
            } else { 
                AllTracks = append(AllTracks, Track { 
                    ID : track.UniqueID, H_date : track.Date.String(),
                    Pilot : track.Pilot, Glider : track.GliderType,
                    Glider_id : track.GliderID, Track_length : track.FlightRecorder,
                })
                AllIds = append(AllIds, TrackId { Id : track.UniqueID })
                json.NewEncoder(w).Encode(url)
            }

        } else if r.Method == "GET" {
            //fmt.Println("Status: ", http.statusCode)
            s := make([]string, len(AllIds))
            for index, ids := range AllIds {
                s[index] = ids.Id
            }
            json.NewEncoder(w).Encode(s)
        } else {
            http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
        }  

    } else {
        fmt.Fprintf(w, "Something went wrong! Status: " + string(http.StatusOK))
    }
      
}

/*
** Retrieves a track by its id, Accepts only GET
*/
func ShowTrackInfo(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    if http.StatusOK == 200{
        if r.Method == "GET" {
            id := r.FormValue("id")
            for _, track := range AllTracks {
                if track.ID == id {
                    json.NewEncoder(w).Encode(track)
                    return
                }
            }
            Err := make(map[string]string)
            Err["Error"] = "The track with the id: " + id + " does not excist!"
            json.NewEncoder(w).Encode(Err)
        } else {
            http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
        }
    }
    
}




func main() {
    http.HandleFunc("/api/", Index)
    http.HandleFunc("/api/igc", RegAndShowTrackIds)
    http.HandleFunc("/api/igc/", ShowTrackInfo)

    log.Fatal(http.ListenAndServe(":32123", nil))
}