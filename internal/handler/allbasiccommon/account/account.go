package account

import (
	"net/http"

	hh "github.com/karincake/apem/handlerhelper"
	td "github.com/karincake/tempe/data"

	m "github.com/munaja/blog-practice-be-using-go/internal/model/user"
	sac "github.com/munaja/blog-practice-be-using-go/internal/service/account"
	sau "github.com/munaja/blog-practice-be-using-go/internal/service/auth"
)

func Register(w http.ResponseWriter, r *http.Request) {
	var input m.RegisterDto
	if hh.ValidateStructByIOR(w, r.Body, &input) == false {
		return
	}

	res, err := sac.Register(input)
	hh.DataResponse(w, res, err, nil, nil)
}

func ConfirmByEmail(w http.ResponseWriter, r *http.Request) {
	var input m.ResendConfirmationEmailDto
	if hh.ValidateStructByURL(w, *r.URL, &input) == false {
		return
	}

	res, err := sac.ConfirmByEmail(input)
	hh.DataResponse(w, res, err, nil, nil)
}

func ResendConfirmationEmail(w http.ResponseWriter, r *http.Request) {
	var input m.ResendEmailConfirmDto
	if hh.ValidateStructByIOR(w, r.Body, &input) == false {
		return
	}

	res, err := sac.ResendConfirmationEmail(input)
	hh.DataResponse(w, res, err, nil, nil)
}

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

func RequestResetPassword(w http.ResponseWriter, r *http.Request) {
	var input m.RequestResetPassDto
	if hh.ValidateStructByIOR(w, r.Body, &input) == false {
		return
	}

	res, err := sac.RequestResetPass(input)
	hh.DataResponse(w, res, err, nil, nil)
}

func CheckResetPassword(w http.ResponseWriter, r *http.Request) {
	var input m.CheckResetPassDto
	if hh.ValidateStructByURL(w, *r.URL, &input) == false {
		return
	}

	res, err := sac.CheckResetPass(input)
	hh.DataResponse(w, res, err, nil, nil)
}

func ResetPassword(w http.ResponseWriter, r *http.Request) {
	var input1 m.CheckResetPassDto
	if hh.ValidateStructByURL(w, *r.URL, &input1) == false {
		return
	}
	var input2 m.ResetPassDto
	if hh.ValidateStructByIOR(w, r.Body, &input2) == false {
		return
	}

	res, err := sac.ResetPass(input1, input2)
	hh.DataResponse(w, res, err, nil, nil)
}

func Check(w http.ResponseWriter, r *http.Request) {
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
