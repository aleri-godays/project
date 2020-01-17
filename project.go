package project

import "context"

type Project struct {
	ID        int     `json:"id"`
	Name      string  `json:"name"`
	User      string  `json:"user"`
	DailyRate float64 `json:"daily_rate"`
}

type Repository interface {
	Add(ctx context.Context, p *Project) (*Project, error)
	Get(ctx context.Context, id int) (*Project, error)
	Delete(ctx context.Context, id int) error
	Update(ctx context.Context, p *Project) error
	All(ctx context.Context) ([]*Project, error)
}
