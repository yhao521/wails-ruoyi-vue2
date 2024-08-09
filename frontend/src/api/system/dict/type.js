// import request from '@/utils/request'

import { ListDictType, GetTypeDict, DeleteDataType, SaveType, UpdateType, RefreshCache, GetOptionSelect } from 'wailsjs/go/system/DictDataAPI'

// 查询字典类型列表
export function listType(query) {
  return new Promise((resolve, reject) => {
    return ListDictType(query).then((res) => {
      resolve(res.data.data)
    }).catch((err) => {
      reject(err)
    });
  })
}

// 查询字典类型详细
export function getType(dictId) {
  return new Promise((resolve, reject) => {
    return GetTypeDict(dictId).then((res) => {
      resolve(res.data)
    }).catch((err) => {
      reject(err)
    });
  })
}

// 新增字典类型
export function addType(data) {
  return SaveType(data)
}

// 修改字典类型
export function updateType(data) {
  return UpdateType(data)
}

// 删除字典类型
export function delType(dictId) {
  return DeleteDataType(dictId)
}

// 刷新字典缓存
export function refreshCache() {
  return RefreshCache()
  // return request({
  //   url: '/system/dict/type/refreshCache',
  //   method: 'delete'
  // })
}

// 获取字典选择框列表
export function optionselect() {
  return GetOptionSelect()
  // return request({
  //   url: '/system/dict/type/optionselect',
  //   method: 'get'
  // })
}