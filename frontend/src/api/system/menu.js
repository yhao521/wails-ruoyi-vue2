// import request from '@/utils/request'
import { GetRoutersHandler, ListMenu, GetMenuInfo, DeleteMenu, SaveMenu, UpdateMenu, GetTreeSelect, TreeSelectByRole } from 'wailsjs/go/system/MenuAPI'

// 获取路由
export const getRouters = () => {
  return GetRoutersHandler()
}
// 查询菜单列表
export function listMenu(query) {
  console.debug('listMenu-query:', query)
  let queryParams = query;
  if (query != undefined) {
    queryParams = {
      other: {
        menuName: query.menuName,
        visible: query.visible,
      }
    }
  }
  console.debug('listMenu-queryParams:', query)
  return new Promise((resolve, reject) => {
    return ListMenu(queryParams).then((res) => {
      resolve(res.data)
    }).catch((err) => {
      reject(err)
    });
  })
}

// 查询菜单详细
export function getMenu(menuId) {
  return new Promise((resolve, reject) => {
    return GetMenuInfo(menuId).then((res) => {
      resolve(res.data)
    }).catch((err) => {
      reject(err)
    });
  })
}

// 查询菜单下拉树结构
export function treeselect() {
  return GetTreeSelect()
  // return request({
  //   url: '/system/menu/treeselect',
  //   method: 'get'
  // })
}

// 根据角色ID查询菜单下拉树结构
export function roleMenuTreeselect(roleId) {
  return TreeSelectByRole(roleId)
  // return request({
  //   url: '/system/menu/roleMenuTreeselect/' + roleId,
  //   method: 'get'
  // })
}

// 新增菜单
export function addMenu(data) {
  return SaveMenu(data)
}

// 修改菜单
export function updateMenu(data) {
  return UpdateMenu(data)
}

// 删除菜单
export function delMenu(menuId) {
  return DeleteMenu(menuId)
}