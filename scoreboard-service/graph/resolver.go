package graph

import (
	"fmt"

	"soa.mafia-game/scoreboard-service/graph/model"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	ScoreboardsStore map[string]model.Scoreboard
	Comments         map[string]model.Comment
}

func (r *Resolver) AddCommentRelation(comment model.Comment) error {
	scoreboard, exists := r.ScoreboardsStore[comment.ScoreboardID]
	if !exists {
		return fmt.Errorf("scoreboard %v not found", comment.ScoreboardID)
	}
	scoreboard.Related = append(scoreboard.Related, &comment)
	r.ScoreboardsStore[comment.ScoreboardID] = scoreboard
	return nil
}
