import { ListOperlog, DelectOperlog, ClearOperlog } from "wailsjs/go/monitor/OperLogAPI";

// 查询操作日志列表
export function list(query) {
  return new Promise((resolve, reject) => {
    return ListOperlog(query).then((res) => {
      console.debug('ListOperlog', res)
      resolve(res)
    }).catch((err) => {
      reject(err)
    });
  })
}

// 删除操作日志
export function delOperlog(operId) {
  return DelectOperlog(operId)
}

// 清空操作日志
export function cleanOperlog() {
  return ClearOperlog()
}
