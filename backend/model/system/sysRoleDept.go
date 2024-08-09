package system

import (
	"mySparkler/pkg/db"
	"mySparkler/pkg/utils/R"
)

type SysRoleDept struct {
	RoleId int `json:"roleId" gorm:"column:role_id"`
	DeptId int `json:"deptId" gorm:"column:dept_id"`
}

// TableName 指定数据库表名称
func (SysRoleDept) TableName() string {
	return "sys_role_dept"
}

// 删除角色与部门关联
func DeleteRoleDept(roleIds string) {
	sql := "delete from sys_role_dept where role_id in ( " + roleIds + " )"
	err := db.Db().Exec(sql).Error
	if err != nil {
		panic(R.ReturnFailMsg(err.Error()))
	}
}
func DeleteRoleDeptByRole(roleId string) {
	err := db.Db().Where("role_id = ?", roleId).Delete(&SysRoleDept{}).Error
	if err != nil {
		panic(R.ReturnFailMsg(err.Error()))
	}
}
