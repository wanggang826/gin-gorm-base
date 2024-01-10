package e

var MsgFlags = map[int]string{
	SUCCESS:        "ok",
	ERROR:          "系统异常",
	INVALID_PARAMS: "请求参数错误",

	ERROR_COMMON_URL_PARAM_ERROR:   "参数错误",
	ERROR_PARAM_GIN_BINDJSON_ERROR: "参数类型错误",
	ERROR_REDIS_ERROR:              "缓存系统异常",
	ERROR_DB_ERROR:                 "数据库错误",
	ERROR_DO_ERROR:                 "操作失败",

	ERROR_LOGIN_FAIL:                          "登录失败",
	ERROR_LOGIN_NONE:                          "未登录",
	ERROR_LOGIN_PC_SCAN:                       "扫码登录失败",
	ERROR_ADMIN_LOGIN_USERNAME_EMPTY:          "用户名不能为空",
	ERROR_ADMIN_LOGIN_USERNAME_PASSWORD_ERROR: "用户名密码错误",
	ERROR_ADMIN_CHANGE_PASSWORD_DIFF_ERROR:    "两次密码不一致",
	ERROR_POWER_NOT_ENOUGH:                    "权限错误",
	ERROR_AUTH_CHECK_TOKEN_FAIL:               "Token鉴权失败",
	ERROR_AUTH_CHECK_TOKEN_TIMEOUT:            "Token已超时",
	ERROR_AUTH_TOKEN:                          "Token生成失败",
	ERROR_AUTH:                                "Token错误",

	ERROR_USER_GETINFO_FAIL:      "用户信息获取失败",
	ERROR_USER_FOLLOW_PARAMS_ERR: "关注参数错误",
	ERROR_USER_HAS_FOLLOWED:      "已关注，不能重复关注",
	ERROR_USER_BANNED:            "账号被禁",
	ERROR_USER_NOT_EXIST:         "用户不存在",
	ERROR_USER_IS_BLACKED:        "用户已在黑名单中",

	ERROR_PAY_CREATE_ORDER: "下单失败，请稍后重试或联系客服",
	ERROR_PAY_REFUND:       "退款失败，请稍后重试或联系客服",

	ERROR_INCOME_STATUS_FINISH:    "收益已确认不可撤回",
	ERROR_INCOME_STATUS_NORELEASE: "仅未发布状态收益可修改",

	ERROR_APP_SPREADER_EXIST: "该流量主已申请推广该应用，不可重复",
}

// GetMsg get error information based on Code
func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}
	return MsgFlags[ERROR]
}
