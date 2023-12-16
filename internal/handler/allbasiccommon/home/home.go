package home

import (
	"net/http"

	hh "github.com/karincake/apem/handlerhelper"
	td "github.com/karincake/tempe/data"
)

func Index(w http.ResponseWriter, r *http.Request) {
	hh.WriteJSON(w, http.StatusOK, td.Message{
		Message: "Welcome to SKGOBO API!!",
	}, nil)
}
