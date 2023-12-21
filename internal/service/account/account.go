package account

import (
	"errors"

	sc "github.com/jinzhu/copier"
	a "github.com/karincake/apem"
	dg "github.com/karincake/apem/databasegorm"
	l "github.com/karincake/apem/lang"
	te "github.com/karincake/tempe/error"
	"gorm.io/gorm"

	m "github.com/munaja/blog-practice-be-using-go/internal/model/user"
	mut "github.com/munaja/blog-practice-be-using-go/internal/model/usertoken"
	sut "github.com/munaja/blog-practice-be-using-go/internal/service/usertoken"
	p "github.com/munaja/blog-practice-be-using-go/pkg/password"
	sh "github.com/munaja/blog-practice-be-using-go/pkg/servicehelper"
)

const source = "user"

func Register(input m.RegisterDto) (any, error) {
	var user m.User
	if err := sc.Copy(&user, &input); err != nil {
		// return nil, td.IS{"message": l.I.Msg("copy-fail") + " user", "code": "fetch-fail"}
		return nil, te.XErrors{"struct": te.XError{Code: "copy-fail", Message: l.I.Msg("data-copy-fail")}}
	}

	password, err := p.Hash(*user.Password)
	if err != nil {
		// return nil, td.IS{"message": l.I.Msg("process-fail") + " password", "code": "fetch-fail"}
		return nil, te.XErrors{"data": te.XError{Code: "process-fail", Message: l.I.Msg("data-process-fail")}}
	} else {
		user.Password = &password
	}

	err = dg.I.Where(&m.RegisterDto{Email: input.Email}).Or(&m.RegisterDto{Name: input.Name}).First(&user).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		// return nil, sh.SetError("request", "create-data", source, "failed", l.I.Msg("process-fail")+": "+err.Error(), nil)
		// return nil, map[string]te.Error{"data": {Message: l.I.Msg("data-update-fail") + ": " + err.Error(), Code: "create-data"}}
		return nil, te.XErrors{"data": te.XError{Code: "update-fail", Message: l.I.Msg("data-update-fail")}}
	}

	xerr := te.XErrors{}
	if user.Id > 0 {
		if user.Name == input.Name {
			xerr["name"] = te.XError{Code: "registered", Message: l.I.Msg("registered")}
		}
		if user.Email == input.Email {
			xerr["email"] = te.XError{Code: "registered", Message: l.I.Msg("registered")}
		}
		if len(xerr) > 0 {
			return nil, xerr
		}
	}

	if err = dg.I.Transaction(func(tx *gorm.DB) error {
		StatusNew := m.USNew
		user.Status = &StatusNew
		if err := tx.Create(&user).Error; err != nil {
			return err
		}

		var userToken mut.UserToken
		if err := sut.Request(&userToken, input.Email, mut.UTConfirmByEmail, 60*24*3, 5); err != nil {
			return err
		} else if err := tx.Save(&userToken).Error; err != nil {
			return err
		}

		// useful for development
		if a.Apem.Env == "development" {
			user.UserToken = append(user.UserToken, userToken)
		}

		return nil
	}); err != nil {
		ed := sh.Event{
			Feature: "register",
			Action:  "create-data",
			Source:  source,
			Status:  "failed",
			ECode:   "data-create-fail",
			EDetail: err.Error(),
		}
		return nil, sh.SetError(ed, nil)
	}

	user.Password = nil
	user.Status = nil
	return user, nil
}

func ConfirmByEmail(input m.ResendConfirmationEmailDto) (any, error) {
	// check user token
	var userToken mut.UserToken
	if localECode, myEDetail := sut.CheckByEmail(&userToken, input.Email, input.Token, mut.UTConfirmByEmail); localECode != "" {
		return nil, te.XErrors{"email-confirmation": te.XError{Code: localECode, Message: l.I.Msg(localECode) + ", " + myEDetail}}
	}

	// check user
	var user m.User
	if err := checkUserStatus(&user, input.Email, m.USNew); err != nil {
		return nil, te.XErrors{"email-confirmation": te.XError{Code: "emailConfirm-fail", Message: l.I.Msg("emailConfirm-fail") + ", " + err.Error()}}
	}

	// go
	if err := dg.I.Transaction(func(tx *gorm.DB) error {
		UserStatus := m.USActive
		user.Status = &UserStatus
		result := tx.Save(&user)
		if result.RowsAffected == 0 {
			return errors.New(l.I.Msg("save-fail") + ": user")
		}

		result = tx.Unscoped().Delete(&userToken)
		if result.RowsAffected == 0 {
			return errors.New(l.I.Msg("save-fail") + ": token")
		}
		return nil
	}); err != nil {
		ed := sh.Event{
			Feature: "confirm-by-email",
			Action:  "process-data",
			Source:  source,
			Status:  "failed",
			ECode:   "save-fail",
			EDetail: err.Error(),
		}
		return nil, sh.SetError(ed, nil)
	}

	user.Password = nil
	user.Status = nil
	return user, nil
}

func ResendConfirmationEmail(input m.ResendEmailConfirmDto) (any, error) {
	// check user
	var user m.User
	if err := checkUserStatus(&user, input.Email, m.USNew); err != nil {
		return nil, te.XErrors{"user": te.XError{Code: "emailConfirm-resend-fail", Message: l.I.Msg("emailConfirm-resend-fail") + ", " + err.Error()}}
	}

	// check user token
	const tokenType = mut.UTConfirmByEmail
	var userToken mut.UserToken
	if err := sut.Request(&userToken, input.Email, tokenType, 60*24*3, 5); err != nil {
		return nil, te.XErrors{"user": te.XError{Code: "emailConfirm-resend-fail", Message: l.I.Msg("emailConfirm-resend-fail") + ", " + err.Error()}}
	}

	// next go
	ed := sh.Event{
		Feature: "confirm-by-email",
		Action:  "create-data",
		Source:  source,
		Status:  "failed",
		ECode:   "generate-fail",
	}
	if err := sut.Request(&userToken, input.Email, tokenType, 3*24*60, 5); err != nil {
		ed.EDetail = "token"
		return nil, sh.SetError(ed, userToken)
	} else if err := dg.I.Save(&userToken).Error; err != nil {
		ed.EDetail = "confirmation email"
		return nil, sh.SetError(ed, userToken)
	}
	return "request is done successfully", nil
}

func ChangePassword(id int, input m.ChangePassDto) (any, error) {
	var user m.User
	result := dg.I.First(&user, id)
	if result.RowsAffected == 0 {
		return nil, nil
	} else if result.Error != nil {
		return nil, sh.SetError(sh.Event{
			Feature: "change-password",
			Action:  "fetch-data",
			Source:  source,
			Status:  "failed",
			ECode:   "fetch-data-fail",
		}, user)
	} else if !p.Check(input.OldPassword, *user.Password) {
		return nil, te.XErrors{"oldPassword": te.XError{Code: "invalid", Message: l.I.Msg("invalid")}}
	}

	password, err := p.Hash(input.NewPassword)
	if err != nil {
		return nil, te.XErrors{"password": te.XError{Code: "data-generate-fail", Message: l.I.Msg("data-generate-fail")}}
	}

	user.Password = &password
	if result := dg.I.Save(&user); result.Error != nil {
		return nil, sh.SetError(sh.Event{
			Feature: "change-password",
			Action:  "save-data",
			Source:  source,
			Status:  "failed",
			ECode:   "save-data-fail",
		}, user)
	}

	return l.I.Msg("request-ok"), nil
}

func RequestResetPass(input m.RequestResetPassDto) (any, error) {
	// check user
	var user m.User
	if err := checkUserStatus(&user, input.Email, m.USActive); err != nil {
		return nil, te.XErrors{"email": te.XError{Code: "data-process-fail", Message: l.I.Msg("data-process-fail") + ", " + err.Error()}}
	}

	// chek user token
	var userToken mut.UserToken
	const tokenType = mut.UTResetPass
	if err := sut.Request(&userToken, input.Email, tokenType, 30, 5); err != nil {
		return nil, te.XErrors{"email": te.XError{Code: "data-process-fail", Message: l.I.Msg("data-process-fail") + ", " + err.Error()}}
	} else if err := dg.I.Save(&userToken).Error; err != nil {
		ed := sh.Event{
			Feature: "reset-password",
			Action:  "request",
			Source:  source,
			Status:  "failed",
			ECode:   "save-fail",
			EDetail: err.Error(),
		}
		return nil, sh.SetError(ed, userToken)
	}

	return l.I.Msg("request-ok"), nil
}

func CheckResetPass(input m.CheckResetPassDto) (any, error) {
	var userToken mut.UserToken
	if errCode, _ := sut.CheckByEmail(&userToken, input.Email, input.Token, mut.UTResetPass); errCode != "" {
		return nil, te.XErrors{"email": te.XError{Code: errCode, Message: l.I.Msg(errCode)}}
	}

	return userToken, nil
}

func ResetPass(input1 m.CheckResetPassDto, input2 m.ResetPassDto) (any, error) {
	var userToken mut.UserToken
	if errCode, _ := sut.CheckByEmail(&userToken, input1.Email, input1.Token, mut.UTResetPass); errCode != "" {
		return nil, te.XErrors{"email": te.XError{Code: errCode, Message: l.I.Msg(errCode)}}
	}

	// TODO: PINDAH KE VALIDATOR
	if input2.NewPassword != input2.RePassword {
		return nil, te.XErrors{"rePassword": te.XError{Code: "equalToField", Message: l.I.Msg("equalToField") + " new password"}}
	}

	ed := sh.Event{
		Feature: "reset-password",
		Action:  "request",
		Source:  source,
		Status:  "failed",
		ECode:   "save-fail",
	}

	var user m.User
	result := dg.I.Where(&m.User{Email: input1.Email}).First(&user)
	if result.RowsAffected == 0 {
		return nil, nil
	} else if result.Error != nil {
		ed.ECode = "data-fetch-fail"
		ed.EDetail = result.Error.Error()
		return nil, sh.SetError(ed, user)
	}

	password, err := p.Hash(input2.NewPassword)
	if err != nil {
		ed.ECode = "data-process-fail"
		ed.EDetail = err.Error()
		return nil, sh.SetError(ed, user)
	} else {
		user.Password = &password
	}

	if result := dg.I.Save(&user); result.Error != nil {
		ed.ECode = "data-update-fail"
		ed.EDetail = result.Error.Error()
		return nil, sh.SetError(ed, user)
	}

	dg.I.Unscoped().Delete(&userToken)
	return l.I.Msg("request-ok"), nil
}
