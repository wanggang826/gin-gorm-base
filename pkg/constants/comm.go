package constants

// 公用常量
const (
	TimeLayout       = "2006-01-02 15:04:05"
	DateLayout       = "2006-01-02"
	DefaultPageLimit = 20

	RoleIdAdmin   = 1 //管理员
	RoleIdPartner = 2 //流量主

	DefaultPassword = "123456"

	StatusEnable  = 1 //启用
	StatusDisable = 2 //禁用

	IncomeStatusNoRelease = 1 //未发布
	IncomeStatusReleased  = 2 //已发布
	IncomeStatusFinish    = 3 //已完成

	AppSpreadStatusApply = 1 //1申请中
	IncomeStatusSuccess  = 2 //2推广中
	IncomeStatusFail     = 3 //3申请失败
	IncomeStatusDisable  = 4 //4已下架
)
