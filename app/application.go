package app

import (
	"pk/components/uploader"
	"pk/components/xlsxer"
	"pk/engine/engines"
	"sync"
)

// 应用程序实例
var applicationOnce sync.Once
var applicationInstance *Application
// 文件处理引擎实例
var fileEngine *uploader.FileUploadReceiver

func init() {
	// 初始化文件处理引擎
	fileEngineInstance, err := uploader.New("./assets/storage")
	if err != nil {
		panic(err)
	}
	fileEngine = fileEngineInstance
	fileEngine.SetFailureAllowed(true)
}

// 应用程序实例结构
type Application struct {
	FileEngine					*uploader.FileUploadReceiver // 文件处理引擎
	PlayerEngine				*engines.PlayerEngine          // 玩家引擎
	UnitEngine 					*engines.UnitEngine          // 单位引擎
	XlsxEngine 					*xlsxer.Xlsxer               // 表格文件处理引擎
}

// 获取到当前应用程序实例
func Get() *Application {
	applicationOnce.Do(func() {
		applicationInstance = &Application{
			FileEngine:   fileEngine,
			PlayerEngine: engines.NewPlayerEngine(),
			UnitEngine:   engines.NewUnitEngine(),
			XlsxEngine:   xlsxer.NewXlsxer(),
		}
	})
	return applicationInstance
}