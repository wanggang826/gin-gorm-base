package models

import (
	"gin-gorm-base/pkg/logging"
	"github.com/gin-gonic/gin"
)

const TableNameYyqAdmin = "admin"

// Admin 账号表
type Admin struct {
	Id            int          `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	Username      string       `gorm:"column:username;not null;comment:用户名" json:"username"`                    // 用户名
	Password      string       `gorm:"column:password;not null;comment:密码" json:"password"`                     // 密码
	RoleID        int          `gorm:"column:role_id;not null;comment:角色ID" json:"role_id"`                     // 角色ID
	Name          string       `gorm:"column:name;not null;comment:名称" json:"name"`                             // 姓名
	Nickname      string       `gorm:"column:nickname;not null;comment:名称" json:"nickname"`                     // 名称
	Mobile        string       `gorm:"column:mobile;not null;comment:名称" json:"mobile"`                         // 手机号
	LastLoginTime int          `gorm:"column:last_login_time;not null;comment:最后一次登录时间" json:"last_login_time"` // 最后一次登录时间
	LastLoginIP   string       `gorm:"column:last_login_ip;not null;comment:最后登录ip" json:"last_login_ip"`       // 最后登录ip
	Status        int          `gorm:"column:status;not null;default:1;comment:状态：2禁用，1启用" json:"status"`       // 状态：2禁用，1启用
	IsDelete      bool         `gorm:"column:is_delete;not null;comment:是否删除 0否1是" json:"is_delete"`            // 是否删除 0否1是
	CreatedAt     int          `gorm:"column:create_time" json:"create_time"`
	UpdatedAt     int          `gorm:"column:update_time" json:"update_time"`
	Ctx           *gin.Context `gorm:"-"`
}

func (*Admin) TableName() string {
	return TableNameYyqAdmin
}

func (item *Admin) Add() (int, error) {
	res := db.Create(item)
	if res.Error != nil {
		logging.WithCtxError(item.Ctx, "DB: Admin.Add: ", "msg", res.Error.Error(), "item", item)
	}
	return item.Id, res.Error
}

func (item *Admin) Edit(data interface{}) (int64, error) {
	res := db.Model(item).Updates(data)

	if res.Error != nil {
		logging.WithCtxError(item.Ctx, "DB: Admin.Edit: ", "msg", res.Error.Error(), "item", item, "data", item)
	}
	return res.RowsAffected, res.Error
}

func (item *Admin) Del() (int64, error) {
	res := db.Where(item).Delete(item)
	if res.Error != nil {
		logging.WithCtxError(item.Ctx, "DB: Admin.Delete: ", "msg", res.Error.Error(), "item", item)
	}
	return res.RowsAffected, res.Error
}

func (item *Admin) Get() (bool, error) {
	res := db.Where(item).First(&item)
	if res.Error != nil {
		logging.WithCtxError(item.Ctx, "DB: Admin.Get: ", "msg", res.Error.Error(), "item", item)
		return false, res.Error
	}
	return true, nil
}

func GetAdminList(ctx *gin.Context, nickname string, status, roleId int) ([]*Admin, error) {
	var items []*Admin
	dbItem := db.Model(&Admin{}).Where("is_delete = 0")
	if nickname != "" {
		dbItem = dbItem.Where("nickname like ?", "%"+nickname+"%")
	}
	if status > 0 {
		dbItem = dbItem.Where("status = ?", status)
	}
	if roleId > 0 {
		dbItem = dbItem.Where("role_id = ?", roleId)
	}
	res := dbItem.Order("id DESC").Find(&items)
	if res.Error != nil {
		logging.WithCtxError(ctx, "DB: Admin.GetAdminList: ", "msg", res.Error.Error())
		return nil, res.Error
	}
	return items, nil
}

func GetAdminCount(ctx *gin.Context, nickname string, status int) (int, error) {
	var count int64
	dbItem := db.Model(&Admin{}).Where("is_delete = 0 and id > 1")
	if nickname != "" {
		dbItem = dbItem.Where("nickname like ?", "%"+nickname+"%")
	}
	if status > 0 {
		dbItem = dbItem.Where("status = ?", status)
	}
	res := dbItem.Count(&count)
	if res.Error != nil {
		logging.WithCtxError(ctx, "DB: Admin.GetAdminCount: ", "msg", res.Error.Error())
		return 0, res.Error
	}
	return int(count), nil
}

func GetAdminPage(ctx *gin.Context, nickname string, status int, limit int, offset int) ([]*Admin, error) {
	var items []*Admin
	dbItem := db.Model(&Admin{}).Where("is_delete = 0 and id > 1")
	if nickname != "" {
		dbItem = dbItem.Where("nickname like ?", "%"+nickname+"%")
	}
	if status > 0 {
		dbItem = dbItem.Where("status = ?", status)
	}
	res := dbItem.Limit(limit).Offset(offset).Order("id DESC").Find(&items)
	if res.Error != nil {
		logging.WithCtxError(ctx, "DB: Admin.GetAdminPage: ", "msg", res.Error.Error())
		return nil, res.Error
	}
	return items, nil
}

func GetAdminById(ctx *gin.Context, id int) *Admin {
	var item *Admin
	res := db.Model(&Admin{}).Where("id = ?", id).First(&item)
	if res.Error != nil {
		logging.WithCtxError(ctx, "DB: Admin.GetAdminById: ", "msg", res.Error.Error())
		return nil
	}
	return item
}

func GeAdminByUsername(ctx *gin.Context, username string) *Admin {
	var item *Admin
	res := db.Model(&Admin{}).Where("username = ? and is_delete = 0", username).First(&item)
	if res.Error != nil {
		logging.WithCtxError(ctx, "DB: Admin.GeAdminByUsername: ", "msg", res.Error.Error())
		return nil
	}
	return item
}

func GetAdminMap(ctx *gin.Context) (map[int]*Admin, error) {
	var items []*Admin
	dbItem := db.Model(&Admin{})
	res := dbItem.Order("id DESC").Find(&items)
	if res.Error != nil {
		logging.Error("DB: Book.GetAdminMap: ", "msg", res.Error.Error())
		return nil, res.Error
	}
	result := make(map[int]*Admin)
	for _, v := range items {
		result[v.Id] = v
	}
	return result, nil
}
