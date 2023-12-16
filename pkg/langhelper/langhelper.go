package langhelper

import (
	"fmt"

	l "github.com/karincake/apem/lang"
	t "github.com/karincake/tempe/error"
)

func ErrorMsgGen(errCode string, errDetail ...string) string {
	errMsg := ""
	if len(errDetail) == 0 || errDetail[0] == "" {
		errMsg = l.I.Msg(errCode)
	} else {
		errMsg = fmt.Sprintf(l.I.Msg(errCode), errDetail[0])
	}
	return errMsg
}

func ErrorBundler(errCode string, errDetail ...string) t.XError {
	if len(errDetail) == 0 {
		return t.XError{Code: errCode, Message: ErrorMsgGen(errCode)}
	} else {
		return t.XError{Code: errCode, Message: ErrorMsgGen(errCode, errDetail[0])}
	}
}
