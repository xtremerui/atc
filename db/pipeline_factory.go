package db

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/concourse/atc/db/lock"
)

//go:generate counterfeiter . PipelineFactory

type PipelineFactory interface {
	TeamPipelines(names ...string) ([]Pipeline, error)
	PublicPipelines() ([]Pipeline, error)
	AllPipelines() ([]Pipeline, error)
}

type pipelineFactory struct {
	conn        Conn
	lockFactory lock.LockFactory
}

func NewPipelineFactory(conn Conn, lockFactory lock.LockFactory) PipelineFactory {
	return &pipelineFactory{
		conn:        conn,
		lockFactory: lockFactory,
	}
}

func (f *pipelineFactory) TeamPipelines(names ...string) ([]Pipeline, error) {
	rows, err := pipelinesQuery.
		Where(sq.Eq{"name": names}).
		OrderBy("team_id ASC", "ordering ASC").
		RunWith(f.conn).
		Query()
	if err != nil {
		return nil, err
	}

	teamPipelines, err := scanPipelines(f.conn, f.lockFactory, rows)
	if err != nil {
		return nil, err
	}

	rows, err = pipelinesQuery.
		Where(sq.NotEq{"name": names}).
		Where(sq.Eq{"public": true}).
		OrderBy("team_id ASC", "ordering ASC").
		RunWith(f.conn).
		Query()
	if err != nil {
		return nil, err
	}

	otherTeamPublicPipelines, err := scanPipelines(f.conn, f.lockFactory, rows)
	if err != nil {
		return nil, err
	}

	return append(teamPipelines, otherTeamPublicPipelines...), nil
}

func (f *pipelineFactory) PublicPipelines() ([]Pipeline, error) {
	rows, err := pipelinesQuery.
		Where(sq.Eq{"p.public": true}).
		OrderBy("t.name, ordering").
		RunWith(f.conn).
		Query()
	if err != nil {
		return nil, err
	}

	pipelines, err := scanPipelines(f.conn, f.lockFactory, rows)
	if err != nil {
		return nil, err
	}

	return pipelines, nil
}

func (f *pipelineFactory) AllPipelines() ([]Pipeline, error) {
	rows, err := pipelinesQuery.
		OrderBy("ordering").
		RunWith(f.conn).
		Query()
	if err != nil {
		return nil, err
	}

	return scanPipelines(f.conn, f.lockFactory, rows)
}
