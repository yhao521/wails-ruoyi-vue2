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

type DictDataAPI struct {
	ctx context.Context
	baseAPI.Base
}

var dictDataAPI *DictDataAPI
var onceDictData sync.Once

// NewApp creates a new App application struct
func NewDictDataAPI() *DictDataAPI {
	if dictDataAPI == nil {
		onceDictData.Do(func() {
			dictDataAPI = &DictDataAPI{}
		})
	}
	return dictDataAPI
}

func (s *DictDataAPI) ListDict(params map[string]interface{}) R.Result {
	if utils.CheckTypeByReflectNil(params) {
		fmt.Printf("params: %v\n", params)
		return R.ReturnFailMsg("参数不能为空")
	}
	var param = tools.SearchTableDataParam{
		Other: system.SysDictData{},
	}
	err := mapstructure.Decode(params, &param)
	if err != nil {
		fmt.Println(err.Error())
	}

	result := system.SelectDictDataList(param, false)
	return R.ReturnSuccess(gin.H{
		"rows":  result.Rows,
		"total": result.Total,
	})
}

func (s *DictDataAPI) ExportDict(params map[string]interface{}, ctx context.Context) {
	if utils.CheckTypeByReflectNil(params) {
		fmt.Printf("params: %v\n", params)
		// return R.ReturnFailMsg("参数不能为空")
	}

	var param = tools.SearchTableDataParam{
		Other: system.SysDictData{},
	}
	err := mapstructure.Decode(params, &param)
	if err != nil {
		fmt.Println(err.Error())
	}

	result := system.SelectDictDataList(param, false)
	var list = result.Rows.([]system.SysDictData)
	//定义首行标题
	dataKey := make([]map[string]string, 0)
	dataKey = append(dataKey, map[string]string{
		"key":    "dictCode",
		"title":  "字典编码",
		"width":  "10",
		"is_num": "0",
	})
	dataKey = append(dataKey, map[string]string{
		"key":    "dictSort",
		"title":  "字典排序",
		"width":  "10",
		"is_num": "0",
	})
	dataKey = append(dataKey, map[string]string{
		"key":    "dictLabel",
		"title":  "字典标签",
		"width":  "10",
		"is_num": "0",
	})
	dataKey = append(dataKey, map[string]string{
		"key":    "dictValue",
		"title":  "字典键值",
		"width":  "10",
		"is_num": "0",
	})
	dataKey = append(dataKey, map[string]string{
		"key":    "dictType",
		"title":  "字典类型",
		"width":  "10",
		"is_num": "0",
	})
	dataKey = append(dataKey, map[string]string{
		"key":    "isDefault",
		"title":  "是否默认",
		"width":  "10",
		"is_num": "0",
	})
	dataKey = append(dataKey, map[string]string{
		"key":    "status",
		"title":  "状态",
		"width":  "10",
		"is_num": "0",
	})
	//填充数据
	data := make([]map[string]interface{}, 0)
	if len(list) > 0 {
		for _, v := range list {
			defaults := v.IsDefault
			var df = ""
			if defaults == "Y" {
				df = "是"
			}
			if defaults == "N" {
				df = "否"
			}
			var status = v.Status
			statusStr := ""
			if status == "0" {
				statusStr = "正常"
			}
			if status == "1" {
				statusStr = "停用"
			}
			data = append(data, map[string]interface{}{
				"dictCode":  v.DictCode,
				"dictSort":  v.DictSort,
				"dictLabel": v.DictLabel,
				"dictValue": v.DictValue,
				"dictType":  v.DictType,
				"isDefault": df,
				"status":    statusStr,
			})
		}
	}
	ex := tools.NewMyExcel()
	ex.ExportToWeb(dataKey, data, ctx)
}

func (s *DictDataAPI) GetDictCode(dictCode int) R.Result {
	// dictCode := context.Param("dictCode")
	result := system.FindDictCodeById(dictCode)
	return R.ReturnSuccess(gin.H{
		"msg":  "操作成功",
		"code": http.StatusOK,
		"data": result,
	})
}

func (s *DictDataAPI) DictTypeHandler(dictType string) R.Result {
	// dictType := context.Param("dictType")
	result := system.SelectDictDataByType(dictType)
	return R.ReturnSuccess(gin.H{
		"msg":  "操作成功",
		"code": http.StatusOK,
		"data": result,
	})
}

func (s *DictDataAPI) SaveDictData(param map[string]interface{}) R.Result {

	userId := jwt.CacheGetUserId()
	var dictDataParam system.SysDictData
	// if utils.CheckTypeByReflectNil(dictDataParam) {
	// 	return R.ReturnFailMsg("参数不能为空")
	// }
	err := mapstructure.Decode(param, &dictDataParam)
	if err != nil {
		fmt.Println(err.Error())
	}
	user := system.FindUserById(userId)
	dictDataParam.CreateBy = user.UserName
	dictDataParam.CreateTime = time.Now()
	dictDataParam.UpdateTime = time.Now()
	result := system.SaveDictData(dictDataParam)
	return result
}

func (s *DictDataAPI) UpDictData(param map[string]interface{}) R.Result {

	userId := jwt.CacheGetUserId()
	var dictDataParam system.SysDictData
	if utils.CheckTypeByReflectNil(dictDataParam) {
		return R.ReturnFailMsg("参数不能为空")

	}
	err := mapstructure.Decode(param, &dictDataParam)
	if err != nil {
		fmt.Println(err.Error())
	}
	user := system.FindUserById(userId)
	dictDataParam.UpdateBy = user.UserName
	dictDataParam.UpdateTime = time.Now()
	result := system.EditDictData(dictDataParam)
	return result
}

func (s *DictDataAPI) DeleteDictData(dictCodes string) R.Result {
	// var dictCodes = context.Param("dictCodes")
	result := system.DeleteDictData(dictCodes)
	return result
}

// ListDictType ---------------------------------------------
func (s *DictDataAPI) ListDictType(params map[string]interface{}) R.Result {
	if utils.CheckTypeByReflectNil(params) {
		return R.ReturnFailMsg("参数不能为空")
	}
	var param = tools.SearchTableDataParam{
		Other: system.SysDictType{},
		// Params:   params.Params,
	}
	err := mapstructure.Decode(params, &param)
	if err != nil {
		fmt.Println(err.Error())
	}

	var result = system.SelectSysDictTypeList(param, true)
	// context.JSON(http.StatusOK, result)
	return R.ReturnSuccess(gin.H{
		"msg":  "操作成功",
		"code": http.StatusOK,
		"data": result,
	})
}

func (s *DictDataAPI) ExportType(params map[string]interface{}, ctx context.Context) {
	var param = tools.SearchTableDataParam{
		Other: system.SysDictType{},
		// Params:   params.Params,
	}
	err := mapstructure.Decode(params, &param)
	if err != nil {
		fmt.Println(err.Error())
	}
	var result = system.SelectSysDictTypeList(param, false)

	var list = result.Rows.([]system.SysDictType)
	//定义首行标题
	dataKey := make([]map[string]string, 0)
	dataKey = append(dataKey, map[string]string{
		"key":    "dictId",
		"title":  "字典主键",
		"width":  "10",
		"is_num": "0",
	})
	dataKey = append(dataKey, map[string]string{
		"key":    "dictName",
		"title":  "字典名称",
		"width":  "10",
		"is_num": "0",
	})
	dataKey = append(dataKey, map[string]string{
		"key":    "dictType",
		"title":  "字典类型",
		"width":  "10",
		"is_num": "0",
	})
	dataKey = append(dataKey, map[string]string{
		"key":    "status",
		"title":  "状态",
		"width":  "10",
		"is_num": "0",
	})
	//填充数据
	data := make([]map[string]interface{}, 0)
	if len(list) > 0 {
		for _, v := range list {
			status := v.Status
			var statusStr = ""
			if status == "0" {
				statusStr = "正常"
			}
			if status == "1" {
				statusStr = "停用"
			}
			data = append(data, map[string]interface{}{
				"dictId":   v.DictId,
				"dictName": v.DictName,
				"dictType": v.DictType,
				"status":   statusStr,
			})
		}
	}
	ex := tools.NewMyExcel()
	ex.ExportToWeb(dataKey, data, ctx)
}

func (s *DictDataAPI) GetTypeDict(dictId string) R.Result {
	// dictId := context.Param("dictId")
	result := system.FindTypeDictById(utils.GetInterfaceToInt(dictId))
	return R.ReturnSuccess(gin.H{
		"msg":  "操作成功",
		"code": http.StatusOK,
		"data": result,
	})
}

func (s *DictDataAPI) SaveType(param map[string]interface{}) R.Result {
	var dictTypeParam system.SysDictType

	if utils.CheckTypeByReflectNil(param) {
		return R.ReturnFailMsg("参数不能为空")
	}
	fmt.Println("SaveType", param)
	err := mapstructure.Decode(param, &dictTypeParam)
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println("SaveType", dictTypeParam)
	result := system.SaveType(dictTypeParam)
	// context.JSON(http.StatusOK, result)
	return result
}

func (s *DictDataAPI) UpdateType(param map[string]interface{}) R.Result {
	var dictTypeParam system.SysDictType
	if utils.CheckTypeByReflectNil(param) {
		return R.ReturnFailMsg("参数不能为空")
	}
	err := mapstructure.Decode(param, &dictTypeParam)
	if err != nil {
		fmt.Println(err.Error())
	}
	result := system.UploadType(dictTypeParam)
	return result
}

func (s *DictDataAPI) DeleteDataType(dictIds string) R.Result {
	// dictIds := context.Param("dictIds")
	result := system.DeleteDataType(dictIds)
	return result
}

func (s *DictDataAPI) RefreshCache() R.Result {
	result := system.RefreshCache()
	return result
}

func (s *DictDataAPI) GetOptionSelect() R.Result {
	result := system.GetOptionSelect()
	return result
}
