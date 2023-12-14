package profile

import (
	"github.com/jinzhu/copier"
	sc "github.com/jinzhu/copier"
	dg "github.com/karincake/apem/databasegorm"
	"gorm.io/gorm"

	l "github.com/karincake/apem/lang"
	te "github.com/karincake/tempe/error"

	m "github.com/munaja/blog-practice-be-using-go/internal/model/person"
	su "github.com/munaja/blog-practice-be-using-go/internal/service/user"
	sh "github.com/munaja/blog-practice-be-using-go/pkg/servicehelper"
)

const source = "profile"

func GetDetail(user_name string, opt ...string) (any, error) {
	// get user first
	user, err := su.GetDetailByName(user_name)
	if err != nil {
		return nil, te.XErrors{"data-notFound": te.XError{Code: "data-fetch-fail", Message: l.I.Msg("data-fetch-fail") + ", " + err.Error()}}
	}

	var person *m.Person
	err = dg.I.Where("User_Id = ?", user.Id).First(&person).Error
	// special case, profile will still showing empty data if no data found
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, te.XErrors{"data-notFound": te.XError{Code: "process-fail", Message: l.I.Msg("data-fetch-fail")}}
	}

	if len(opt) > 0 && opt[0] == "of-mine" {
		var profil m.ProfilOMDto = m.ProfilOMDto{}
		profil.User_Id = user.Id
		profil.User_Name = user.Name
		profil.User_CreatedAt = user.CreatedAt
		if err = sc.CopyWithOption(&profil, person, copier.Option{IgnoreEmpty: true}); err != nil {
			return nil, te.XErrors{"data": te.XError{Code: "data-copy-fail"}}
		}
		return profil, nil
	} else {
		var profil m.ProfilDto = m.ProfilDto{}
		profil.User_Id = user.Id
		profil.User_Name = user.Name
		profil.User_CreatedAt = user.CreatedAt
		if err = sc.CopyWithOption(&profil, person, copier.Option{IgnoreEmpty: true}); err != nil {
			return nil, te.XErrors{"data": te.XError{Code: "data-copy-fail"}}
		}
		return profil, nil
	}
}

func Update(user_name string, data m.UpdateDto) (any, error) {
	// prepare event
	event := sh.Event{
		Feature: source,
		Action:  "get-resource",
		Source:  source,
		Status:  "failed",
		ECode:   "data-fetch-fail",
	}

	// get user first
	user, err := su.GetIdByName(user_name)
	if err != nil {
		return nil, te.XErrors{"data": te.XError{Code: "process-fail", Message: l.I.Msg("data-fetch-fail") + ", " + err.Error()}}
	}

	// old data
	var oldData *m.Person = &m.Person{}
	err = dg.I.Where("User_Id = ?", user["Id"]).First(&oldData).Error
	// special case, profile will still showing empty data if no data found
	if err != nil && err != gorm.ErrRecordNotFound {
		event.EDetail = err.Error()
		return nil, te.XErrors{"data": te.XError{Code: "process-fail", Message: l.I.Msg("data-fetch-fail") + ", " + err.Error()}}
	}

	new := false
	if err == gorm.ErrRecordNotFound {
		new = true
		user_id, _ := user["Id"].(int)
		oldData.User_Id = user_id
	}

	if err := sc.CopyWithOption(&oldData, data, sc.Option{IgnoreEmpty: true}); err != nil {
		return nil, sh.SetError(event, data)
	}

	if !new {
		err = dg.I.Save(oldData).Error
	} else {
		err = dg.I.Create(oldData).Error
	}

	if err != nil {
		return nil, err
	}
	return oldData, nil
}
