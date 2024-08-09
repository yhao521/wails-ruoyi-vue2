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

type RoleAPI struct {
	ctx context.Context
	baseAPI.Base
}

var roleAPI *RoleAPI
var onceRole sync.Once

// NewApp creates a new App application struct
func NewRoleAPI() *RoleAPI {
	if roleAPI == nil {
		onceRole.Do(func() {
			roleAPI = &RoleAPI{}
		})
	}
	return roleAPI
}

func (a *RoleAPI) ListRole(params map[string]interface{}) tools.TableDataInfo {

	var param = tools.SearchTableDataParam{
		Other: system.SysRoles{
			// RoleName: roleName,
			// RoleKey:  roleKey,
			// Status:   status,
		},
		Params: tools.Params{
			// BeginTime: beginTime,
			// EndTime:   endTime,
		},
	}

	err := mapstructure.Decode(params, &param)
	if err != nil {
		fmt.Println(err.Error())
	}

	var result = system.SelectRoleList(param, true)
	return result
}

func (a *RoleAPI) ExportRole(params map[string]interface{}) {

	var param = tools.SearchTableDataParam{
		// PageNum:  pageNum,
		// PageSize: pageSize,
		Other: system.SysRoles{
			// RoleName: roleName,
			// RoleKey:  roleKey,
			// Status:   status,
		},
		Params: tools.Params{
			// BeginTime: beginTime,
			// EndTime:   endTime,
		},
	}
	err := mapstructure.Decode(params, &param)
	if err != nil {
		fmt.Println(err.Error())
	}
	var result = system.SelectRoleList(param, false)
	var list = result.Rows.([]system.SysRoles)
	//定义首行标题
	dataKey := make([]map[string]string, 0)
	dataKey = append(dataKey, map[string]string{
		"key":    "roleId",
		"title":  "角色序号",
		"width":  "10",
		"is_num": "0",
	})
	dataKey = append(dataKey, map[string]string{
		"key":    "roleName",
		"title":  "角色名称",
		"width":  "10",
		"is_num": "0",
	})
	dataKey = append(dataKey, map[string]string{
		"key":    "roleKey",
		"title":  "角色权限",
		"width":  "10",
		"is_num": "0",
	})
	dataKey = append(dataKey, map[string]string{
		"key":    "roleSort",
		"title":  "角色排序",
		"width":  "10",
		"is_num": "0",
	})
	dataKey = append(dataKey, map[string]string{
		"key":    "dataScope",
		"title":  "数据范围",
		"width":  "10",
		"is_num": "0",
	})
	dataKey = append(dataKey, map[string]string{
		"key":    "status",
		"title":  "角色状态",
		"width":  "10",
		"is_num": "0",
	})
	//填充数据
	data := make([]map[string]interface{}, 0)
	if len(list) > 0 {
		for _, v := range list {
			var dataScope = v.DataScope
			dataScopeStr := ""
			if dataScope == "1" {
				dataScopeStr = "所有数据权限"
			}
			if dataScope == "2" {
				dataScopeStr = "自定义数据权限"
			}
			if dataScope == "3" {
				dataScopeStr = "本部门数据权限"
			}
			if dataScope == "4" {
				dataScopeStr = "本部门及以下数据权限"
			}
			if dataScope == "5" {
				dataScopeStr = "仅本人数据权限"
			}
			status := v.Status
			statusStr := ""
			if status == "1" {
				statusStr = "正常"
			}
			if status == "0" {
				statusStr = "停用"
			}
			data = append(data, map[string]interface{}{
				"roleId":    v.RoleId,
				"roleName":  v.RoleName,
				"roleKey":   v.RoleKey,
				"roleSort":  v.RoleSort,
				"dataScope": dataScopeStr,
				"status":    statusStr,
			})
		}
	}
	ex := tools.NewMyExcel()
	ex.ExportToWeb(dataKey, data, a.ctx)
}

func (a *RoleAPI) GetRoleInfo(roleId string) R.Result {
	// roleId := context.Param("roleId")

	userId := jwt.CacheGetUserId()
	if !system.IsAdminById(utils.GetInterfaceToInt(userId)) {
		var isCheck = system.CheckRoleDataScope(roleId)
		if isCheck {
			return R.ReturnFailMsg("没有权限访问角色数据！")

		}
	}
	result := system.FindRoleInfoById(roleId)
	return R.ReturnSuccess(result)
}

func (a *RoleAPI) SaveRole(params map[string]interface{}) R.Result {

	userId := jwt.CacheGetUserId()
	var rolesParam = system.SysRolesParam{}
	if utils.CheckTypeByReflectNil(params) {
		return R.ReturnFailMsg("参数不能为空")
	}
	err := mapstructure.Decode(params, &rolesParam)
	if err != nil {
		fmt.Println(err.Error())
	}
	user := system.FindUserById(userId)
	rolesParam.CreateBy = user.UserName
	rolesParam.CreateTime = time.Now()
	result := system.SaveRole(rolesParam)
	return result
}

/*修改权限*/
func (a *RoleAPI) UploadRole(params map[string]interface{}) R.Result {

	userId := jwt.CacheGetUserId()
	var rolesParam = system.SysRolesParam{}
	if utils.CheckTypeByReflectNil(params) {
		return R.ReturnFailMsg("参数不能为空")
	}
	err := mapstructure.Decode(params, &rolesParam)
	if err != nil {
		fmt.Println(err.Error())
	}
	user := system.FindUserById(userId)
	rolesParam.UpdateBy = user.UserName
	rolesParam.UpdateTime = time.Now()
	result := system.UploadRole(rolesParam, utils.GetInterfaceToInt(userId))
	return result
}

func (a *RoleAPI) PutDataScope(params map[string]interface{}) R.Result {

	userId := jwt.CacheGetUserId()
	var rolesParam = system.SysRolesParam{}
	if utils.CheckTypeByReflectNil(params) {
		return R.ReturnFailMsg("参数不能为空")
	}
	err := mapstructure.Decode(params, &rolesParam)
	if err != nil {
		fmt.Println(err.Error())
	}
	result := system.PutDataScope(rolesParam, utils.GetInterfaceToInt(userId))
	return result
}

func (a *RoleAPI) ChangeRoleStatus(params map[string]interface{}) R.Result {

	userId := jwt.CacheGetUserId()
	var rolesParam = system.SysRoles{}
	if utils.CheckTypeByReflectNil(params) {
		return R.ReturnFailMsg("参数不能为空")
	}
	err := mapstructure.Decode(params, &rolesParam)
	if err != nil {
		fmt.Println(err.Error())
	}
	user := system.FindUserById(userId)
	rolesParam.UpdateBy = user.UserName
	rolesParam.UpdateTime = time.Now()
	result := system.ChangeRoleStatus(rolesParam, utils.GetInterfaceToInt(userId))
	return result
}

func (a *RoleAPI) DeleteRole(roleIds string) R.Result {

	userId := jwt.CacheGetUserId()
	// var roleIds = context.Param("roleIds")
	system.DeleteRolesById(roleIds, utils.GetInterfaceToInt(userId))
	return R.ReturnSuccess(gin.H{
		"msg":  "操作成功",
		"code": http.StatusOK,
		"data": "",
	})
}

/*不需要分页*/
func (a *RoleAPI) GetRoleOptionSelect() R.Result {
	result := system.GetRoleOptionSelect()
	return result
}

/*分页*/
func (a *RoleAPI) GetAllocatedList(params map[string]interface{}) tools.TableDataInfo {

	var param = tools.SearchTableDataParam{
		// PageNum:  pageNum,
		// PageSize: pageSize,
		Other: system.SysUserParm{
			// RoleId:      roleId,
			// UserName:    userName,
			// Phonenumber: phonenumber,
		},
		Params: tools.Params{
			// BeginTime: beginTime,
			// EndTime:   endTime,
		},
	}
	err := mapstructure.Decode(params, &param)
	if err != nil {
		fmt.Println(err.Error())
	}

	result := system.GetAllocatedList(param)
	return result
}

/*分页*/
func (a *RoleAPI) GetUnAllocatedList(params map[string]interface{}) tools.TableDataInfo {

	var param = tools.SearchTableDataParam{
		Other: system.SysUserParm{
			// RoleId:      roleId,
			// UserName:    userName,
			// Phonenumber: phonenumber,
		},
		Params: tools.Params{
			// BeginTime: beginTime,
			// EndTime:   endTime,
		},
	}
	err := mapstructure.Decode(params, &param)
	if err != nil {
		fmt.Println(err.Error())
	}

	result := system.GetUnAllocatedList(param)
	return result
}

func (a *RoleAPI) CancelRole(params map[string]interface{}) R.Result {
	var rolesParam = system.SysUserRolesParam{}
	if utils.CheckTypeByReflectNil(params) {
		return R.ReturnFailMsg("参数不能为空")
	}
	err := mapstructure.Decode(params, &rolesParam)
	if err != nil {
		fmt.Println(err.Error())
	}

	result := system.CancelRole(rolesParam.UserId, rolesParam.RoleId)
	return result
}
func (a *RoleAPI) CancelAllRole(params map[string]string) R.Result {
	roleId := params["roleId"]
	userIds := params["userIds"]

	result := system.CancelAllRole(roleId, userIds)
	return result
}

func (a *RoleAPI) SelectRoleAll(params map[string]string) R.Result {
	roleId := params["roleId"]
	userIds := params["userIds"]
	userId := params["userId"]
	result := system.SelectRoleAll(roleId, userIds, utils.GetInterfaceToInt(userId))
	return result
}

func (a *RoleAPI) GetDeptTreeRole(roleId string) R.Result {
	// roleId := context.Param("roleId")
	checkedKeys := system.GetDeptTreeRole(roleId)
	depts := system.SelectDeptTreeList()
	return R.ReturnSuccess(gin.H{
		"msg":         "操作成功",
		"code":        http.StatusOK,
		"checkedKeys": checkedKeys,
		"depts":       depts,
	})
}
