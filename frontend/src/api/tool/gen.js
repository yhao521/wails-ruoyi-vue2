// import request from '@/utils/request'
import { GenList, GenDbList, PreviewGenTable, GenDelete, GenBatchCode, SynchDb, ImportTable, GenInfo, GenEdit, CreateTable } from "wailsjs/go/tools/GenAPI";

// 查询生成表数据
export function listTable(query) {
  return new Promise((resolve, reject) => {
    return GenList(query).then((res) => {
      console.debug('GenList', res)
      resolve(res.data)
    }).catch((err) => {
      reject(err)
    });
  })
  // return request({
  //   url: '/tool/gen/list',
  //   method: 'get',
  //   params: query
  // })
}
// 查询db数据库列表
export function listDbTable(query) {

  return new Promise((resolve, reject) => {
    return GenDbList(query).then((res) => {
      console.debug('GenDbList', res)
      resolve(res.data)
    }).catch((err) => {
      reject(err)
    });
  })
}

// 查询表详细信息
export function getGenTable(tableId) {
  return GenInfo(tableId)
  // return request({
  //   url: '/tool/gen/' + tableId,
  //   method: 'get'
  // })
}

// 修改代码生成信息
export function updateGenTable(data) {
  return GenEdit(data)
  // return request({
  //   url: '/tool/gen',
  //   method: 'put',
  //   data: data
  // })
}

// 导入表
export function importTable(data) {
  return ImportTable(data.tables)
  // return request({
  //   url: '/tool/gen/importTable',
  //   method: 'post',
  //   params: data
  // })
}

// 创建表
export function createTable(data) {
  return CreateTable(data)
}

// 预览生成代码
export function previewTable(tableId) {
  return PreviewGenTable(tableId)
  // return request({
  //   url: '/tool/gen/preview/' + tableId,
  //   method: 'get'
  // })
}

// 删除表数据
export function delTable(tableId) {
  return GenDelete(tableId)
  // return request({
  //   url: '/tool/gen/' + tableId,
  //   method: 'delete'
  // })
}

// 生成代码（自定义路径）
export function genCode(tableName, filePath) {
  return GenBatchCode(tableName, filePath)
  // return request({
  //   url: '/tool/gen/genCode/' + tableName,
  //   method: 'get'
  // })
}

// 同步数据库
export function synchDb(tableName) {
  return SynchDb(tableName)
  // return request({
  //   url: '/tool/gen/synchDb/' + tableName,
  //   method: 'get'
  // })
}
