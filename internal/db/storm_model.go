package db

import "github.com/aleri-godays/project"

type StormProject struct {
	ID        int    `storm:"id,increment"`
	Name      string `storm:"index"`
	User      string `storm:"index"`
	DailyRate float64
}

func projectToStormProject(p *project.Project) *StormProject {
	return &StormProject{
		ID:        p.ID,
		Name:      p.Name,
		User:      p.User,
		DailyRate: p.DailyRate,
	}
}

func stormProjectToProject(p *StormProject) *project.Project {
	return &project.Project{
		ID:        p.ID,
		Name:      p.Name,
		User:      p.User,
		DailyRate: p.DailyRate,
	}
}
