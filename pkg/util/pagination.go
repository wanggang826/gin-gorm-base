package util

import (
	"gin-gorm-base/pkg/constants"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
	"math"
)

// GetAppPageParams GetPage get page parameters
func GetAppPageParams(c *gin.Context) (limit int, offset int) {
	limit = com.StrTo(c.Query("limit")).MustInt()
	offset = com.StrTo(c.Query("offset")).MustInt()
	if limit < 1 {
		limit = constants.DefaultPageLimit
	}

	return
}

// GetPageParams GetPage get page parameters
func GetPageParams(c *gin.Context) (limit int, offset int, pageNo int) {
	offset = 0
	limit = com.StrTo(c.Query("page_size")).MustInt()
	pageNo = com.StrTo(c.Query("page_no")).MustInt()
	if pageNo > 0 {
		offset = (pageNo - 1) * limit
	}
	if limit < 1 {
		limit = constants.DefaultPageLimit
	}

	return
}

func ResponsePageData(total int, data interface{}, pageSize int, pageNo int) map[string]interface{} {
	var responseMap = make(map[string]interface{})

	totalPage := int(math.Ceil(float64(total) / float64(pageSize)))

	responseMap["data"] = data
	responseMap["page_size"] = pageSize
	responseMap["page_no"] = pageNo
	responseMap["total_page"] = totalPage
	responseMap["total_count"] = total

	return responseMap
}
