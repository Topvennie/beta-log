// Package model contains all databank models
package model

import "github.com/Topvennie/beta-log/pkg/sqlc"

type User struct {
	ID   int
	UID  string
	Name string
}

func UserModel(user sqlc.User) *User {
	return &User{
		ID:   int(user.ID),
		UID:  user.Uid,
		Name: user.Name,
	}
}
