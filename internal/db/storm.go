package db

import (
	"context"
	"fmt"
	"github.com/aleri-godays/project"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	log "github.com/sirupsen/logrus"
)
import "github.com/asdine/storm/v3"

type stormDB struct {
	db *storm.DB
}

func NewStormDB(dbPath string) *storm.DB {
	p := fmt.Sprintf("%s/project.db", dbPath)
	db, err := storm.Open(p)
	if err != nil {
		log.WithFields(log.Fields{
			"db_path": p,
			"error":   err,
		}).Fatal("could not open storm db")
	}
	return db
}

func NewStormRepository(db *storm.DB) project.Repository {
	if err := db.Init(&StormProject{}); err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("could not initialize StormProject bucket")
	}

	sdb := &stormDB{
		db: db,
	}

	return sdb
}

func (s *stormDB) Add(ctx context.Context, p *project.Project) (*project.Project, error) {
	span := newSpanFromContext(ctx, "Add")
	defer span.Finish()

	sp := projectToStormProject(p)
	if err := s.db.Save(sp); err != nil {
		return nil, fmt.Errorf("could not save project: %w", err)
	}
	return stormProjectToProject(sp), nil
}

func (s *stormDB) Get(ctx context.Context, id int) (*project.Project, error) {
	span := newSpanFromContext(ctx, "Get")
	defer span.Finish()

	var sp StormProject
	if err := s.db.One("ID", id, &sp); err != nil {
		if err == storm.ErrNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("could not fetch project '%d': %w", id, err)
	}
	return stormProjectToProject(&sp), nil
}

func (s *stormDB) Delete(ctx context.Context, id int) error {
	span := newSpanFromContext(ctx, "Delete")
	defer span.Finish()

	sp := StormProject{ID: id}
	if err := s.db.DeleteStruct(&sp); err != nil {
		return fmt.Errorf("could not delete project '%d': %w", id, err)
	}
	return nil
}

func (s *stormDB) Update(ctx context.Context, p *project.Project) error {
	span := newSpanFromContext(ctx, "Update")
	defer span.Finish()

	sp := projectToStormProject(p)
	if err := s.db.Update(sp); err != nil {
		return fmt.Errorf("could not update project '%d': %w", p.ID, err)
	}
	return nil
}

func (s *stormDB) All(ctx context.Context) ([]*project.Project, error) {
	span := newSpanFromContext(ctx, "All")
	defer span.Finish()

	var sps []StormProject
	if err := s.db.All(&sps); err != nil {
		return nil, fmt.Errorf("could not fetch all projects: %w", err)
	}
	if len(sps) == 0 {
		return nil, nil
	}
	ps := make([]*project.Project, 0, len(sps))
	for _, sp := range sps {
		p := stormProjectToProject(&sp)
		ps = append(ps, p)
	}
	return ps, nil
}

func newSpanFromContext(ctx context.Context, opName string) opentracing.Span {
	span, _ := opentracing.StartSpanFromContext(ctx, "db-"+opName)
	ext.DBInstance.Set(span, "project")
	ext.DBType.Set(span, "storm")
	ext.Component.Set(span, "database")
	return span
}
