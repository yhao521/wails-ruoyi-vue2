package baseAPI

import (
	"bufio"
	"context"
	"embed"
	"fmt"
	"io"
	"log"
	"mySparkler/backend/model"
	"mySparkler/backend/model/monitor"
	"mySparkler/backend/model/system"
	"mySparkler/backend/model/tools"
	"mySparkler/pkg/db"
	"mySparkler/pkg/file"
	"mySparkler/pkg/utils/R"
	"mySparkler/pkg/ylog"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strings"

	runtime2 "github.com/wailsapp/wails/v2/pkg/runtime"
)

// Base 控制器基类
type Base struct {
	ctx context.Context
}

func NewBase() *Base {
	return &Base{}
}

const (
	SuccessCode = 200
	ErrorCode   = 101
	NoLoginCode = 501
)

// Res 返回结果
type Res struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// res 返回
func (b *Base) Res(code int, msg string, data interface{}) Res {
	return Res{
		Code: code,
		Msg:  msg,
		Data: data,
	}
}

// success 返回成功
func (b *Base) Success(data interface{}) Res {
	return b.Res(SuccessCode, "操作成功", data)
}

// error 返回错误
func (b *Base) Error(message string) Res {
	return b.Res(ErrorCode, message, nil)
}

func (b *Base) GetCtx() context.Context {
	return b.ctx
}

// setCtx 设置上下文对象
func (b *Base) SetCtx(ctx context.Context) {
	b.ctx = ctx
}

// log 增加日志记录
func (b *Base) Log(content any) {
	ylog.Log(content)
}

// schema 根据model自动建立数据表
func (b *Base) Schema(sqlfiles embed.FS) {
	db := db.Dbp(b.GetAppPath())
	db.AutoMigrate(
		&model.Config{},
		&system.SysConfig{},
		&system.SysDept{},
		&system.SysDictData{},
		&system.SysDictType{},
		&system.SysMenu{},
		&system.SysNotice{},

		&system.SysPost{},
		&system.SysRoleDept{},
		&system.SysRoleMenu{},
		&system.SysRoles{},

		&system.SysUser{},
		&system.SysUserPost{},
		&system.SysUserRole{},

		&monitor.SysJob{},
		&monitor.SysJobLog{},
		&monitor.SysLogininfor{},
		&monitor.SysOperLog{},

		&tools.GenTable{},
		&tools.GenTableColumn{},
	)
	res := system.SelectConfigByKey("sys.data.init")

	b.Log("init-data")
	b.Log(res)

	if res != "true" {
		//数据初始化
		data, _ := sqlfiles.ReadFile(path.Join("sql", "sqlite_data.sql"))
		sqlarr := strings.Split(string(data), ";")
		for _, v := range sqlarr {
			if v == "" {
				continue
			}
			db.Exec(v)
		}
	}

}

// getAppPath 获取应用主目录
func (b *Base) GetAppPath() string {
	//获取我的文档目录
	return file.GetAppPath()
}

// pathExist 判断文件目录是否存在，不存在创建
func (b Base) PathExist(path string) string {
	_, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		_ = os.MkdirAll(path, os.ModePerm)
	}
	return path
}

// shellCMD 以shell方式运行cmd命令
func (b *Base) Command(cmdStr string, paramStr string) {
	exec.Command(cmdStr, paramStr).Start()
}

// setSystemBackground 设置系统壁纸
func (b *Base) SetSystemBackground(path string) {
	// 设置壁纸
	cmd := exec.Command("gsettings", "set", "org.gnome.desktop.background", "picture-uri", "file://"+path)
	err := cmd.Run()
	if err != nil {
		log.Fatalf("Error setting wallpaper: %s", err)
	}
}

// httpGet get请求
func (b *Base) HttpGet(url string) []byte {
	res, err := http.Get(url)
	if err != nil {
		panic(any("获取壁纸失败"))
	}
	defer res.Body.Close()
	bytes, _ := io.ReadAll(res.Body)
	return bytes
}

func (b *Base) OpenMacDir(path string) {
	b.Command("open", path)
}

// openDir 打开ubuntu目录
func (b *Base) OpenDir(path string) {
	b.Command("xdg-open", fmt.Sprintf("file://%s", path))
}

// SelectFile 选择需要处理的文件
func (b *Base) SelectFile(title string, filetype string) string {
	if title == "" {
		title = "选择文件"
	}
	if filetype == "" {
		filetype = "*.txt;*.json"
	}
	selection, err := runtime2.OpenFileDialog(b.ctx, runtime2.OpenDialogOptions{
		Title: title,
		Filters: []runtime2.FileFilter{
			{
				DisplayName: "文本数据",
				Pattern:     filetype,
			},
		},
	})
	if err != nil {
		return fmt.Sprintf("err %s!", err)
	}
	return selection
}

// SaveFile 选择需要处理的文件
func (b *Base) SaveFilePath(fileName string, filetype string) string {

	if filetype == "" {
		filetype = "*.txt;*.json;*.zip"
	}
	filepath, err := runtime2.SaveFileDialog(b.ctx, runtime2.SaveDialogOptions{
		ShowHiddenFiles: true,
		DefaultFilename: fileName,
		Title:           "选择文件",
		Filters: []runtime2.FileFilter{
			{
				DisplayName: "文本数据",
				Pattern:     filetype,
			},
		},
	})
	if err != nil {
		return fmt.Sprintf("err %s!", err)
	}

	return filepath
}

// SaveFile 选择需要处理的文件
func (b *Base) SaveFile(data string, fileName string, filetype string) (resp R.Result) {

	if filetype == "" {
		filetype = "*.txt;*.json;*.zip"
	}
	filepath, err := runtime2.SaveFileDialog(b.ctx, runtime2.SaveDialogOptions{
		ShowHiddenFiles: true,
		DefaultFilename: fileName,
		Title:           "选择文件",
		Filters: []runtime2.FileFilter{
			{
				DisplayName: "文本数据",
				Pattern:     filetype,
			},
		},
	})
	if err != nil {
		resp.Msg = fmt.Sprintf("err %s!", err)
		resp.Code = 500
		return
	}

	file, err := os.Create(filepath)
	if err != nil {
		resp.Msg = err.Error()
		return
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	_, _ = writer.WriteString(data)

	writer.Flush()
	resp.Data = filepath
	resp.Code = 200
	return
}
