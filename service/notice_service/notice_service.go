package notice_service

import (
	"gin-gorm-base/models"
	"gin-gorm-base/pkg/e"
	"gin-gorm-base/pkg/logging"
	"gin-gorm-base/pkg/util"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
)

type NoticeRes struct {
	Id         int    `json:"id"`
	Title      string `json:"title"`   // 标题
	Image      string `json:"image"`   // 图片
	CateID     int    `json:"cate_id"` // 种类 1功能更新 2活动通知
	UpdateTime string `json:"update_time"`
}

func HandleNotice(dbNotice *models.Notice) *NoticeRes {
	if dbNotice == nil {
		return nil
	}
	resp := &NoticeRes{
		Id:         dbNotice.Id,
		Title:      dbNotice.Title,
		Image:      dbNotice.Image,
		CateID:     dbNotice.CateID,
		UpdateTime: util.TimeToStrTime(dbNotice.UpdatedAt),
	}
	return resp
}

func GetNoticePage(c *gin.Context) map[string]interface{} {
	limit, offset, pageNo := util.GetPageParams(c)
	title := c.Query("title")
	cateId := com.StrTo(c.Query("cate_id")).MustInt()
	total, _ := models.GetNoticeCount(c, title, cateId)
	list, _ := models.GetNoticePage(c, title, cateId, limit, offset)
	retList := make([]*NoticeRes, len(list))
	for k, v := range list {
		retList[k] = HandleNotice(v)
	}
	return util.ResponsePageData(total, retList, limit, pageNo)
}

type SaveNoticeParams struct {
	Id     int    `json:"id"`
	Title  string `json:"title"`   // 标题
	Image  string `json:"image"`   // 图片
	CateID int    `json:"cate_id"` // 种类 1功能更新 2活动通知
	Sort   int    `json:"sort"`    // 排序
}

func AddNotice(c *gin.Context) (int, int) {
	var params SaveNoticeParams
	err := c.BindJSON(&params)
	if err != nil {
		return 0, e.ERROR_PARAM_GIN_BINDJSON_ERROR
	}
	if params.Title == "" {
		return 0, e.ERROR_COMMON_URL_PARAM_ERROR
	}
	Notice := &models.Notice{
		Ctx:    c,
		Title:  params.Title,
		Image:  params.Image,
		CateID: params.CateID,
		Sort:   params.Sort,
	}
	_, err = Notice.Add()
	if err != nil {
		logging.WithCtxError(c, "AddNotice Add fail", "err", err)
		return 0, e.ERROR_DO_ERROR
	}
	return Notice.Id, e.SUCCESS
}

func EditNotice(c *gin.Context) (int, int) {
	var params SaveNoticeParams
	err := c.BindJSON(&params)
	if err != nil {
		return 0, e.ERROR_PARAM_GIN_BINDJSON_ERROR
	}
	if params.Id < 1 {
		return 0, e.ERROR_COMMON_URL_PARAM_ERROR
	}
	if params.Title == "" {
		return 0, e.ERROR_COMMON_URL_PARAM_ERROR
	}

	oldNotice := models.GetNoticeById(c, params.Id)
	if oldNotice == nil {
		return 0, e.ERROR_COMMON_URL_PARAM_ERROR
	}

	editNotice := &models.Notice{
		Id: oldNotice.Id,
	}
	editData := make(map[string]interface{})
	editData["title"] = params.Title
	editData["image"] = params.Image
	editData["cate_id"] = params.CateID
	editData["sort"] = params.Sort
	_, err = editNotice.Edit(editData)
	if err != nil {
		logging.WithCtxError(c, "EditNotice Edit fail", "err", err)
		return 0, e.ERROR_DO_ERROR
	}
	return oldNotice.Id, e.SUCCESS
}

type DelNoticeParams struct {
	Id int `json:"id"`
}

func DelNotice(c *gin.Context) int {
	var params DelNoticeParams
	err := c.BindJSON(&params)
	if err != nil {
		return e.ERROR_PARAM_GIN_BINDJSON_ERROR
	}
	if params.Id < 1 {
		return e.ERROR_COMMON_URL_PARAM_ERROR
	}
	oldNotice := models.GetNoticeById(c, params.Id)
	if oldNotice == nil {
		return e.ERROR_COMMON_URL_PARAM_ERROR
	}

	editNotice := &models.Notice{
		Ctx: c,
		Id:  oldNotice.Id,
	}
	editData := make(map[string]interface{})
	editData["is_delete"] = 1
	_, err = editNotice.Edit(editData)
	if err != nil {
		logging.WithCtxError(c, "DelNotice Edit fail", "err", err)
		return e.ERROR_DO_ERROR
	}
	return e.SUCCESS
}
