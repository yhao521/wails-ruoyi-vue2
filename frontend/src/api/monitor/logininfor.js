import { LoginInformListHandler, DeleteByIdHandler, CleanHandler, UnlockHandler } from "wailsjs/go/monitor/LoginInforAPI";

// 查询登录日志列表
export function list(query) {
  return new Promise((resolve, reject) => {
    return LoginInformListHandler(query).then((res) => {
      console.debug('LoginInformListHandler', res)
      resolve(res)
    }).catch((err) => {
      reject(err)
    });
  })
}

// 删除登录日志
export function delLogininfor(infoId) {
  return DeleteByIdHandler(infoId)
}

// 解锁用户登录状态
export function unlockLogininfor(userName) {
  return UnlockHandler(userName)
}

// 清空登录日志
export function cleanLogininfor() {
  return CleanHandler()
}
