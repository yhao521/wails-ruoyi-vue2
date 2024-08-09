// import request from '@/utils/request'

import { ListNotice, GetNotice, SaveNotice, UploadNotice, DeleteNotice } from 'wailsjs/go/system/NoticeAPI'
// 查询公告列表
export function listNotice(query) {
  let queryParams = query;
  queryParams = {
    pageNum: query.pageNum,
    pageSize: query.pageSize,
    other: {
      noticeTitle: query.noticeTitle,
      createBy: query.createBy,
      noticeType: query.noticeType,
    }
  }
  return new Promise((resolve, reject) => {
    return ListNotice(queryParams).then((res) => {
      resolve(res)
    }).catch((err) => {
      reject(err)
    });
  })
}

// 查询公告详细
export function getNotice(noticeId) {
  return GetNotice(noticeId)
}

// 新增公告
export function addNotice(data) {
  return SaveNotice(data)
}

// 修改公告
export function updateNotice(data) {
  return UploadNotice(data)
}

// 删除公告
export function delNotice(noticeId) {
  return DeleteNotice(noticeId)
}