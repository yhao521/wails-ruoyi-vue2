package system

import (
	"context"
	"fmt"
	"mySparkler/backend/api/baseAPI"
	"mySparkler/backend/model/monitor"
	"mySparkler/backend/model/system"
	"mySparkler/backend/model/tools"
	"mySparkler/pkg/constants"
	"mySparkler/pkg/db"
	"mySparkler/pkg/utils"
	"mySparkler/pkg/utils/R"
	"mySparkler/pkg/utils/jwt"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	"github.com/xuri/excelize/v2"
)

type UserAPI struct {
	ctx context.Context
	baseAPI.Base
}

var userAPI *UserAPI
var onceUser sync.Once

// NewApp creates a new App application struct
func NewUserAPI() *UserAPI {
	if userAPI == nil {
		onceUser.Do(func() {
			userAPI = &UserAPI{}
		})
	}
	return userAPI
}

func (a *UserAPI) LoginHandler(username string, password string, code string, uuid string) R.Result {
	var param = system.LoginParam{
		Code:     code,
		Password: password,
		UserName: username,
		Uuid:     uuid,
	}
	if reflect.DeepEqual(param, system.LoginParam{}) {
		monitor.LoginInfoAdd(param, "登录失败，参数为空", false)
		return R.ReturnFailMsg("参数不能为空")
	} else {
		fmt.Println("param.UserName: " + param.UserName)
		var captchaEnabled = system.SelectCaptchaEnabled()

		if captchaEnabled {
			isVerify := utils.VerifyCaptcha(param.Uuid, param.Code)
			if isVerify {
				isExist := system.IsExistUser(param.UserName)
				if isExist {
					return findUser(param)
				} else {
					monitor.LoginInfoAdd(param, "登录失败，用户不存在", false)
					return R.ReturnFailMsg("用户不存在")
				}
			} else {
				monitor.LoginInfoAdd(param, "登录失败，验证码错误", false)
				return R.ReturnFailMsg("请输入正确的验证码")
			}
		} else {
			isExist := system.IsExistUser(param.UserName)
			if isExist {
				return findUser(param)
			} else {
				monitor.LoginInfoAdd(param, "登录失败，用户不存在", false)
				return R.ReturnFailMsg("用户不存在")
			}
		}
	}
}

// 判断用户是否存在 返回bool类型
func findUser(param system.LoginParam) R.Result {
	var loginName = param.UserName
	var pass = param.Password
	var user = system.SysUser{}
	cacheUserId := jwt.CacheGetUserId()
	if cacheUserId != 0 {
		user = system.FindUserPassById(cacheUserId)
	} else {
		user = system.FindUserByName(loginName)
	}
	fmt.Println(user.UserName)
	if user.UserId != 0 {
		if user.Status == "1" {
			monitor.LoginInfoAdd(param, "登录失败，账号已停用", false)
			return R.ReturnFailMsg("账号已停用")

		}
		fmt.Println("pass", pass, user.Password)
		// 验证 密码是否正确
		if utils.PasswordVerify(pass, user.Password) {
			// tokenString, err := jwt.CreateToken(user.UserName, user.UserId, user.DeptId)
			// if err != nil {
			// 	monitor.LoginInfoAdd(param, "登录失败，"+err.Error(), false)
			// 	return R.ReturnFailMsg("登录失败")

			// }
			if cacheUserId == 0 {
				jwt.CacheUserId(user.UserId)
			}
			monitor.LoginInfoAdd(param, "登录成功", true)
			return R.ReturnSuccess(gin.H{
				"msg":   "登录成功",
				"token": "99999",
				"id":    user.UserId,
			})
		} else {
			monitor.LoginInfoAdd(param, "登录失败，密码错误", false)
			return R.ReturnFailMsg("登录失败，密码错误")
		}

	} else {
		monitor.LoginInfoAdd(param, "登录失败，用户不存在", false)
		return R.ReturnFailMsg("用户不存在")
	}
}

func (a *UserAPI) GetInfoHandler() R.Result {
	userId := jwt.CacheGetUserId()
	fmt.Println("GetInfoHandler", userId)
	var user = system.FindUserById(userId)
	roles := system.GetRolePermission(user)
	permissions := system.GetMenuPermission(user)
	dept := system.GetDeptInfo(strconv.Itoa(user.DeptId))
	return R.ReturnSuccess(gin.H{
		"msg":    "获取成功",
		"code":   http.StatusOK,
		"userId": userId,
		"user": gin.H{
			"userName":    user.UserName,
			"nickName":    user.NickName,
			"phonenumber": user.Phonenumber,
			"email":       user.Email,
			"avatar":      user.Avatar,
			"sex":         user.Sex,
			"createTime":  user.CreateTime,
			"dept":        dept,
		},
		"roles":       roles,
		"permissions": permissions,
	})
}
func (a *UserAPI) GetCookie() R.Result {

	userId := jwt.CacheGetUserId()
	// if userId != nil {
	var user = system.FindUserById(userId)
	fmt.Println("GetCookie", user)
	if user.UserId != 0 {
		return R.ReturnSuccess(user)
	} else {
		return R.ReturnFailMsg("获取失败")
	}
	// }
	// return R.ReturnSuccess("退出成功")
}
func (a *UserAPI) LogoutHandler() R.Result {

	userId := jwt.CacheGetUserId()
	// if userId != nil {
	var user = system.FindUserById(userId)
	fmt.Println("LogoutHandler", user)
	if user.UserId != 0 {
		// 开始删除缓存
		jwt.CacheCleanAll()
		return R.ReturnSuccess("退出成功")
	} else {

		return R.ReturnFailMsg("退出失败")
	}
	// }
	// return R.ReturnSuccess("退出成功")
}

// CaptchaImageHandler 验证码 输出
func (a *UserAPI) CaptchaImageHandler() R.Result {
	var captchaEnabled = system.SelectCaptchaEnabled()
	if captchaEnabled {
		id, b64s, _, err := utils.CreateImageCaptcha("Number")
		if err != nil {
			return R.ReturnFailMsg("创建二维码失败，请联系管理员")

		}
		return R.ReturnSuccess(gin.H{
			"msg":            "操作成功",
			"img":            strings.ReplaceAll(b64s, "data:image/png;base64,", ""),
			"code":           http.StatusOK,
			"captchaEnabled": captchaEnabled,
			"uuid":           id,
		})
	} else {
		return R.ReturnSuccess(gin.H{
			"msg":            "操作成功",
			"code":           http.StatusOK,
			"captchaEnabled": captchaEnabled,
		})
	}
}

// UpdatePwdHandler 修改密码
func (a *UserAPI) UpdatePwdHandler(params map[string]string) R.Result {

	userId := jwt.CacheGetUserId()
	var user = system.FindUserById(userId)
	// var newPassword1 = params["newPassword"]
	// println(newPassword1)
	// 没有这个，下面的为空很奇怪

	// from := context.Request.Form
	OldPassword := params["oldPassword"]
	NewPassword := params["newPassword"]

	if OldPassword == "" {
		return R.ReturnFailMsg("参数不能为空")
	}

	if NewPassword == "" {
		return R.ReturnFailMsg("参数不能为空")

	}

	// 验证旧密码
	if utils.PasswordVerify(OldPassword, user.Password) {
		// 验证新密码
		if utils.PasswordVerify(NewPassword, user.Password) {
			return R.ReturnFailMsg("新密码不能与旧密码相同")

		}
		// 加密
		passString, err3 := utils.PasswordHash(NewPassword)
		if err3 != nil {
			return R.ReturnFailMsg("加密失败")

		}
		// 更新 密码
		err2 := db.Db().Model(&user).Update("password", passString)
		if err2.Error != nil {
			return R.ReturnFailMsg("修改密码失败")

		}
		return R.ReturnSuccess("修改密码成功")
	} else {
		return R.ReturnFailMsg("旧密码错误")

	}
}

// ProfileHandler 查询个人信息
func (a *UserAPI) ProfileHandler() R.Result {

	userId := jwt.CacheGetUserId()
	var user system.SysUser
	err1 := db.Db().Where("user_id = ?", userId).First(&user)

	if err1.Error != nil {
		return R.ReturnFailMsg("未找到用户")

	}
	user.Password = ""
	return R.ReturnSuccess(gin.H{
		"msg":  "获取成功",
		"code": http.StatusOK,
		"data": gin.H{
			"userName":    user.UserName,
			"nickName":    user.NickName,
			"phonenumber": user.Phonenumber,
			"email":       user.Email,
			"sex":         user.Sex,
			"createTime":  user.CreateTime.Format(constants.TimeFormat),
		},
		"roleGroup": system.SelectRolesByUserName(user.UserName), // 目前暂时没有用
		"postGroup": system.SelectUserPostGroup(user.UserName),   // 目前暂时没有用
	})
}

// PostProfileHandler 修改个人信息
func (a *UserAPI) PostProfileHandler(params map[string]interface{}) R.Result {

	userId := jwt.CacheGetUserId()
	param := system.Userparam{}
	if utils.CheckTypeByReflectNil(params) {
		return R.ReturnFailMsg("参数不能为空")
	}
	err := mapstructure.Decode(params, &param)
	if err != nil {
		fmt.Println(err.Error())
	}
	var result = system.EditProfileUserInfo(utils.GetInterfaceToInt(userId), param)
	return result
}

// AvatarHandler 上传头像 并更新
func (a *UserAPI) AvatarHandler(context *gin.Context) R.Result {

	userId := jwt.CacheGetUserId()
	var user system.SysUser
	err1 := db.Db().Where("user_id = ?", userId).First(&user)

	if err1.Error != nil {
		return R.ReturnFailMsg("未找到用户")
	}

	file, errLoad := context.FormFile("avatarfile")
	if errLoad != nil {
		msg := "获取上传文件错误:" + errLoad.Error()
		return R.ReturnFailMsg(msg)
	}

	//fileExt := strings.ToLower(path.Ext(file.Filename))
	//if fileExt != ".png" && fileExt != ".jpg" && fileExt != ".gif" && fileExt != ".jpeg" {
	//	return R.ReturnSuccess(gin.H{
	//		"code": http.StatusInternalServerError,
	//		"msg":  "上传失败!只允许png,jpg,gif,jpeg文件",
	//	})
	//	return
	//}

	//上传图片
	errFile := context.SaveUploadedFile(file, "./static/images/"+file.Filename)
	if errFile != nil {
		return R.ReturnFailMsg("上传图片异常，请联系管理员")
	}

	// 更新 头像
	err2 := db.Db().Model(&user).Update("avatar", "/profile/"+file.Filename)
	if err2.Error != nil {
		return R.ReturnFailMsg("上传图片异常，请联系管理员")
	}

	return R.ReturnSuccess(gin.H{
		"code":   http.StatusOK,
		"msg":    "上传头像成功!",
		"imgUrl": "/profile/" + file.Filename,
	})
}

/*-----------用户管理----------------------------*/

func (a *UserAPI) ListUser(params map[string]interface{}) tools.TableDataInfo {

	var param = tools.SearchTableDataParam{
		// PageNum:  pageNum,
		// PageSize: pageSize,
		Other: system.SysUser{
			// DeptId:      deptId,
			// UserName:    userName,
			// Phonenumber: phonenumber,
			// Status:      status,
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

	return system.SelectUserList(param, true)
}

func (a *UserAPI) ExportExport(params map[string]interface{}) {

	var param = tools.SearchTableDataParam{
		// PageNum:  pageNum,
		// PageSize: pageSize,
		Other: system.SysUser{
			// DeptId:      deptId,
			// UserName:    userName,
			// Phonenumber: phonenumber,
			// Status:      status,
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

	var data1 = system.SelectUserParmList(param, true)
	var list = data1.Rows.([]system.SysUserExcel)

	//定义首行标题
	dataKey := make([]map[string]string, 0)
	dataKey = append(dataKey, map[string]string{
		"key":    "userId",
		"title":  "用户序号",
		"width":  "10",
		"is_num": "0",
	})
	dataKey = append(dataKey, map[string]string{
		"key":    "deptId",
		"title":  "部门编号",
		"width":  "15",
		"is_num": "0",
	})
	dataKey = append(dataKey, map[string]string{
		"key":    "userName",
		"title":  "登录名称",
		"width":  "20",
		"is_num": "0",
	})
	dataKey = append(dataKey, map[string]string{
		"key":    "nickName",
		"title":  "用户名称",
		"width":  "20",
		"is_num": "0",
	})
	dataKey = append(dataKey, map[string]string{
		"key":    "email",
		"title":  "用户邮箱",
		"width":  "20",
		"is_num": "0",
	})
	dataKey = append(dataKey, map[string]string{
		"key":    "phonenumber",
		"title":  "手机号码",
		"width":  "20",
		"is_num": "0",
	})
	dataKey = append(dataKey, map[string]string{
		"key":    "sex",
		"title":  "用户性别",
		"width":  "10",
		"is_num": "0",
	})
	dataKey = append(dataKey, map[string]string{
		"key":    "status",
		"title":  "帐号状态",
		"width":  "10",
		"is_num": "0",
	})
	dataKey = append(dataKey, map[string]string{
		"key":    "loginIp",
		"title":  "最后登录IP",
		"width":  "30",
		"is_num": "0",
	})
	dataKey = append(dataKey, map[string]string{
		"key":    "loginDate",
		"title":  "最后登录时间",
		"width":  "60",
		"is_num": "0",
	})
	dataKey = append(dataKey, map[string]string{
		"key":    "deptName",
		"title":  "部门名称",
		"width":  "50",
		"is_num": "0",
	})
	dataKey = append(dataKey, map[string]string{
		"key":    "leader",
		"title":  "部门负责人",
		"width":  "30",
		"is_num": "0",
	})

	//填充数据
	data := make([]map[string]interface{}, 0)
	if len(list) > 0 {
		for _, v := range list {
			var sexStatus = v.Sex
			var sex = ""
			if sexStatus == "0" {
				sex = "男"
			} else if sexStatus == "1" {
				sex = "女"
			} else {
				sex = "未知"
			}
			var status = v.Status
			var statusStr = ""
			if status == "0" {
				statusStr = "正常"
			} else {
				statusStr = "停用"
			}
			//timeObj, _ := time.Parse(time.RFC3339, v.LoginDate)
			var loginData = v.LoginDate.Format(constants.TimeFormat)
			data = append(data, map[string]interface{}{
				"userId":      v.UserId,
				"deptId":      v.DeptId,
				"userName":    v.UserName,
				"nickName":    v.NickName,
				"email":       v.Email,
				"phonenumber": v.Phonenumber,
				"sex":         sex,
				"status":      statusStr,
				"loginIp":     v.LoginIp,
				"loginDate":   loginData,
				"deptName":    v.DeptName,
				"leader":      v.Leader,
			})
		}
	}
	ex := tools.NewMyExcel()
	ex.ExportToWeb(dataKey, data, a.ctx)
}

func (a *UserAPI) ImportUserData(context *gin.Context) R.Result {
	file, _, errLoad := context.Request.FormFile("file")
	if errLoad != nil {
		msg := "获取上传文件错误:" + errLoad.Error()
		return R.ReturnFailMsg(msg)
	}

	xlsx, err := excelize.OpenReader(file)
	if err != nil {
		return R.ReturnFailMsg("请选择文件")
	}
	var users []system.SysUserParm

	var updateSupport = context.DefaultQuery("updateSupport", "")

	rows, _ := xlsx.GetRows("Sheet1")
	for irow, row := range rows {
		if irow > 0 {
			var data []string
			for i := range row {
				data = append(data, row[i])
			}
			users = append(users, system.SysUserParm{
				UserName:    data[2],
				NickName:    data[3],
				Phonenumber: data[5],
				Sex:         data[6],
				Email:       data[4],
				Status:      data[7],
				CreateTime:  time.Now(),
				DeptId:      utils.GetInterfaceToInt(data[1]),
				PostIds:     utils.Split(data[8]),
				RoleIds:     utils.Split(data[9]),
			})
		}
	}
	if len(users) == 0 {
		return R.ReturnFailMsg("请在表格中添加数据")
	}

	var error, message = system.ImportUserData(users, updateSupport)

	if error == "" {
		return R.ReturnSuccess(message)
	} else {
		return R.ReturnFailMsg(message)
	}

}

// 下载模版
func (a *UserAPI) ImportTemplate(context *gin.Context) {
	//定义首行标题
	dataKey := make([]map[string]string, 0)
	dataKey = append(dataKey, map[string]string{
		"key":    "userId",
		"title":  "用户序号",
		"width":  "10",
		"is_num": "0",
	})
	dataKey = append(dataKey, map[string]string{
		"key":    "deptId",
		"title":  "部门编号",
		"width":  "15",
		"is_num": "0",
	})
	dataKey = append(dataKey, map[string]string{
		"key":    "userName",
		"title":  "登录名称",
		"width":  "20",
		"is_num": "0",
	})
	dataKey = append(dataKey, map[string]string{
		"key":    "nickName",
		"title":  "用户名称",
		"width":  "20",
		"is_num": "0",
	})
	dataKey = append(dataKey, map[string]string{
		"key":    "email",
		"title":  "用户邮箱",
		"width":  "20",
		"is_num": "0",
	})
	dataKey = append(dataKey, map[string]string{
		"key":    "phonenumber",
		"title":  "手机号码",
		"width":  "20",
		"is_num": "0",
	})
	dataKey = append(dataKey, map[string]string{
		"key":    "sex",
		"title":  "用户性别",
		"width":  "10",
		"is_num": "0",
	})
	dataKey = append(dataKey, map[string]string{
		"key":    "status",
		"title":  "帐号状态",
		"width":  "10",
		"is_num": "0",
	})
	dataKey = append(dataKey, map[string]string{
		"key":    "status",
		"title":  "岗位",
		"width":  "11",
		"is_num": "0",
	})
	dataKey = append(dataKey, map[string]string{
		"key":    "status",
		"title":  "权限",
		"width":  "12",
		"is_num": "0",
	})
	//填充数据
	data := make([]map[string]interface{}, 0)

	ex := tools.NewMyExcel()
	ex.ExportToWeb(dataKey, data, context)
}

func (a *UserAPI) GetUserInfo(useridP int) R.Result {
	/*参数用户*/
	// userIdStr := context.Param("userId")
	//登录用户
	userId := jwt.CacheGetUserId()
	// userId := utils.GetInterfaceToInt(loginUserId)
	fmt.Println("GetUserInfo-userIdStr:", useridP)

	if useridP == 0 {
		useridP = userId
	}
	// useridP, _ := strconv.Atoi(userIdStr)

	system.CheckUserDataScope(userId, useridP)

	user := system.FindUserById(useridP)

	var roles []system.SysRoles
	var roleIds []int
	// 登录者的权限
	roles = system.SelectRolePermissionByUserId(userId)

	posts := system.SelectSysPostList(tools.SearchTableDataParam{
		Other: system.SysPost{},
	}, false).Rows

	var result = gin.H{
		"msg":   "操作成功",
		"code":  http.StatusOK,
		"data":  user,
		"roles": roles,
		"posts": posts,
	}
	// 判断当期是否为管理员
	if useridP != 0 {
		postIds := system.SelectPostListByUserId(useridP)
		roles2 := system.SelectRolePermissionByUserId(useridP)
		for _, sysRoles := range roles2 {
			roleIds = append(roleIds, sysRoles.RoleId)
		}
		result["postIds"] = postIds
		result["roleIds"] = roleIds
	}

	return R.ReturnSuccess(result)
}

func (a *UserAPI) SaveUser(params map[string]interface{}) R.Result {

	userId := jwt.CacheGetUserId()
	var user = system.SysUserParm{}
	if utils.CheckTypeByReflectNil(params) {
		return R.ReturnFailMsg("参数不能为空")
	}

	err := mapstructure.Decode(params, &user)
	if err != nil {
		fmt.Println(err.Error())
	}

	var username = user.UserName
	if system.IsExistUser(username) {
		return R.ReturnSuccess(gin.H{
			"msg": "用户名已存在",
		})
	}
	var phonenumber = user.Phonenumber
	if system.IsExistUserByPhoneNumber(phonenumber) {
		return R.ReturnSuccess(gin.H{
			"msg":  "手机号已存在",
			"code": http.StatusInternalServerError,
		})
	}
	var email = user.Email
	if system.IsExistUserByEmail(email) {
		return R.ReturnSuccess(gin.H{
			"msg":  "邮箱已存在",
			"code": http.StatusInternalServerError,
		})
	}
	var password = user.Password
	var pwd, _ = utils.PasswordHash(password)
	var user1 = system.FindUserById(userId)
	user.CreateBy = user1.UserName
	user.CreateTime = time.Now()
	user.Password = pwd

	/*用户名、手机号、邮箱不能重复验证*/
	var message = system.SaveUser(user)
	return R.ReturnSuccess(gin.H{
		"msg":  "操作成功",
		"code": http.StatusOK,
		"data": message,
	})
}

func (a *UserAPI) UpdateUser(params map[string]interface{}) R.Result {

	userId := jwt.CacheGetUserId()
	var userParm = system.SysUserParm{}
	if utils.CheckTypeByReflectNil(params) {
		return R.ReturnFailMsg("参数不能为空")
	}
	err := mapstructure.Decode(params, &userParm)
	if err != nil {
		fmt.Println(err.Error())
	}
	var uId = userParm.UserId
	if uId == 0 {
		return R.ReturnFailMsg("参数不能为空")
	}
	var user = system.FindUserById(userId)
	var userSql = system.FindUserById(uId)

	var username = userParm.UserName
	if !system.IsExistUser(username) {
		if userSql.UserName != username {
			return R.ReturnSuccess(gin.H{
				"msg":  "用户名:" + username + "已存在",
				"code": http.StatusInternalServerError,
			})
		}
	}
	var phonenumber = userParm.Phonenumber
	if !system.IsExistUserByPhoneNumber(phonenumber) {
		if userSql.Phonenumber != phonenumber {
			return R.ReturnSuccess(gin.H{
				"msg":  "手机号:" + phonenumber + "已存在",
				"code": http.StatusInternalServerError,
			})
		}
	}
	var email = userParm.Email
	if !system.IsExistUserByEmail(email) {
		if userSql.Email != email {
			return R.ReturnSuccess(gin.H{
				"msg":  "邮箱:" + email + "已存在",
				"code": http.StatusInternalServerError,
			})
		}
	}
	userParm.UpdateBy = user.UserName
	userParm.UpdateTime = time.Now()
	system.UploadUser(userParm)
	return R.ReturnSuccess(gin.H{
		"msg":  "操作成功",
		"code": http.StatusOK,
	})
}

func (a *UserAPI) DeleteUserById(userIds string) R.Result {
	// var userIds = context.Param("userIds")

	userIdA := utils.Split(userIds)
	// len(userIdA)
	var uIds = []int{}
	for i, v := range userIdA {
		uIds[i] = utils.GetInterfaceToInt(v)
	}
	return system.DeleteUser(uIds)
}

// 重设密码
func (a *UserAPI) ResetPwd(params map[string]interface{}) R.Result {
	/*获取为空有可能内部参数错误*/
	var user = system.SysUserParm{}
	if utils.CheckTypeByReflectNil(params) {
		return R.ReturnFailMsg("参数不能为空")
	}
	err := mapstructure.Decode(params, &user)
	if err != nil {
		fmt.Println(err.Error())
	}
	result := system.ResetPwd(user)
	return result
}

func (a *UserAPI) ChangeUserStatus(params map[string]interface{}) R.Result {
	var user = system.SysUserParm{}
	if utils.CheckTypeByReflectNil(params) {
		return R.ReturnFailMsg("参数不能为空")
	}
	if user.UserId == constants.AdminId {
		return R.ReturnFailMsg("管理员数据不开操作")
	}
	err := mapstructure.Decode(params, &user)
	if err != nil {
		fmt.Println(err.Error())
	}
	result := system.ChangeUserStatus(user)
	return result
}

func (a *UserAPI) GetAuthUserRole(userId int) R.Result {
	// 登录者的权限
	var roles []system.SysRoles
	if system.IsAdminById(userId) {
		roles = system.SelectRolePermissionByUserId(userId)
	} else {
		roles = system.SelectRolePermissionByUserId(userId)
	}

	user := system.FindUserById(userId)

	return R.ReturnSuccess(gin.H{
		"user":  user,
		"roles": roles,
	})
}

func (a *UserAPI) PutAuthUser(uId int, roleIds string) R.Result {
	uIds := []int{uId}
	system.DeleteRolesByUser(uIds)
	result := system.InsertRolesByUser(uId, utils.Split(roleIds))
	return result
}

// 登录获取菜单
func (a *UserAPI) GetUserDeptTree() R.Result {
	var list = system.GetUserDeptTree()
	return R.ReturnSuccess(list)
}
