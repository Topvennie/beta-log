package dto

import "github.com/Topvennie/beta-log/internal/database/model"

type Gym struct {
	ID       int          `json:"id"`
	Name     string       `json:"name"`
	IconPath string       `json:"icon_path"`
	Source   model.Source `json:"source"`
}

func GymDTO(g *model.Gym) Gym {
	return Gym{
		ID:       g.ID,
		Name:     g.Name,
		IconPath: g.IconPath,
		Source:   g.Source,
	}
}

type GymCreate struct {
	Name     string `json:"name" validate:"required"`
	IconPath string `json:"icon_path"`
}

func (g GymCreate) ToModel() model.Gym {
	return model.Gym{
		Name:     g.Name,
		IconPath: g.IconPath,
	}
}

type GymUpdate struct {
	ID       int    `json:"id" validate:"required,min=1"`
	Name     string `json:"name" validate:"required"`
	IconPath string `json:"icon_path"`
}

func (g GymUpdate) ToModel() model.Gym {
	return model.Gym{
		ID:       g.ID,
		Name:     g.Name,
		IconPath: g.IconPath,
	}
}

type GymStatsVisits struct {
	Gym    string `json:"gym"`
	Amount int    `json:"amount"`
}

type GymStatsTop struct {
	Gym   string `json:"gym"`
	Top   int    `json:"top"`
	Flash int    `json:"flash"`
}

type GymStatsDistribution struct {
	Gym          string      `json:"gym"`
	Distribution map[int]int `json:"distribution"`
}

type GymStats struct {
	Total             int                    `json:"total"`
	MostVisited       string                 `json:"most_visited"`
	MostVisitedAmount int                    `json:"most_visited_amount"`
	Sessions          int                    `json:"sessions"`
	GraphVisits       []GymStatsVisits       `json:"graph_visits"`
	GraphTop          []GymStatsTop          `json:"graph_top"`
	GraphDistribution []GymStatsDistribution `json:"graph_distribution"`
}
