package usertoken

import (
	"errors"
	"time"

	"github.com/google/uuid"
	dg "github.com/karincake/apem/databasegorm"
	"gorm.io/gorm"

	m "github.com/munaja/blog-be-in-go/internal/model/usertoken"
)

// non db-storing
func Request(userToken *m.UserToken, email string, tokenType m.Type, duration uint16, maxAttempt byte) error {
	err := dg.I.Where("User_Email = ? AND Type = ?", email, tokenType).First(&userToken).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if userToken.Id > 0 {
		if errCode := CheckExpiration(*userToken.ExpiredAt); errCode == "" {
			if userToken.AttemptCount >= maxAttempt {
				return errors.New("terlalu banyak percobaan request " + m.GetTypeText(tokenType))
			} else {
				userToken.AttemptCount = userToken.AttemptCount + 1
			}
		} else {
			userToken.AttemptCount = 0
		}
	} else {
		userToken.User_Email = email
		userToken.Type = tokenType
	}
	token, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	expiredAt := time.Now().Local().Add(time.Minute * time.Duration(duration))
	userToken.Token = token
	userToken.ExpiredAt = &expiredAt
	return nil
}

func CheckByEmail(userToken *m.UserToken, email string, token string, tokenType m.Type) (eCode, eDetail string) {
	eDb := dg.I.Where("User_Email = ? AND Token = ? AND Type = ?", email, token, tokenType).First(&userToken).Error
	if errors.Is(eDb, gorm.ErrRecordNotFound) {
		return "userToken-incorrect", ""
	} else if eDb != nil {
		return "data-fetch-fail", eDb.Error()
		// } else if userToken.Type != tokenType {
		// 	return "token-invalidType", ""
	} else if errCode := CheckExpiration(*userToken.ExpiredAt); errCode != "" {
		return errCode, "token is expired"
	}
	return "", ""
}

func CheckExpiration(input time.Time) (errCode string) {
	now := time.Now()
	if now.Before(input) {
		return ""
	} else {
		return "token-expired"
	}
}
