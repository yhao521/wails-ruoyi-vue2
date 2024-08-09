// import request from '@/utils/request'

import { ListDept, ExcludeDept, GetDept, SaveDept, UpDataDept, DeleteDept } from 'wailsjs/go/system/DeptAPI'

// 查询部门列表
export function listDept(query) {
  console.debug(query)
  let queryParams = query;
  queryParams = {
    pageNum: query.pageNum,
    pageSize: query.pageSize,
    other: {
      deptName: query.deptName,
      status: query.status,
    }
  }
  return new Promise((resolve, reject) => {
    return ListDept(queryParams).then((res) => {
      console.debug(res)
      resolve(res)
    }).catch((err) => {
      reject(err)
    });
  })
}

// 查询部门列表（排除节点）
export function listDeptExcludeChild(deptId) {
  return ExcludeDept(deptId)
}

// 查询部门详细
export function getDept(deptId) {
  return GetDept(deptId)
}

// 新增部门
export function addDept(data) {
  return SaveDept(data)
}

// 修改部门
export function updateDept(data) {
  return UpDataDept(data)
}

// 删除部门
export function delDept(deptId) {
  return DeleteDept(deptId)
}