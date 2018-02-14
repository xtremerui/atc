package db

import (
	"context"
	"errors"
)

type UserFactory interface {
	GetUser(context.Context) (User, error)
}

type userFactory struct {
	conn Conn
}

func NewUserFactory(conn Conn) UserFactory {
	return &userFactory{
		conn: conn,
	}
}

func (f *userFactory) GetUser(c context.Context) (User, error) {

	return nil, errors.New("no user")
}
