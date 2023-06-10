package scoreboard

import (
	"encoding/json"
	"time"
)

type Scoreboard struct {
	Players   []string  `json:"players"`
	StartedAt time.Time `json:"startedAt"`
	Id        string    `json:"id"`
	Winner    string    `json:"winner"`
}

func (s Scoreboard) MarshalJSON() ([]byte, error) {
	if len(s.Id) == 0 {
		return json.Marshal(struct {
			Players   []string  `json:"players"`
			StartedAt time.Time `json:"startedAt"`
			Winner    string    `json:"winner"`
		}{
			Players:   s.Players,
			StartedAt: s.StartedAt,
			Winner:    s.Winner,
		})
	}
	return json.Marshal(struct {
		Players   []string  `json:"players"`
		StartedAt time.Time `json:"startedAt"`
		Id        string    `json:"id"`
		Winner    string    `json:"winner"`
	}{
		Players:   s.Players,
		StartedAt: s.StartedAt,
		Id:        s.Id,
		Winner:    s.Winner,
	})
}
