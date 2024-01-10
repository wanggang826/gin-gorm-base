package user_service

import (
	"gin-gorm-base/models"
	"gin-gorm-base/pkg/douyin"
	"gin-gorm-base/pkg/logging"
	"strconv"
	"time"
)

func SaveUser(code, anonymousCode string) (int, error) {
	dyService := douyin.NewDyService()
	session, err := dyService.Code2Session(code, anonymousCode)
	if err != nil {
		logging.Error("Code2Session error:", err)
		return 0, err
	}
	var userExist *models.User
	if session.Openid != "" {
		userExist = models.GetUserByUnionId(session.Unionid)
	} else if session.AnonymousOpenid != "" {
		userExist = models.GetUserByAnonymousOpenid(session.AnonymousOpenid)
	}

	userId := 0
	if userExist != nil {
		if session.SessionKey != "" {
			editUser := &models.User{
				Id: userExist.Id,
			}
			editData := make(map[string]interface{})
			editData["session_key"] = session.SessionKey
			editData["update_time"] = int(time.Now().Unix())
			_, err := editUser.Edit(editData)
			if err != nil {
				return userExist.Id, err
			}
		}
		userId = userExist.Id
	} else {
		newUser := &models.User{
			Openid:          session.Openid,
			Unionid:         session.Unionid,
			AnonymousOpenid: session.AnonymousOpenid,
			SessionKey:      session.SessionKey,
			CreatedAt:       int(time.Now().Unix()),
			UpdatedAt:       int(time.Now().Unix()),
		}
		_, err := newUser.Add()
		if err != nil {
			logging.Error("Add user error:", err)
			return 0, err
		}
		userId = newUser.Id
	}
	return userId, nil
}

func GetUserOpenId(userId int) (string, string) {
	user, _ := models.GetUserById(userId)
	if user == nil {
		return "", ""
	}
	return user.Openid, user.AnonymousOpenid
}

type ResUser struct {
	Id              int    `json:"id"`
	Nickname        string `json:"nickname"`
	VipExpiresTime  int    `json:"vip_expires_time"`
	VipDay          int    `json:"vip_day"`
	BtBalance       int    `json:"bt_balance"`
	LoginType       int    `json:"login_type"`
	AnonymousOpenid string `json:"anonymous_openid"`
	Openid          string `json:"openid"`
	CreatedAt       int    `json:"created_at"`
}

func GetUserInfo(uid int) *ResUser {
	dbUser, _ := models.GetUserById(uid)
	user := &ResUser{}
	nowTime := int(time.Now().Unix())
	if dbUser != nil {
		user.Id = dbUser.Id
		user.Nickname = "游客" + strconv.Itoa(dbUser.Id)
		user.VipExpiresTime = dbUser.VipExpiresTime
		user.Openid = dbUser.Openid
		user.AnonymousOpenid = dbUser.AnonymousOpenid
		if nowTime < dbUser.VipExpiresTime {
			user.VipDay = (user.VipExpiresTime-nowTime)/86400 + 1
		}
		user.CreatedAt = dbUser.CreatedAt
		if dbUser.Openid == "" {
			user.LoginType = 1
		}
	}
	return user
}
