// import SaveRolerequest from '@/utils/request'

import {
  ListRole, ExportRole, GetRoleInfo, SaveRole, UploadRole, PutDataScope, ChangeRoleStatus, DeleteRole
  , GetRoleOptionSelect, GetAllocatedList, GetUnAllocatedList, CancelRole, CancelAllRole, SelectRoleAll
  , GetDeptTreeRole
} from 'wailsjs/go/system/RoleAPI'

// 查询角色列表
export function listRole(query) {
  let queryParams = query;
  queryParams = {
    pageNum: query.pageNum,
    pageSize: query.pageSize,
    other: {
      roleName: query.roleName,
      roleKey: query.roleKey,
      status: query.status,
    }
  }
  return new Promise((resolve, reject) => {
    return ListRole(queryParams).then((res) => {
      resolve(res)
    }).catch((err) => {
      reject(err)
    });
  })
}

// 查询角色详细
export function getRole(roleId) {
  return GetRoleInfo(roleId)
}

// 新增角色
export function addRole(data) {
  return SaveRole(data)
}

// 修改角色
export function updateRole(data) {
  return UploadRole(
    data
  )
}

// 角色数据权限
export function dataScope(data) {
  return PutDataScope(
    data
  )
}

// 角色状态修改
export function changeRoleStatus(roleId, status) {
  const data = {
    roleId,
    status
  }
  return ChangeRoleStatus(data)
}

// 删除角色
export function delRole(roleId) {
  return DeleteRole(roleId)
}

// 查询角色已授权用户列表
export function allocatedUserList(query) {
  return GetRoleOptionSelect(query)
}

// 查询角色未授权用户列表
export function unallocatedUserList(query) {
  return GetUnAllocatedList(query)
}

// 取消用户授权角色
export function authUserCancel(data) {
  return CancelRole(data)
}

// 批量取消用户授权角色
export function authUserCancelAll(data) {
  return CancelAllRole(data)
}

// 授权用户选择
export function authUserSelectAll(data) {
  return SelectRoleAll(data)
}

// 根据角色ID查询部门树结构
export function deptTreeSelect(roleId) {
  return GetDeptTreeRole(roleId)
}
