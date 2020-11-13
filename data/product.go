package data

import (
	"database/sql"
)

type Player struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Score int    `json:"score"`
	IsCut bool   `json:is_cut`
}

func (p *Player) GetPlayer(db *sql.DB) error {
	return db.QueryRow("SELECT name, score, is_cut FROM players WHERE id=$1", p.ID).
		Scan(&p.Name, &p.Score, &p.IsCut)
}

func (p *Player) UpdatePlayer(db *sql.DB) error {
	_, err := db.Exec("UPDATE players SET name=$1, score=$2 is_cut=$3 WHERE id=$4", p.Name, p.Score, p.IsCut, p.ID)
	return err
}

func (p *Player) DeletePlayer(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM players WHERE id=$1", p.ID)
	return err
}

func (p *Player) CreatePlayer(db *sql.DB) error {
	return db.QueryRow(
		"INSERT INTO players(name, price) VALUES($1, $2, $3) RETURNING id", p.Name, p.Score, p.IsCut).
		Scan(&p.ID)

}

func GetPlayers(db *sql.DB, start, count int) ([]Player, error) {
	rows, err := db.Query("SELECT id, name, score, is_cut FROM players LIMIT $1 OFFSET $2", count, start)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	players := []Player{}

	for rows.Next() {
		var p Player
		if err := rows.Scan(&p.ID, &p.Name, &p.Score, &p.IsCut); err != nil {
			return nil, err
		}
		players = append(players, p)
	}
	return players, nil
}
