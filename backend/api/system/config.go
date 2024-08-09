package system

import (
	"context"
	"fmt"
	"mySparkler/backend/api/baseAPI"
	"mySparkler/backend/model/system"
	"mySparkler/backend/model/tools"
	"mySparkler/pkg/utils"
	"mySparkler/pkg/utils/R"
	"mySparkler/pkg/utils/jwt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
)

type ConfigAPI struct {
	ctx context.Context
	baseAPI.Base
}

var configAPI *ConfigAPI
var onceConfig sync.Once

// NewApp creates a new App application struct
func NewConfigAPI() *ConfigAPI {
	if configAPI == nil {
		onceConfig.Do(func() {
			configAPI = &ConfigAPI{}
		})
	}
	return configAPI
}

func (a *ConfigAPI) ListConfig(params map[string]interface{}) tools.TableDataInfo {

	var param = tools.SearchTableDataParam{
		Other: system.SysConfig{
		},
		Params: tools.Params{
		},
	}

	err := mapstructure.Decode(params, &param)
	if err != nil {
		fmt.Println(err.Error())
	}
	var result = system.SelectConfigList(param, true)
	return result
}

func (a *ConfigAPI) ExportConfig(context *gin.Context) R.Result {

	var configParam system.SysConfig
	if err := context.ShouldBind(&configParam); err != nil {
		return R.ReturnFailMsg("参数不能为空")
	}
	var param = tools.SearchTableDataParam{
		PageNum:  0,
		PageSize: 10,
		Other: system.SysConfig{
			ConfigName: configParam.ConfigName,
			ConfigKey:  configParam.ConfigKey,
			ConfigType: configParam.ConfigType,
		},
	}
	tab := system.SelectConfigList(param, false)
	list := tab.Rows.([]system.SysConfig)
	//定义首行标题
	dataKey := make([]map[string]string, 0)
	dataKey = append(dataKey, map[string]string{
		"key":    "configId",
		"title":  "参数主键",
		"width":  "10",
		"is_num": "0",
	})
	dataKey = append(dataKey, map[string]string{
		"key":    "configName",
		"title":  "参数名称",
		"width":  "10",
		"is_num": "0",
	})
	dataKey = append(dataKey, map[string]string{
		"key":    "configKey",
		"title":  "参数键名",
		"width":  "10",
		"is_num": "0",
	})
	dataKey = append(dataKey, map[string]string{
		"key":    "configValue",
		"title":  "参数键值",
		"width":  "10",
		"is_num": "0",
	})
	dataKey = append(dataKey, map[string]string{
		"key":    "configType",
		"title":  "系统内置",
		"width":  "10",
		"is_num": "0",
	})

	//填充数据
	data := make([]map[string]interface{}, 0)
	if len(list) > 0 {
		for _, v := range list {
			configType := v.ConfigType
			var configTypeStr = ""
			if configType == "Y" {
				configTypeStr = "是"
			} else if configType == "N" {
				configTypeStr = "否"
			}
			data = append(data, map[string]interface{}{
				"configId":    v.ConfigId,
				"configName":  v.ConfigName,
				"configKey":   v.ConfigKey,
				"configValue": v.ConfigValue,
				"configType":  configTypeStr,
			})
		}
	}
	ex := tools.NewMyExcel()
	ex.ExportToWeb(dataKey, data, context)
	return R.ReturnSuccess(gin.H{
		"msg":  "操作成功",
		"code": http.StatusOK,
	})
}

func (a *ConfigAPI) GetConfigInfo(configId int) R.Result {
	// configId := context.Param("configId")
	result := system.GetConfigInfo(configId)
	return R.ReturnSuccess(gin.H{
		"msg":  "操作成功",
		"code": http.StatusOK,
		"data": result,
	})
}

func (a *ConfigAPI) GetConfigKey(configKey string) R.Result {
	// configKey := context.Param("configKey")
	var config = system.SysConfig{ConfigKey: configKey}
	var result = system.SelectConfig(config)
	return R.ReturnSuccess(gin.H{
		"msg": result.ConfigValue,
	})
}

func (a *ConfigAPI) SaveConfig(params map[string]interface{}) R.Result {

	userId := jwt.CacheGetUserId()
	var configParam = system.SysConfig{}
	if utils.CheckTypeByReflectNil(params) {
		return R.ReturnFailMsg("参数不能为空")
	}
	err := mapstructure.Decode(params, &configParam)
	if err != nil {
		fmt.Println(err.Error())
	}
	user := system.FindUserById(userId)
	configParam.CreateBy = user.UserName
	configParam.CreateTime = time.Now()
	result := system.SaveConfig(configParam)
	return result
}

func (a *ConfigAPI) UploadConfig(params map[string]interface{}) R.Result {

	userId := jwt.CacheGetUserId()
	var configParam system.SysConfig
	if utils.CheckTypeByReflectNil(params) {
		return R.ReturnFailMsg("参数不能为空")
	}
	err := mapstructure.Decode(params, &configParam)
	if err != nil {
		fmt.Println(err.Error())
	}
	user := system.FindUserById(userId)
	configParam.UpdateBy = user.UserName
	configParam.UpdateTime = time.Now()
	result := system.EditConfig(configParam)
	return result
}

func (a *ConfigAPI) DetectConfig(configIds string) R.Result {

	userId := jwt.CacheGetUserId()
	println(userId)
	// var configIds = context.Param("configIds")
	result := system.DelConfig(configIds)
	return result
}

func (a *ConfigAPI) DeleteCacheConfig(refreshCache string) R.Result {

	userId := jwt.CacheGetUserId()
	println(userId)
	// var refreshCache = context.Param("refreshCache")
	result := system.DelCacheConfig(refreshCache)
	return result
}
