package accountom

import (
	"net/http"

	hh "github.com/karincake/apem/handlerhelper"
	l "github.com/karincake/apem/lang"
	td "github.com/karincake/tempe/data"
	te "github.com/karincake/tempe/error"
)

func checkRepassword(password, repassword string, w http.ResponseWriter, r *http.Request) bool {
	if password != repassword {
		hh.WriteJSON(w, http.StatusUnauthorized, td.II{
			"Meta":   td.IS{"count": "1"},
			"Errors": te.XErrors{"repassword": te.XError{Code: "equalToField", Message: l.I.Msg("equalToField") + ": password"}},
		}, nil)
		return false
	}
	return true
}
