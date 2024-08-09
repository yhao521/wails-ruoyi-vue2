// import request from '@/utils/request'

import { ListConfig, GetConfigKey, UploadConfig, DetectConfig, SaveConfig, GetConfigInfo, DeleteCacheConfig } from 'wailsjs/go/system/ConfigAPI.js'

// 查询参数列表
export function listConfig(query) {
  let queryParams = query;
  queryParams = {
    pageNum: query.pageNum,
    pageSize: query.pageSize,
    other: {
      configName: query.configName,
      configKey: query.configKey,
      configType: query.configType
    }
  }
  return new Promise((resolve, reject) => {
    return ListConfig(queryParams).then((res) => {
      resolve(res)
    }).catch((err) => {
      reject(err)
    });
  })
}

// 查询参数详细
export function getConfig(configId) {
  return new Promise((resolve, reject) => {
    return GetConfigInfo(configId).then((res) => {
      resolve(res.data)
    }).catch((err) => {
      reject(err)
    });
  })
}

// 根据参数键名查询参数值
export function getConfigKey(configKey) {
  return GetConfigKey(configKey)
}

// 新增参数配置
export function addConfig(data) {
  return SaveConfig(data)
}

// 修改参数配置
export function updateConfig(data) {
  return UploadConfig(data)
}

// 删除参数配置
export function delConfig(configId) {
  return DetectConfig(configId)
}

// 刷新参数缓存
export function refreshCache() {
  return DeleteCacheConfig()
}
