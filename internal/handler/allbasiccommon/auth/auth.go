package auth

import (
	"context"
	"net/http"

	hh "github.com/karincake/apem/handlerhelper"
	td "github.com/karincake/tempe/data"
	te "github.com/karincake/tempe/error"

	m "github.com/munaja/blog-practice-be-using-go/internal/model/user"
	s "github.com/munaja/blog-practice-be-using-go/internal/service/auth"
)

var Position m.Position

func Login(w http.ResponseWriter, r *http.Request) {
	var input m.LoginDto
	if hh.ValidateStructByIOR(w, r.Body, &input) == false {
		return
	}

	input.Position = Position
	boolTrue := true
	if Position == m.UPOperator {
		input.OptStatus = &boolTrue
	} else if Position == m.UPOwner {
		input.OwnerStatus = &boolTrue
	}
	res, err := s.GenToken(input)
	if err != nil {
		hh.WriteJSON(w, http.StatusUnauthorized, td.II{"errors": err}, nil)
	} else {
		hh.DataResponse(w, res, nil, nil, nil)
	}
}

func Logout(w http.ResponseWriter, r *http.Request) {
	authInfo := context.Context.Value(r.Context(), "authInfo").(*s.AuthInfo)
	s.RevokeToken(authInfo.Uuid)
	hh.WriteJSON(w, http.StatusOK, td.IS{"message": "logged out"}, nil)
}

func GuardMW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accessDetail, err := s.ExtractToken(r, s.AccessToken)
		if err != nil {
			hh.WriteJSON(w, http.StatusUnauthorized, err.(te.XError), nil)
			return
		}
		ctx := context.WithValue(r.Context(), "authInfo", accessDetail)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
