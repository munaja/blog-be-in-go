package user

import (
	dg "github.com/karincake/apem/databasegorm"
	"gorm.io/gorm"

	mu "github.com/munaja/blog-be-in-go/internal/model/user"
)

func GetDetailByName(name string) (*mu.User, error) {
	var user *mu.User
	err := dg.I.Model(&mu.User{}).Where("Name = ?", name).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		return nil, err
	}
	return user, nil
}

func GetIdByName(name string) (map[string]interface{}, error) {
	user := map[string]interface{}{}
	err := dg.I.Model(&mu.User{}).Select("Id").Where("Name = ?", name).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		return nil, err
	}
	return user, nil
}
