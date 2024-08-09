package monitor

import (
	"context"
	"fmt"
	"mySparkler/backend/api/baseAPI"
	"mySparkler/backend/model/monitor"
	"mySparkler/backend/model/tools"
	"mySparkler/pkg/constants"
	"mySparkler/pkg/utils"
	"mySparkler/pkg/utils/R"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
)

type LoginInforAPI struct {
	ctx context.Context
	baseAPI.Base
}

var loginInforAPI *LoginInforAPI
var onceLoginInforAPI sync.Once

// NewApp creates a new App application struct
func NewLoginInforAPI() *LoginInforAPI {
	if loginInforAPI == nil {
		onceLoginInforAPI.Do(func() {
			loginInforAPI = &LoginInforAPI{}
		})
	}
	return loginInforAPI
}

func (g *LoginInforAPI) LoginInformListHandler(params map[string]interface{}) tools.TableDataInfo {
	rows, total := getListLoginLog(params)
	if rows == nil {
		return tools.Fail()
	} else {
		return tools.Success(rows, total)
	}
}

func getListLoginLog(params map[string]interface{}) ([]monitor.SysLogininfor, int64) {
	/*分页*/

	var param = tools.SearchTableDataParam{
		PageNum:  1,
		PageSize: 10,
		Other: monitor.SysLogininfor{
			Ipaddr:   "",
			UserName: "",
			Status:   "",
		},
		OrderByColumn: "",
		IsAsc:         "",
		Params: tools.Params{
			BeginTime: "",
			EndTime:   "",
		},
	}
	err := mapstructure.Decode(params, &param)
	if err != nil {
		fmt.Println(err.Error())
	}
	return monitor.SelectLogininforList(param)
}

func (g *LoginInforAPI) ExportHandler(params map[string]interface{}) {
	list, _ := getListLoginLog(params)
	//定义首行标题
	dataKey := make([]map[string]string, 0)
	dataKey = append(dataKey, map[string]string{
		"key":    "infoId",
		"title":  "序号",
		"width":  "10",
		"is_num": "0",
	})
	dataKey = append(dataKey, map[string]string{
		"key":    "userName",
		"title":  "用户账号",
		"width":  "10",
		"is_num": "0",
	})
	dataKey = append(dataKey, map[string]string{
		"key":    "status",
		"title":  "登录状态",
		"width":  "10",
		"is_num": "0",
	})
	dataKey = append(dataKey, map[string]string{
		"key":    "ipaddr",
		"title":  "登录地址",
		"width":  "20",
		"is_num": "0",
	})
	dataKey = append(dataKey, map[string]string{
		"key":    "loginLocation",
		"title":  "登录地点",
		"width":  "20",
		"is_num": "0",
	})
	dataKey = append(dataKey, map[string]string{
		"key":    "browser",
		"title":  "浏览器",
		"width":  "20",
		"is_num": "0",
	})
	dataKey = append(dataKey, map[string]string{
		"key":    "os",
		"title":  "操作系统",
		"width":  "10",
		"is_num": "0",
	})
	dataKey = append(dataKey, map[string]string{
		"key":    "msg",
		"title":  "提示消息",
		"width":  "30",
		"is_num": "0",
	})
	dataKey = append(dataKey, map[string]string{
		"key":    "loginTime",
		"title":  "访问时间",
		"width":  "50",
		"is_num": "0",
	})
	//填充数据
	data := make([]map[string]interface{}, 0)
	if len(list) > 0 {
		for _, v := range list {
			var status = v.Status
			var statusStr = ""
			if status == "0" {
				statusStr = "成功"
			} else {
				statusStr = "失败"
			}

			var loginTime = v.LoginTime.Format(constants.TimeFormat)
			data = append(data, map[string]interface{}{
				"infoId":        v.InfoId,
				"userName":      v.UserName,
				"status":        statusStr,
				"ipaddr":        v.Ipaddr,
				"loginLocation": v.LoginLocation,
				"browser":       v.Browser,
				"os":            v.Os,
				"msg":           v.Msg,
				"loginTime":     loginTime,
			})
		}
	}
	ex := tools.NewMyExcel()
	ex.ExportToWeb(dataKey, data, g.ctx)
}

func (g *LoginInforAPI) DeleteByIdHandler(operId string) R.Result {
	// var operId = context.Param("infoIds")
	var result = monitor.DelectLoginlog(utils.Split(operId))
	return result
}

func (g *LoginInforAPI) CleanHandler() R.Result {
	return monitor.ClearLoginlog()
}

func (g *LoginInforAPI) UnlockHandler(userName string) R.Result {
	// var userName = context.Param("userName")
	monitor.UnlockByUserName(userName)
	return R.ReturnSuccess(gin.H{
		"msg":  "操作成功",
		"code": http.StatusOK,
	})
}
