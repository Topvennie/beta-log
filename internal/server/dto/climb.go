package dto

import (
	"time"

	"github.com/Topvennie/beta-log/internal/database/model"
	"github.com/Topvennie/beta-log/internal/grade"
	"github.com/Topvennie/beta-log/pkg/utils"
)

type climb struct {
	ID         int               `json:"id"`
	Grade      string            `json:"grade"`
	HoldColor  string            `json:"hold_color"`
	ClimbType  model.ClimbType   `json:"climb_type"`
	FinishType model.ClimbFinish `json:"finish_type"`
	Source     model.Source      `json:"source"`
}

func climbDTO(c *model.Climb, setting model.Setting) climb {
	return climb{
		ID:         c.ID,
		Grade:      grade.Grade(c.Grade).Format(setting.GradeSystem),
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

func ClimbDayDTO(d *model.ClimbDay, setting model.Setting) ClimbDay {
	return ClimbDay{
		ID:     d.ID,
		Date:   d.Date,
		Gym:    GymDTO(&d.Gym),
		Climbs: utils.SliceMap(d.Climbs, func(c model.Climb) climb { return climbDTO(&c, setting) }),
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

type ClimbCreate struct {
	Grade      string            `json:"grade"`
	HoldColor  string            `json:"hold_color" validate:"hexcolor"`
	ClimbType  model.ClimbType   `json:"climb_type" validate:"required"`
	FinishType model.ClimbFinish `json:"finish_type" validate:"required"`
}

func (c ClimbCreate) ToModel(setting model.Setting) model.Climb {
	return model.Climb{
		Grade:      int(grade.FromString(c.Grade, setting.GradeSystem)),
		HoldColor:  c.HoldColor,
		ClimbType:  c.ClimbType,
		FinishType: c.FinishType,
	}
}

type ClimbDayCreate struct {
	Date   time.Time     `json:"date" validate:"required"`
	GymID  int           `json:"gym_id" validate:"required,min=1"`
	Climbs []ClimbCreate `json:"climbs" validate:"required,min=1,dive"`
}

func (c ClimbDayCreate) ToModel(setting model.Setting) model.ClimbDay {
	return model.ClimbDay{
		Date:   c.Date,
		GymID:  c.GymID,
		Climbs: utils.SliceMap(c.Climbs, func(c ClimbCreate) model.Climb { return c.ToModel(setting) }),
	}
}

type ClimbDayUpdate struct {
	ID     int           `json:"id" validate:"required"`
	Date   time.Time     `json:"date" validate:"required"`
	GymID  int           `json:"gym_id" validate:"required,min=1"`
	Climbs []ClimbCreate `json:"climbs" validate:"required,min=1,dive"`
}

func (c ClimbDayUpdate) ToModel(setting model.Setting) model.ClimbDay {
	return model.ClimbDay{
		ID:     c.ID,
		Date:   c.Date,
		GymID:  c.GymID,
		Climbs: utils.SliceMap(c.Climbs, func(c ClimbCreate) model.Climb { return c.ToModel(setting) }),
	}
}

type ClimbStatsProgress struct {
	Date   string `json:"date"`
	Grade  string `json:"grade"`
	Volume int    `json:"volume"`
}

type ClimbStatsGrade struct {
	Grade  string `json:"grade"`
	Flash  int    `json:"flash"`
	Top    int    `json:"top"`
	Repeat int    `json:"repeat"`
}

type ClimbStats struct {
	Total                  int                  `json:"total"`
	TotalUnique            int                  `json:"total_unique"`
	Flash                  int                  `json:"flash"`
	Top                    int                  `json:"top"`
	Repeat                 int                  `json:"repeat"`
	Best                   string               `json:"best"`
	BestAmount             int                  `json:"best_amount"`
	BestFlash              string               `json:"best_flash"`
	BestFlashAmount        int                  `json:"best_flash_amount"`
	Sessions               int                  `json:"sessions"`
	Boulder                int                  `json:"boulder"`
	Lead                   int                  `json:"lead"`
	MedianClimbsPerSession int                  `json:"median_climbs_per_session"`
	GraphProgress          []ClimbStatsProgress `json:"graph_progress"`
	GraphPerGrade          []ClimbStatsGrade    `json:"graph_per_grade"`
}
