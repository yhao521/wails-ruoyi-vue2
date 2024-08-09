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
	"strconv"
	"sync"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/wxnacy/wgo/arrays"
)

type DeptAPI struct {
	ctx context.Context
	baseAPI.Base
}

var deptAPI *DeptAPI
var onceDept sync.Once

// NewApp creates a new App application struct
func NewDeptAPI() *DeptAPI {
	if deptAPI == nil {
		onceDept.Do(func() {
			deptAPI = &DeptAPI{}
		})
	}
	return deptAPI
}

func (a *DeptAPI) ListDept(params map[string]interface{}) R.Result {
	log.Default().Println("ListDept-params:", params)

	var param = tools.SearchTableDataParam{
		PageNum:  1,
		PageSize: 10,
		Other: system.SysDept{
			DeptName: "",
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
	log.Default().Println("ListDept-param:", param)
	var rows, total = system.GetDeptList(param, false)

	if total > 0 {
		for _, dept := range rows {
			dept.Children = []system.SysDeptResult{}
		}
	}
	return R.ReturnSuccess(rows)
}

/*排除节点*/
func (a *DeptAPI) ExcludeDept(deptId string) R.Result {
	// deptId := context.Param("deptId")
	var param = tools.SearchTableDataParam{}
	var list, _ = system.GetDeptList(param, false)
	var ExcludeList []system.SysDeptResult
	for i := 0; i < len(list); i++ {
		bean := list[i]
		ancestors := bean.Ancestors
		ancestors1 := utils.SplitStr(ancestors)
		index := arrays.ContainsString(ancestors1, deptId)
		if deptId != strconv.Itoa(bean.DeptId) || index == -1 {
			ExcludeList = append(ExcludeList, bean)
		}
	}
	return R.ReturnSuccess(ExcludeList)
}

func (a *DeptAPI) GetDept(deptId string) R.Result {
	// deptId := context.Param("deptId")
	result := system.GetDeptInfo(deptId)
	return R.ReturnSuccess(result)
}

func (a *DeptAPI) SaveDept(params map[string]interface{}) R.Result {

	userId := jwt.CacheGetUserId()
	var deptParam = system.SysDept{}
	if utils.CheckTypeByReflectNil(params) {
		return R.ReturnFailMsg("参数不能为空")
	}
	err := mapstructure.Decode(params, &deptParam)
	if err != nil {
		fmt.Println(err.Error())
	}
	user := system.FindUserById(userId)
	deptParam.CreateBy = user.UserName
	deptParam.CreateTime = time.Now()
	result := system.SaveDept(deptParam)
	return result
}

func (a *DeptAPI) UpDataDept(params map[string]interface{}) R.Result {
	userId := jwt.CacheGetUserId()
	var deptParam = system.SysDept{}
	if utils.CheckTypeByReflectNil(params) {
		return R.ReturnFailMsg("参数不能为空")
	}
	err := mapstructure.Decode(params, &deptParam)
	if err != nil {
		fmt.Println(err.Error())
	}
	user := system.FindUserById(userId)
	deptParam.UpdateBy = user.UserName
	deptParam.UpdateTime = time.Now()
	result := system.UpDataDept(deptParam)
	return result
}

func (a *DeptAPI) DeleteDept(deptId string) R.Result {
	// var deptId = context.Param("deptId")
	result := system.DeleteDept(deptId)
	return result
}
