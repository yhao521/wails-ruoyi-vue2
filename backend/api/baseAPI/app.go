package baseAPI

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"mySparkler/config"
	"mySparkler/pkg/utils/R"
	"os"
	"path/filepath"
	"runtime"

	"github.com/vrischmann/userdir"
	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/menu/keys"
	runtime2 "github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx context.Context
	// baseAPI.Base
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

var configFile = flag.String("f", "./config.yaml", "")

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) Startup(ctx context.Context, ymlDefault string) {
	a.ctx = ctx
	// 初始化配置文件
	config.InitAppConfig(*configFile, ymlDefault)
}
func (a *App) Shutdown(ctx context.Context) {
	a.ctx = ctx
}

// 自定义菜单
func (a *App) ApplicationMenu() *menu.Menu {

	AppMenu := menu.NewMenu()
	if runtime.GOOS == "darwin" {
		AppMenu.Append(menu.AppMenu())
	}
	FileMenu := AppMenu.AddSubmenu("文件")
	// FileMenu.AddText("&Open", keys.CmdOrCtrl("o"), openFile)
	FileMenu.AddSeparator()
	FileMenu.AddText("退出", keys.CmdOrCtrl("q"), func(_ *menu.CallbackData) {
		runtime2.Quit(a.ctx)
	})

	if runtime.GOOS == "darwin" {
		AppMenu.Append(menu.EditMenu()) // on macos platform, we should append EditMenu to enable Cmd+C,Cmd+V,Cmd+Z... shortcut
		AppMenu.Append(menu.WindowMenu())
	}

	return AppMenu
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

// SelectFile 选择需要处理的文件
func (a *App) SelectFile(title string, filetype string) string {
	if title == "" {
		title = "选择文件"
	}
	if filetype == "" {
		filetype = "*.txt;*.json"
	}
	selection, err := runtime2.OpenFileDialog(a.ctx, runtime2.OpenDialogOptions{
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

func (a *App) SelectPathDownload() string {
	filePath := filepath.Join(os.Getenv("HOME"), "Downloads")

	selection, err := runtime2.OpenDirectoryDialog(a.ctx, runtime2.OpenDialogOptions{
		Title:            "选择目录",
		DefaultDirectory: filePath,
	})
	if err != nil {
		return fmt.Sprintf("err %s!", err)
	}
	return selection
}

func (a *App) SelectPath(filePath string) string {
	if filePath == "" {
		filePath = userdir.GetDataHome()
	}
	selection, err := runtime2.OpenDirectoryDialog(a.ctx, runtime2.OpenDialogOptions{
		Title:            "选择目录",
		DefaultDirectory: filePath,
	})
	if err != nil {
		return fmt.Sprintf("err %s!", err)
	}
	return selection
}

// SaveFile 选择需要处理的文件
func (a *App) SaveFilePath(fileName string, filetype string) string {

	if filetype == "" {
		filetype = "*.txt;*.json;*.zip"
	}
	filepath, err := runtime2.SaveFileDialog(a.ctx, runtime2.SaveDialogOptions{
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
func (a *App) SaveFile(data string, fileName string, filetype string) (resp R.Result) {

	if filetype == "" {
		filetype = "*.txt;*.json;*.zip"
	}
	filepath, err := runtime2.SaveFileDialog(a.ctx, runtime2.SaveDialogOptions{
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
