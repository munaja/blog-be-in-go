package auth

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/spf13/viper"

	dg "github.com/karincake/apem/databasegorm"
	l "github.com/karincake/apem/lang"
	ms "github.com/karincake/apem/memstorageredis"
	td "github.com/karincake/tempe/data"
	te "github.com/karincake/tempe/error"
	lh "github.com/munaja/blog-practice-be-using-go/pkg/langhelper"
	p "github.com/munaja/blog-practice-be-using-go/pkg/password"

	mu "github.com/munaja/blog-practice-be-using-go/internal/model/user"
)

//	type TokenDetails struct {
//		AccessToken  string
//		RefreshToken string
//		AccessUuid   string
//		RefreshUuid  string
//		AtExpires    int64
//		RtExpires    int64
//	}

// Generates token and store in redis at one place
// just return the error code
func GenToken(input mu.LoginDto) (any, error) {
	// Get User
	var user mu.User
	if errCode := getAndCheck(&user, mu.User{Name: input.Name}); errCode != "" {
		return nil, te.XErrors{"authentication": te.XError{Code: errCode, Message: lh.ErrorMsgGen(errCode)}}
	}

	if *&user.LoginAttemptCount > 5 {
		if &user.LastSuccessLogin != nil {
			now := time.Now()
			lastAllowdLogin := user.LastAllowdLogin
			if lastAllowdLogin.After(now.Add(-time.Hour * 1)) {
				return nil, te.XErrors{"authentication": te.XError{Code: "auth-login-tooMany", Message: lh.ErrorMsgGen("auth-login-tooMany")}}
			} else {
				user.LastAllowdLogin = time.Now()
				user.LoginAttemptCount = 0
				dg.I.Save(&user)
			}
		} else {
			user.LastAllowdLogin = time.Now()
			dg.I.Save(&user)
			return nil, te.XErrors{"authentication": te.XError{Code: "auth-login-tooMany", Message: lh.ErrorMsgGen("auth-login-tooMany")}}
		}
	}

	if p.Check(input.Password, *user.Password) == false {
		user.LoginAttemptCount = user.LoginAttemptCount + 1
		dg.I.Save(&user)
		return nil, te.XErrors{"authentication": te.XError{Code: "auth-login-incorrect", Message: lh.ErrorMsgGen("auth-login-incorrect")}}
	} else if *user.Status == mu.USBlocked {
		return nil, te.XErrors{"authentication": te.XError{Code: "auth-login-blocked", Message: lh.ErrorMsgGen("auth-login-blocked")}}
	} else if *user.Status == mu.USNew {
		return nil, te.XErrors{"authentication": te.XError{Code: "auth-login-unverified", Message: lh.ErrorMsgGen("auth-login-unverified")}}
	}

	// Access token prep
	id, err := uuid.NewRandom()
	if err != nil {
		panic(fmt.Sprintf(l.I.Msg("uuid-gen-fail"), err))
	}
	aUuid := id.String()

	// calculate
	durations := strings.Split(strings.ToLower(input.Duration), "-")
	duration := time.Hour * 24
	if len(durations) == 2 {
		val, err := strconv.Atoi(durations[0])
		if err == nil {
			if durations[1] == "m" {
				duration = time.Minute * time.Duration(val)
			} else if durations[1] == "h" {
				duration = time.Hour * time.Duration(val)
			} else if durations[1] == "d" {
				duration = time.Hour * 24 * time.Duration(val)
			}
		}
	}
	atExpires := time.Now().Add(duration).Unix()

	// key
	atSecretKey := viper.GetString("authConf.atSecretKey")

	// Creating Access Token
	atClaims := jwt.MapClaims{}
	atClaims["user_id"] = user.Id
	atClaims["user_name"] = user.Name
	atClaims["user_email"] = user.Email
	atClaims["exp"] = atExpires
	atClaims["uuid"] = aUuid
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	ats, err := at.SignedString([]byte(atSecretKey))
	if err != nil {
		return nil, te.XErrors{"user": te.XError{Code: "token-sign-err", Message: lh.ErrorMsgGen("token-sign-err")}}
	}
	// Save to redis
	now := time.Now()
	atx := time.Unix(atExpires, 0) //converting Unix to UTC(to Time object)
	err = ms.I.Set(aUuid, strconv.Itoa(user.Id), atx.Sub(now)).Err()
	if err != nil {
		panic(fmt.Sprintf(l.I.Msg("redis-store-fail"), err.Error()))
	}

	user.LoginAttemptCount = 0
	user.LastSuccessLogin = time.Now()
	user.LastAllowdLogin = time.Now()
	dg.I.Save(&user)

	// Current data
	return td.II{
		"id":          strconv.Itoa(user.Id),
		"name":        user.Name,
		"email":       user.Email,
		"accessToken": ats,
	}, nil
}

func RevokeToken(uuid string) {
	ms.I.Del(uuid)
}

func VerifyToken(r *http.Request, tokenType TokenType) (data *jwt.Token, errCode, errDetail string) {
	auth := r.Header.Get("Authorization")
	if auth == "" {
		return nil, "auth-missingHeader", ""
	}
	authArr := strings.Split(auth, " ")
	if len(authArr) == 2 {
		auth = authArr[1]
	}

	token, err := jwt.Parse(auth, func(token *jwt.Token) (any, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf(l.I.Msg("token-sign-unexcpeted"), token.Header["alg"])
		}
		if tokenType == AccessToken {
			return []byte(viper.GetString("authConf.atSecretKey")), nil
		} else {
			return []byte(viper.GetString("authConf.rtSecretKey")), nil
		}
	})
	if err != nil {
		return nil, "token-parse-fail", err.Error()
	}
	return token, "", ""
}

func ExtractToken(r *http.Request, tokenType TokenType) (data *AuthInfo, err error) {
	token, errCode, errDetail := VerifyToken(r, tokenType)
	if errCode != "" {
		return nil, te.XError{Code: errCode, Message: lh.ErrorMsgGen(errCode, errDetail)}
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		accessUuid, ok := claims["uuid"].(string)
		if !ok {
			return nil, te.XError{Code: "token-invalid", Message: lh.ErrorMsgGen("token-invalid", "uuid not available")}
		}
		user_id, myErr := strconv.ParseInt(fmt.Sprintf("%.f", claims["user_id"]), 10, 64)
		if myErr != nil {
			return nil, te.XError{Code: "token-invalid", Message: lh.ErrorMsgGen("token-invalid", "uuid is not available")}
		}
		accessUuidRedis := ms.I.Get(accessUuid)
		if accessUuidRedis.String() == "" {
			return nil, te.XError{Code: "token-unidentified", Message: lh.ErrorMsgGen("token-unidentified")}
		}
		user_name := fmt.Sprintf("%v", claims["user_name"])
		user_email := fmt.Sprintf("%.f", claims["user_email"])
		data = &AuthInfo{
			Uuid:       accessUuid,
			User_Id:    int(user_id),
			User_Name:  user_name,
			User_Email: user_email,
		}
		return
	}
	return nil, te.XError{Code: "token", Message: "token-invalid"}
}
