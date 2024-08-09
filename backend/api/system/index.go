package system

import (
	"mySparkler/pkg/utils/R"

	"github.com/gin-gonic/gin"
)

// 后台获取 首页数据

func IndexData() {
	//
}

// IndexHandler 测试代码
func IndexHandler(context *gin.Context) R.Result {
	return R.ReturnSuccess("Hello ruoyi go")
}
