// import request from '@/utils/request'

import { DictTypeHandler, ListDict, GetDictCode, DeleteDictData, SaveDictData, UpDictData } from 'wailsjs/go/system/DictDataAPI'

// 查询字典数据列表
export function listData(query) {
  let queryParams = query;
  queryParams = {
    pageNum: query.pageNum,
    pageSize: query.pageSize,
    other: {
      dictType: query.dictType,
      dictLabel: query.dictLabel,
      status: query.status,
    }
  }
  return new Promise((resolve, reject) => {
    return ListDict(queryParams).then((res) => {
      console.debug('ListDict', res)
      resolve(res.data)
    }).catch((err) => {
      reject(err)
    });
  })
}

// 查询字典数据详细
export function getData(dictCode) {
  return GetDictCode(dictCode)
}

// 根据字典类型查询字典数据信息
export function getDicts(dictType) {
  return new Promise((resolve, reject) => {
    return DictTypeHandler(dictType).then((res) => {
      resolve(res.data)
    }).catch((err) => {
      reject(err)
    });
  })
}

// 新增字典数据
export function addData(data) {
  return SaveDictData(data)
}

// 修改字典数据
export function updateData(data) {
  return UpDictData(data)
}

// 删除字典数据
export function delData(dictCode) {
  return DeleteDictData(dictCode)
}
