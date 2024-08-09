// import request from '@/utils/request'
// import { parseStrEmpty } from "@/utils/ruoyi";

import { GetUserDeptTree, AvatarHandler, ProfileHandler, ListUser, PostProfileHandler, SaveUser, GetUserInfo, UpdateUser, DeleteUserById, ResetPwd, ChangeUserStatus, GetAuthUserRole, PutAuthUser } from 'wailsjs/go/system/UserAPI'

// 查询用户列表
export function listUser(query) {
  let queryParams = query;
  queryParams = {
    pageNum: query.pageNum,
    pageSize: query.pageSize,
    other: {
      userName: query.userName,
      phonenumber: query.phonenumber,
      status: query.status,
      deptId: query.deptId,
    }
  }
  console.debug(queryParams)
  return new Promise((resolve, reject) => {
    return ListUser(queryParams).then((res) => {
      console.debug(res)
      resolve(res)
    }).catch((err) => {
      reject(err)
    });
  })
}

// 查询用户详细
export function getUser(userId) {
  return new Promise((resolve, reject) => {
    return GetUserInfo(userId).then((res) => {
      resolve(res.data)
    }).catch((err) => {
      reject(err)
    });
  })
}

// 新增用户
export function addUser(data) {
  return SaveUser(data)
}

// 修改用户
export function updateUser(data) {
  return UpdateUser(data)
}

// 删除用户
export function delUser(userId) {
  return DeleteUserById(userId)
}

// 用户密码重置
export function resetUserPwd(userId, password) {
  const data = {
    userId,
    password
  }
  return ResetPwd(data)
}

// 用户状态修改
export function changeUserStatus(userId, status) {
  const data = {
    userId,
    status
  }
  return ChangeUserStatus(data)
}

// 查询用户个人信息
export function getUserProfile() {
  return ProfileHandler()
}

// 修改用户个人信息
export function updateUserProfile(data) {
  return PostProfileHandler(data)
}

// 用户密码重置
export function updateUserPwd(oldPassword, newPassword) {
  const data = {
    oldPassword,
    newPassword
  }
  return UpdatePwdHandler(data)
}

// 用户头像上传
export function uploadAvatar(data) {
  return AvatarHandler(data)
}

// 查询授权角色
export function getAuthRole(userId) {
  return GetAuthUserRole(userId)
}

// 保存授权角色
export function updateAuthRole(data) {
  return PutAuthUser(data)
}

// 查询部门下拉树结构
export function deptTreeSelect() {
  return new Promise((resolve, reject) => {
    return GetUserDeptTree().then((res) => {
      console.debug('deptTreeSelect', res)
      resolve(res)
    }).catch((err) => {
      reject(err)
    });
  })
}
