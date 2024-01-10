package models

import "gin-gorm-base/pkg/logging"

const TableNameUser = "user"

// User 用户信息表
type User struct {
	Id              int    `gorm:"column:id;primaryKey;autoIncrement:true;comment:主键ID" json:"id"`               // 主键ID
	Openid          string `gorm:"column:openid;not null;comment:openid" json:"openid"`                          // openid
	AnonymousOpenid string `gorm:"column:anonymous_openid;not null;comment:接口返回的匿名登录凭证" json:"anonymous_openid"` // 接口返回的匿名登录凭证
	Unionid         string `gorm:"column:unionid;not null;comment:unionId" json:"unionid"`                       // unionId
	SessionKey      string `gorm:"column:session_key;not null;comment:session_key" json:"session_key"`           // session_key
	Nickname        string `gorm:"column:nickname;not null;comment:昵称" json:"nickname"`                          // 昵称
	Gender          int    `gorm:"column:gender;not null;comment:用户的性别，1=男性，2=女性，0=未知" json:"gender"`            // 用户的性别，1=男性，2=女性，0=未知
	City            string `gorm:"column:city;not null;comment:用户所在城市" json:"city"`                              // 用户所在城市
	Country         string `gorm:"column:country;not null;comment:用户所在国家" json:"country"`                        // 用户所在国家
	Province        string `gorm:"column:province;not null;comment:用户所在省份" json:"province"`                      // 用户所在省份
	Language        string `gorm:"column:language;not null;comment:用户的语言，简体中文为zh_CN" json:"language"`            // 用户的语言，简体中文为zh_CN
	AvatarURL       string `gorm:"column:avatar_url;not null;comment:头像" json:"avatar_url"`                      // 头像
	Mobile          string `gorm:"column:mobile;not null;comment:手机号，带区号" json:"mobile"`                         // 手机号，带区号
	CountryCode     string `gorm:"column:country_code;not null;comment:国家码" json:"country_code"`                 // 国家码
	Email           string `gorm:"column:email;not null;comment:邮箱" json:"email"`                                // 邮箱
	VipExpiresTime  int    `gorm:"column:vip_expires_time;not null;comment:会员到期时间" json:"vip_expires_time"`      // 会员到期时间
	IsBanned        int    `gorm:"column:is_banned;not null;comment:用户是否被禁用，0否，1是" json:"is_banned"`             // 用户是否被禁用，0否，1是
	CreateTime      int    `gorm:"column:create_time;not null;comment:创建时间" json:"create_time"`                  // 创建时间
	UpdateTime      int    `gorm:"column:update_time;not null;comment:最后更新时间" json:"update_time"`                // 最后更新时间
}

// TableName User's table name
func (*User) TableName() string {
	return TableNameUser
}

func (item *User) Add() (int, error) {
	res := db.Create(item)
	if res.Error != nil {
		logging.Error("DB: User.Add: ", "msg", res.Error.Error(), "item", item)
	}
	return item.Id, res.Error
}

func (item *User) Edit(data interface{}) (int64, error) {
	res := db.Model(item).Updates(data)

	if res.Error != nil {
		logging.Error("DB: User.Edit: ", "msg", res.Error.Error(), "item", item, "data", item)
	}
	return res.RowsAffected, res.Error
}

func (item *User) Del() (int64, error) {
	res := db.Where(item).Delete(item)
	if res.Error != nil {
		logging.Error("DB: User.Delete: ", "msg", res.Error.Error(), "item", item)
	}
	return res.RowsAffected, res.Error
}

func (item *User) Get() (bool, error) {
	res := db.Where(item).First(&item)
	if res.Error != nil {
		logging.Error("DB: User.Get: ", "msg", res.Error.Error(), "item", item)
		return false, res.Error
	}
	return true, nil
}

func GetUserById(id int) (*User, error) {
	user := &User{}
	result := db.Model(&User{}).Where("id = ? ", id).First(user)
	return user, result.Error
}

func GetUserByUnionId(unionId string) *User {
	var item *User
	dbItem := db.Model(&User{}).Where("unionid = ? ", unionId)
	res := dbItem.First(&item)
	if res.Error != nil {
		logging.Error("DB: User.GetUserByUnionId: ", "msg", res.Error.Error())
		return nil
	}
	return item
}

func GetUserByAnonymousOpenid(anonymousOpenid string) *User {
	var item *User
	dbItem := db.Model(&User{}).Where("anonymous_openid = ? ", anonymousOpenid)
	res := dbItem.First(&item)
	if res.Error != nil {
		logging.Error("DB: User.GetUserByAnonymousOpenid: ", "msg", res.Error.Error())
		return nil
	}
	return item
}
