package dto

import "github.com/Topvennie/beta-log/internal/database/model"

type User struct {
	ID   int    `json:"id"`
	UID  string `json:"uid"`
	Name string `json:"name"`
}

func UserDTO(user model.User) User {
	return User(user)
}

func (u *User) ToModel() model.User {
	return model.User(*u)
}
