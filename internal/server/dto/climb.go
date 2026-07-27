package dto

import (
	"time"

	"github.com/Topvennie/beta-log/internal/database/model"
	"github.com/Topvennie/beta-log/pkg/utils"
)

type climb struct {
	ID         int               `json:"id"`
	Grade      int               `json:"grade"`
	Color      string            `json:"color"`
	HoldColor  string            `json:"hold_color"`
	ClimbType  model.ClimbType   `json:"climb_type"`
	FinishType model.ClimbFinish `json:"finish_type"`
	Source     model.Source      `json:"source"`
}

func climbDTO(c *model.Climb) climb {
	return climb{
		ID:         c.ID,
		Grade:      c.Grade,
		HoldColor:  c.HoldColor,
		ClimbType:  c.ClimbType,
		FinishType: c.FinishType,
		Source:     c.Source,
	}
}

type ClimbDay struct {
	ID     int          `json:"id"`
	Date   time.Time    `json:"date"`
	Gym    Gym          `json:"gym"`
	Climbs []climb      `json:"climbs"`
	Source model.Source `json:"source"`
}

func ClimbDayDTO(d *model.ClimbDay) ClimbDay {
	return ClimbDay{
		ID:     d.ID,
		Date:   d.Date,
		Gym:    GymDTO(&d.Gym),
		Climbs: utils.SliceMap(d.Climbs, func(c model.Climb) climb { return climbDTO(&c) }),
		Source: d.Source,
	}
}

type ClimbDayFilter struct {
	UserID int
	Limit  int
	Offset int
}

func (c ClimbDayFilter) ToModel() model.ClimbDayFilter {
	return model.ClimbDayFilter(c)
}

type ClimbStatsGraphProgress struct {
	Date   string `json:"date"`
	Grade  int    `json:"grade"`
	Volume int    `json:"volume"`
}

type ClimbStatsGraphGrade struct {
	Grade  int `json:"grade"`
	Flash  int `json:"flash"`
	Top    int `json:"top"`
	Repeat int `json:"repeat"`
}

type ClimbStats struct {
	Total                  int                       `json:"total"`
	TotalUnique            int                       `json:"total_unique"`
	Flash                  int                       `json:"flash"`
	Top                    int                       `json:"top"`
	Repeat                 int                       `json:"repeat"`
	Best                   int                       `json:"best"`
	BestAmount             int                       `json:"best_amount"`
	BestFlash              int                       `json:"best_flash"`
	BestFlashAmount        int                       `json:"best_flash_amount"`
	Sessions               int                       `json:"sessions"`
	MedianClimbsPerSession int                       `json:"median_climbs_per_session"`
	GraphProgress          []ClimbStatsGraphProgress `json:"graph_progress"`
	GraphPerGrade          []ClimbStatsGraphGrade    `json:"graph_per_grade"`
}
