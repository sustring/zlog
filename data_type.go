package zlog

// QueueMaxNumber 日志缓冲队列组大长度
var QueueMaxNumber = 1000000 // 100w

// LevelMap 日志级别码表
var LevelMap = make(map[int]string)

// 日志级别
const (
	Debug    = 0x0001 // 1
	Info     = 0x0002 // 2
	Notice   = 0x0004 // 4
	Warn     = 0x0008 // 8
	Error    = 0x0010 // 16
	Critical = 0x0020 // 32

	All = 0x003f // 63
)

// 日志的相关默认配置参数
const (
	ModulesAll  = "all"  // display all modules log
	ModulesNone = "none" // hide all modules log

	DefaultLogFilePath   = "/tmp"
	DefaultLogFilePrefix = "server"

	DefaultSplitPeriod       = 60 * 60 * 24
	DefaultSplitSize         = 1024 * 1024 * 10
	DefaultLogFileSavePeriod = 7
)

// Printer 日志执行对象接口
type Printer interface {
	Printf(level int, module string, format string, v ...interface{})
}

type oneLog struct {
	level      int
	module     string
	callerFile string
	callerLine int
	callerName string
	format     string
	args       []interface{}
}

func initLevelMap() {
	LevelMap[Debug] = "[ DEBUG  ] - "
	LevelMap[Info] = "[ INFO   ] - "
	LevelMap[Notice] = "[ NOTICE ] - "
	LevelMap[Warn] = "[ WARN   ] - "
	LevelMap[Error] = "[ ERROR  ] - "
	LevelMap[Critical] = "[CRITICAL] - "
}
