package accountom

import (
	"net/http"

	hh "github.com/karincake/apem/handlerhelper"
	td "github.com/karincake/tempe/data"

	m "github.com/munaja/blog-practice-be-using-go/internal/model/user"
	sac "github.com/munaja/blog-practice-be-using-go/internal/service/account"
	sau "github.com/munaja/blog-practice-be-using-go/internal/service/auth"
)

func ChangePassword(w http.ResponseWriter, r *http.Request) {
	var input m.ChangePassDto
	if hh.ValidateStructByIOR(w, r.Body, &input) == false {
		return
	}

	ctx := r.Context()
	authInfo, ok := ctx.Value("authInfo").(*sau.AuthInfo)
	if !ok {
		hh.WriteJSON(w, http.StatusUnauthorized, nil, nil)
		return
	}

	res, err := sac.ChangePassword(authInfo.User_Id, input)
	hh.DataResponse(w, res, err, nil, nil)
}

func GetDetail(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	authInfo, ok := ctx.Value("authInfo").(*sau.AuthInfo)
	if !ok {
		hh.WriteJSON(w, http.StatusUnauthorized, nil, nil)
		return
	}

	data := td.II{
		"data": authInfo,
	}
	hh.WriteJSON(w, http.StatusOK, data, nil)
}
