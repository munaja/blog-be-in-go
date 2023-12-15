package blogom

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	hh "github.com/karincake/apem/handlerhelper"

	m "github.com/munaja/blog-practice-be-using-go/internal/model/blog"
	sau "github.com/munaja/blog-practice-be-using-go/internal/service/auth"
	s "github.com/munaja/blog-practice-be-using-go/internal/service/blog"
)

func Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	authInfo, ok := ctx.Value("authInfo").(*sau.AuthInfo)
	if !ok {
		hh.WriteJSON(w, http.StatusUnauthorized, nil, nil)
		return
	}

	var res *m.Blog
	var err error
	status := r.URL.Query().Get("status")
	if status == "publish" {
		var input m.CreatePublishDto
		if hh.ValidateStructByIOR(w, r.Body, &input) == false {
			return
		}
		input.User_Id = &authInfo.User_Id
		input.Status = m.BSActive
		res, err = s.Create(input)
	} else {
		var input m.CreateDraftDto
		if hh.ValidateStructByIOR(w, r.Body, &input) == false {
			return
		}
		input.User_Id = &authInfo.User_Id
		input.Status = m.BSDraft
		res, err = s.Create(input)
	}

	hh.DataResponse(w, res, err, nil, nil)
}

func Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	authInfo, ok := ctx.Value("authInfo").(*sau.AuthInfo)
	if !ok {
		hh.WriteJSON(w, http.StatusUnauthorized, nil, nil)
		return
	}

	id := hh.ValidateInt(w, "id", chi.URLParam(r, "id"))
	if id <= 0 {
		return
	}

	var res *m.Blog
	var err error
	status := r.URL.Query().Get("status")
	if status == "publish" {
		var input m.UpdatePublishDto
		if hh.ValidateStructByIOR(w, r.Body, &input) == false {
			return
		}
		input.User_Id = &authInfo.User_Id
		input.Status = m.BSActive
		res, err = s.Update(id, input)
	} else {
		var input m.UpdateDraftDto
		if hh.ValidateStructByIOR(w, r.Body, &input) == false {
			return
		}
		input.User_Id = &authInfo.User_Id
		input.Status = m.BSDraft
		res, err = s.Update(id, input)
	}

	hh.DataResponse(w, res, err, nil, nil)
}

func Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	authInfo, ok := ctx.Value("authInfo").(*sau.AuthInfo)
	if !ok {
		hh.WriteJSON(w, http.StatusUnauthorized, nil, nil)
		return
	}

	id := hh.ValidateInt(w, "id", chi.URLParam(r, "id"))
	if id <= 0 {
		return
	}

	input := m.DeleteDto{User_Id: &authInfo.User_Id}
	res, err := s.Delete(id, input)

	hh.DataResponse(w, res, err, nil, nil)
}

func GetList(w http.ResponseWriter, r *http.Request) {
}

func GetDetail(w http.ResponseWriter, r *http.Request) {
}
