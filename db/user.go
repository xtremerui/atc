package db

import sq "github.com/Masterminds/squirrel"

type User interface {
	ID() string
	TeamNames() []string
	Workers() ([]Worker, error)
}

type user struct {
	conn Conn

	id    string
	teams []string
}

func (u *user) ID() string {
	return u.id
}

func (u *user) TeamNames() []string {
	return u.teams
}

func (u *user) Workers() ([]Worker, error) {
	return getWorkers(u.conn, workersQuery.Where(sq.Or{
		sq.Eq{"t.name": u.teams},
		sq.Eq{"w.team_id": nil},
	}))
}
