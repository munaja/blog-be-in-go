package customer

import (
	"github.com/munaja/blog-practice-be-using-go/internal/model/user"
	"github.com/munaja/blog-practice-be-using-go/internal/model/usertoken"
)

func GetModelList() (data []interface{}) {
	tableList := []interface{}{
		&user.User{},
		&usertoken.UserToken{},
	}
	data = append(data, tableList...)

	return data
}
