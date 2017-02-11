package db

import "database/sql"

//go:generate counterfeiter . TeamDBFactory
type TeamDBFactory interface {
	GetTeamDBById(int) TeamDB
	GetTeamDBByName(string) (TeamDB, bool, error)
}

type teamDBFactory struct {
	conn        Conn
	bus         *notificationsBus
	lockFactory LockFactory
}

func NewTeamDBFactory(conn Conn, bus *notificationsBus, lockFactory LockFactory) TeamDBFactory {
	return &teamDBFactory{
		conn:        conn,
		bus:         bus,
		lockFactory: lockFactory,
	}
}

func (f *teamDBFactory) GetTeamDBById(teamId int) TeamDB {
	return &teamDB{
		teamId:       teamId,
		conn:         f.conn,
		buildFactory: newBuildFactory(f.conn, f.bus, f.lockFactory),
	}
}

func (f *teamDBFactory) GetTeamDBByName(teamName string) (TeamDB, bool, error) {
	var id int
	row := f.conn.QueryRow(`
	SELECT id FROM teams WHERE LOWER(name) = LOWER($1)
`, teamName)

	err := row.Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, false, nil
		}
		return nil, false, err
	}

	return &teamDB{
		teamId:       id,
		conn:         f.conn,
		buildFactory: newBuildFactory(f.conn, f.bus, f.lockFactory),
	}, true, nil
}
