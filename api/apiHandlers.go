package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/marni/goigc"
	"igcinfo/model"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

/*   Global vars   */
var Id = 1
var AllIds []model.TrackId
var AllTracks []model.Track
var Start time.Time

/*
** Basepoint of the API. Gives basic info about the API
 */
func Index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	reg := regexp.MustCompile("^/(api/)$")
	parts := reg.FindStringSubmatch(r.URL.Path)

	uptime := model.GetUptime(Start)

	if parts != nil {
		if r.Method == "GET" {
			json.NewEncoder(w).Encode(model.Information{
				Uptime: uptime, Info: "Service for IGC tracks.", Version: "version 1.0",
			})
		} else {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		}
	} else {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	}

}

/*
** Accepts POST or GET request
** Restores a track when the right igc- url is sent with POST
** Shows slices of IDs of tracks restored in the memory when GET are used
 */
func RegAndShowTrackIds(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	switch r.Method {
	case "POST":
		var url model.URL
		defer r.Body.Close()
		if err := json.NewDecoder(r.Body).Decode(&url); err != nil {
			http.Error(w, "The api only accepts JSON data", 404)
		}

		track, err := igc.ParseLocation(url.Url)
		if err != nil {
			http.Error(w, "Empty or wrong igc url provided", 404)
		} else {
			var trackLen float64
			for i := 0; i < len(track.Points)-1; i++ {
				trackLen += track.Points[i].Distance(track.Points[i+1])
			}
			AllTracks = append(AllTracks, model.Track{
				ID: Id, HDate: track.Date,
				Pilot: track.Pilot, Glider: track.GliderType,
				GliderId: track.GliderID, TrackLength: trackLen,
			})
			newTrack := model.TrackId{Id: Id}
			AllIds = append(AllIds, newTrack)
			Id++
			json.NewEncoder(w).Encode(newTrack)
		}

	case "GET":
		s := make([]int, len(AllIds))
		for index, ids := range AllIds {
			s[index] = ids.Id
		}
		json.NewEncoder(w).Encode(s)

	default:
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

/*
** Retrieves a track by its id, Accepts only GET
 */
func ShowTrackInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	var field string
	path := strings.Split(r.URL.Path, "/")
	id, conErr := strconv.Atoi(path[3])
	if conErr != nil {
		http.Error(w, "Wrong or empty id provided!", http.StatusNotFound)
		return
	}
	if len(path) < 3 || len(path) > 5 {
		http.Error(w, "Not implimented yet", http.StatusNotImplemented)
		return
	} else if len(path) == 5 {
		field = path[4]
	}

	if r.Method == "GET" {
		track, err := GetTrackById(id)
		if err == nil && field == "" {
			json.NewEncoder(w).Encode(track)
		} else if err == nil && field != "" {
			ShowTrackField(w, r, track, field)
		} else {
			http.Error(w, "Did't find the track with id ("+path[3]+")", 404)
		}
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

/*
** Retrieves the track field, Accepts only GET
 */
func ShowTrackField(w http.ResponseWriter, r *http.Request, obj model.Track, field string) {
	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")

	switch field {
	case "pilot":
		fmt.Fprint(w, obj.Pilot)
	case "glider":
		fmt.Fprint(w, obj.Glider)
	case "glider_id":
		fmt.Fprint(w, obj.GliderId)
	case "calculated total track length":
		fmt.Fprint(w, obj.TrackLength)
	case "H_date":
		fmt.Fprint(w, obj.HDate)
	default:
		http.Error(w, "Wrong field provided", http.StatusNotFound)
	}
}

// Returns a struct of Track that matches the param id.
func GetTrackById(id int) (T model.Track, err error) {
	for _, T = range AllTracks {
		if id == T.ID {
			return T, nil
		}
	}
	return T, errors.New("no track found")
}
