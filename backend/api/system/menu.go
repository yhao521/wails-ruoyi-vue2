package system

import (
	"context"
	"fmt"
	"mySparkler/backend/api/baseAPI"
	"mySparkler/backend/model/system"
	"mySparkler/backend/model/tools"
	"mySparkler/pkg/db"
	"mySparkler/pkg/utils"
	"mySparkler/pkg/utils/R"
	"mySparkler/pkg/utils/jwt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
)

// App struct
type MenuAPI struct {
	ctx context.Context
	baseAPI.Base
}

var menuAPI *MenuAPI
var onceMenu sync.Once

// NewApp creates a new App application struct
func NewMenuAPI() *MenuAPI {
	if menuAPI == nil {
		onceMenu.Do(func() {
			menuAPI = &MenuAPI{}
		})
	}
	return menuAPI
}

// GetRoutersHandler /*
func (a *MenuAPI) GetRoutersHandler() R.Result {

	userId := jwt.CacheGetUserId()
	var user system.SysUser
	err1 := db.Db().Where("user_id = ?", userId).First(&user)

	if err1.Error != nil {
		return R.ReturnFailMsg("未找到用户")

	}

	menu := system.SelectMenuTreeByUserId(utils.GetInterfaceToInt(userId))
	var data = system.BuildMenus(menu)
	return R.ReturnSuccess(data)
}

func (a *MenuAPI) ListMenu(params map[string]interface{}) R.Result {

	// menuName := context.DefaultQuery("menuName", "")
	// status := context.DefaultQuery("status", "")
	var param = tools.SearchTableDataParam{
		Other:  system.SysMenu{},
		Params: tools.Params{},
	}
	err := mapstructure.Decode(params, &param)
	if err != nil {
		fmt.Println(err.Error())
	}

	userId := jwt.CacheGetUserId()
	if system.IsAdminById(utils.GetInterfaceToInt(userId)) {
		var result = system.SelectSysMenuList(param)
		return R.ReturnSuccess(gin.H{
			"msg":  "操作成功",
			"code": http.StatusOK,
			"data": result,
		})
	} else {
		var result = system.SelectSysMenuListByUserId(int(userId), param)
		return R.ReturnSuccess(gin.H{
			"msg":  "操作成功",
			"code": http.StatusOK,
			"data": result,
		})
	}
}

func (a *MenuAPI) GetMenuInfo(menuId int) R.Result {
	// userId := jwt.CacheGetUserId()

	println(menuId)
	var date = system.FindMenuInfoById(menuId)
	return R.ReturnSuccess(gin.H{
		"msg":  "操作成功",
		"code": http.StatusOK,
		"data": date,
	})
}

func (a *MenuAPI) GetTreeSelect(params map[string]interface{}) R.Result {
	userId := jwt.CacheGetUserId()

	var menuParm = system.SysMenu{}
	err := mapstructure.Decode(params, &menuParm)
	if err != nil {
		fmt.Println(err.Error())
	}
	if utils.CheckTypeByReflectNil(menuParm) {
		return R.ReturnFailMsg("参数不能为空")
	}
	menu := system.SelectMenuTree(utils.GetInterfaceToInt(userId), menuParm)
	var result = system.BuildMenuTreeSelect(menu)

	return R.ReturnSuccess(result)
}

func (a *MenuAPI) TreeSelectByRole(roleId string, params map[string]interface{}) R.Result {
	userId := jwt.CacheGetUserId()

	var menuPerms = system.SysMenu{}
	err := mapstructure.Decode(params, &menuPerms)
	if err != nil {
		fmt.Println(err.Error())
	}
	menu := system.SelectMenuTree(utils.GetInterfaceToInt(userId), menuPerms)
	var result = system.BuildMenuTreeSelect(menu)
	roles := system.FindRoleInfoById(roleId)
	var checkedKeys = system.SelectMenuListByRoleId(roleId, roles.MenuCheckStrictly)
	return R.ReturnSuccess(gin.H{
		"menus":       result,
		"checkedKeys": checkedKeys,
	})
}

func (a *MenuAPI) SaveMenu(params map[string]interface{}) R.Result {

	var menuParm = system.SysMenu{}
	userId := jwt.CacheGetUserId()
	if utils.CheckTypeByReflectNil(params) {
		return R.ReturnFailMsg("参数不能为空")
	}
	err := mapstructure.Decode(params, &menuParm)
	if err != nil {
		fmt.Println(err.Error())
	}
	user := system.FindUserById(userId)
	menuParm.CreateBy = user.UserName
	menuParm.CreateTime = time.Now()
	result := system.AddMenu(menuParm)
	return R.ReturnSuccess(result)
}

func (a *MenuAPI) UpdateMenu(params map[string]interface{}) R.Result {
	userId := jwt.CacheGetUserId()

	var menuParm = system.SysMenu{}
	if utils.CheckTypeByReflectNil(params) {
		return R.ReturnFailMsg("参数不能为空")
	}

	err := mapstructure.Decode(params, &menuParm)
	if err != nil {
		fmt.Println(err.Error())
	}

	user := system.FindUserById(userId)
	menuParm.UpdateBy = user.UserName
	menuParm.UpdateTime = time.Now()
	result := system.UpdateMenu(menuParm)
	return R.ReturnSuccess(result)
}

func (a *MenuAPI) DeleteMenu(menuId int) R.Result {
	// var menuId = context.Param("menuId")
	result := system.DeleteMenu(menuId)
	return R.ReturnSuccess(result)
}
