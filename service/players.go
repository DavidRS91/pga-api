package service

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"

	"github.com/DavidRS91/pga-api/data"
	"github.com/gorilla/mux"
)

func (s *Server) GetPlayer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid player ID")
		return
	}
	p := data.Player{ID: id}
	if err := p.GetPlayer(s.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Player not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())

		}
		return
	}
	respondWithJSON(w, http.StatusOK, p)
}

func (s *Server) GetPlayers(w http.ResponseWriter, r *http.Request) {
	count, _ := strconv.Atoi(r.FormValue("count"))
	start, _ := strconv.Atoi(r.FormValue("start"))
	if count > 10 || count < 1 {
		count = 10
	}
	if start < 0 {
		start = 0
	}

	players, err := data.GetPlayers(s.DB, start, count)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, players)
}

func (s *Server) CreatePlayer(w http.ResponseWriter, r *http.Request) {
	var p data.Player
	d := json.NewDecoder(r.Body)
	if err := d.Decode(&p); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	if err := p.CreatePlayer(s.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusCreated, p)
}

func (s *Server) UpdatePlayer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid player ID")
		return
	}

	var p data.Player
	d := json.NewDecoder(r.Body)
	if err := d.Decode(&p); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()
	p.ID = id

	if err := p.UpdatePlayer(s.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, p)
}

func (s *Server) DeletePlayer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid player ID")
		return
	}

	p := data.Player{ID: id}
	if err := p.DeletePlayer(s.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func (s *Server) SyncPlayers(w http.ResponseWriter, r *http.Request) {
	resp, err := http.Get("https://datacrunch.9c9media.ca/statsapi/sports/golf/leagues/golf/pga/scoreboard?brand=tsn")
	if err != nil {
		fmt.Printf("err: %w", err)
		respondWithError(w, http.StatusBadRequest, "Invalid player ID")
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	var result []map[string]interface{}
	json.Unmarshal(body, &result)

	eventList := result[0]["eventList"]
	events := InterfaceSlice(eventList)

	// TODO:
	// 1. generate event map
	// 2. grab playerEventStatsList key
	// 3. parse playerEventStatsList for player data
	v := reflect.ValueOf(events[0])
	if v.Kind() == reflect.Map {
		for _, key := range v.MapKeys() {
			strct := v.MapIndex(key)
			fmt.Println(key.Interface(), strct.Interface())
		}
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
	return
}

// borrowed from https://stackoverflow.com/questions/12753805/type-converting-slices-of-interfaces
func InterfaceSlice(slice interface{}) []interface{} {
	s := reflect.ValueOf(slice)
	if s.Kind() != reflect.Slice {
		panic("InterfaceSlice() given a non-slice type")
	}

	ret := make([]interface{}, s.Len())

	for i := 0; i < s.Len(); i++ {
		ret[i] = s.Index(i).Interface()
	}

	return ret
}
