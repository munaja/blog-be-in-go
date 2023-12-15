package blogcategory

import "github.com/karincake/getuk"

type Status byte

type BlogCategory struct {
	Id int `json:"id" gorm:"primaryKey;autoIncrement"`
	getuk.DateModel
	Code int     `json:"code" gorm:"size:20;"`
	Name *string `json:"name" gorm:"size:50"`
}

type CreateDto struct {
	Code *string `json:"code" validate:"required"`
	Name *string `json:"name" validate:"required"`
}

// actually can use the same DTO of creating data
type UpdateDto struct {
	Code *string `json:"code" validate:"required"`
	Name *string `json:"name" validate:"required"`
}

const (
	BSDraft    Status = 0
	BSActive   Status = 1
	BSInactive Status = 2
	BSBlocked  Status = 3
)
