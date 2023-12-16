package blog

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	hh "github.com/karincake/apem/handlerhelper"

	m "github.com/munaja/blog-practice-be-using-go/internal/model/blog"
	s "github.com/munaja/blog-practice-be-using-go/internal/service/blog"
)

func GetList(w http.ResponseWriter, r *http.Request) {
	input := m.ListFilterDto{
		Title_Opt: "left",
	}
	if hh.ValidateStructByURL(w, *r.URL, &input) == false {
		return
	}

	res, err := s.GetList(input)
	hh.DataResponse(w, res, err)
}

func GetDetail(w http.ResponseWriter, r *http.Request) {
	id := hh.ValidateInt(w, "id", chi.URLParam(r, "id"))
	if id <= 0 {
		return
	}

	input := m.DetailFilterDto{}
	res, err := s.GetDetail(id, input)
	hh.DataResponse(w, res, err)
}
