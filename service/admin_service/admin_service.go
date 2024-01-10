package admin_service

import (
	"gin-gorm-base/models"
	"gin-gorm-base/pkg/app"
	"gin-gorm-base/pkg/constants"
	"gin-gorm-base/pkg/e"
	"gin-gorm-base/pkg/logging"
	"gin-gorm-base/pkg/util"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
	"time"
)

type AdminLoginParams struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func AdminLogin(c *gin.Context) (map[string]interface{}, int) {
	var params AdminLoginParams
	err := c.BindJSON(&params)
	if err != nil {
		logging.Error("GetRequest: Login Gin BindJSON:", "msg", err.Error())
		return nil, e.ERROR_PARAM_GIN_BINDJSON_ERROR
	}

	valid := validation.Validation{}
	valid.Required(params.Username, "Username")
	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		return nil, e.ERROR_ADMIN_LOGIN_USERNAME_EMPTY
	}
	valid.Required(params.Password, "Password")
	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		return nil, e.INVALID_PARAMS
	}

	admin := models.GeAdminByUsername(c, params.Username)
	if admin == nil {
		return nil, e.ERROR_ADMIN_LOGIN_USERNAME_PASSWORD_ERROR
	}
	if admin.Status == constants.StatusDisable {
		return nil, e.ERROR_USER_BANNED
	}
	if admin.Password != params.Password {
		return nil, e.ERROR_ADMIN_LOGIN_USERNAME_PASSWORD_ERROR
	}
	editAdmin := &models.Admin{
		Ctx: c,
		Id:  admin.Id,
	}
	editData := make(map[string]interface{})
	editData["last_login_time"] = int(time.Now().Unix())
	editData["last_login_ip"] = c.ClientIP()
	_, err = editAdmin.Edit(editData)
	if err != nil {
		logging.WithCtxError(c, "AdminLogin Edit fail", "err", err)
	}
	token, tErr := util.GenerateAdminToken(admin.Id, admin.RoleID)
	if tErr != nil {
		logging.WithCtxError(c, "AdminLogin GenerateAdminToken fail", "err", tErr)
		return nil, e.ERROR_AUTH_TOKEN
	}
	result := make(map[string]interface{})
	result["token"] = token
	result["adminInfo"] = HandleAdmin(admin)
	return result, e.SUCCESS
}

type AdminRes struct {
	Id         int    `json:"id"`
	Username   string `json:"username"`
	RoleId     int    `json:"role_id"`
	Name       string `json:"name"`
	Nickname   string `json:"nickname"`
	Mobile     string `json:"mobile"`
	Status     int    `json:"status"`
	CreateTime string `json:"create_time"`
}

func HandleAdmin(dbAdmin *models.Admin) *AdminRes {
	if dbAdmin == nil {
		return nil
	}
	resp := &AdminRes{
		Id:         dbAdmin.Id,
		Username:   dbAdmin.Username,
		RoleId:     dbAdmin.RoleID,
		Name:       dbAdmin.Name,
		Nickname:   dbAdmin.Nickname,
		Mobile:     dbAdmin.Mobile,
		Status:     dbAdmin.Status,
		CreateTime: util.TimeToStrTime(dbAdmin.CreatedAt),
	}
	return resp
}

func GetAdminPage(c *gin.Context) map[string]interface{} {
	limit, offset, pageNo := util.GetPageParams(c)
	nickname := c.Query("nickname")
	status := com.StrTo(c.Query("status")).MustInt()
	total, _ := models.GetAdminCount(c, nickname, status)
	list, _ := models.GetAdminPage(c, nickname, status, limit, offset)
	retList := make([]*AdminRes, len(list))
	for k, v := range list {
		retList[k] = HandleAdmin(v)
	}
	return util.ResponsePageData(total, retList, limit, pageNo)
}

func GetAdminInfo(c *gin.Context) *AdminRes {
	info := models.GetAdminById(c, c.GetInt("uid"))
	return HandleAdmin(info)
}

func SearchAdmin(c *gin.Context) []*AdminRes {
	nickname := c.Query("nickname")
	list, _ := models.GetAdminList(c, nickname, constants.StatusEnable, constants.RoleIdPartner)
	retList := make([]*AdminRes, len(list))
	for k, v := range list {
		retList[k] = HandleAdmin(v)
	}
	return retList
}

type SaveAdminParams struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	RoleId   int    `json:"role_id"`
	Name     string `json:"name"`
	Nickname string `json:"nickname"`
	Mobile   string `json:"mobile"`
}

func AddAdmin(c *gin.Context) (int, int) {
	roleId := c.GetInt("roleId")
	if roleId != constants.RoleIdAdmin {
		return 0, e.ERROR_POWER_NOT_ENOUGH
	}
	var params SaveAdminParams
	err := c.BindJSON(&params)
	if err != nil {
		return 0, e.ERROR_PARAM_GIN_BINDJSON_ERROR
	}
	if params.Username == "" || params.Password == "" {
		return 0, e.ERROR_COMMON_URL_PARAM_ERROR
	}
	admin := &models.Admin{
		Ctx:      c,
		Username: params.Username,
		Password: params.Password,
		RoleID:   params.RoleId,
		Name:     params.Name,
		Nickname: params.Nickname,
		Mobile:   params.Mobile,
		Status:   constants.StatusEnable,
	}
	_, err = admin.Add()
	if err != nil {
		logging.WithCtxError(c, "AddAdmin Add fail", "err", err)
		return 0, e.ERROR_DO_ERROR
	}
	return admin.Id, e.SUCCESS
}

func EditAdmin(c *gin.Context) (int, int) {
	roleId := c.GetInt("roleId")
	if roleId != constants.RoleIdAdmin {
		return 0, e.ERROR_POWER_NOT_ENOUGH
	}
	var params SaveAdminParams
	err := c.BindJSON(&params)
	if err != nil {
		return 0, e.ERROR_PARAM_GIN_BINDJSON_ERROR
	}
	if params.Id < 2 {
		return 0, e.ERROR_COMMON_URL_PARAM_ERROR
	}

	oldAdmin := models.GetAdminById(c, params.Id)
	if oldAdmin == nil {
		return 0, e.ERROR_COMMON_URL_PARAM_ERROR
	}

	editAdmin := &models.Admin{
		Ctx: c,
		Id:  oldAdmin.Id,
	}
	editData := make(map[string]interface{})
	editData["role_id"] = params.RoleId
	editData["name"] = params.Name
	editData["nickname"] = params.Nickname
	editData["mobile"] = params.Mobile
	_, err = editAdmin.Edit(editData)
	if err != nil {
		logging.WithCtxError(c, "EditAdmin Edit fail", "err", err)
		return 0, e.ERROR_DO_ERROR
	}
	return oldAdmin.Id, e.SUCCESS
}

type DelAdminParams struct {
	Id int `json:"id"`
}

func DelAdmin(c *gin.Context) int {
	roleId := c.GetInt("roleId")
	if roleId != constants.RoleIdAdmin {
		return e.ERROR_POWER_NOT_ENOUGH
	}
	var params DelAdminParams
	err := c.BindJSON(&params)
	if err != nil {
		return e.ERROR_PARAM_GIN_BINDJSON_ERROR
	}
	if params.Id < 2 {
		return e.ERROR_COMMON_URL_PARAM_ERROR
	}
	oldAdmin := models.GetAdminById(c, params.Id)
	if oldAdmin == nil {
		return e.ERROR_COMMON_URL_PARAM_ERROR
	}

	editAdmin := &models.Admin{
		Ctx: c,
		Id:  oldAdmin.Id,
	}
	editData := make(map[string]interface{})
	editData["is_delete"] = 1
	_, err = editAdmin.Edit(editData)
	if err != nil {
		logging.WithCtxError(c, "DelAdmin Edit fail", "err", err)
		return e.ERROR_DO_ERROR
	}
	return e.SUCCESS
}

type ChangeAdminStatusParams struct {
	Id     int `json:"id"`
	Status int `json:"status"`
}

func ChangeAdminStatus(c *gin.Context) int {
	roleId := c.GetInt("roleId")
	if roleId != constants.RoleIdAdmin {
		return e.ERROR_POWER_NOT_ENOUGH
	}
	var params ChangeAdminStatusParams
	err := c.BindJSON(&params)
	if err != nil {
		return e.ERROR_PARAM_GIN_BINDJSON_ERROR
	}
	if params.Id < 2 {
		return e.ERROR_COMMON_URL_PARAM_ERROR
	}
	oldAdmin := models.GetAdminById(c, params.Id)
	if oldAdmin == nil {
		return e.ERROR_COMMON_URL_PARAM_ERROR
	}

	editAdmin := &models.Admin{
		Ctx: c,
		Id:  oldAdmin.Id,
	}
	editData := make(map[string]interface{})
	editData["status"] = params.Status
	_, err = editAdmin.Edit(editData)
	if err != nil {
		logging.WithCtxError(c, "ChangeAdminStatus Edit fail", "err", err)
		return e.ERROR_DO_ERROR
	}
	return e.SUCCESS
}

type EditAdminPwdParams struct {
	Id          int    `json:"id"`
	Password    string `json:"password"`
	NewPwd      string `json:"new_pwd"`
	AgainNewPwd string `json:"again_new_pwd"`
}

func EditAdminPwd(c *gin.Context) int {
	var params EditAdminPwdParams
	err := c.BindJSON(&params)
	if err != nil {
		return e.ERROR_PARAM_GIN_BINDJSON_ERROR
	}
	if params.Id < 2 {
		return e.ERROR_COMMON_URL_PARAM_ERROR
	}
	if params.Password == "" || params.NewPwd == "" || params.AgainNewPwd == "" {
		return e.ERROR_COMMON_URL_PARAM_ERROR
	}

	if params.NewPwd != params.AgainNewPwd {
		return e.ERROR_ADMIN_CHANGE_PASSWORD_DIFF_ERROR
	}
	admin := models.GetAdminById(c, params.Id)
	if params.Password != admin.Password {
		return e.ERROR_ADMIN_LOGIN_USERNAME_PASSWORD_ERROR
	}
	editAdmin := &models.Admin{
		Ctx: c,
		Id:  admin.Id,
	}
	editData := make(map[string]interface{})
	editData["password"] = params.NewPwd
	_, err = editAdmin.Edit(editData)
	if err != nil {
		logging.WithCtxError(c, "EditAdminPwd Edit fail", "err", err)
		return e.ERROR_DO_ERROR
	}
	return e.SUCCESS
}

func ResetAdminPwd(c *gin.Context) int {
	roleId := c.GetInt("roleId")
	if roleId != constants.RoleIdAdmin {
		return e.ERROR_POWER_NOT_ENOUGH
	}
	var params EditAdminPwdParams
	err := c.BindJSON(&params)
	if err != nil {
		return e.ERROR_PARAM_GIN_BINDJSON_ERROR
	}
	if params.Id < 1 {
		return e.ERROR_COMMON_URL_PARAM_ERROR
	}
	admin := models.GetAdminById(c, params.Id)
	editAdmin := &models.Admin{
		Ctx: c,
		Id:  admin.Id,
	}
	newPassword := util.Md5String(constants.DefaultPassword)
	editData := make(map[string]interface{})
	editData["password"] = newPassword
	_, err = editAdmin.Edit(editData)
	if err != nil {
		logging.WithCtxError(c, "EditAdminPwd Edit fail", "err", err)
		return e.ERROR_DO_ERROR
	}
	return e.SUCCESS
}
