package system

import (
	"context"
	"fmt"
	"log"
	"mySparkler/backend/api/baseAPI"
	"mySparkler/backend/model/system"
	"mySparkler/backend/model/tools"
	"mySparkler/pkg/utils"
	"mySparkler/pkg/utils/R"
	"mySparkler/pkg/utils/jwt"
	"sync"
	"time"

	"github.com/mitchellh/mapstructure"
)

type PostAPI struct {
	ctx context.Context
	baseAPI.Base
}

var postAPI *PostAPI
var oncePostAPI sync.Once

// NewApp creates a new App application struct
func NewPostAPI() *PostAPI {
	if postAPI == nil {
		oncePostAPI.Do(func() {
			postAPI = &PostAPI{}
		})
	}
	return postAPI
}

func (a *PostAPI) ListPost(params map[string]interface{}) tools.TableDataInfo {

	var param = tools.SearchTableDataParam{
		PageNum:  1,
		PageSize: 10,
		Other: system.SysPost{
			PostName: "",
			PostCode: "",
			Status:   "",
		},
		Params: tools.Params{
			BeginTime: "",
			EndTime:   "",
		},
	}
	log.Default().Println("ListDept-param:", param)
	err := mapstructure.Decode(params, &param)
	if err != nil {
		fmt.Println(err.Error())
	}
	var result = system.SelectSysPostList(param, true)
	return result
}
func (a *PostAPI) ExportPost(params map[string]interface{}) {

	var param = tools.SearchTableDataParam{
		PageNum:  1,
		PageSize: 10,
		Other: system.SysPost{
			PostName: "",
			PostCode: "",
			Status:   "",
		},
		Params: tools.Params{
			BeginTime: "",
			EndTime:   "",
		},
	}
	log.Default().Println("ListDept-param:", param)
	err := mapstructure.Decode(params, &param)
	if err != nil {
		fmt.Println(err.Error())
	}
	var result = system.SelectSysPostList(param, false)
	var list = result.Rows.([]system.SysPost)
	//定义首行标题
	dataKey := make([]map[string]string, 0)
	dataKey = append(dataKey, map[string]string{
		"key":    "postId",
		"title":  "岗位序号",
		"width":  "10",
		"is_num": "0",
	})
	dataKey = append(dataKey, map[string]string{
		"key":    "postCode",
		"title":  "岗位编码",
		"width":  "10",
		"is_num": "0",
	})
	dataKey = append(dataKey, map[string]string{
		"key":    "postName",
		"title":  "岗位名称",
		"width":  "10",
		"is_num": "0",
	})
	dataKey = append(dataKey, map[string]string{
		"key":    "postSort",
		"title":  "岗位排序",
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
			var status = v.Status
			var statusStr = ""
			if status == "0" {
				statusStr = "正常"
			} else {
				statusStr = "停用"
			}
			data = append(data, map[string]interface{}{
				"postId":   v.PostId,
				"postCode": v.PostCode,
				"postName": v.PostName,
				"postSort": v.PostSort,
				"status":   statusStr,
			})
		}
	}
	ex := tools.NewMyExcel()
	ex.ExportToWeb(dataKey, data, a.ctx)
}
func (a *PostAPI) GetPostInfo(postId int) R.Result {
	// var postId = context.Param("postId")
	result := system.FindPostInfoById(postId)
	return R.ReturnSuccess(result)
}

func (a *PostAPI) SavePost(params map[string]interface{}) R.Result {

	userId := jwt.CacheGetUserId()
	var postParam = system.SysPost{}
	if utils.CheckTypeByReflectNil(params) {
		return R.ReturnFailMsg("参数不能为空")
	}
	err := mapstructure.Decode(params, &postParam)
	if err != nil {
		fmt.Println(err.Error())
	}
	user := system.FindUserById(userId)
	postParam.CreateBy = user.UserName
	postParam.CreateTime = time.Now()
	result := system.SavePost(postParam)
	return result
}

func (a *PostAPI) UploadPost(params map[string]interface{}) R.Result {

	userId := jwt.CacheGetUserId()
	var postParam = system.SysPost{}
	if utils.CheckTypeByReflectNil(params) {
		return R.ReturnFailMsg("参数不能为空")
	}
	err := mapstructure.Decode(params, &postParam)
	if err != nil {
		fmt.Println(err.Error())
	}
	user := system.FindUserById(userId)
	postParam.UpdateBy = user.UserName
	postParam.UpdateTime = time.Now()
	result := system.EditPost(postParam)
	return result
}

func (a *PostAPI) DeletePost(postIds string) R.Result {
	// var postIds = context.Param("postIds")
	result := system.DeletePost(postIds)
	return result
}

func (a *PostAPI) GetPostOptionSelect(params map[string]interface{}) R.Result {
	var param = tools.SearchTableDataParam{}
	err := mapstructure.Decode(params, &param)
	if err != nil {
		fmt.Println(err.Error())
	}
	var result = system.SelectSysPostList(param, false)
	return R.ReturnSuccess(result.Rows)
}
