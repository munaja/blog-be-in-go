package account

import (
	"errors"
	"fmt"

	l "github.com/karincake/apem/lang"

	dg "github.com/karincake/apem/databasegorm"
	m "github.com/munaja/blog-practice-be-using-go/internal/model/user"
	"gorm.io/gorm"
)

func checkUserStatus(user *m.User, email string, status m.Status) error {
	err := dg.I.Where("Email = ?", email).First(user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New(fmt.Sprintf(l.I.Msg("data-notFound-condition"), "user", "email", email))
	} else if err != nil {
		return errors.New(l.I.Msg("data-fetch-fail"))
	} else if user.Status != nil && *(user.Status) != status {
		return errors.New(fmt.Sprintf(l.I.Msg("data-state-mismatch"), "user", m.GetUSText(status)))
	}
	return nil
}
