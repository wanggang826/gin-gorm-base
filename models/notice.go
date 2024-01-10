package models

import (
	"gin-gorm-base/pkg/logging"
	"github.com/gin-gonic/gin"
)

const TableNameNotice = "notice"

// Notice 公告
type Notice struct {
	Id        int          `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	Title     string       `gorm:"column:title;not null;comment:标题" json:"title"`                           // 标题
	Image     string       `gorm:"column:image;not null;comment:图片" json:"image"`                           // 图片
	CateID    int          `gorm:"column:cate_id;not null;default:1;comment:种类 1功能更新 2活动通知" json:"cate_id"` // 种类 1功能更新 2活动通知
	IsDelete  bool         `gorm:"column:is_delete;not null;comment:是否删除 0否1是" json:"is_delete"`            // 是否删除 0否1是
	Sort      int          `gorm:"column:sort;not null;default:1;comment:排序" json:"sort"`                   // 排序
	CreatedAt int          `gorm:"column:create_time" json:"create_time"`
	UpdatedAt int          `gorm:"column:update_time" json:"update_time"`
	Ctx       *gin.Context `gorm:"-"`
}

// TableName Notice's table name
func (*Notice) TableName() string {
	return TableNameNotice
}

func (item *Notice) Add() (int, error) {
	res := db.Create(item)
	if res.Error != nil {
		logging.WithCtxError(item.Ctx, "DB: Notice.Add: ", "msg", res.Error.Error(), "item", item)
	}
	return item.Id, res.Error
}

func (item *Notice) Edit(data interface{}) (int64, error) {
	res := db.Model(item).Updates(data)

	if res.Error != nil {
		logging.WithCtxError(item.Ctx, "DB: Notice.Edit: ", "msg", res.Error.Error(), "item", item, "data", item)
	}
	return res.RowsAffected, res.Error
}

func (item *Notice) Del() (int64, error) {
	res := db.Where(item).Delete(item)
	if res.Error != nil {
		logging.WithCtxError(item.Ctx, "DB: Notice.Delete: ", "msg", res.Error.Error(), "item", item)
	}
	return res.RowsAffected, res.Error
}

func (item *Notice) Get() (bool, error) {
	res := db.Where(item).First(&item)
	if res.Error != nil {
		logging.WithCtxError(item.Ctx, "DB: Notice.Get: ", "msg", res.Error.Error(), "item", item)
		return false, res.Error
	}
	return true, nil
}

func GetNoticeList(ctx *gin.Context, cateId int) ([]*Notice, error) {
	var items []*Notice
	dbItem := db.Model(&Notice{}).Where("is_delete = 0 and cate_id = ?", cateId)
	res := dbItem.Order("id DESC").Find(&items)
	if res.Error != nil {
		logging.WithCtxError(ctx, "DB: Notice.GetNoticeList: ", "msg", res.Error.Error())
		return nil, res.Error
	}
	return items, nil
}

func GetNoticeCount(ctx *gin.Context, title string, cateId int) (int, error) {
	var count int64
	dbItem := db.Model(&Notice{}).Where("is_delete = 0")
	if cateId > 0 {
		dbItem = dbItem.Where("cate_id = ?", cateId)
	}
	if title != "" {
		dbItem = dbItem.Where("title like ?", "%"+title+"%")
	}
	res := dbItem.Count(&count)
	if res.Error != nil {
		logging.WithCtxError(ctx, "DB: Notice.GetNoticeCount: ", "msg", res.Error.Error())
		return 0, res.Error
	}
	return int(count), nil
}

func GetNoticePage(ctx *gin.Context, title string, cateId int, limit int, offset int) ([]*Notice, error) {
	var items []*Notice
	dbItem := db.Model(&Notice{}).Where("is_delete = 0")
	if cateId > 0 {
		dbItem = dbItem.Where("cate_id = ?", cateId)
	}
	if title != "" {
		dbItem = dbItem.Where("title like ?", "%"+title+"%")
	}
	res := dbItem.Limit(limit).Offset(offset).Order("id DESC").Find(&items)
	if res.Error != nil {
		logging.WithCtxError(ctx, "DB: Notice.GetNoticePage: ", "msg", res.Error.Error())
		return nil, res.Error
	}
	return items, nil
}

func GetNoticeById(ctx *gin.Context, id int) *Notice {
	var item *Notice
	res := db.Model(&Notice{}).Where("id = ? and is_delete = 0", id).First(&item)
	if res.Error != nil {
		logging.WithCtxError(ctx, "DB: Notice.GetNoticeById: ", "msg", res.Error.Error())
		return nil
	}
	return item
}
