package blog

import (
	"errors"
	"strconv"

	sc "github.com/jinzhu/copier"
	dg "github.com/karincake/apem/databasegorm"
	l "github.com/karincake/apem/lang"
	g "github.com/karincake/getuk"
	td "github.com/karincake/tempe/data"
	te "github.com/karincake/tempe/error"
	"gorm.io/gorm"

	m "github.com/munaja/blog-practice-be-using-go/internal/model/blog"
	sh "github.com/munaja/blog-practice-be-using-go/pkg/servicehelper"
)

const source = "blog"

func Create[T m.CreateDraftDto | m.CreatePublishDto](dto T) (*m.Blog, error) {
	data := &m.Blog{}
	switch dt := any(dto).(type) {
	case m.CreateDraftDto:
		if err := sc.Copy(data, dt); err != nil {
			return nil, te.XErrors{"struct": te.XError{Code: "copy-fail", Message: l.I.Msg("data-copy-fail")}}
		}
	case m.CreatePublishDto:
		if err := sc.Copy(data, dt); err != nil {
			return nil, te.XErrors{"struct": te.XError{Code: "copy-fail", Message: l.I.Msg("data-copy-fail")}}
		}
	}

	if err := dg.I.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&data).Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		ed := sh.Event{
			Feature: "blog",
			Action:  "create-data",
			Source:  source,
			Status:  "failed",
			ECode:   "data-create-fail",
			EDetail: err.Error(),
		}
		return nil, sh.SetError(ed, nil)
	}

	return data, nil
}

func Update[T m.UpdateDraftDto | m.UpdatePublishDto](id int, dto T) (*m.Blog, error) {
	data := &m.Blog{}

	switch dt := any(dto).(type) {
	case m.UpdateDraftDto:
		if err := dg.I.Where("Id = ? AND User_Id = ?", id, dt.User_Id).First(data).Error; err != nil {
			return nil, returnFetchError(err)
		}
		if err := sc.Copy(data, dt); err != nil {
			return nil, te.XErrors{"struct": te.XError{Code: "copy-fail", Message: l.I.Msg("data-copy-fail")}}
		}
	case m.UpdatePublishDto:
		if err := dg.I.Where("Id = ? AND User_Id = ?", id, dt.User_Id).First(data).Error; err != nil {
			return nil, returnFetchError(err)
		}
		if err := sc.Copy(data, dt); err != nil {
			return nil, te.XErrors{"struct": te.XError{Code: "copy-fail", Message: l.I.Msg("data-copy-fail")}}
		}
	}

	if err := dg.I.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(&data).Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		ed := sh.Event{
			Feature: "blog",
			Action:  "create-data",
			Source:  source,
			Status:  "failed",
			ECode:   "data-create-fail",
			EDetail: err.Error(),
		}
		return nil, sh.SetError(ed, nil)
	}

	return data, nil
}

func Delete(id int, dto m.DeleteDto) (*string, error) {
	data := &m.Blog{}

	if err := dg.I.Where("Id = ? AND User_Id = ?", id, dto.User_Id).First(data).Error; err != nil {
		return nil, returnFetchError(err)
	}
	if err := dg.I.Delete(data).Error; err != nil {
		return nil, err
	}

	msg := l.I.Msg("data-delete-success")
	return &msg, nil
}

func GetList(dto m.ListFilterDto) (*td.Data, error) {
	var data []m.Blog
	var count int64

	var pagination g.Pagination
	result := dg.I.
		Model(&m.Blog{}).
		Scopes(g.Filter(dto)).
		Count(&count).
		Scopes(g.Paginate(dto, &pagination)).
		Find(&data)

	if result.Error != nil {
		ed := sh.Event{
			Feature: "blog",
			Action:  "get-list",
			Source:  source,
			Status:  "failed",
			ECode:   "data-create-fail",
			EDetail: result.Error.Error(),
		}
		return nil, sh.SetError(ed, data)
	}

	return &td.Data{
		Meta: td.IS{
			"totalCount":   strconv.Itoa(int(count)),
			"currentCount": strconv.Itoa(int(result.RowsAffected)),
			"page":         strconv.Itoa(pagination.Page),
			"pageSize":     strconv.Itoa(pagination.PageSize),
		},
		Data: data,
	}, nil
}

func GetDetail(id int, dto m.DetailFilterDto) (*td.Data, error) {
	data := &m.Blog{}

	result := dg.I.
		Where("Id = ? ", id).
		Scopes(g.Filter(dto)).
		First(data)
	if result.Error != nil {
		return nil, returnFetchError(result.Error)
	}

	return &td.Data{
		Data: data,
	}, nil

}

func returnFetchError(err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return te.XError{Code: "data-notFound", Message: l.I.Msg("data-notFound")}
	} else {
		return err
	}
}
