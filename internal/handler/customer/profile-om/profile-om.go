package profileom

import (
	"net/http"

	hh "github.com/karincake/apem/handlerhelper"

	m "github.com/munaja/blog-practice-be-using-go/internal/model/person"
	sau "github.com/munaja/blog-practice-be-using-go/internal/service/auth"
	s "github.com/munaja/blog-practice-be-using-go/internal/service/profile"
)

func GetDetail(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	authInfo, ok := ctx.Value("authInfo").(*sau.AuthInfo)
	if !ok {
		hh.WriteJSON(w, http.StatusUnauthorized, nil, nil)
		return
	}

	res, err := s.GetDetail(authInfo.User_Name, "of-mine")
	hh.DataResponse(w, res, err, nil, nil)
}

func Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	authInfo, ok := ctx.Value("authInfo").(*sau.AuthInfo)
	if !ok {
		hh.WriteJSON(w, http.StatusUnauthorized, nil, nil)
		return
	}

	var payload m.UpdateDto
	if hh.ValidateStructByIOR(w, r.Body, &payload) == false {
		return
	}

	res, err := s.Update(authInfo.User_Name, payload)
	hh.DataResponse(w, res, err, nil, nil)
}
