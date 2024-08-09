// import request from '@/utils/request'
import { ListPost, GetPostInfo, SavePost, UploadPost, DeletePost } from 'wailsjs/go/system/PostAPI'

// 查询岗位列表
export function listPost(query) {
  let queryParams = query;
  queryParams = {
    pageNum: query.pageNum,
    pageSize: query.pageSize,
    other: {
      postName: query.postName,
      postCode: query.postCode,
      status: query.status,
    }
  }
  return new Promise((resolve, reject) => {
    return ListPost(queryParams).then((res) => {
      resolve(res)
    }).catch((err) => {
      reject(err)
    });
  })
}

// 查询岗位详细
export function getPost(postId) {
  return GetPostInfo(postId)
}

// 新增岗位
export function addPost(data) {
  return SavePost(data)
}

// 修改岗位
export function updatePost(data) {
  return UploadPost(data)
}

// 删除岗位
export function delPost(postId) {
  return DeletePost(postId)
}
