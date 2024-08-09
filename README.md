# README
 wails-demo  加入了 ruoyi-go ruoyi-vue

## About

This is the official Wails Vue template.

You can configure the project by editing `wails.json`. More information about the project settings can be found
here: https://wails.io/docs/reference/project-config

## Live Development

To run in live development mode, run `wails dev` in the project directory. This will run a Vite development
server that will provide very fast hot reload of your frontend changes. If you want to develop in a browser
and have access to your Go methods, there is also a dev server that runs on http://localhost:34115. Connect
to this in your browser, and you can call your Go code from devtools.

## Building

To build a redistributable, production mode package, use `wails build`.
>Linux 或 Mac 系统构建时，须要保证目录文件有操作权限。模板中 Windows 系统相关代码将会失效，如运行时Go相关代码报错，请移除这些代码后再进行编译。
```shell
#如果 frontend/wailsjs 目录不存在，先执行如下命令
wails generate module

#如果首次打包提示 dist目录不存在，先执行下面命令
cd frontend
npm i
npm run build

#检测
wails doctor

#调试
wails dev

#打包 .exe
wails build 

#打包 .exe 带安装步骤
wails build -nsis -upx

# mac编译ext 
./build-win.sh
# mac打包app
./build-mac.sh
```


## 功能列表




| 功能     | 说明                                                         |
| :------: | :----------------------------------------------------------- |
| 代码编辑器 | 表单构建中复杂编辑器 |
|   表单构建       |       拖动表单元素生成相应的  vue +element-ui代码。     |
|代码生成|前后端代码的生成（java、html、xml、sql）支持CRUD下载 。 |
| 用户管理 | 用户是系统操作者，该功能主要完成系统用户配置。               |
| 部门管理 | 配置系统组织机构（公司、部门、小组），树结构展现支持数据权限。 |
| 岗位管理 | 配置系统用户所属担任职务。                                   |
|菜单管理|配置系统菜单，操作权限，按钮权限标识等。 |
|角色管理|角色菜单权限分配、设置角色按机构进行数据范围权限划分。 |
|字典管理|对系统中经常使用的一些较为固定的数据进行维护。 |
|参数管理|对系统动态配置常用参数。 |
|通知公告|系统通知公告信息发布维护。 |
|操作日志|系统正常操作日志记录和查询；系统异常信息日志记录和查询。 |
|登录日志|系统登录日志记录查询包含登录异常。 |
|服务监控|监视当前系统CPU、内存、磁盘、堆栈等相关信息。 |

感谢以下开源项目
https://gitee.com/xbuntu/godesk

 ruoyi-go

 https://gitee.com/y_project/RuoYi-Vue