
import { ServerData } from "wailsjs/go/monitor/ServerAPI";

// 获取服务信息
export function getServer() {
  return new Promise((resolve, reject) => {
    return ServerData().then((res) => {
      console.debug('ServerData', res)
      resolve(res)
    }).catch((err) => {
      reject(err)
    });
  })
}