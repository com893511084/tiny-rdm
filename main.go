package main

import (
	"context"
	"embed"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/linux"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
	runtime2 "github.com/wailsapp/wails/v2/pkg/runtime"
	"runtime"
	"tinyrdm/backend/consts"
	"tinyrdm/backend/services"
)

//go:embed all:frontend/dist
var assets embed.FS

//go:embed build/appicon.png
var icon []byte

var version = "0.0.0"
var gaMeasurementID, gaSecretKey string

func main() {
	// Create an instance of the app structure
	sysSvc := services.System()
	connSvc := services.Connection()
	browserSvc := services.Browser()
	cliSvc := services.Cli()
	monitorSvc := services.Monitor()
	pubsubSvc := services.Pubsub()
	prefSvc := services.Preferences()
	prefSvc.SetAppVersion(version)
	// 读取窗口大小
	windowWidth, windowHeight, maximised := prefSvc.GetWindowSize()
	windowStartState := options.Normal
	if maximised {
		windowStartState = options.Maximised
	}

	// 创建 menu 菜单
	appMenu := menu.NewMenu()
	if runtime.GOOS == "darwin" {
		appMenu.Append(menu.AppMenu())
		appMenu.Append(menu.EditMenu())
		appMenu.Append(menu.WindowMenu())
	}

	// 使用选项创建应用程序
	err := wails.Run(&options.App{
		// 标题
		Title: "Tiny RDM",
		// 窗口的初始宽度。
		Width: windowWidth,
		// 窗口的初始高度。
		Height: windowHeight,
		// 最小宽度
		MinWidth: consts.MIN_WINDOW_WIDTH,
		// 最小高度
		MinHeight: consts.MIN_WINDOW_HEIGHT,
		// 窗口启动状态 全屏, 最大化,最小化
		WindowStartState: windowStartState,
		// 是否无边框
		Frameless: runtime.GOOS != "darwin",
		// 应用程序要使用的菜单。 菜单参考 中有关菜单的更多详细信息。
		Menu: appMenu,
		// 在生产环境中启用浏览器的默认上下文菜单。
		EnableDefaultContextMenu: true,
		// 资源服务
		// 这定义了资产服务特定的选项。
		// 它允许使用静态资产自定义资产服务，
		// 使用 http.Handler 动态地提供资产或使用 assetsserver.Middleware 钩到请求链。
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		// 窗口的默认背景颜色
		BackgroundColour: options.NewRGBA(27, 38, 54, 0),
		// 启动时隐藏窗口 设置为 true 时，应用程序将被隐藏，直到调用显示窗口。
		StartHidden: true,
		// 启动时调用的方法 此回调在前端创建之后调用，但在 index.html 加载之前调用。 它提供了应用程序上下文。
		OnStartup: func(ctx context.Context) {
			sysSvc.Start(ctx, version)
			connSvc.Start(ctx)
			browserSvc.Start(ctx)
			cliSvc.Start(ctx)
			monitorSvc.Start(ctx)
			pubsubSvc.Start(ctx)

			services.GA().SetSecretKey(gaMeasurementID, gaSecretKey)
			services.GA().Startup(version)
		},
		// 在前端加载完毕 index.html 及其资源后调用此回调。 它提供了应用程序上下文。
		OnDomReady: func(ctx context.Context) {
			x, y := prefSvc.GetWindowPosition(ctx)
			runtime2.WindowSetPosition(ctx, x, y)
			runtime2.WindowShow(ctx)
		},
		// 应用关闭前回调
		OnBeforeClose: func(ctx context.Context) (prevent bool) {
			x, y := runtime2.WindowGetPosition(ctx)
			prefSvc.SaveWindowPosition(x, y)
			return false
		},
		// 应用退出回调
		OnShutdown: func(ctx context.Context) {
			browserSvc.Stop()
			cliSvc.CloseAll()
			monitorSvc.StopAll()
			pubsubSvc.StopAll()
		},
		// 定义需要绑定到前端的方法的结构实例切片。
		Bind: []interface{}{
			sysSvc,
			connSvc,
			browserSvc,
			cliSvc,
			monitorSvc,
			pubsubSvc,
			prefSvc,
		},
		// MAC 的特定选项
		Mac: &mac.Options{
			TitleBar: mac.TitleBarHiddenInset(),
			About: &mac.AboutInfo{
				Title:   "Tiny RDM " + version,
				Message: "A modern lightweight cross-platform Redis desktop client.\n\nCopyright © 2024",
				Icon:    icon,
			},
			WebviewIsTransparent: false,
			WindowIsTranslucent:  false,
			DisableZoom:          true,
		},
		// Windows 的特定选项
		Windows: &windows.Options{
			WebviewIsTransparent:              true,
			WindowIsTranslucent:               true,
			DisableFramelessWindowDecorations: true,
		},
		// Linux 的特定选项
		Linux: &linux.Options{
			ProgramName:         "Tiny RDM",
			Icon:                icon,
			WebviewGpuPolicy:    linux.WebviewGpuPolicyOnDemand,
			WindowIsTranslucent: true,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
