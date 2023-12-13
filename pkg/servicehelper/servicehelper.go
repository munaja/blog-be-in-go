package servicehelper

import (
	"encoding/json"

	l "github.com/karincake/apem/lang"
	lz "github.com/karincake/apem/loggerzap"
	te "github.com/karincake/tempe/error"
	"go.uber.org/zap"
)

// To standardize the error logging format
type Event struct {
	Feature string
	Action  string
	Source  string
	Status  string
	ECode   string
	EDetail string
}

func SetError(d Event, data any) te.XError {
	dataString, _ := json.Marshal(data)
	msg := l.I.Msg(d.ECode)
	lz.I.Error(msg,
		zap.String("feature", d.Feature),
		zap.String("source", d.Source),
		zap.String("action", d.Action),
		zap.String("status", d.Status),
		zap.String("data", string(dataString)))
	return te.XError{Code: d.ECode, Message: l.I.Msg(d.ECode)}
}
