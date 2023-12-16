package person

import (
	"time"

	"github.com/karincake/getuk"
	gdt "gorm.io/datatypes"
)

type Person struct {
	Id int `json:"id" gorm:"primaryKey;autoIncrement"`
	getuk.DateModel
	User_Id   int       `json:"user_id "`
	Name      string    `json:"name" gorm:"size:100"`
	Birthdate *gdt.Date `json:"birthDate"`
	Gender    *byte     `json:"gender"`
	Education *string   `json:"education" gorm:"size:50"`
	Country   *string   `json:"country" gorm:"size:50"`
	Province  *string   `json:"province" gorm:"size:50"`
	Regency   *string   `json:"regency" gorm:"size:50"`
	Address   string    `json:"address" gorm:"size:200"`
	Phone     *string   `json:"phone" gorm:"size:50"`
}

type UpdateDto struct {
	Name      string    `json:"name" validate:"required;alphaSpace;maxLength=100"`
	Birthdate *gdt.Date `json:"birthDate" validate:"required"`
	Gender    *byte     `json:"gender" validate:"required"`
	Education *string   `json:"education" validate:"required;maxLength=50"`
	Country   *string   `json:"country" validate:"required;maxLength=50"`
	Province  *string   `json:"province" validate:"required;maxLength=50"`
	Regency   *string   `json:"regency" validate:"required;maxLength=50"`
	Address   string    `json:"address" validate:"required;maxLength=200"`
	Phone     *string   `json:"phone" validate:"required;maxLength=50"`
}

type ProfilDto struct {
	// user
	User_Id        int       `json:"user_id"`
	User_Name      string    `json:"user_name"`
	User_CreatedAt time.Time `json:"user_createdAt"`
	// person
	Name      string `json:"name"`
	Education string `json:"education"`
	Country   string `json:"country"`
	Province  string `json:"province"`
	Regency   string `json:"regency"`
}

type ProfilOMDto struct {
	// user
	User_Id        int       `json:"user_id"`
	User_Name      string    `json:"user_name"`
	User_CreatedAt time.Time `json:"user_createdAt"`
	// person
	Name      string  `json:"name"`
	Education string  `json:"education"`
	Country   string  `json:"country"`
	Province  string  `json:"province"`
	Regency   string  `json:"regency"`
	Address   string  `json:"address"`
	Phone     *string `json:"phone"`
}
