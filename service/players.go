package service

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
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
	// resp, err := http.Get("https://datacrunch.9c9media.ca/statsapi/sports/golf/leagues/golf/pga/scoreboard?brand=tsn")
	resp, err := http.Get("https://datacrunch.9c9media.ca/statsapi/sports/golf/leagues/golf/pga/tournament/?brand=tsn&type=json")
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Error fetching data")
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	var result map[string]interface{}
	json.Unmarshal(body, &result)

	eventList := result["events"]
	events := InterfaceSlice(eventList)
	eventMap := InterfaceMap(events[0])
	players := InterfaceSlice(eventMap["playerEventStatsList"])
	for _, p := range players {
		pm := InterfaceMap(p)
		info := InterfaceMap(pm["player"])
		stats := InterfaceMap(pm["stats"])
		fmt.Println("ID: ", info["playerId"], info["displayName"], ": ", stats["scoreTotal"])
		name := info["displayName"].(string)
		score, err := InterfaceInt(stats["scoreTotal"])
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		tsnID, err := InterfaceInt(info["playerId"])
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		player := data.Player{
			Name:  name,
			Score: score,
			IsCut: false,
			TsnID: tsnID,
		}
		if err := player.CreatePlayer(s.DB); err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}
	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
	return
}
