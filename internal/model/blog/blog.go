package blog

import (
	"time"

	"github.com/karincake/getuk"
)

type Status byte

type Blog struct {
	Id int `json:"id" gorm:"primaryKey;autoIncrement"`
	getuk.DateModel
	User_Id      int     `json:"user_id"`
	Title        *string `json:"title" gorm:"size:100"`
	CategoryTags *string `json:"categoryTags" gorm:"size:2000"`
	CoverImg     *string `json:"coverImg" gorm:"size:5000"`
	Content      *string `json:"content" gorm:"type:text"`
	Status       *Status `json:"status,omitempty"`
}

type CreateDraftDto struct {
	User_Id      *int    `json:"-"`
	Title        string  `json:"title" validate:"maxLength=50"`
	CategoryTags *string `json:"categoryTags"`
	CoverImg     *byte   `json:"coverImg"`
	Content      *string `json:"content"`
	Status       Status  `json:"-"`
}

type CreatePublishDto struct {
	User_Id      *int    `json:"-"`
	Title        string  `json:"title" validate:"required;maxLength=50"`
	CategoryTags *string `json:"categoryTags" validate:"required"`
	CoverImg     *byte   `json:"coverImg"`
	Content      *string `json:"content" validate:"required"`
	Status       Status  `json:"-"`
}

type UpdateDraftDto struct {
	User_Id      *int    `json:"-"`
	Title        string  `json:"title" validate:"maxLength=50"`
	CategoryTags *string `json:"categoryTags"`
	CoverImg     *byte   `json:"coverImg"`
	Content      *string `json:"content"`
	Status       Status  `json:"-"`
}

type UpdatePublishDto struct {
	User_Id      *int    `json:"-"`
	Title        string  `json:"title" validate:"required;maxLength=50"`
	CategoryTags *string `json:"categoryTags" validate:"required"`
	CoverImg     *byte   `json:"coverImg"`
	Content      *string `json:"content" validate:"required"`
	Status       Status  `json:"-"`
}

type DeleteDto struct {
	User_Id *int `json:"-"`
}

type ListFilterDto struct {
	User_Id          *int   `json:"user_id"`
	Title            string `json:"title"`
	Title_Opt        string `json:"title_opt"`
	CategoryTags     string `json:"categoryTags"`
	CategoryTags_Opt string `json:"categoryTags_opt"`
	Status           Status `json:"-"`
	CreateAt         time.Time
	Page             int `json:"page"`
	PageSize         int `json:"pagesize"`
}

type DetailFilterDto struct {
	User_Id *int `json:"-"`
}

const (
	BSDraft    Status = 0
	BSActive   Status = 1
	BSInactive Status = 2
	BSBlocked  Status = 3
)
