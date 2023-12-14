package profile

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	hh "github.com/karincake/apem/handlerhelper"

	s "github.com/munaja/blog-practice-be-using-go/internal/service/profile"
)

func GetDetail(w http.ResponseWriter, r *http.Request) {
	user_name := hh.ValidateString(w, "user_name", chi.URLParam(r, "user_name"))
	if user_name == "" {
		return
	}

	res, err := s.GetDetail(user_name)
	hh.DataResponse(w, res, err, nil, nil)
}
