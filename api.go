package zlog

func example() {

}

// Config 日志配置.
type Config struct {
	LogFilePath   string   // 日志文件路径(默认：/tmp,即系统/tmp目录下).
	LogLevels     int      // 将要记录下来的日志的日志等级(默认：63 ,即全部).
	LogFilePrefix string   // 日志文件前缀(默认：server).
	Modules       []string // 将要记录下来的日志的模块(默认：all).
	IsTimeSplit   bool     // 是否周期性的分割日志文件(默认：true).
	SplitPeriod   int64    // 周期性分割日志文件的周期，单位－秒(默认：60*60*24).
	IsSizeSplit   bool     // 是否根据日志文件大小分割日志文件(默认：true).
	SplitSize     int64    // 日志文件超过此大小进行分割，单位－byte(默认：1024*1024*10).
	IsClear       bool     // 是否开启清理日志文件功能(默认：true).
	SavePeriod    int64    // 日志文件保存周期，单位 - 天(默认：7).
}

// Init 初始化日志管理器，否则将使用默认配置.
func Init(config Config) {
	getScreenLogImplInstance().SetScreenLogConfig(config)
	getSysFileLogInstance().SetFileLogConfig(config)
	return
}

// Prints 新增一条屏幕打印日志.
func Prints(level int, module string, format string, v ...interface{}) {
	screenAddOneLog().Printf(level, module, format, v...)
	return
}

// Printf 新增一条文件日志.
func Printf(level int, module string, format string, v ...interface{}) {
	fileAddOneLog().Printf(level, module, format, v...)
	return
}

func init() {
	initLevelMap()
	initScreenLoger()
	initFileLoger()
}
