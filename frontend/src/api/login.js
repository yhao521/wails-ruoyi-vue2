import request from '@/utils/request'

import { CaptchaImageHandler, LoginHandler, GetInfoHandler, LogoutHandler, GetCookie } from 'wailsjs/go/system/UserAPI'

// 登录方法
export function login(username, password, code, uuid) {
  // const data = {
  //   username,
  //   password,
  //   code,
  //   uuid
  // }
  return LoginHandler(username, password, code, uuid)
}

// 注册方法
export function register(data) {
  return request({
    url: '/register',
    headers: {
      isToken: false
    },
    method: 'post',
    data: data
  })
}

// 获取用户详细信息
export function getCookies() {
  return GetCookie()
}
// 获取用户详细信息
export function getInfo() {
  return GetInfoHandler()
}

// 退出方法
export function logout() {
  return LogoutHandler()
}

// 获取验证码
export function getCodeImg() {
  return new Promise((resolve, reject) => {
    return CaptchaImageHandler().then((res) => {
      console.debug(res)
      resolve(res.data)
    }).catch((err) => {
      reject(err)
    });
  })
}