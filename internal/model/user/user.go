package user

import (
	"fmt"
	"time"

	"github.com/karincake/getuk"
)

type Status byte
type Position int16

type User struct {
	Id int `json:"id" gorm:"primaryKey;autoIncrement"`
	getuk.DateModel
	Name                    string  `json:"name" gorm:"size:100;unique"`
	Email                   string  `json:"email" gorm:"size:100;unique"`
	Password                *string `json:"password,omitempty" gorm:"size:200"`
	OptStatus               *bool   `json:"optSatus,omitempty"`
	OwnerStatus             *bool   `json:"ownerSatus,omitempty"`
	Status                  *Status `json:"status,omitempty"`
	FailedLoginAttemptCount int     `json:"-"`
	LastSuccessLogin        time.Time
	LastAllowedLogin        time.Time `json:"-"`
}

type CreateDto struct {
	Name        string    `json:"name" validate:"required"`
	Email       string    `json:"email" validate:"email"`
	Password    string    `json:"password" validate:"required;minLength=8"`
	Status      int16     `json:"status"`
	ValidPeriod time.Time `json:"validPeriod"`
}

type UpdateDto struct {
	Email string `json:"email"`
}

type RegisterDto struct {
	Name       string `json:"name" validate:"required"`
	Email      string `json:"email" validate:"required;email"`
	Password   string `json:"password" validate:"required;minLength=8"`
	RePassword string `json:"repassword" validate:"required;eqField=Password"`
}

type ResendConfirmationEmailDto struct {
	Email string `json:"email" validate:"required;email"`
	Token string `json:"token" validate:"required"`
}

type ResendEmailConfirmDto struct {
	Email string `json:"email" validate:"required;email"`
}

type LoginDto struct {
	Name        string   `json:"name" validate:"required"`
	Password    string   `json:"password" validate:"required"`
	Duration    string   `json:"duration"` // Format: "[val]-[unit]", where val is int, and accepted units are: m (minute), h (hour), d (day), month. Any error format will results in 24 hours of duration
	OptStatus   *bool    `json:"-"`
	OwnerStatus *bool    `json:"-"`
	Position    Position `json:"-"`
}

type FilterDto struct {
	Name   *string `json:"name"`
	Type   *int16  `json:"type"`
	Email  *string `json:"email"`
	Status *int16  `json:"status"`
	// fixed fields
	Page     int   `json:"page"`
	PageSize int64 `json:"page_size"`
}

type ChangePassDto struct {
	OldPassword string `json:"oldPassword,omitempty" validate:"required;minLength=8"`
	NewPassword string `json:"newPassword,omitempty" validate:"required;minLength=8"`
	RePassword  string `json:"rePassword,omitempty" validate:"required;minLength=8;eqField=NewPassword"`
}

type RequestResetPassDto struct {
	Email string `json:"email" validate:"required;email"`
}

type CheckResetPassDto struct {
	Email string `json:"email" validate:"required;email"`
	Token string `json:"token" validate:"required"`
}

type ResetPassDto struct {
	NewPassword string `json:"newPassword,omitempty" validate:"required;minLength=8"`
	RePassword  string `json:"rePassword,omitempty" validate:"required;minLength=8"`
}

const (
	UPCustomer Position = 1
	UPOperator Position = 2
	UPAdmin    Position = 3
	UPOwner    Position = 4

	USNew       Status = 0
	USActive    Status = 1
	USBlocked   Status = 2
	USSuspended Status = 3
)

var usList map[Status]string = map[Status]string{
	0: "New",
	1: "Active",
	2: "Blocked",
	3: "Suspended",
}

func GetUSText(code Status) string {
	status, _ := usList[code]
	fmt.Println(code)
	fmt.Println(status)
	return status
}
