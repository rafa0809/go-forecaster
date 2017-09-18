package manager

import (
	"github.com/server-forecaster/model/entity"
)

type UserManager struct {
	BaseManager
}

func (manager UserManager) AddUser(user *entity.User) bool {
	manager.DB.Create(&user)
	return !manager.DB.NewRecord(user)
}

func (manager UserManager) GetUserByAlias(alias string) *entity.User {
	user := entity.User{}
	manager.DB.Where("Alias = ?", alias).First(&user)
	return &user
}

func CreateUserManager() UserManager {
	return UserManager{BaseManager: Create()}
}
