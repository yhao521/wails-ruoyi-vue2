package run

import (
	"context"
	"embed"
	"runtime"

	"mySparkler/backend/api/baseAPI"
	"mySparkler/backend/api/monitor"
	"mySparkler/backend/api/system"
	"mySparkler/backend/api/tools"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/linux"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

var icon []byte

// WailsRun 初始化
func WailsRun(assets embed.FS, ymlDefault string, sqlfiles embed.FS, templates embed.FS) {
	// 创建控制器实例
	// Create an instance of the app structure
	app := baseAPI.NewApp()
	user := system.NewUserAPI()
	menuAPI := system.NewMenuAPI()
	dictDataAPI := system.NewDictDataAPI()
	genApi := tools.NewGenAPI()
	configApi := system.NewConfigAPI()
	deptApi := system.NewDeptAPI()
	roleAPI := system.NewRoleAPI()
	noticeAPI := system.NewNoticeAPI()
	postAPI := system.NewPostAPI()
	operLogAPI := monitor.NewOperLogAPI()
	serverAPI := monitor.NewServerAPI()

	loginInforAPI := monitor.NewLoginInforAPI()

	// Create application with options
	err := wails.Run(&options.App{
		Title:         "我的工具",
		Width:         1300,
		Height:        768,
		DisableResize: false,
		// Fullscreen:    false,

		// StartHidden:       false,
		// HideWindowOnClose: false,
		// AlwaysOnTop:       false,
		Frameless: runtime.GOOS != "darwin",
		// Frameless:        false,
		EnableDefaultContextMenu: true,
		BackgroundColour:         &options.RGBA{R: 255, G: 255, B: 255, A: 1},
		// BackgroundColour:  &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		AssetServer:        &assetserver.Options{Assets: assets},
		Menu:               app.ApplicationMenu(),
		Logger:             nil,
		LogLevel:           logger.INFO,
		LogLevelProduction: logger.ERROR,
		OnStartup: func(ctx context.Context) {
			app.Startup(ctx, ymlDefault)
			//设置 context 对象
			user.SetCtx(ctx)
			menuAPI.SetCtx(ctx)
			dictDataAPI.SetCtx(ctx)
			genApi.SetCtx(ctx, templates)
			configApi.SetCtx(ctx)
			deptApi.SetCtx(ctx)
			roleAPI.SetCtx(ctx)
			noticeAPI.SetCtx(ctx)
			postAPI.SetCtx(ctx)

			operLogAPI.SetCtx(ctx)
			serverAPI.SetCtx(ctx)
			loginInforAPI.SetCtx(ctx)


			menuAPI.Schema(sqlfiles)
		},
		OnDomReady: func(ctx context.Context) {
		},
		OnShutdown: func(ctx context.Context) {
			app.Shutdown(ctx)
		},
		// OnBeforeClose: func(ctx context.Context) (prevent bool) {
		// 	dialog, err := runtime.MessageDialog(ctx, runtime.MessageDialogOptions{Type: runtime.QuestionDialog, Title: "退出?", Message: "你确定要退出吗?"})
		// 	println("dialog:", dialog)
		// 	if err != nil {
		// 		return false
		// 	}
		// 	return dialog != "Yes"
		// },
		Bind: []interface{}{
			app,
			user,
			menuAPI,
			dictDataAPI,
			genApi,
			configApi,
			deptApi,
			roleAPI,
			postAPI,
			noticeAPI,
			operLogAPI,
			serverAPI,
			loginInforAPI,
		},
		// EnumBind:         []interface{}{},
		// WindowStartState: 0,
		ErrorFormatter: func(err error) any {
			return err.Error()
		},
		// CSSDragProperty:                  "",
		// CSSDragValue:                     "",
		// EnableDefaultContextMenu:         false,
		// EnableFraudulentWebsiteDetection: false,
		// SingleInstanceLock:               &options.SingleInstanceLock{},
		Windows: &windows.Options{
			// Webview 透明
			WebviewIsTransparent: false,
			// // 窗口半透明 true
			WindowIsTranslucent:               false,
			DisableFramelessWindowDecorations: true,

			BackdropType:        windows.Mica,
			DisableWindowIcon:   false,
			WebviewUserDataPath: "",
			WebviewBrowserPath:  "",
			Theme:               windows.SystemDefault,
			CustomTheme: &windows.ThemeSettings{
				DarkModeTitleBar:   windows.RGB(20, 20, 20),
				DarkModeTitleText:  windows.RGB(200, 200, 200),
				DarkModeBorder:     windows.RGB(20, 0, 20),
				LightModeTitleBar:  windows.RGB(200, 200, 200),
				LightModeTitleText: windows.RGB(20, 20, 20),
				LightModeBorder:    windows.RGB(200, 200, 200),
			},
			// User messages that can be customised
			// Messages * windows.Messages
			// OnSuspend is called when Windows enters low power mode
			// OnSuspend func(),
			// OnResume is called when Windows resumes from low power mode
			// OnResume func(),
			// WebviewGpuDisabled: false,

		},
		Mac: &mac.Options{
			// TitleBar: &mac.TitleBar{
			// 	TitlebarAppearsTransparent: false,
			// 	HideTitle:                  false,
			// 	HideTitleBar:               false,
			// 	FullSizeContent:            false,
			// 	UseToolbar:                 false,
			// 	HideToolbarSeparator:       false,
			// },
			TitleBar:   mac.TitleBarDefault(),
			Appearance: mac.NSAppearanceNameDarkAqua,
			// Webview 透明
			WebviewIsTransparent: false,
			// 窗口半透明 true
			WindowIsTranslucent: false,
			About: &mac.AboutInfo{
				Title:   "我的工具箱",
				Message: "© 2023 Me",
				Icon:    icon,
			},
			Preferences: &mac.Preferences{
				TabFocusesLinks:        mac.Enabled,
				TextInteractionEnabled: mac.Enabled,
				FullscreenEnabled:      mac.Enabled,
			},
		},
		Linux: &linux.Options{
			ProgramName:         "我的工具箱",
			Icon:                icon,
			WebviewGpuPolicy:    linux.WebviewGpuPolicyOnDemand,
			WindowIsTranslucent: true,
		},
		Experimental: &options.Experimental{},
		Debug: options.Debug{
			// 设置为 true 将在应用程序启动时打开 Web 检查器。
			OpenInspectorOnStartup: true,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}

}
